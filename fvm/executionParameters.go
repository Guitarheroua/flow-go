package fvm

import (
	"context"
	"fmt"
	"math"

	"github.com/rs/zerolog"

	"github.com/onflow/flow-go/fvm/storage/snapshot"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/common"

	"github.com/onflow/flow-go/fvm/blueprints"
	"github.com/onflow/flow-go/fvm/environment"
	"github.com/onflow/flow-go/fvm/errors"
	"github.com/onflow/flow-go/fvm/meter"
	"github.com/onflow/flow-go/fvm/storage"
	"github.com/onflow/flow-go/fvm/storage/derived"
	"github.com/onflow/flow-go/fvm/storage/state"
	"github.com/onflow/flow-go/fvm/systemcontracts"
)

func ProcedureStateParameters(
	ctx Context,
	proc Procedure,
) state.StateParameters {
	return state.DefaultParameters().
		WithMeterParameters(getBasicMeterParameters(ctx, proc)).
		WithMaxKeySizeAllowed(ctx.MaxStateKeySize).
		WithMaxValueSizeAllowed(ctx.MaxStateValueSize)
}

// getBasicMeterParameters returns the set of meter parameters used for
// general procedure execution.  Subparts of the procedure execution may
// specify custom meter parameters via nested transactions.
func getBasicMeterParameters(
	ctx Context,
	proc Procedure,
) meter.MeterParameters {
	params := meter.DefaultParameters().
		WithComputationLimit(uint(proc.ComputationLimit(ctx))).
		WithMemoryLimit(proc.MemoryLimit(ctx)).
		WithEventEmitByteLimit(ctx.EventCollectionByteSizeLimit).
		WithStorageInteractionLimit(ctx.MaxStateInteractionSize)

	// NOTE: The memory limit (and interaction limit) may be overridden by the
	// environment.  We need to ignore the override in that case.
	if proc.ShouldDisableMemoryAndInteractionLimits(ctx) {
		params = params.WithMemoryLimit(math.MaxUint64).
			WithStorageInteractionLimit(math.MaxUint64)
	}

	return params
}

// getExecutionParameters returns the set of meter parameters used for
// transaction/script body execution and the minimum required version as defined by the
// NodeVersionBeacon contract.
func getExecutionParameters(
	log zerolog.Logger,
	ctx Context,
	proc Procedure,
	txnState storage.TransactionPreparer,
) (state.ExecutionParameters, *snapshot.ExecutionSnapshot, error) {
	meterParams := getBasicMeterParameters(ctx, proc)

	executionParams, executionParamsStateRead, err := txnState.GetStateExecutionParameters(
		txnState,
		NewExecutionParametersComputer(log, ctx, txnState))
	if err != nil {
		return state.ExecutionParameters{
			MeterParameters: meterParams,
		}, nil, err
	}

	if executionParams.ComputationWeights != nil {
		meterParams = meterParams.WithComputationWeights(
			executionParams.ComputationWeights)
	}

	if executionParams.MemoryWeights != nil {
		meterParams = meterParams.WithMemoryWeights(executionParams.MemoryWeights)
	}

	if executionParams.MemoryLimit != nil {
		meterParams = meterParams.WithMemoryLimit(*executionParams.MemoryLimit)
	}

	// NOTE: The memory limit (and interaction limit) may be overridden by the
	// environment.  We need to ignore the override in that case.
	if proc.ShouldDisableMemoryAndInteractionLimits(ctx) {
		meterParams = meterParams.WithMemoryLimit(math.MaxUint64).
			WithStorageInteractionLimit(math.MaxUint64)
	}

	return state.ExecutionParameters{
		MeterParameters: meterParams,
	}, executionParamsStateRead, nil
}

type ExecutionParametersComputer struct {
	log      zerolog.Logger
	ctx      Context
	txnState storage.TransactionPreparer
}

func NewExecutionParametersComputer(
	log zerolog.Logger,
	ctx Context,
	txnState storage.TransactionPreparer,
) ExecutionParametersComputer {
	return ExecutionParametersComputer{
		log:      log,
		ctx:      ctx,
		txnState: txnState,
	}
}

func (computer ExecutionParametersComputer) Compute(
	_ state.NestedTransactionPreparer,
	_ struct{},
) (
	derived.StateExecutionParameters,
	error,
) {
	var overrides derived.StateExecutionParameters
	var err error
	computer.txnState.RunWithAllLimitsDisabled(func() {
		overrides, err = computer.getExecutionParameters()
	})

	if err != nil {
		return overrides, fmt.Errorf(
			"error getting environment meter parameter overrides: %w",
			err)
	}

	return overrides, nil
}

func (computer ExecutionParametersComputer) getExecutionParameters() (
	derived.StateExecutionParameters,
	error,
) {
	sc := systemcontracts.SystemContractsForChain(computer.ctx.Chain.ChainID())

	// The execution parameters are stored in the ExecutionParametersAccount. This is
	// just the service account for all networks except mainnet and testnet.
	// For mainnet and testnet, the execution parameters are stored in a separate
	// account, so that they are separated from the frequently changing data on the
	// service account.
	service := common.Address(sc.ExecutionParametersAccount.Address)

	env := environment.NewScriptEnv(
		context.Background(),
		computer.ctx.TracerSpan,
		computer.ctx.EnvironmentParams,
		computer.txnState)

	overrides := derived.StateExecutionParameters{}

	// set the property if no error, but if the error is a fatal error then
	// return it
	setIfOk := func(prop string, err error, setter func()) (fatal error) {
		err, fatal = errors.SplitErrorTypes(err)
		if fatal != nil {
			// this is a fatal error. return it
			computer.ctx.Logger.
				Error().
				Err(fatal).
				Msgf("error getting %s", prop)
			return fatal
		}
		if err != nil {
			// this is a general error.
			// could be that no setting was present in the state,
			// or that the setting was not parseable,
			// or some other deterministic thing.
			computer.ctx.Logger.
				Debug().
				Err(err).
				Msgf("could not set %s. Using defaults", prop)
			return nil
		}
		// everything is ok. do the setting
		setter()
		return nil
	}

	computationWeights, err := GetExecutionEffortWeights(env, service)
	err = setIfOk(
		"execution effort weights",
		err,
		func() { overrides.ComputationWeights = computationWeights })
	if err != nil {
		return overrides, err
	}

	memoryWeights, err := GetExecutionMemoryWeights(env, service)
	err = setIfOk(
		"execution memory weights",
		err,
		func() { overrides.MemoryWeights = memoryWeights })
	if err != nil {
		return overrides, err
	}

	memoryLimit, err := GetExecutionMemoryLimit(env, service)
	err = setIfOk(
		"execution memory limit",
		err,
		func() { overrides.MemoryLimit = &memoryLimit })
	if err != nil {
		return overrides, err
	}

	return overrides, nil
}

func getExecutionWeights[KindType common.ComputationKind | common.MemoryKind](
	env environment.Environment,
	service common.Address,
	path cadence.Path,
	defaultWeights map[KindType]uint64,
) (
	map[KindType]uint64,
	error,
) {
	runtime := env.BorrowCadenceRuntime()
	defer env.ReturnCadenceRuntime(runtime)

	value, err := runtime.ReadStored(service, path)

	if err != nil {
		// this might be fatal, return as is
		return nil, err
	}

	weightsRaw, ok := cadenceValueToWeights(value)
	if !ok {
		// this is a non-fatal error. It is expected if the weights are not
		// set up on the network yet.
		return nil, errors.NewCouldNotGetExecutionParameterFromStateError(
			service.Hex(),
			path.String())
	}

	// Merge the default weights with the weights from the state.
	// This allows for weights that are not set in the state, to be set by default.
	// In case the network is stuck because of a transaction using an FVM feature that has 0 weight
	// (or is not metered at all), the defaults can be changed and the network restarted
	// instead of trying to change the weights with a transaction.
	weights := make(map[KindType]uint64, len(defaultWeights))
	for k, v := range defaultWeights {
		weights[k] = v
	}
	for k, v := range weightsRaw {
		weights[KindType(k)] = v
	}

	return weights, nil
}

// cadenceValueToWeights converts a cadence value to a map of weights used for
// metering
func cadenceValueToWeights(value cadence.Value) (map[uint]uint64, bool) {
	dict, ok := value.(cadence.Dictionary)
	if !ok {
		return nil, false
	}

	result := make(map[uint]uint64, len(dict.Pairs))
	for _, p := range dict.Pairs {
		key, ok := p.Key.(cadence.UInt64)
		if !ok {
			return nil, false
		}

		value, ok := p.Value.(cadence.UInt64)
		if !ok {
			return nil, false
		}

		result[uint(key)] = uint64(value)
	}

	return result, true
}

// GetExecutionEffortWeights reads stored execution effort weights from the service account
func GetExecutionEffortWeights(
	env environment.Environment,
	service common.Address,
) (
	computationWeights meter.ExecutionEffortWeights,
	err error,
) {
	return getExecutionWeights(
		env,
		service,
		blueprints.TransactionFeesExecutionEffortWeightsPath,
		meter.DefaultComputationWeights)
}

// GetExecutionMemoryWeights reads stored execution memory weights from the service account
func GetExecutionMemoryWeights(
	env environment.Environment,
	service common.Address,
) (
	memoryWeights meter.ExecutionMemoryWeights,
	err error,
) {
	return getExecutionWeights(
		env,
		service,
		blueprints.TransactionFeesExecutionMemoryWeightsPath,
		meter.DefaultMemoryWeights)
}

// GetExecutionMemoryLimit reads the stored execution memory limit from the service account
func GetExecutionMemoryLimit(
	env environment.Environment,
	service common.Address,
) (
	memoryLimit uint64,
	err error,
) {
	runtime := env.BorrowCadenceRuntime()
	defer env.ReturnCadenceRuntime(runtime)

	value, err := runtime.ReadStored(
		service,
		blueprints.TransactionFeesExecutionMemoryLimitPath)
	if err != nil {
		// this might be fatal, return as is
		return 0, err
	}

	memoryLimitRaw, ok := value.(cadence.UInt64)
	if value == nil || !ok {
		// this is a non-fatal error. It is expected if the weights are not set up on the network yet.
		return 0, errors.NewCouldNotGetExecutionParameterFromStateError(
			service.Hex(),
			blueprints.TransactionFeesExecutionMemoryLimitPath.String())
	}

	return uint64(memoryLimitRaw), nil
}
