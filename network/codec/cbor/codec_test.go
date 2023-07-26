package cbor_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	blockProposalData := unittest.ProposalFixture()
	b.Run(fmt.Sprintf("cbor_encode_block_proposal"), func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := cborCodec.Encode(blockProposalData)
			if err != nil {
				b.Error(err)
			}
		}
	})

	execReceiptData := unittest.ExecutionReceiptFixture(unittest.WithResult(unittest.ExecutionResultFixture()))
	b.Run(fmt.Sprintf("cbor_encode_execution_receipt"), func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err := cborCodec.Encode(execReceiptData)
			if err != nil {
				b.Error(err)
			}
		}
	})

	//TODO: Add another inputs for benchmarking if needed
}

func BenchmarkCodec_Decode(b *testing.B) {
	cborCodec := cbor.NewCodec()
	blockProposalData := unittest.ProposalFixture()
	blockProposalDataEncoded, err := cborCodec.Encode(blockProposalData)
	if err != nil {
		b.Error(err)
	}

	b.Run(fmt.Sprintf("cbor_decode_block_proposal"), func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err = cborCodec.Decode(blockProposalDataEncoded)
			if err != nil {
				b.Error(err)
			}
		}
	})

	execReceiptData := unittest.ExecutionReceiptFixture(unittest.WithResult(unittest.ExecutionResultFixture()))
	execReceiptDataEncoded, err := cborCodec.Encode(execReceiptData)
	if err != nil {
		b.Error(err)
	}

	b.Run(fmt.Sprintf("cbor_decode_execution_receipt"), func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_, err = cborCodec.Decode(execReceiptDataEncoded)
			if err != nil {
				b.Error(err)
			}
		}
	})

	//TODO: Add another inputs for benchmarking  if needed
}
