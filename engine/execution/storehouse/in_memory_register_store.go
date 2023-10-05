package storehouse

import (
	"fmt"
	"sync"

	"github.com/onflow/flow-go/engine/execution"
	"github.com/onflow/flow-go/model/flow"
)

var _ execution.InMemoryRegisterStore = (*InMemoryRegisterStore)(nil)

var ErrPruned = fmt.Errorf("block is pruned")
var ErrNotExecuted = fmt.Errorf("block is not executed")

type InMemoryRegisterStore struct {
	sync.RWMutex
	registersByBlockID map[flow.Identifier]map[flow.RegisterID]flow.RegisterValue // for storing the registers
	parentByBlockID    map[flow.Identifier]flow.Identifier                        // for register updates to be fork-aware
	blockIDsByHeight   map[uint64]map[flow.Identifier]struct{}                    // for pruning
	prunedHeight       uint64
	prunedID           flow.Identifier // to ensure all blocks are extending from pruned block (last finalized and executed block)
}

func NewInMemoryRegisterStore(lastHeight uint64, lastID flow.Identifier) *InMemoryRegisterStore {
	return &InMemoryRegisterStore{
		registersByBlockID: make(map[flow.Identifier]map[flow.RegisterID]flow.RegisterValue),
		parentByBlockID:    make(map[flow.Identifier]flow.Identifier),
		blockIDsByHeight:   make(map[uint64]map[flow.Identifier]struct{}),
		prunedHeight:       lastHeight,
		prunedID:           lastID,
	}
}

// SaveRegisters saves the registers of a block to InMemoryRegisterStore
// It needs to ensure the block is above the pruned height and is connected to the pruned block
func (s *InMemoryRegisterStore) SaveRegisters(
	height uint64,
	blockID flow.Identifier,
	parentID flow.Identifier,
	registers []flow.RegisterEntry,
) error {
	// preprocess data before acquiring the lock
	regs := make(map[flow.RegisterID]flow.RegisterValue)
	for _, reg := range registers {
		regs[reg.Key] = reg.Value
	}

	s.Lock()
	defer s.Unlock()

	// ensure all saved registers are above the pruned height
	if height <= s.prunedHeight {
		return fmt.Errorf("saving pruned registers height %v <= pruned height %v", height, s.prunedHeight)
	}

	// ensure the block is not already saved
	_, ok := s.registersByBlockID[blockID]
	if ok {
		// already exist
		return fmt.Errorf("saving registers for block %s, but it already exists", blockID)
	}

	// make sure parent is a known block or the pruned block, which forms a fork
	_, ok = s.registersByBlockID[parentID]
	if !ok {
		// parent doesn't exist, check if it is the pruned block
		if parentID != s.prunedID {
			return fmt.Errorf("saving registers for block %s, but its parent %s is not saved", blockID, parentID)
		}
	}

	// update registers for the block
	s.registersByBlockID[blockID] = regs

	// update index on parent
	s.parentByBlockID[blockID] = parentID

	// update index on height
	sameHeight, ok := s.blockIDsByHeight[height]
	if !ok {
		sameHeight = make(map[flow.Identifier]struct{})
		s.blockIDsByHeight[height] = sameHeight
	}

	sameHeight[blockID] = struct{}{}
	return nil
}

// GetRegister will return the latest updated value of the given register
// since the pruned height.
// It returns ErrPruned if the register is unknown or not updated since the pruned height
// Can't return ErrNotFound, since we can't distinguish between not found or not updated since the pruned height
func (s *InMemoryRegisterStore) GetRegister(height uint64, blockID flow.Identifier, register flow.RegisterID) (flow.RegisterValue, error) {
	s.RLock()
	defer s.RUnlock()

	if height <= s.prunedHeight {
		return flow.RegisterValue{}, fmt.Errorf("cannot get register at height %d, it is pruned (prunedHeight: %v): %w", height, s.prunedHeight, ErrPruned)
	}

	_, ok := s.registersByBlockID[blockID]
	if !ok {
		return flow.RegisterValue{}, fmt.Errorf("cannot get register at height %d, block %v is not saved: %w", height, blockID, ErrNotExecuted)
	}

	// traverse the fork to find the latest updated value of the given register
	// if not found, it means the register is not updated from the pruned block to the given block
	block := blockID
	for {
		// TODO: do not hold the read lock when reading register from the updated register map
		reg, ok := s.readRegisterAtBlockID(block, register)
		if ok {
			return reg, nil
		}

		// the register didn't get updated at this block, so check its parent

		parent, ok := s.parentByBlockID[block]
		if !ok {
			// if parent doesn't exist, we check if we've reached the pruned block,
			// if so, the register is not found.
			// otherwise, it means the parent block index is not consistent, which is a bug
			// we've reached the pruned block, so the register is not found
			if block == s.prunedID {
				return flow.RegisterValue{}, fmt.Errorf("cannot get register at height %d, block %v is pruned: %w", height, blockID, ErrPruned)
			}

			return flow.RegisterValue{},
				fmt.Errorf("inconsistent parent block index in in-memory-register-store, ancient block %v is not found when getting register at block %v",
					block, blockID)
		}

		block = parent
	}
}

func (s *InMemoryRegisterStore) readRegisterAtBlockID(blockID flow.Identifier, register flow.RegisterID) (flow.RegisterValue, bool) {
	registers, ok := s.registersByBlockID[blockID]
	if !ok {
		return flow.RegisterValue{}, false
	}

	reg, ok := registers[register]
	return reg, ok
}

// GetUpdatedRegisters returns the updated registers of a block
func (s *InMemoryRegisterStore) GetUpdatedRegisters(height uint64, blockID flow.Identifier) ([]flow.RegisterEntry, error) {
	s.RLock()
	defer s.RUnlock()
	if height <= s.prunedHeight {
		return nil, fmt.Errorf("cannot get register at height %d, it is pruned %v", height, s.prunedHeight)
	}

	registerUpdates, ok := s.registersByBlockID[blockID]
	if !ok {
		return nil, fmt.Errorf("cannot get register at height %d, block %s is not found", height, blockID)
	}

	// convert from map to into slice
	registers := make([]flow.RegisterEntry, 0, len(registerUpdates))
	for regID, reg := range registerUpdates {
		registers = append(registers, flow.RegisterEntry{
			Key:   regID,
			Value: reg,
		})
	}

	return registers, nil
}

// Prune prunes the register store to the given height
// The pruned height must be an executed block, the caller should ensure that by calling SaveRegisters before.
// TODO: It does not block the caller, the pruning work is done async
func (s *InMemoryRegisterStore) Prune(height uint64, blockID flow.Identifier) error {
	finalizedFork, err := s.findFinalizedFork(height, blockID)
	if err != nil {
		return fmt.Errorf("cannot find finalized fork: %w", err)
	}

	s.Lock()
	defer s.Unlock()

	for i := len(finalizedFork) - 1; i >= 0; i-- {
		// traverse from lower height to higher height
		blockID := finalizedFork[i]

		err := s.pruneByHeight(s.prunedHeight+1, blockID)
		if err != nil {
			return fmt.Errorf("could not prune by height %v: %w", s.prunedHeight+1, err)
		}
	}

	return nil
}

func (s *InMemoryRegisterStore) PrunedHeight() uint64 {
	s.RLock()
	defer s.RUnlock()
	return s.prunedHeight
}

func (s *InMemoryRegisterStore) IsBlockExecuted(height uint64, blockID flow.Identifier) (bool, error) {
	s.RLock()
	defer s.RUnlock()

	// finalized and executed blocks are pruned
	if height <= s.prunedHeight {
		return false, fmt.Errorf("below pruned height")
	}

	_, ok := s.registersByBlockID[blockID]
	return ok, nil
}

// findFinalizedFork returns the finalized fork from higher height to lower height
// the last block's height is s.prunedHeight + 1
func (s *InMemoryRegisterStore) findFinalizedFork(height uint64, blockID flow.Identifier) ([]flow.Identifier, error) {
	s.RLock()
	defer s.RUnlock()

	if height <= s.prunedHeight {
		return nil, fmt.Errorf("cannot find finalized fork at height %d, it is pruned (prunedHeight: %v)", height, s.prunedHeight)
	}
	prunedHeight := height
	block := blockID

	// finalized fork from pruned height to the last finalized height
	fork := make([]flow.Identifier, 0, height-s.prunedHeight)
	for {
		fork = append(fork, block)
		prunedHeight--

		parent, ok := s.parentByBlockID[block]
		if !ok {
			return nil, fmt.Errorf("inconsistent parent block index in in-memory-register-store, ancient block %s is not found when finding finalized fork at height %v", block, height)
		}
		if parent == s.prunedID {
			break
		}
		block = parent
	}

	if prunedHeight != s.prunedHeight {
		return nil, fmt.Errorf("inconsistent parent block index in in-memory-register-store, pruned height %d is not equal to %d", prunedHeight, s.prunedHeight)
	}

	return fork, nil
}

func (s *InMemoryRegisterStore) pruneByHeight(height uint64, finalized flow.Identifier) error {
	s.removeBlock(height, finalized)

	// remove conflicting forks
	for blockID := range s.blockIDsByHeight[height] {
		s.pruneFork(height, blockID)
	}

	if len(s.blockIDsByHeight[height]) > 0 {
		return fmt.Errorf("all forks on the same height should have been pruend, but actually not: %v", len(s.blockIDsByHeight[height]))
	}

	delete(s.blockIDsByHeight, height)
	s.prunedHeight = height
	s.prunedID = finalized
	return nil
}

func (s *InMemoryRegisterStore) removeBlock(height uint64, blockID flow.Identifier) {
	delete(s.registersByBlockID, blockID)
	delete(s.parentByBlockID, blockID)
	delete(s.blockIDsByHeight[height], blockID)
}

func (s *InMemoryRegisterStore) pruneFork(height uint64, blockID flow.Identifier) {
	s.removeBlock(height, blockID)
	// all its children must be at height + 1, whose parent is blockID

	nextHeight := height + 1
	blocksAtNextHeight, ok := s.blockIDsByHeight[nextHeight]
	if !ok {
		return
	}

	for block := range blocksAtNextHeight {
		isChild := s.parentByBlockID[block] == blockID
		if isChild {
			s.pruneFork(nextHeight, block)
		}
	}
}
