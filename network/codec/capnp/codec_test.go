package capnp

import (
	"fmt"
	"github.com/onflow/flow-go/model/flow"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3"
	clusterblockproposal "github.com/onflow/flow-go/model/encoding/capnp/messages/clusterblockproposal"
	transactionbody "github.com/onflow/flow-go/model/encoding/capnp/messages/transactionbody"
	"github.com/onflow/flow-go/utils/unittest"
)

func BenchmarkCodec_Encode(b *testing.B) {
	codec := NewCapnpCodec(false)

	tbMsg, tbSeg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	require.NoError(b, err)

	TransactionBodyMessageFixture(b, tbSeg)

	transactionBodyEncoded, err := codec.Encode(tbMsg)
	require.NoError(b, err)
	b.Logf("Size of TransactionBody encoded: %d bytes\n", len(transactionBodyEncoded))

	b.Run(fmt.Sprintf("encode_transaction_body"), func(b *testing.B) {
		defer runCleanTimer(b)
		for n := 0; n < b.N; n++ {
			_, err := codec.Encode(tbMsg)
			if err != nil {
				b.Error(err)
			}
		}
	})

	numberOfPayloads := [3]int{5, 15, 30}
	for _, n := range numberOfPayloads {
		cbpMsg, cbpSeg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		require.NoError(b, err)

		ClusterBlockFixture(b, cbpSeg, n)

		cbpMsgEncoded, err := codec.Encode(cbpMsg)
		require.NoError(b, err)
		b.Logf("Size of ClusterBlockProposal with %d payloads encoded: %d bytes\n", n, len(cbpMsgEncoded))

		b.Run(fmt.Sprintf("encode_cluster_block_with_payload_number_%d", n), func(b *testing.B) {
			defer runCleanTimer(b)
			for n := 0; n < b.N; n++ {
				_, err = codec.Encode(cbpMsg)
				if err != nil {
					b.Error(err)
				}
			}
		})
	}
}

func BenchmarkCodec_Decode(b *testing.B) {
	codec := NewCapnpCodec(false)

	tbMsg, tbSeg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	require.NoError(b, err)

	TransactionBodyMessageFixture(b, tbSeg)

	transactionBodyEncoded, err := codec.Encode(tbMsg)
	require.NoError(b, err)
	b.Logf("Size of TransactionBody encoded: %d bytes\n", len(transactionBodyEncoded))

	//tbMsgDecoded, _, err := capnp.NewMessage(capnp.SingleSegment(nil))
	//require.NoError(b, err)

	b.Run(fmt.Sprintf("decode_transaction_body"), func(b *testing.B) {
		defer runCleanTimer(b)
		for n := 0; n < b.N; n++ {
			_, err := codec.Decode(transactionBodyEncoded)
			if err != nil {
				b.Error(err)
			}

			//err = capnpMarshaler.Unmarshal(capnpBuffer.Bytes(), tbMsgDecoded)
			//if err != nil {
			//	b.Error(err)
			//}
		}
	})

	numberOfPayloads := [3]int{5, 15, 30}
	for _, n := range numberOfPayloads {
		cbpMsg, cbpSeg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		require.NoError(b, err)

		ClusterBlockFixture(b, cbpSeg, n)

		cbpMsgEncoded, err := codec.Encode(cbpMsg)
		require.NoError(b, err)
		b.Logf("Size of ClusterBlockProposal with %d payloads encoded: %d bytes\n", n, len(cbpMsgEncoded))

		//cbpMsgDecoded, _, err := capnp.NewMessage(capnp.SingleSegment(nil))
		//require.NoError(b, err)

		b.Run(fmt.Sprintf("decode_cluster_block_with_payload_number_%d", n), func(b *testing.B) {
			defer runCleanTimer(b)
			for n := 0; n < b.N; n++ {
				_, err = codec.Decode(cbpMsgEncoded)
				if err != nil {
					b.Error(err)
				}

				//err = capnpMarshaler.Unmarshal(capnpBuffer.Bytes(), cbpMsgDecoded)
				//if err != nil {
				//	b.Error(err)
				//}
			}
		})
	}
}

func runCleanTimer(b *testing.B) func() {
	b.StopTimer()
	b.ResetTimer()
	b.StartTimer()

	return b.StopTimer
}

func ClusterBlockFixture(b *testing.B, s *capnp.Segment, n int) clusterblockproposal.ClusterBlockProposalMessage {
	if n <= 0 {
		n = 5
	}

	clusterBlockProposalMsg, err := clusterblockproposal.NewRootClusterBlockProposalMessage(s)
	require.NoError(b, err)

	height := 1 + uint64(rand.Uint32())
	parentView := height + uint64(rand.Intn(1000))
	view := parentView + 1 + uint64(rand.Intn(10))

	var lastViewTC clusterblockproposal.TimeoutCertificateMessage

	if view != parentView+1 {
		newestQC, err := clusterblockproposal.NewRootQuorumCertificateMessage(s)
		require.NoError(b, err)
		newestQC.SetView(parentView)
		blockId := unittest.IdentifierFixture()
		err = newestQC.SetBlockID(blockId[:])
		require.NoError(b, err)
		err = newestQC.SetSignerIndices(unittest.SignerIndicesFixture(3))
		require.NoError(b, err)
		err = newestQC.SetSigData(unittest.QCSigDataFixture())
		require.NoError(b, err)

		lastViewTC, err = clusterblockproposal.NewRootTimeoutCertificateMessage(s)
		require.NoError(b, err)
		lastViewTC.SetView(view - 1)
		err = lastViewTC.SetNewestQC(newestQC)
		require.NoError(b, err)
		newestQCViews, err := capnp.NewUInt64List(s, 1)
		require.NoError(b, err)
		newestQCViews.Set(0, newestQC.View())
		err = lastViewTC.SetNewestQCViews(newestQCViews)
		require.NoError(b, err)
		err = lastViewTC.SetSignerIndices(unittest.SignerIndicesFixture(4))
		require.NoError(b, err)
		err = lastViewTC.SetSigData(unittest.QCSigDataFixture())
		require.NoError(b, err)
	}

	headerMessage, err := clusterblockproposal.NewRootHeaderMessage(s)
	require.NoError(b, err)
	err = headerMessage.SetChainID(string(flow.Emulator))
	require.NoError(b, err)
	parentBlockId := unittest.IdentifierFixture()
	err = headerMessage.SetParentID(parentBlockId[:])
	require.NoError(b, err)
	headerMessage.SetHeight(height + 1)
	payloadHash := unittest.IdentifierFixture()
	err = headerMessage.SetPayloadHash(payloadHash[:])
	require.NoError(b, err)
	headerMessage.SetTimestamp(time.Now().UTC().Unix())
	headerMessage.SetView(view)
	headerMessage.SetParentView(parentView)
	err = headerMessage.SetParentVoterIndices(unittest.SignerIndicesFixture(4))
	require.NoError(b, err)
	err = headerMessage.SetParentVoterSigData(unittest.QCSigDataFixture())
	require.NoError(b, err)
	proposerID := unittest.IdentifierFixture()
	err = headerMessage.SetProposerID(proposerID[:])
	require.NoError(b, err)
	err = headerMessage.SetProposerSigData(unittest.SignatureFixture())
	require.NoError(b, err)
	err = headerMessage.SetLastViewTC(lastViewTC)
	require.NoError(b, err)

	err = clusterBlockProposalMsg.SetHeader(headerMessage)
	require.NoError(b, err)

	clusterBlockPayloadMessage, err := clusterblockproposal.NewRootUntrustedClusterBlockPayloadMessage(s)
	require.NoError(b, err)

	collection, err := transactionbody.NewTransactionBodyMessage_List(s, int32(n))
	require.NoError(b, err)
	for i := 0; i < n; i++ {
		tbm := TransactionBodyMessageFixture(b, s)

		err = collection.Set(i, tbm)
		require.NoError(b, err)
	}
	err = clusterBlockPayloadMessage.SetCollection(collection)
	require.NoError(b, err)

	err = clusterBlockPayloadMessage.SetReferenceBlockID(flow.ZeroID[:])
	require.NoError(b, err)

	clusterBlockProposalMsg.SetPayload(clusterBlockPayloadMessage)

	payload := unittest.ClusterPayloadFixture(n)
	header := unittest.BlockHeaderFixture()
	header.PayloadHash = payload.Hash()

	return clusterBlockProposalMsg
}

func TransactionBodyMessageFixture(b *testing.B, s *capnp.Segment) transactionbody.TransactionBodyMessage {
	transactionBodyMsg, err := transactionbody.NewRootTransactionBodyMessage(s)
	require.NoError(b, err)

	refBlockId := unittest.IdentifierFixture()
	err = transactionBodyMsg.SetReferenceBlockID(refBlockId[:])
	require.NoError(b, err)

	err = transactionBodyMsg.SetScript([]byte("pub fun main() {}"))
	require.NoError(b, err)

	transactionBodyMsg.SetGasLimit(10)

	proposalKeyMessage, err := transactionbody.NewRootProposalKeyMessage(s)
	require.NoError(b, err)

	address := unittest.AddressFixture()
	err = proposalKeyMessage.SetAddress(address[:])
	require.NoError(b, err)
	proposalKeyMessage.SetKeyIndex(1)
	proposalKeyMessage.SetSequenceNumber(0)

	err = transactionBodyMsg.SetProposalKey(proposalKeyMessage)
	require.NoError(b, err)

	address = unittest.AddressFixture()
	err = transactionBodyMsg.SetPayer(address[:])
	require.NoError(b, err)

	dataList, err := capnp.NewDataList(s, 1)
	require.NoError(b, err)

	address = unittest.AddressFixture()
	err = dataList.Set(0, address[:])
	require.NoError(b, err)
	err = transactionBodyMsg.SetAuthorizers(dataList)
	require.NoError(b, err)

	transactionSignatureFixture := unittest.TransactionSignatureFixture()
	transactionSignatureMessage, err := transactionbody.NewRootTransactionSignatureMessage(s)
	err = transactionSignatureMessage.SetAddress(transactionSignatureFixture.Address[:])
	require.NoError(b, err)
	transactionSignatureMessage.SetSignerIndex(int32(transactionSignatureFixture.SignerIndex))
	transactionSignatureMessage.SetKeyIndex(transactionSignatureFixture.KeyIndex)
	err = transactionSignatureMessage.SetSignature(transactionSignatureFixture.Signature)
	require.NoError(b, err)

	transSigList, err := transactionbody.NewTransactionSignatureMessage_List(s, 1)
	require.NoError(b, err)
	transSigList.Set(0, transactionSignatureMessage)
	transactionBodyMsg.SetEnvelopeSignatures(transSigList)

	return transactionBodyMsg
}
