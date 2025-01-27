package mocks

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/state/protocol"
	"github.com/onflow/flow-go/state/protocol/invalid"
)

// EpochQuery implements protocol.EpochQuery for testing purposes.
// Safe for concurrent use by multiple goroutines.
// Only supports committed epochs.
type EpochQuery struct {
	t         *testing.T
	mu        sync.RWMutex
	counter   uint64                             // represents the current epoch
	byCounter map[uint64]protocol.CommittedEpoch // all epochs
}

func NewEpochQuery(t *testing.T, counter uint64, epochs ...protocol.CommittedEpoch) *EpochQuery {
	mock := &EpochQuery{
		t:         t,
		counter:   counter,
		byCounter: make(map[uint64]protocol.CommittedEpoch),
	}

	for _, epoch := range epochs {
		mock.Add(epoch)
	}

	return mock
}

func (mock *EpochQuery) Current() protocol.CommittedEpoch {
	mock.mu.RLock()
	defer mock.mu.RUnlock()
	return mock.byCounter[mock.counter]
}

func (mock *EpochQuery) NextUnsafe() protocol.TentativeEpoch {
	mock.mu.RLock()
	defer mock.mu.RUnlock()
	epoch, exists := mock.byCounter[mock.counter+1]
	if !exists {
		return invalid.NewEpoch(protocol.ErrNextEpochNotSetup)
	}
	return epoch
}

func (mock *EpochQuery) NextCommitted() protocol.CommittedEpoch {
	mock.mu.RLock()
	defer mock.mu.RUnlock()
	epoch, exists := mock.byCounter[mock.counter+1]
	if !exists {
		return invalid.NewEpoch(protocol.ErrNextEpochNotSetup)
	}
	return epoch
}

func (mock *EpochQuery) Previous() protocol.CommittedEpoch {
	mock.mu.RLock()
	defer mock.mu.RUnlock()
	epoch, exists := mock.byCounter[mock.counter-1]
	if !exists {
		return invalid.NewEpoch(protocol.ErrNoPreviousEpoch)
	}
	return epoch
}

// Phase returns a phase consistent with the current epoch state.
func (mock *EpochQuery) Phase() flow.EpochPhase {
	mock.mu.RLock()
	defer mock.mu.RUnlock()
	_, exists := mock.byCounter[mock.counter+1]
	if exists {
		return flow.EpochPhaseCommitted
	}
	return flow.EpochPhaseStaking
}

func (mock *EpochQuery) ByCounter(counter uint64) protocol.CommittedEpoch {
	mock.mu.RLock()
	defer mock.mu.RUnlock()
	return mock.byCounter[counter]
}

func (mock *EpochQuery) Transition() {
	mock.mu.Lock()
	defer mock.mu.Unlock()
	mock.counter++
}

func (mock *EpochQuery) Add(epoch protocol.CommittedEpoch) {
	mock.mu.Lock()
	defer mock.mu.Unlock()
	counter, err := epoch.Counter()
	require.NoError(mock.t, err, "cannot add epoch with invalid counter")
	mock.byCounter[counter] = epoch
}
