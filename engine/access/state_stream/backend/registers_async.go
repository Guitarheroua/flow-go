package backend

import (
	"fmt"
	"sync"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/state_synchronization"
	"github.com/onflow/flow-go/storage"
	"go.uber.org/atomic"
)

// RegistersAsyncStore has the same basic structure as access/backend.ScriptExecutor
type RegistersAsyncStore struct {
	registerIndex storage.RegisterIndex
	reporter      state_synchronization.IndexReporter
	initialized   *atomic.Bool
	init          sync.Once
}

func NewRegistersAsyncStore() *RegistersAsyncStore {
	return &RegistersAsyncStore{
		initialized: atomic.NewBool(false),
	}
}

// InitDataAvailable follows the same pattern of backend.ScriptExecutor
// This method can be called at any time after the RegistersAsyncStore object is created and
// calls to GetRegisterValues will return a storage.ErrHeightNotIndexed,
// since we can't disambiguate between the underlying store before bootstrapping or just simply being behind sync
func (r *RegistersAsyncStore) InitDataAvailable(
	indexReporter state_synchronization.IndexReporter,
	registers storage.RegisterIndex,
) {
	r.init.Do(func() {
		defer r.initialized.Store(true)
		r.reporter = indexReporter
		r.registerIndex = registers
	})
}

// RegisterValues atomically gets the register values from the underlying storage.RegisterIndex
// Expected errors:
//   - storage.ErrHeightNotIndexed if the store is still bootstrapping or if the values at the height is not indexed yet
func (r *RegistersAsyncStore) RegisterValues(ids flow.RegisterIDs, height uint64) ([]flow.RegisterValue, error) {
	if !r.isDataAvailable(height) {
		return nil, storage.ErrHeightNotIndexed
	}
	result := make([]flow.RegisterValue, len(ids))
	for idx, regId := range ids {
		val, err := r.registerIndex.Get(regId, height)
		if err != nil {
			return nil, fmt.Errorf("failed to get register value for id %s: %w", regId.String(), err)
		}
		result[idx] = val
	}
	return result, nil
}

func (r *RegistersAsyncStore) isDataAvailable(height uint64) bool {
	return r.initialized.Load() &&
		height <= r.reporter.HighestIndexedHeight() && height >= r.reporter.LowestIndexedHeight()
}
