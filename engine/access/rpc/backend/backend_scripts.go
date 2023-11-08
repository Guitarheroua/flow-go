package backend

import (
	"bytes"
	"context"
	"crypto/md5" //nolint:gosec
	"errors"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	execproto "github.com/onflow/flow/protobuf/go/flow/execution"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onflow/flow-go/engine/access/rpc/connection"
	"github.com/onflow/flow-go/engine/common/rpc"
	fvmerrors "github.com/onflow/flow-go/fvm/errors"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/execution"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/storage"
	"github.com/onflow/flow-go/utils/logging"
)

// uniqueScriptLoggingTimeWindow is the duration for checking the uniqueness of scripts sent for execution
const uniqueScriptLoggingTimeWindow = 10 * time.Minute

type backendScripts struct {
	log               zerolog.Logger
	headers           storage.Headers
	executionReceipts storage.ExecutionReceipts
	state             protocol.State
	connFactory       connection.ConnectionFactory
	metrics           module.BackendScriptsMetrics
	loggedScripts     *lru.Cache[[md5.Size]byte, time.Time]
	nodeCommunicator  Communicator
	scriptExecutor    execution.ScriptExecutor
	scriptExecMode    ScriptExecutionMode
}

// ExecuteScriptAtLatestBlock executes provided script at the latest sealed block.
func (b *backendScripts) ExecuteScriptAtLatestBlock(
	ctx context.Context,
	script []byte,
	arguments [][]byte,
) ([]byte, error) {
	latestHeader, err := b.state.Sealed().Head()
	if err != nil {
		// the latest sealed header MUST be available
		return nil, status.Errorf(codes.Internal, "failed to get latest sealed header: %v", err)
	}

	return b.executeScript(ctx, latestHeader.ID(), latestHeader.Height, script, arguments)
}

// ExecuteScriptAtBlockID executes provided script at the provided block ID.
func (b *backendScripts) ExecuteScriptAtBlockID(
	ctx context.Context,
	blockID flow.Identifier,
	script []byte,
	arguments [][]byte,
) ([]byte, error) {
	header, err := b.headers.ByBlockID(blockID)
	if err != nil {
		return nil, rpc.ConvertStorageError(err)
	}

	return b.executeScript(ctx, blockID, header.Height, script, arguments)
}

// ExecuteScriptAtBlockHeight executes provided script at the provided block height.
func (b *backendScripts) ExecuteScriptAtBlockHeight(
	ctx context.Context,
	blockHeight uint64,
	script []byte,
	arguments [][]byte,
) ([]byte, error) {
	header, err := b.headers.ByHeight(blockHeight)
	if err != nil {
		return nil, rpc.ConvertStorageError(err)
	}

	return b.executeScript(ctx, header.ID(), blockHeight, script, arguments)
}

// executeScript executes the provided script using either the local execution state or the execution
// nodes depending on the node's configuration and the availability of the data.
func (b *backendScripts) executeScript(
	ctx context.Context,
	blockID flow.Identifier,
	height uint64,
	script []byte,
	arguments [][]byte,
) ([]byte, error) {
	// encode to MD5 as low compute/memory lookup key
	// CAUTION: cryptographically insecure md5 is used here, but only to de-duplicate logs.
	// *DO NOT* use this hash for any protocol-related or cryptographic functions.
	insecureScriptHash := md5.Sum(script) //nolint:gosec

	switch b.scriptExecMode {
	case ScriptExecutionModeExecutionNodesOnly:
		result, _, err := b.executeScriptOnAvailableExecutionNodes(ctx, blockID, script, arguments, insecureScriptHash)
		return result, err

	case ScriptExecutionModeLocalOnly:
		result, _, err := b.executeScriptLocally(ctx, blockID, height, script, arguments, insecureScriptHash)
		return result, err

	case ScriptExecutionModeFailover:
		localResult, localDuration, localErr := b.executeScriptLocally(ctx, blockID, height, script, arguments, insecureScriptHash)
		if localErr == nil || isInvalidArgumentError(localErr) {
			return localResult, localErr
		}
		execResult, execDuration, execErr := b.executeScriptOnAvailableExecutionNodes(ctx, blockID, script, arguments, insecureScriptHash)

		b.compareScriptExecutionResults(execResult, execDuration, execErr, localResult, localDuration, localErr, blockID, script, arguments, insecureScriptHash)

		return execResult, execErr

	case ScriptExecutionModeCompare:
		execResult, execDuration, execErr := b.executeScriptOnAvailableExecutionNodes(ctx, blockID, script, arguments, insecureScriptHash)
		// we can only compare the results if there were either no errors or a cadence error
		// since we cannot distinguish the EN error as caused by the block being pruned or some other reason,
		// which may produce a valid RN output but an error for the EN
		if execErr != nil && !isInvalidArgumentError(execErr) {
			return nil, execErr
		}
		localResult, localDuration, localErr := b.executeScriptLocally(ctx, blockID, height, script, arguments, insecureScriptHash)

		b.compareScriptExecutionResults(execResult, execDuration, execErr, localResult, localDuration, localErr, blockID, script, arguments, insecureScriptHash)

		// always return EN results
		return execResult, execErr

	default:
		return nil, status.Errorf(codes.Internal, "unknown script execution mode: %v", b.scriptExecMode)
	}
}

// executeScriptLocally executes the provided script using the local execution state.
func (b *backendScripts) executeScriptLocally(
	ctx context.Context,
	blockID flow.Identifier,
	height uint64,
	script []byte,
	arguments [][]byte,
	insecureScriptHash [md5.Size]byte,
) ([]byte, time.Duration, error) {
	execStartTime := time.Now()

	result, err := b.scriptExecutor.ExecuteAtBlockHeight(ctx, script, arguments, height)

	execEndTime := time.Now()
	execDuration := execEndTime.Sub(execStartTime)

	lg := b.log.With().
		Str("script_executor_addr", "localhost").
		Hex("block_id", logging.ID(blockID)).
		Uint64("height", height).
		Hex("script_hash", insecureScriptHash[:]).
		Dur("execution_dur_ms", execDuration).
		Logger()

	if err != nil {
		convertedErr := convertScriptExecutionError(err, height)

		if status.Code(convertedErr) == codes.InvalidArgument {
			lg.Debug().Err(err).
				Str("script", string(script)).
				Msg("script failed to execute locally")
		} else {
			lg.Error().Err(err).Msg("script execution failed")
			b.metrics.ScriptExecutionErrorLocal()
		}

		return nil, execDuration, convertedErr
	}

	if b.log.GetLevel() == zerolog.DebugLevel && b.shouldLogScript(execEndTime, insecureScriptHash) {
		lg.Debug().
			Str("script", string(script)).
			Msg("Successfully executed script")
		b.loggedScripts.Add(insecureScriptHash, execEndTime)
	}

	// log execution time
	b.metrics.ScriptExecuted(execDuration, len(script))

	return result, execDuration, nil
}

// executeScriptOnAvailableExecutionNodes executes the provided script using available execution nodes.
func (b *backendScripts) executeScriptOnAvailableExecutionNodes(
	ctx context.Context,
	blockID flow.Identifier,
	script []byte,
	arguments [][]byte,
	insecureScriptHash [md5.Size]byte,
) ([]byte, time.Duration, error) {
	// find few execution nodes which have executed the block earlier and provided an execution receipt for it
	executors, err := executionNodesForBlockID(ctx, blockID, b.executionReceipts, b.state, b.log)
	if err != nil {
		return nil, 0, status.Errorf(codes.Internal, "failed to find script executors at blockId %v: %v", blockID.String(), err)
	}

	lg := b.log.With().
		Hex("block_id", logging.ID(blockID)).
		Hex("script_hash", insecureScriptHash[:]).
		Logger()

	var result []byte
	var execDuration time.Duration
	errToReturn := b.nodeCommunicator.CallAvailableNode(
		executors,
		func(node *flow.Identity) error {
			execStartTime := time.Now()

			result, err = b.tryExecuteScriptOnExecutionNode(ctx, node.Address, blockID, script, arguments)

			executionTime := time.Now()
			execDuration = executionTime.Sub(execStartTime)

			if err != nil {
				return err
			}

			if b.log.GetLevel() == zerolog.DebugLevel {
				if b.shouldLogScript(executionTime, insecureScriptHash) {
					lg.Debug().
						Str("script_executor_addr", node.Address).
						Str("script", string(script)).
						Dur("execution_dur_ms", execDuration).
						Msg("Successfully executed script")
					b.loggedScripts.Add(insecureScriptHash, executionTime)
				}
			}

			// log execution time
			b.metrics.ScriptExecuted(time.Since(execStartTime), len(script))

			return nil
		},
		func(node *flow.Identity, err error) bool {
			if status.Code(err) == codes.InvalidArgument {
				lg.Debug().Err(err).
					Str("script_executor_addr", node.Address).
					Str("script", string(script)).
					Msg("script failed to execute on the execution node")
				return true
			}
			return false
		},
	)

	if errToReturn != nil {
		if status.Code(errToReturn) != codes.InvalidArgument {
			b.metrics.ScriptExecutionErrorOnExecutionNode()
			b.log.Error().Err(errToReturn).Msg("script execution failed for execution node internal reasons")
		}
		return nil, execDuration, rpc.ConvertError(errToReturn, "failed to execute script on execution nodes", codes.Internal)
	}

	return result, execDuration, nil
}

// tryExecuteScriptOnExecutionNode attempts to execute the script on the given execution node.
func (b *backendScripts) tryExecuteScriptOnExecutionNode(
	ctx context.Context,
	executorAddress string,
	blockID flow.Identifier,
	script []byte,
	arguments [][]byte,
) ([]byte, error) {
	execRPCClient, closer, err := b.connFactory.GetExecutionAPIClient(executorAddress)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create client for execution node %s: %v",
			executorAddress, err)
	}
	defer closer.Close()

	execResp, err := execRPCClient.ExecuteScriptAtBlockID(ctx, &execproto.ExecuteScriptAtBlockIDRequest{
		BlockId:   blockID[:],
		Script:    script,
		Arguments: arguments,
	})
	if err != nil {
		return nil, status.Errorf(status.Code(err), "failed to execute the script on the execution node %s: %v", executorAddress, err)
	}
	return execResp.GetValue(), nil
}

// isInvalidArgumentError checks if the error is from an invalid argument
func isInvalidArgumentError(scriptExecutionErr error) bool {
	return status.Code(scriptExecutionErr) == codes.InvalidArgument
}

func (b *backendScripts) compareScriptExecutionResults(
	execNodeResult []byte,
	execDuration time.Duration,
	execErr error,
	localResult []byte,
	localDuration time.Duration,
	localErr error,
	blockID flow.Identifier,
	script []byte,
	arguments [][]byte,
	insecureScriptHash [md5.Size]byte,
) {
	// record errors caused by missing local data
	if localErr != nil && status.Code(localErr) == codes.OutOfRange {
		b.metrics.ScriptExecutionNotIndexed()
		b.logScriptExecutionComparison(execNodeResult, execDuration, execErr, localResult, localDuration, localErr, blockID, script, arguments, insecureScriptHash,
			"script execution results do not match EN because data is not indexed yet")
		return
	}

	// check errors first
	if execErr != nil {
		if execErr == localErr {
			b.metrics.ScriptExecutionErrorMatch()
			return
		}

		// error strings generally won't match since the code paths are slightly different
		// check if the underlying error is the same
		if status.Code(execErr) == status.Code(localErr) {
			// both script execution implementations use the same engine, which adds
			// "failed to execute script at block" to the message before returning. Any characters
			// before this can be ignored. The string that comes after is the original error and
			// should match.
			execParts := strings.Split(execErr.Error(), "failed to execute script at block")
			localParts := strings.Split(localErr.Error(), "failed to execute script at block")
			if len(execParts) == 2 && len(localParts) == 2 {
				if execParts[1] == localParts[1] {
					b.metrics.ScriptExecutionErrorMatch()
					return
				}

				// by default ENs are configured with longer script error size limits. check if the
				// non-truncated parts match
				localSubParts := strings.Split(localParts[1], " ... ")
				if len(localSubParts) == 2 &&
					strings.HasPrefix(execParts[1], localSubParts[0]) &&
					strings.HasSuffix(execParts[1], localSubParts[1]) {
					b.metrics.ScriptExecutionErrorMatch()
					return
				}
			}
		}

		b.metrics.ScriptExecutionErrorMismatch()
		b.logScriptExecutionComparison(execNodeResult, execDuration, execErr, localResult, localDuration, localErr, blockID, script, arguments, insecureScriptHash,
			"cadence errors from local execution do not match and EN")
		return
	}

	if bytes.Equal(execNodeResult, localResult) {
		b.metrics.ScriptExecutionResultMatch()
		return
	}

	b.metrics.ScriptExecutionResultMismatch()
	b.logScriptExecutionComparison(execNodeResult, execDuration, execErr, localResult, localDuration, localErr, blockID, script, arguments, insecureScriptHash,
		"script execution results from local execution do not match EN")
}

// logScriptExecutionComparison logs the script execution comparison between local execution and execution node
func (b *backendScripts) logScriptExecutionComparison(
	execNodeResult []byte,
	execDuration time.Duration,
	execErr error,
	localResult []byte,
	localDuration time.Duration,
	localErr error,
	blockID flow.Identifier,
	script []byte,
	arguments [][]byte,
	insecureScriptHash [md5.Size]byte,
	msg string,
) {
	args := make([]string, len(arguments))
	for i, arg := range arguments {
		args[i] = string(arg)
	}

	lgCtx := b.log.With().
		Hex("block_id", blockID[:]).
		Hex("script_hash", insecureScriptHash[:]).
		Str("script", string(script)).
		Strs("args", args)

	if execErr != nil {
		lgCtx = lgCtx.AnErr("execution_node_error", execErr)
	} else {
		lgCtx = lgCtx.Hex("execution_node_result", execNodeResult)
	}
	lgCtx = lgCtx.Dur("execution_node_duration_ms", execDuration)

	if localErr != nil {
		lgCtx = lgCtx.AnErr("local_error", localErr)
	} else {
		lgCtx = lgCtx.Hex("local_result", localResult)
	}
	lgCtx = lgCtx.Dur("local_duration_ms", localDuration)

	lg := lgCtx.Logger()
	lg.Debug().Msg(msg)
}

// shouldLogScript checks if the script hash is unique in the time window
func (b *backendScripts) shouldLogScript(execTime time.Time, scriptHash [md5.Size]byte) bool {
	timestamp, seen := b.loggedScripts.Get(scriptHash)
	if seen {
		return execTime.Sub(timestamp) >= uniqueScriptLoggingTimeWindow
	}
	return true
}

// convertScriptExecutionError converts the script execution error to a gRPC error
func convertScriptExecutionError(err error, height uint64) error {
	if err == nil {
		return nil
	}

	var coded fvmerrors.CodedError
	if fvmerrors.As(err, &coded) {
		// general FVM/ledger errors
		if coded.Code().IsFailure() {
			return rpc.ConvertError(err, "failed to execute script", codes.Internal)
		}

		// runtime errors
		return status.Errorf(codes.InvalidArgument, "failed to execute script: %v", err)
	}

	return convertIndexError(err, height, "failed to execute script")
}

// convertIndexError converts errors related to index to a gRPC error
func convertIndexError(err error, height uint64, defaultMsg string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, execution.ErrDataNotAvailable) {
		return status.Errorf(codes.OutOfRange, "data for block height %d is not available", height)
	}

	if errors.Is(err, storage.ErrNotFound) {
		return status.Errorf(codes.NotFound, "data not found: %v", err)
	}

	return rpc.ConvertError(err, defaultMsg, codes.Internal)
}
