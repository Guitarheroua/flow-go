package cbor_test

import (
	"fmt"
	"github.com/onflow/flow-go/model/cluster"
	"github.com/onflow/flow-go/model/messages"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/onflow/flow-go/network/codec"
	"github.com/onflow/flow-go/network/codec/cbor"
	"github.com/onflow/flow-go/utils/unittest"
)

func TestCodec_Decode(t *testing.T) {
	t.Parallel()

	c := cbor.NewCodec()

	t.Run("decodes message successfully", func(t *testing.T) {
		t.Parallel()

		data := unittest.ProposalFixture()
		encoded, err := c.Encode(data)
		require.NoError(t, err)

		decoded, err := c.Decode(encoded)
		require.NoError(t, err)
		require.Equal(t, data, decoded)
	})

	t.Run("returns error when data is empty", func(t *testing.T) {
		t.Parallel()

		decoded, err := c.Decode(nil)
		assert.Nil(t, decoded)
		assert.True(t, codec.IsErrInvalidEncoding(err))

		decoded, err = c.Decode([]byte{})
		assert.Nil(t, decoded)
		assert.True(t, codec.IsErrInvalidEncoding(err))
	})

	t.Run("returns error when message code is invalid", func(t *testing.T) {
		t.Parallel()

		decoded, err := c.Decode([]byte{codec.CodeMin.Uint8()})
		assert.Nil(t, decoded)
		assert.True(t, codec.IsErrUnknownMsgCode(err))

		decoded, err = c.Decode([]byte{codec.CodeMax.Uint8()})
		assert.Nil(t, decoded)
		assert.True(t, codec.IsErrUnknownMsgCode(err))

		decoded, err = c.Decode([]byte{codec.CodeMin.Uint8() - 1})
		assert.Nil(t, decoded)
		assert.True(t, codec.IsErrUnknownMsgCode(err))

		decoded, err = c.Decode([]byte{codec.CodeMax.Uint8() + 1})
		assert.Nil(t, decoded)
		assert.True(t, codec.IsErrUnknownMsgCode(err))
	})

	t.Run("returns error when unmarshalling fails - empty", func(t *testing.T) {
		t.Parallel()

		decoded, err := c.Decode([]byte{codec.CodeBlockProposal.Uint8()})
		assert.Nil(t, decoded)
		assert.True(t, codec.IsErrMsgUnmarshal(err))
	})

	t.Run("returns error when unmarshalling fails - wrong type", func(t *testing.T) {
		t.Parallel()

		data := unittest.ProposalFixture()
		encoded, err := c.Encode(data)
		require.NoError(t, err)

		encoded[0] = codec.CodeCollectionGuarantee.Uint8()

		decoded, err := c.Decode(encoded)
		assert.Nil(t, decoded)
		assert.True(t, codec.IsErrMsgUnmarshal(err))
	})

	t.Run("returns error when unmarshalling fails - corrupt", func(t *testing.T) {
		t.Parallel()

		data := unittest.ProposalFixture()
		encoded, err := c.Encode(data)
		require.NoError(t, err)

		encoded[2] = 0x20 // corrupt payload

		decoded, err := c.Decode(encoded)
		assert.Nil(t, decoded)
		assert.True(t, codec.IsErrMsgUnmarshal(err))
	})
}

func BenchmarkCodec_Encode(b *testing.B) {
	cborCodec := cbor.NewCodec()

	transactionBody := unittest.TransactionBodyFixture()
	transactionBodyEncoded, err := cborCodec.Encode(&transactionBody)
	require.NoError(b, err)
	b.Logf("Size of TransactionBody encoded: %d bytes\n", len(transactionBodyEncoded))

	b.Run(fmt.Sprintf("encode_transaction_body"), func(b *testing.B) {
		defer runCleanTimer(b)
		for n := 0; n < b.N; n++ {
			_, err := cborCodec.Encode(&transactionBody)
			if err != nil {
				b.Error(err)
			}
		}
	})

	numberOfPayloads := [3]int{5, 15, 30}
	for _, n := range numberOfPayloads {
		clusterBlockProposalData := ClusterBlockFixture(n)
		clusterBlockProposalEncodedData, err := cborCodec.Encode(clusterBlockProposalData)
		require.NoError(b, err)
		b.Logf("Size of ClusterBlockProposal with %d payloads encoded: %d bytes\n", n, len(clusterBlockProposalEncodedData))

		b.Run(fmt.Sprintf("encode_cluster_block_with_payload_number_%d", n), func(b *testing.B) {
			defer runCleanTimer(b)
			for n := 0; n < b.N; n++ {
				_, err := cborCodec.Encode(clusterBlockProposalData)
				if err != nil {
					b.Error(err)
				}
			}
		})
	}
}

func BenchmarkCodec_Decode(b *testing.B) {
	cborCodec := cbor.NewCodec()

	transactionBody := unittest.TransactionBodyFixture()
	transactionBodyEncoded, err := cborCodec.Encode(&transactionBody)
	require.NoError(b, err)
	b.Logf("Size of TransactionBody encoded: %d bytes\n", len(transactionBodyEncoded))

	b.Run(fmt.Sprintf("decode_transaction_body"), func(b *testing.B) {
		defer runCleanTimer(b)
		for n := 0; n < b.N; n++ {
			_, err := cborCodec.Decode(transactionBodyEncoded)
			if err != nil {
				b.Error(err)
			}
		}
	})

	numberOfPayloads := [3]int{5, 15, 30}
	for _, n := range numberOfPayloads {
		clusterBlockProposalData := ClusterBlockFixture(n)
		clusterBlockProposalEncodedData, err := cborCodec.Encode(clusterBlockProposalData)
		require.NoError(b, err)
		b.Logf("Size of ClusterBlockProposal with %d payloads encoded: %d bytes\n", n, len(clusterBlockProposalEncodedData))

		b.Run(fmt.Sprintf("decode_cluster_block_with_payload_number_%d", n), func(b *testing.B) {
			defer runCleanTimer(b)
			for n := 0; n < b.N; n++ {
				_, err := cborCodec.Decode(clusterBlockProposalEncodedData)
				if err != nil {
					b.Error(err)
				}
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

func ClusterBlockFixture(n int) *messages.ClusterBlockProposal {
	if n <= 0 {
		n = 5
	}

	payload := unittest.ClusterPayloadFixture(n)
	header := unittest.BlockHeaderFixture()
	header.PayloadHash = payload.Hash()
	block := cluster.Block{
		Header:  header,
		Payload: payload,
	}
	return unittest.ClusterProposalFromBlock(&block)
}
