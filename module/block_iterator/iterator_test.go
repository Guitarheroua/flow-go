package block_iterator

import (
	"fmt"
	"testing"

	"github.com/dgraph-io/badger/v2"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/metrics"
	storagebadger "github.com/onflow/flow-go/storage/badger"
	"github.com/onflow/flow-go/storage/badger/operation"
	"github.com/onflow/flow-go/utils/unittest"
)

func TestIterateHeight(t *testing.T) {
	unittest.RunWithBadgerDB(t, func(db *badger.DB) {
		// create blocks with siblings
		b1 := &flow.Header{Height: 1}
		b2 := &flow.Header{Height: 2}
		b3 := &flow.Header{Height: 3}
		bs := []*flow.Header{b1, b2, b3}

		// index height
		for _, b := range bs {
			require.NoError(t, db.Update(operation.IndexBlockHeight(b.Height, b.ID())))
		}

		progress := &saveNextHeight{}

		// create iterator
		// b0 is the root block, iterate from b1 to b3
		iterRange := module.IterateRange{Start: b1.Height, End: b3.Height}
		headers := storagebadger.NewHeaders(&metrics.NoopCollector{}, db)
		getBlockIDByIndex := func(height uint64) (flow.Identifier, bool, error) {
			blockID, err := headers.BlockIDByHeight(height)
			if err != nil {
				return flow.ZeroID, false, err
			}

			return blockID, true, nil
		}
		iter := NewIndexedBlockIterator(getBlockIDByIndex, progress, iterRange)

		// iterate through all blocks
		visited := make(map[flow.Identifier]struct{})
		for {
			id, ok, err := iter.Next()
			require.NoError(t, err)
			if !ok {
				break
			}

			// preventing duplicate visit
			_, ok = visited[id]
			require.False(t, ok, fmt.Sprintf("block %v is visited twice", id))

			visited[id] = struct{}{}
		}

		// verify all blocks are visited
		for _, b := range bs {
			_, ok := visited[b.ID()]
			require.True(t, ok, fmt.Sprintf("block %v is not visited", b.ID()))
			delete(visited, b.ID())
		}
		require.Empty(t, visited)

		// save the next to iterate height and verify

		require.NoError(t, iter.Checkpoint())

		savedNextHeight, err := progress.LoadState()
		require.NoError(t, err)

		require.Equal(t, b3.Height+1, savedNextHeight,
			fmt.Sprintf("saved next height should be %v, but got %v", b3.Height, savedNextHeight))

	})
}

type saveNextHeight struct {
	savedNextHeight uint64
}

var _ module.IterateProgressWriter = (*saveNextHeight)(nil)
var _ module.IterateProgressReader = (*saveNextHeight)(nil)

func (s *saveNextHeight) SaveState(height uint64) error {
	s.savedNextHeight = height
	return nil
}

func (s *saveNextHeight) LoadState() (uint64, error) {
	return s.savedNextHeight, nil
}
