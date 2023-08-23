package avro

import (
	"fmt"
	"github.com/onflow/flow-go/model/cluster"
	codecAvro "github.com/onflow/flow-go/model/encoding/avro"
	"github.com/onflow/flow-go/model/messages"
	"github.com/onflow/flow-go/utils/unittest"
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkCodec_Encode(b *testing.B) {
	tbCodec := NewCodec(codecAvro.TransactionBodySchema)
	transactionBody := unittest.TransactionBodyFixture()
	transactionBodyMap := codecAvro.ConvertStructToMap(transactionBody)
	transactionBodyEncoded, err := tbCodec.Encode(transactionBodyMap)
	require.NoError(b, err)
	b.Logf("Size of TransactionBody encoded: %d bytes\n", len(transactionBodyEncoded))

	b.Run(fmt.Sprintf("encode_transaction_body"), func(b *testing.B) {
		defer runCleanTimer(b)
		for n := 0; n < b.N; n++ {
			_, err := tbCodec.Encode(transactionBodyMap)
			if err != nil {
				b.Error(err)
			}
		}
	})

	numberOfPayloads := [3]int{5, 15, 30}
	cbCodec := NewCodec(codecAvro.ClusterBlockProposal)
	for _, n := range numberOfPayloads {
		clusterBlockProposalData := ClusterBlockFixture(n)
		clusterBlockProposalDataMap := codecAvro.ConvertStructToMap(*clusterBlockProposalData)
		clusterBlockProposalEncodedData, err := cbCodec.Encode(clusterBlockProposalDataMap)
		require.NoError(b, err)
		b.Logf("Size of ClusterBlockProposal with %d payloads encoded: %d bytes\n", n, len(clusterBlockProposalEncodedData))

		b.Run(fmt.Sprintf("encode_cluster_block_with_payload_number_%d", n), func(b *testing.B) {
			defer runCleanTimer(b)
			for n := 0; n < b.N; n++ {
				_, err := cbCodec.Encode(clusterBlockProposalDataMap)
				if err != nil {
					b.Error(err)
				}
			}
		})
	}
}

func BenchmarkCodec_Decode(b *testing.B) {
	tbCodec := NewCodec(codecAvro.TransactionBodySchema)
	transactionBody := unittest.TransactionBodyFixture()
	transactionBodyMap := codecAvro.ConvertStructToMap(transactionBody)
	transactionBodyEncoded, err := tbCodec.Encode(transactionBodyMap)
	require.NoError(b, err)
	b.Logf("Size of TransactionBody encoded: %d bytes\n", len(transactionBodyEncoded))

	b.Run(fmt.Sprintf("decode_transaction_body"), func(b *testing.B) {
		defer runCleanTimer(b)
		for n := 0; n < b.N; n++ {
			_, err := tbCodec.Decode(transactionBodyEncoded)
			if err != nil {
				b.Error(err)
			}
		}
	})

	numberOfPayloads := [3]int{5, 15, 30}
	cbCodec := NewCodec(codecAvro.ClusterBlockProposal)
	for _, n := range numberOfPayloads {
		clusterBlockProposalData := ClusterBlockFixture(n)
		clusterBlockProposalDataMap := codecAvro.ConvertStructToMap(*clusterBlockProposalData)
		clusterBlockProposalEncodedData, err := cbCodec.Encode(clusterBlockProposalDataMap)
		require.NoError(b, err)
		b.Logf("Size of ClusterBlockProposal with %d payloads encoded: %d bytes\n", n, len(clusterBlockProposalEncodedData))

		b.Run(fmt.Sprintf("decode_cluster_block_with_payload_number_%d", n), func(b *testing.B) {
			defer runCleanTimer(b)
			for n := 0; n < b.N; n++ {
				_, err := cbCodec.Decode(clusterBlockProposalEncodedData)
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
