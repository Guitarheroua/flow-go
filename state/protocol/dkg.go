package protocol

import (
	"github.com/onflow/crypto"

	"github.com/onflow/flow-go/model/flow"
)

// DKG represents the result of running the distributed key generation
// procedure for the random beacon.
type DKG interface {

	// Size is the number of members in the DKG.
	Size() uint

	// GroupKey is the group public key.
	GroupKey() crypto.PublicKey

	// Index returns the index for the given node.
	// Error Returns:
	// * protocol.IdentityNotFoundError if nodeID is not a valid DKG participant.
	Index(nodeID flow.Identifier) (uint, error)

	// KeyShare returns the public key share for the given node.
	// Error Returns:
	// * protocol.IdentityNotFoundError if nodeID is not a valid DKG participant.
	KeyShare(nodeID flow.Identifier) (crypto.PublicKey, error)

	// KeyShares returns all public key shares that are result of the distributed key generation.
	KeyShares() []crypto.PublicKey

	// NodeID returns the node identifier for the given index.
	// An exception is returned if the index is ≥ Size().
	// Intended for use outside the hotpath, with runtime
	// scaling linearly in the number of DKG participants (ie. Size()).
	NodeID(index uint) (flow.Identifier, error)
}
