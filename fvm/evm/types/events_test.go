package types_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/onflow/go-ethereum/core/vm"

	"github.com/onflow/cadence/encoding/ccf"
	cdcCommon "github.com/onflow/cadence/runtime/common"
	gethCommon "github.com/onflow/go-ethereum/common"
	gethTypes "github.com/onflow/go-ethereum/core/types"
	"github.com/onflow/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	flowSdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go/fvm/evm/testutils"
	"github.com/onflow/flow-go/fvm/evm/types"
)

func TestEVMBlockExecutedEventCCFEncodingDecoding(t *testing.T) {
	t.Parallel()

	block := &types.Block{
		Height:          2,
		Timestamp:       100,
		TotalSupply:     big.NewInt(1500),
		ParentBlockHash: gethCommon.HexToHash("0x2813452cff514c3054ac9f40cd7ce1b016cc78ab7f99f1c6d49708837f6e06d1"),
		ReceiptRoot:     gethCommon.Hash{},
		TotalGasUsed:    15,
		TransactionHashes: []gethCommon.Hash{
			gethCommon.HexToHash("0x70b67ce6710355acf8d69b2ea013d34e212bc4824926c5d26f189c1ca9667246"),
		},
	}

	event := types.NewBlockEvent(block)
	ev, err := event.Payload.ToCadence()
	require.NoError(t, err)

	bep, err := types.DecodeBlockEventPayload(ev)
	require.NoError(t, err)

	assert.Equal(t, bep.Height, block.Height)

	blockHash, err := block.Hash()
	require.NoError(t, err)
	assert.Equal(t, bep.Hash, blockHash.Hex())

	assert.Equal(t, bep.TotalSupply.Value, block.TotalSupply)
	assert.Equal(t, bep.Timestamp, block.Timestamp)
	assert.Equal(t, bep.TotalGasUsed, block.TotalGasUsed)
	assert.Equal(t, bep.ParentBlockHash, block.ParentBlockHash.Hex())
	assert.Equal(t, bep.ReceiptRoot, block.ReceiptRoot.Hex())

	hashes := make([]gethCommon.Hash, len(bep.TransactionHashes))
	for i, h := range bep.TransactionHashes {
		hashes[i] = gethCommon.HexToHash(h.ToGoValue().(string))
	}
	assert.Equal(t, hashes, block.TransactionHashes)

	v, err := ccf.Encode(ev)
	require.NoError(t, err)
	assert.Equal(t, ccf.HasMsgPrefix(v), true)

	evt, err := ccf.Decode(nil, v)
	require.NoError(t, err)

	assert.Equal(t, evt.Type().ID(), "evm.BlockExecuted")

	location, qualifiedIdentifier, err := cdcCommon.DecodeTypeID(nil, "evm.BlockExecuted")
	require.NoError(t, err)

	assert.Equal(t, flowSdk.EVMLocation{}, location)
	assert.Equal(t, "BlockExecuted", qualifiedIdentifier)
}

func TestEVMTransactionExecutedEventCCFEncodingDecoding(t *testing.T) {
	t.Parallel()

	txEncoded := "fff83b81ff0194000000000000000000000000000000000000000094000000000000000000000000000000000000000180895150ae84a8cdf00000825208"
	txBytes, err := hex.DecodeString(txEncoded)
	require.NoError(t, err)
	txHash := testutils.RandomCommonHash(t)
	blockHash := testutils.RandomCommonHash(t)
	data := "000000000000000000000000000000000000000000000000000000000000002a"
	dataBytes, err := hex.DecodeString(data)
	require.NoError(t, err)
	blockHeight := uint64(2)
	deployedAddress := types.NewAddress(gethCommon.HexToAddress("0x99466ed2e37b892a2ee3e9cd55a98b68f5735db2"))
	log := &gethTypes.Log{
		Index:       1,
		BlockNumber: blockHeight,
		BlockHash:   blockHash,
		TxHash:      txHash,
		TxIndex:     3,
		Address:     gethCommon.HexToAddress("0x99466ed2e37b892a2ee3e9cd55a98b68f5735db2"),
		Data:        dataBytes,
		Topics: []gethCommon.Hash{
			gethCommon.HexToHash("0x24abdb5865df5079dcc5ac590ff6f01d5c16edbc5fab4e195d9febd1114503da"),
		},
	}
	vmError := vm.ErrOutOfGas
	txResult := &types.Result{
		VMError:                 vmError,
		TxType:                  255,
		GasConsumed:             23200,
		DeployedContractAddress: &deployedAddress,
		ReturnedValue:           dataBytes,
		Logs:                    []*gethTypes.Log{log},
		TxHash:                  txHash,
	}

	t.Run("evm.TransactionExecuted with failed status", func(t *testing.T) {
		event := types.NewTransactionEvent(txResult, txBytes, blockHeight, blockHash)
		ev, err := event.Payload.ToCadence()
		require.NoError(t, err)

		tep, err := types.DecodeTransactionEventPayload(ev)
		require.NoError(t, err)

		assert.Equal(t, tep.BlockHeight, blockHeight)
		assert.Equal(t, tep.BlockHash, blockHash.Hex())
		assert.Equal(t, tep.Hash, txHash.Hex())
		assert.Equal(t, tep.Payload, txEncoded)
		assert.Equal(t, types.ErrorCode(tep.ErrorCode), types.ExecutionErrCodeOutOfGas)
		assert.Equal(t, tep.TransactionType, txResult.TxType)
		assert.Equal(t, tep.GasConsumed, txResult.GasConsumed)
		assert.Equal(
			t,
			tep.ContractAddress,
			txResult.DeployedContractAddress.ToCommon().Hex(),
		)

		encodedLogs, err := rlp.EncodeToBytes(txResult.Logs)
		require.NoError(t, err)
		assert.Equal(t, tep.Logs, hex.EncodeToString(encodedLogs))

		v, err := ccf.Encode(ev)
		require.NoError(t, err)
		assert.Equal(t, ccf.HasMsgPrefix(v), true)

		evt, err := ccf.Decode(nil, v)
		require.NoError(t, err)

		assert.Equal(t, evt.Type().ID(), "evm.TransactionExecuted")

		location, qualifiedIdentifier, err := cdcCommon.DecodeTypeID(nil, "evm.TransactionExecuted")
		require.NoError(t, err)

		assert.Equal(t, flowSdk.EVMLocation{}, location)
		assert.Equal(t, "TransactionExecuted", qualifiedIdentifier)
	})

	t.Run("evm.TransactionExecuted with non-failed status", func(t *testing.T) {
		txResult.VMError = nil

		event := types.NewTransactionEvent(txResult, txBytes, blockHeight, blockHash)
		ev, err := event.Payload.ToCadence()
		require.NoError(t, err)

		tep, err := types.DecodeTransactionEventPayload(ev)
		require.NoError(t, err)

		assert.Equal(t, tep.BlockHeight, blockHeight)
		assert.Equal(t, tep.BlockHash, blockHash.Hex())
		assert.Equal(t, tep.Hash, txHash.Hex())
		assert.Equal(t, tep.Payload, txEncoded)
		assert.Equal(t, types.ErrCodeNoError, types.ErrorCode(tep.ErrorCode))
		assert.Equal(t, tep.TransactionType, txResult.TxType)
		assert.Equal(t, tep.GasConsumed, txResult.GasConsumed)
		assert.NotNil(t, txResult.DeployedContractAddress)
		assert.Equal(
			t,
			tep.ContractAddress,
			txResult.DeployedContractAddress.ToCommon().Hex(),
		)

		encodedLogs, err := rlp.EncodeToBytes(txResult.Logs)
		require.NoError(t, err)
		assert.Equal(t, tep.Logs, hex.EncodeToString(encodedLogs))

		v, err := ccf.Encode(ev)
		require.NoError(t, err)
		assert.Equal(t, ccf.HasMsgPrefix(v), true)

		evt, err := ccf.Decode(nil, v)
		require.NoError(t, err)

		assert.Equal(t, evt.Type().ID(), "evm.TransactionExecuted")

		location, qualifiedIdentifier, err := cdcCommon.DecodeTypeID(nil, "evm.TransactionExecuted")
		require.NoError(t, err)

		assert.Equal(t, flowSdk.EVMLocation{}, location)
		assert.Equal(t, "TransactionExecuted", qualifiedIdentifier)
	})
}
