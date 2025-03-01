package votecollector

import (
	"errors"
	"sync"
	"testing"

	"github.com/onflow/crypto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/atomic"

	"github.com/onflow/flow-go/consensus/hotstuff"
	"github.com/onflow/flow-go/consensus/hotstuff/committees"
	"github.com/onflow/flow-go/consensus/hotstuff/helper"
	mockhotstuff "github.com/onflow/flow-go/consensus/hotstuff/mocks"
	"github.com/onflow/flow-go/consensus/hotstuff/model"
	hotstuffvalidator "github.com/onflow/flow-go/consensus/hotstuff/validator"
	"github.com/onflow/flow-go/consensus/hotstuff/verification"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/local"
	modulemock "github.com/onflow/flow-go/module/mock"
	"github.com/onflow/flow-go/module/signature"
	"github.com/onflow/flow-go/utils/unittest"
)

func TestStakingVoteProcessor(t *testing.T) {
	suite.Run(t, new(StakingVoteProcessorTestSuite))
}

// StakingVoteProcessorTestSuite is a test suite that holds mocked state for isolated testing of StakingVoteProcessor.
type StakingVoteProcessorTestSuite struct {
	VoteProcessorTestSuiteBase

	processor       *StakingVoteProcessor
	allParticipants flow.IdentityList
}

func (s *StakingVoteProcessorTestSuite) SetupTest() {
	s.VoteProcessorTestSuiteBase.SetupTest()
	s.allParticipants = unittest.IdentityListFixture(14)
	s.processor = &StakingVoteProcessor{
		log:               unittest.Logger(),
		block:             s.proposal.Block,
		stakingSigAggtor:  s.stakingAggregator,
		onQCCreated:       s.onQCCreated,
		minRequiredWeight: s.minRequiredWeight,
		done:              *atomic.NewBool(false),
		allParticipants:   s.allParticipants,
	}
}

// TestInitialState tests that Block() and Status() return correct values after calling constructor
func (s *StakingVoteProcessorTestSuite) TestInitialState() {
	require.Equal(s.T(), s.proposal.Block, s.processor.Block())
	require.Equal(s.T(), hotstuff.VoteCollectorStatusVerifying, s.processor.Status())
}

// TestProcess_VoteNotForProposal tests that vote should pass to validation only if it has correct
// view and block ID matching proposal that is locked in StakingVoteProcessor
func (s *StakingVoteProcessorTestSuite) TestProcess_VoteNotForProposal() {
	err := s.processor.Process(unittest.VoteFixture(unittest.WithVoteView(s.proposal.Block.View)))
	require.ErrorAs(s.T(), err, &VoteForIncompatibleBlockError)
	require.False(s.T(), model.IsInvalidVoteError(err))

	err = s.processor.Process(unittest.VoteFixture(unittest.WithVoteBlockID(s.proposal.Block.BlockID)))
	require.ErrorAs(s.T(), err, &VoteForIncompatibleViewError)
	require.False(s.T(), model.IsInvalidVoteError(err))

	s.stakingAggregator.AssertNotCalled(s.T(), "Verify")
}

// TestProcess_InvalidSignature tests that StakingVoteProcessor doesn't collect signatures for votes with invalid signature.
// Checks are made for cases where both staking and threshold signatures were submitted.
func (s *StakingVoteProcessorTestSuite) TestProcess_InvalidSignature() {
	exception := errors.New("unexpected-exception")

	// sentinel error from `InvalidSignerError` should be wrapped as `InvalidVoteError`
	voteA := unittest.VoteForBlockFixture(s.proposal.Block, unittest.VoteWithStakingSig())
	s.stakingAggregator.On("Verify", voteA.SignerID, mock.Anything).Return(model.NewInvalidSignerErrorf("")).Once()
	err := s.processor.Process(voteA)
	require.Error(s.T(), err)
	require.True(s.T(), model.IsInvalidVoteError(err))
	require.True(s.T(), model.IsInvalidSignerError(err))

	// sentinel error from `ErrInvalidSignature` should be wrapped as `InvalidVoteError`
	voteB := unittest.VoteForBlockFixture(s.proposal.Block, unittest.VoteWithStakingSig())
	s.stakingAggregator.On("Verify", voteB.SignerID, mock.Anything).Return(model.ErrInvalidSignature).Once()
	err = s.processor.Process(voteB)
	require.Error(s.T(), err)
	require.True(s.T(), model.IsInvalidVoteError(err))
	require.ErrorAs(s.T(), err, &model.ErrInvalidSignature)

	// unexpected errors from `Verify` should be propagated, but should _not_ be wrapped as `InvalidVoteError`
	voteC := unittest.VoteForBlockFixture(s.proposal.Block, unittest.VoteWithStakingSig())
	s.stakingAggregator.On("Verify", voteC.SignerID, mock.Anything).Return(exception)
	err = s.processor.Process(voteC)
	require.ErrorIs(s.T(), err, exception)              // unexpected errors from verifying the vote signature should be propagated
	require.False(s.T(), model.IsInvalidVoteError(err)) // but not interpreted as an invalid vote

	s.stakingAggregator.AssertNotCalled(s.T(), "TrustedAdd")
}

// TestProcess_TrustedAdd_Exception tests that unexpected exceptions returned by
// WeightedSignatureAggregator.TrustedAdd(..) are _not_ interpreted as invalid votes
func (s *StakingVoteProcessorTestSuite) TestProcess_TrustedAdd_Exception() {
	exception := errors.New("unexpected-exception")
	stakingVote := unittest.VoteForBlockFixture(s.proposal.Block, unittest.VoteWithStakingSig())
	*s.stakingAggregator = mockhotstuff.WeightedSignatureAggregator{}
	s.stakingAggregator.On("Verify", stakingVote.SignerID, mock.Anything).Return(nil).Once()
	s.stakingAggregator.On("TrustedAdd", stakingVote.SignerID, mock.Anything).Return(uint64(0), exception).Once()
	err := s.processor.Process(stakingVote)
	require.ErrorIs(s.T(), err, exception)
	require.False(s.T(), model.IsInvalidVoteError(err))
	s.stakingAggregator.AssertExpectations(s.T())
}

// TestProcess_BuildQCError tests error path during process of building QC.
// Building QC is a one time operation, we need to make sure that failing in one of the steps leads to exception.
func (s *StakingVoteProcessorTestSuite) TestProcess_BuildQCError() {
	// In this test we will mock all dependencies for happy path, and replace some branches with unhappy path
	// to simulate errors along the branches.
	vote := unittest.VoteForBlockFixture(s.proposal.Block)

	// in this test case we aren't able to aggregate staking signature
	exception := errors.New("staking-aggregate-exception")
	stakingSigAggregator := &mockhotstuff.WeightedSignatureAggregator{}
	stakingSigAggregator.On("Verify", mock.Anything, mock.Anything).Return(nil).Once()
	stakingSigAggregator.On("TrustedAdd", mock.Anything, mock.Anything).Return(s.minRequiredWeight, nil).Once()
	stakingSigAggregator.On("Aggregate").Return(nil, nil, exception).Once()

	s.processor.stakingSigAggtor = stakingSigAggregator
	err := s.processor.Process(vote)
	require.ErrorIs(s.T(), err, exception)
	stakingSigAggregator.AssertExpectations(s.T())
}

// TestProcess_NotEnoughStakingWeight tests a scenario where we first don't have enough weight,
// then we iteratively increase it to the point where we have enough staking weight. No QC should be created.
func (s *StakingVoteProcessorTestSuite) TestProcess_NotEnoughStakingWeight() {
	for i := s.sigWeight; i < s.minRequiredWeight; i += s.sigWeight {
		vote := unittest.VoteForBlockFixture(s.proposal.Block)
		s.stakingAggregator.On("Verify", vote.SignerID, crypto.Signature(vote.SigData)).Return(nil).Once()
		err := s.processor.Process(vote)
		require.NoError(s.T(), err)
	}
	require.False(s.T(), s.processor.done.Load())
	s.onQCCreatedState.AssertNotCalled(s.T(), "onQCCreated")
	s.stakingAggregator.AssertExpectations(s.T())
}

// TestProcess_CreatingQC tests a scenario when we have collected enough staking weight
// and proceed to build QC. Created QC has to have all signatures and identities aggregated by
// aggregator.
func (s *StakingVoteProcessorTestSuite) TestProcess_CreatingQC() {
	// prepare test setup: 13 votes with staking sigs
	stakingSigners := s.allParticipants[:14].NodeIDs()
	signerIndices, err := signature.EncodeSignersToIndices(stakingSigners, stakingSigners)
	require.NoError(s.T(), err)

	// setup aggregator
	*s.stakingAggregator = mockhotstuff.WeightedSignatureAggregator{}
	expectedSigData := unittest.RandomBytes(128)
	s.stakingAggregator.On("Aggregate").Return(stakingSigners, expectedSigData, nil).Once()

	// expected QC
	s.onQCCreatedState.On("onQCCreated", mock.Anything).Run(func(args mock.Arguments) {
		qc := args.Get(0).(*flow.QuorumCertificate)
		// ensure that QC contains correct field
		expectedQC := &flow.QuorumCertificate{
			View:          s.proposal.Block.View,
			BlockID:       s.proposal.Block.BlockID,
			SignerIndices: signerIndices,
			SigData:       expectedSigData,
		}
		require.Equal(s.T(), expectedQC, qc)
	}).Return(nil).Once()

	// add votes
	for _, signer := range stakingSigners {
		vote := unittest.VoteForBlockFixture(s.proposal.Block)
		vote.SignerID = signer
		expectedSig := crypto.Signature(vote.SigData)
		s.stakingAggregator.On("Verify", vote.SignerID, expectedSig).Return(nil).Once()
		s.stakingAggregator.On("TrustedAdd", vote.SignerID, expectedSig).Run(func(args mock.Arguments) {
			s.stakingTotalWeight += s.sigWeight
		}).Return(s.stakingTotalWeight, nil).Once()
		err := s.processor.Process(vote)
		require.NoError(s.T(), err)
	}

	require.True(s.T(), s.processor.done.Load())
	s.onQCCreatedState.AssertExpectations(s.T())
	s.stakingAggregator.AssertExpectations(s.T())

	// processing extra votes shouldn't result in creating new QCs
	vote := unittest.VoteForBlockFixture(s.proposal.Block)
	err = s.processor.Process(vote)
	require.NoError(s.T(), err)

	s.onQCCreatedState.AssertExpectations(s.T())
}

// TestProcess_ConcurrentCreatingQC tests a scenario where multiple goroutines process vote at same time,
// we expect only one QC created in this scenario.
func (s *StakingVoteProcessorTestSuite) TestProcess_ConcurrentCreatingQC() {
	stakingSigners := s.allParticipants[:10].NodeIDs()
	mockAggregator := func(aggregator *mockhotstuff.WeightedSignatureAggregator) {
		aggregator.On("Verify", mock.Anything, mock.Anything).Return(nil)
		aggregator.On("TrustedAdd", mock.Anything, mock.Anything).Return(s.minRequiredWeight, nil)
		aggregator.On("TotalWeight").Return(s.minRequiredWeight)
		aggregator.On("Aggregate").Return(stakingSigners, unittest.RandomBytes(128), nil)
	}

	// mock aggregators, so we have enough weight and shares for creating QC
	*s.stakingAggregator = mockhotstuff.WeightedSignatureAggregator{}
	mockAggregator(s.stakingAggregator)

	// at this point sending any vote should result in creating QC.
	s.onQCCreatedState.On("onQCCreated", mock.Anything).Return(nil).Once()

	var startupWg, shutdownWg sync.WaitGroup

	vote := unittest.VoteForBlockFixture(s.proposal.Block)
	startupWg.Add(1)
	// prepare goroutines, so they are ready to submit a vote at roughly same time
	for i := 0; i < 5; i++ {
		shutdownWg.Add(1)
		go func() {
			defer shutdownWg.Done()
			startupWg.Wait()
			err := s.processor.Process(vote)
			require.NoError(s.T(), err)
		}()
	}

	startupWg.Done()

	// wait for all routines to finish
	shutdownWg.Wait()

	s.onQCCreatedState.AssertNumberOfCalls(s.T(), "onQCCreated", 1)
}

// TestStakingVoteProcessorV2_BuildVerifyQC tests a complete path from creating votes to collecting votes and then
// building & verifying QC.
// We start with leader proposing a block, then new leader collects votes and builds a QC.
// Need to verify that QC that was produced is valid and can be embedded in new proposal.
func TestStakingVoteProcessorV2_BuildVerifyQC(t *testing.T) {
	epochCounter := uint64(3)
	epochLookup := &modulemock.EpochLookup{}
	proposerView := uint64(20)
	epochLookup.On("EpochForView", proposerView).Return(epochCounter, nil)

	// signers hold objects that are created with private key and can sign votes and proposals
	signers := make(map[flow.Identifier]*verification.StakingSigner)
	// prepare staking signers, each signer has its own private/public key pair
	stakingSigners := unittest.IdentityListFixture(7, func(identity *flow.Identity) {
		stakingPriv := unittest.StakingPrivKeyFixture()
		identity.StakingPubKey = stakingPriv.PublicKey()

		me, err := local.New(identity.IdentitySkeleton, stakingPriv)
		require.NoError(t, err)

		signers[identity.NodeID] = verification.NewStakingSigner(me)
	}).Sort(flow.Canonical[flow.Identity])

	leader := stakingSigners[0]
	block := helper.MakeBlock(helper.WithBlockView(proposerView), helper.WithBlockProposer(leader.NodeID))

	committee := &mockhotstuff.DynamicCommittee{}
	committee.On("IdentitiesByEpoch", block.View).Return(stakingSigners.ToSkeleton(), nil)
	committee.On("IdentitiesByBlock", block.BlockID).Return(stakingSigners, nil)
	committee.On("QuorumThresholdForView", mock.Anything).Return(committees.WeightThresholdToBuildQC(stakingSigners.ToSkeleton().TotalWeight()), nil)

	votes := make([]*model.Vote, 0, len(stakingSigners))

	// first staking signer will be leader collecting votes for proposal
	// prepare votes for every member of committee except leader
	for _, signer := range stakingSigners[1:] {
		vote, err := signers[signer.NodeID].CreateVote(block)
		require.NoError(t, err)
		votes = append(votes, vote)
	}

	// create and sign proposal
	leaderVote, err := signers[leader.NodeID].CreateVote(block)
	require.NoError(t, err)
	proposal := helper.MakeSignedProposal(helper.WithProposal(
		helper.MakeProposal(helper.WithBlock(block))), helper.WithSigData(leaderVote.SigData))

	qcCreated := false
	onQCCreated := func(qc *flow.QuorumCertificate) {
		// create verifier that will do crypto checks of created QC
		verifier := verification.NewStakingVerifier()
		// create validator which will do compliance and crypto checked of created QC
		validator := hotstuffvalidator.New(committee, verifier)
		// check if QC is valid against parent
		err := validator.ValidateQC(qc)
		require.NoError(t, err)

		qcCreated = true
	}

	voteProcessorFactory := NewStakingVoteProcessorFactory(committee, onQCCreated)
	voteProcessor, err := voteProcessorFactory.Create(unittest.Logger(), proposal)
	require.NoError(t, err)

	// process votes by new leader, this will result in producing new QC
	for _, vote := range votes {
		err := voteProcessor.Process(vote)
		require.NoError(t, err)
	}

	require.True(t, qcCreated)
}
