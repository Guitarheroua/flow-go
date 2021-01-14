package validation

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/onflow/flow-go/engine"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/storage"
)

type receiptValidator struct {
	state    protocol.State
	index    storage.Index
	results  storage.ExecutionResults
	verifier module.Verifier
}

func NewReceiptValidator(state protocol.State, index storage.Index, results storage.ExecutionResults, verifier module.Verifier) *receiptValidator {
	rv := &receiptValidator{
		state:    state,
		index:    index,
		results:  results,
		verifier: verifier,
	}

	return rv
}

// checkIsStakedNodeWithRole checks whether, at the given block, `nodeID`
//   * has _positive_ weight
//   * and has the expected role
// Returns the following errors:
//   * sentinel engine.InvalidInputError if any of the above-listed conditions are violated.
// Note: the method receives the identity as proof of its existence.
// Therefore, we consider the case where the respective identity is unknown to the
// protocol state as a symptom of a fatal implementation bug.
func (v *receiptValidator) ensureStakedNodeWithRole(identity *flow.Identity, expectedRole flow.Role) error {
	// check that the origin is an expected node
	if identity.Role != expectedRole {
		return engine.NewInvalidInputErrorf("expected node %x to have identity %v but got %v", identity.NodeID, expectedRole, identity.Role)
	}

	// check if the identity has a stake
	if identity.Stake == 0 {
		return engine.NewInvalidInputErrorf("node has zero stake (%x)", identity.NodeID)
	}

	// TODO: check if node was ejected
	return nil
}

// identityForNode ensures that `nodeID` is an authorized member of the network
// at the given block and returns the corresponding node's full identity.
// Error returns:
//   * sentinel engine.InvalidInputError is nodeID is NOT an authorized member of the network
//   * generic error indicating a fatal internal problem
func (v *receiptValidator) identityForNode(blockID flow.Identifier, nodeID flow.Identifier) (*flow.Identity, error) {
	// get the identity of the origin node
	identity, err := v.state.AtBlockID(blockID).Identity(nodeID)
	if err != nil {
		if protocol.IsIdentityNotFound(err) {
			return nil, engine.NewInvalidInputErrorf("unknown node identity: %w", err)
		}
		// unexpected exception
		return nil, fmt.Errorf("failed to retrieve node identity: %w", err)
	}

	return identity, nil
}

func (v *receiptValidator) verifySignature(receipt *flow.ExecutionReceipt, nodeIdentity *flow.Identity) error {
	id := receipt.ID()
	valid, err := v.verifier.Verify(id[:], receipt.ExecutorSignature, nodeIdentity.StakingPubKey)
	if err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	if !valid {
		return engine.NewInvalidInputErrorf("Invalid signature for (%x)", nodeIdentity.NodeID)
	}

	return nil
}

func (v *receiptValidator) verifyChunksFormat(result *flow.ExecutionResult) error {
	for index, chunk := range result.Chunks.Items() {
		if uint(index) != chunk.CollectionIndex {
			return engine.NewInvalidInputErrorf("invalid CollectionIndex, expected %d got %d", index, chunk.CollectionIndex)
		}

		if chunk.BlockID != result.BlockID {
			return engine.NewInvalidInputErrorf("invalid blockID, expected %v got %v", result.BlockID, chunk.BlockID)
		}
	}

	// we create one chunk per collection, plus the
	// system chunk. so we can check if the chunk number matches with the
	// number of guarantees plus one; this will ensure the execution receipt
	// cannot lie about having less chunks and having the remaining ones
	// approved
	requiredChunks := 1 // system chunk: must exist for block's ExecutionResult, even if block payload itself is empty

	index, err := v.index.ByBlockID(result.BlockID)
	if err != nil {
		if !errors.Is(err, storage.ErrNotFound) {
			return err
		}
		// reaching this line means the block is empty, i.e. it has no payload => we expect only the system chunk
	} else {
		requiredChunks += len(index.CollectionIDs)
	}

	if result.Chunks.Len() != requiredChunks {
		return engine.NewInvalidInputErrorf("invalid number of chunks, expected %d got %d",
			requiredChunks, result.Chunks.Len())
	}

	return nil
}

func (v *receiptValidator) previousResult(result *flow.ExecutionResult) (*flow.ExecutionResult, error) {
	prevResult, err := v.results.ByID(result.PreviousResultID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, engine.NewInvalidInputErrorf("receipt's previous result (%x) is unknown", result.PreviousResultID)
		} else {
			return nil, err
		}
	}
	return prevResult, nil
}

// subgraphCheck enforces that result forms a valid sub-graph:
// Let R1 be a result that references block A, and R2 be R1's parent result.
// The execution results form a valid subgraph if and only if R2 references
// A's parent.
func (v *receiptValidator) subgraphCheck(result *flow.ExecutionResult, prevResult *flow.ExecutionResult) error {
	block, err := v.state.AtBlockID(result.BlockID).Head()
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return engine.NewInvalidInputErrorf("no block found %v %w", result.BlockID, err)
		} else {
			return err
		}
	}

	if prevResult.BlockID != block.ParentID {
		return engine.NewInvalidInputErrorf("invalid block for previous result %v", prevResult.BlockID)
	}

	return nil
}

// resultChainCheck enforces that the end state of the parent result
// matches the current result's start state
func (v *receiptValidator) resultChainCheck(result *flow.ExecutionResult, prevResult *flow.ExecutionResult) error {
	finalState, isOk := prevResult.FinalStateCommitment()
	if !isOk {
		return fmt.Errorf("missing final state commitment in execution result %v", prevResult.ID())
	}
	initialState, isOK := result.InitialStateCommit()
	if !isOK {
		return fmt.Errorf("missing initial state commitment in execution result %v", result.ID())
	}
	if !bytes.Equal(initialState, finalState) {
		return engine.NewInvalidInputError("execution results do not form chain")
	}
	return nil
}

// Validate performs checks for ExecutionReceipt being valid or no.
// Checks performed:
// 	* can find stake and stake is positive
//	* signature is correct
//	* chunks are in correct format
// 	* execution result has a valid parent
// Returns nil if all checks passed successfully
func (v *receiptValidator) Validate(receipt *flow.ExecutionReceipt) error {
	identity, err := v.identityForNode(receipt.ExecutionResult.BlockID, receipt.ExecutorID)
	if err != nil {
		return err
	}

	err = v.ensureStakedNodeWithRole(identity, flow.RoleExecution)
	if err != nil {
		return fmt.Errorf("staked node invalid: %w", err)
	}

	err = v.verifySignature(receipt, identity)
	if err != nil {
		return fmt.Errorf("invalid receipt signature: %w", err)
	}

	err = v.verifyChunksFormat(&receipt.ExecutionResult)
	if err != nil {
		return fmt.Errorf("invalid chunks format: %w", err)
	}

	prevResult, err := v.previousResult(&receipt.ExecutionResult)
	if err != nil {
		return err
	}

	err = v.subgraphCheck(&receipt.ExecutionResult, prevResult)
	if err != nil {
		return fmt.Errorf("invalid execution result: %w", err)
	}

	err = v.resultChainCheck(&receipt.ExecutionResult, prevResult)
	if err != nil {
		return fmt.Errorf("invalid execution results chain: %w", err)
	}

	return nil
}
