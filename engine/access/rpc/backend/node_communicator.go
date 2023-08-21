package backend

import (
	"github.com/hashicorp/go-multierror"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onflow/flow-go/model/flow"
)

// maxFailedRequestCount represents the maximum number of failed requests before returning errors.
const maxFailedRequestCount = 3

// NodeAction is a callback function type that represents an action to be performed on a node.
// It takes a node as input and returns an error indicating the result of the action.
type NodeAction func(node *flow.Identity) error

// ErrorTerminator is a callback function that determines whether an error should terminate further execution.
// It takes an error as input and returns a boolean value indicating whether the error should be considered terminal.
type ErrorTerminator func(node *flow.Identity, err error) bool

type Communicator interface {
	CallAvailableNode(
		nodes flow.IdentityList,
		call NodeAction,
		shouldTerminateOnError ErrorTerminator,
	) error
}

var _ Communicator = (*NodeCommunicator)(nil)

// NodeCommunicator is responsible for calling available nodes in the backend.
type NodeCommunicator struct {
	nodeSelectorFactory NodeSelectorFactory
}

// NewNodeCommunicator creates a new instance of NodeCommunicator.
func NewNodeCommunicator(circuitBreakerEnabled bool) *NodeCommunicator {
	return &NodeCommunicator{
		nodeSelectorFactory: NodeSelectorFactory{circuitBreakerEnabled: circuitBreakerEnabled},
	}
}

// CallAvailableNode calls the provided function on the available nodes.
// It iterates through the nodes and executes the function.
// If an error occurs, it applies the custom error terminator (if provided) and keeps track of the errors.
// If the error occurs in circuit breaker, it continues to the next node.
// If the maximum failed request count is reached, it returns the accumulated errors.
func (b *NodeCommunicator) CallAvailableNode(
	nodes flow.IdentityList,
	call NodeAction,
	shouldTerminateOnError ErrorTerminator,
) error {
	var errs *multierror.Error
	nodeSelector, err := b.nodeSelectorFactory.SelectNodes(nodes)
	if err != nil {
		return err
	}

	for node := nodeSelector.Next(); node != nil; node = nodeSelector.Next() {
		err := call(node)
		if err == nil {
			return nil
		}

		if shouldTerminateOnError != nil && shouldTerminateOnError(node, err) {
			return err
		}

		if err == gobreaker.ErrOpenState {
			if !nodeSelector.HasNext() && len(errs.Errors) == 0 {
				errs = multierror.Append(errs, status.Error(codes.Unavailable, "there are no available nodes"))
			}
			continue
		}

		errs = multierror.Append(errs, err)
		if len(errs.Errors) >= maxFailedRequestCount {
			return errs.ErrorOrNil()
		}
	}

	return errs.ErrorOrNil()
}
