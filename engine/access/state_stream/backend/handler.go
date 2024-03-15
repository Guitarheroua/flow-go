package backend

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onflow/flow/protobuf/go/flow/entities"
	"github.com/onflow/flow/protobuf/go/flow/executiondata"

	"github.com/onflow/flow-go/engine/access/state_stream"
	"github.com/onflow/flow-go/engine/access/subscription"
	"github.com/onflow/flow-go/engine/common/rpc"
	"github.com/onflow/flow-go/engine/common/rpc/convert"
	"github.com/onflow/flow-go/model/flow"
)

type Handler struct {
	subscription.StreamingData

	api   state_stream.API
	chain flow.Chain

	eventFilterConfig        state_stream.EventFilterConfig
	defaultHeartbeatInterval uint64
}

var _ executiondata.ExecutionDataAPIServer = (*Handler)(nil)

func NewHandler(api state_stream.API, chain flow.Chain, config Config) *Handler {
	h := &Handler{
		StreamingData:            subscription.NewStreamingData(config.MaxGlobalStreams),
		api:                      api,
		chain:                    chain,
		eventFilterConfig:        config.EventFilterConfig,
		defaultHeartbeatInterval: config.HeartbeatInterval,
	}
	return h
}

func (h *Handler) GetExecutionDataByBlockID(ctx context.Context, request *executiondata.GetExecutionDataByBlockIDRequest) (*executiondata.GetExecutionDataByBlockIDResponse, error) {
	blockID, err := convert.BlockID(request.GetBlockId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not convert block ID: %v", err)
	}

	execData, err := h.api.GetExecutionDataByBlockID(ctx, blockID)
	if err != nil {
		return nil, rpc.ConvertError(err, "could no get execution data", codes.Internal)
	}

	message, err := convert.BlockExecutionDataToMessage(execData)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not convert execution data to entity: %v", err)
	}

	err = convert.BlockExecutionDataEventPayloadsToVersion(message, request.GetEventEncodingVersion())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not convert execution data event payloads to JSON: %v", err)
	}

	return &executiondata.GetExecutionDataByBlockIDResponse{BlockExecutionData: message}, nil
}

func (h *Handler) SubscribeExecutionData(request *executiondata.SubscribeExecutionDataRequest, stream executiondata.ExecutionDataAPI_SubscribeExecutionDataServer) error {
	// check if the maximum number of streams is reached
	if h.StreamCount.Load() >= h.MaxStreams {
		return status.Errorf(codes.ResourceExhausted, "maximum number of streams reached")
	}
	h.StreamCount.Add(1)
	defer h.StreamCount.Add(-1)

	startBlockID := flow.ZeroID
	if request.GetStartBlockId() != nil {
		blockID, err := convert.BlockID(request.GetStartBlockId())
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "could not convert start block ID: %v", err)
		}
		startBlockID = blockID
	}

	sub := h.api.SubscribeExecutionData(stream.Context(), startBlockID, request.GetStartBlockHeight())

	for {
		v, ok := <-sub.Channel()
		if !ok {
			if sub.Err() != nil {
				return rpc.ConvertError(sub.Err(), "stream encountered an error", codes.Internal)
			}
			return nil
		}

		resp, ok := v.(*ExecutionDataResponse)
		if !ok {
			return status.Errorf(codes.Internal, "unexpected response type: %T", v)
		}

		execData, err := convert.BlockExecutionDataToMessage(resp.ExecutionData)
		if err != nil {
			return status.Errorf(codes.Internal, "could not convert execution data to entity: %v", err)
		}

		err = convert.BlockExecutionDataEventPayloadsToVersion(execData, request.GetEventEncodingVersion())
		if err != nil {
			return status.Errorf(codes.Internal, "could not convert execution data event payloads to JSON: %v", err)
		}

		err = stream.Send(&executiondata.SubscribeExecutionDataResponse{
			BlockHeight:        resp.Height,
			BlockExecutionData: execData,
		})
		if err != nil {
			return rpc.ConvertError(err, "could not send response", codes.Internal)
		}
	}
}

func (h *Handler) SubscribeEvents(request *executiondata.SubscribeEventsRequest, stream executiondata.ExecutionDataAPI_SubscribeEventsServer) error {
	// check if the maximum number of streams is reached
	if h.StreamCount.Load() >= h.MaxStreams {
		return status.Errorf(codes.ResourceExhausted, "maximum number of streams reached")
	}
	h.StreamCount.Add(1)
	defer h.StreamCount.Add(-1)

	startBlockID := flow.ZeroID
	if request.GetStartBlockId() != nil {
		blockID, err := convert.BlockID(request.GetStartBlockId())
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "could not convert start block ID: %v", err)
		}
		startBlockID = blockID
	}

	filter := state_stream.EventFilter{}
	if request.GetFilter() != nil {
		var err error
		reqFilter := request.GetFilter()
		filter, err = state_stream.NewEventFilter(
			h.eventFilterConfig,
			h.chain,
			reqFilter.GetEventType(),
			reqFilter.GetAddress(),
			reqFilter.GetContract(),
		)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "invalid event filter: %v", err)
		}
	}

	sub := h.api.SubscribeEvents(stream.Context(), startBlockID, request.GetStartBlockHeight(), filter)

	heartbeatInterval := request.HeartbeatInterval
	if heartbeatInterval == 0 {
		heartbeatInterval = h.defaultHeartbeatInterval
	}

	blocksSinceLastMessage := uint64(0)
	for {
		v, ok := <-sub.Channel()
		if !ok {
			if sub.Err() != nil {
				return rpc.ConvertError(sub.Err(), "stream encountered an error", codes.Internal)
			}
			return nil
		}

		resp, ok := v.(*EventsResponse)
		if !ok {
			return status.Errorf(codes.Internal, "unexpected response type: %T", v)
		}

		// check if there are any events in the response. if not, do not send a message unless the last
		// response was more than HeartbeatInterval blocks ago
		if len(resp.Events) == 0 {
			blocksSinceLastMessage++
			if blocksSinceLastMessage < heartbeatInterval {
				continue
			}
			blocksSinceLastMessage = 0
		}

		// BlockExecutionData contains CCF encoded events, and the Access API returns JSON-CDC events.
		// convert event payload formats.
		// This is a temporary solution until the Access API supports specifying the encoding in the request
		events, err := convert.EventsToMessagesWithEncodingConversion(resp.Events, entities.EventEncodingVersion_CCF_V0, request.GetEventEncodingVersion())
		if err != nil {
			return status.Errorf(codes.Internal, "could not convert events to entity: %v", err)
		}

		err = stream.Send(&executiondata.SubscribeEventsResponse{
			BlockHeight: resp.Height,
			BlockId:     convert.IdentifierToMessage(resp.BlockID),
			Events:      events,
		})
		if err != nil {
			return rpc.ConvertError(err, "could not send response", codes.Internal)
		}
	}
}

func (h *Handler) GetRegisterValues(_ context.Context, request *executiondata.GetRegisterValuesRequest) (*executiondata.GetRegisterValuesResponse, error) {
	// Convert data
	registerIDs, err := convert.MessagesToRegisterIDs(request.GetRegisterIds(), h.chain)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "could not convert register IDs: %v", err)
	}

	// get payload from store
	values, err := h.api.GetRegisterValues(registerIDs, request.GetBlockHeight())
	if err != nil {
		return nil, rpc.ConvertError(err, "could not get register values", codes.Internal)
	}

	return &executiondata.GetRegisterValuesResponse{Values: values}, nil
}

func (h *Handler) createAccountStatusesFilter(statusFilter *executiondata.StatusFilter) (state_stream.EventFilter, error) {
	filter := state_stream.EventFilter{}
	if statusFilter != nil {
		var err error
		filter, err = state_stream.NewEventFilter(
			h.eventFilterConfig,
			h.chain,
			statusFilter.GetEventType(),
			[]string{},
			[]string{},
		)
		if err != nil {
			return filter, status.Errorf(codes.InvalidArgument, "invalid event filter: %v", err)
		}

		addresses := statusFilter.GetAddress()
		if len(addresses) > 0 {
			var fieldFilters []state_stream.FieldFilter
			for _, address := range addresses {
				fieldFilters = append(fieldFilters, state_stream.FieldFilter{
					FieldName:   "address",
					TargetValue: address,
				})
			}

			for eventType := range filter.EventTypes {
				filter.EventFieldFilters[eventType] = fieldFilters
			}
		}
	}

	return filter, nil
}

func convertAccountsStatusesResults(eventVersion entities.EventEncodingVersion, resp *AccountStatusesResponse) ([]*executiondata.SubscribeAccountStatusesResponse_Result, error) {
	var results []*executiondata.SubscribeAccountStatusesResponse_Result
	for address, events := range resp.AccountEvents {
		convertedEvent, err := convert.EventsToMessagesWithEncodingConversion(events, entities.EventEncodingVersion_CCF_V0, eventVersion)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not convert events to entity: %v", err)
		}

		results = append(results, &executiondata.SubscribeAccountStatusesResponse_Result{
			Address: []byte(address),
			Events:  convertedEvent,
		})
	}
	return results, nil
}

func (h *Handler) SubscribeAccountStatusesFromStartBlockID(request *executiondata.SubscribeAccountStatusesFromStartBlockIDRequest, stream executiondata.ExecutionDataAPI_SubscribeAccountStatusesFromStartBlockIDServer) error {
	// check if the maximum number of streams is reached
	if h.streamCount.Load() >= h.maxStreams {
		return status.Errorf(codes.ResourceExhausted, "maximum number of streams reached")
	}

	h.streamCount.Add(1)
	defer h.streamCount.Add(-1)

	startBlockID := flow.ZeroID
	if request.GetStartBlockId() != nil {
		blockID, err := convert.BlockID(request.GetStartBlockId())
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "could not convert start block ID: %v", err)
		}
		startBlockID = blockID
	}

	filter, err := h.createAccountStatusesFilter(request.GetFilter())
	if err != nil {
		return err
	}

	sub := h.api.SubscribeAccountStatuses(stream.Context(), startBlockID, 0, filter)

	heartbeatInterval := request.HeartbeatInterval
	if heartbeatInterval == 0 {
		heartbeatInterval = h.defaultHeartbeatInterval
	}

	blocksSinceLastMessage := uint64(0)
	for {
		v, ok := <-sub.Channel()
		if !ok {
			if sub.Err() != nil {
				return rpc.ConvertError(sub.Err(), "stream encountered an error", codes.Internal)
			}
			return nil
		}

		resp, ok := v.(*AccountStatusesResponse)
		if !ok {
			return status.Errorf(codes.Internal, "unexpected response type: %T", v)
		}

		// check if there are any events in the response. if not, do not send a message unless the last
		// response was more than HeartbeatInterval blocks ago
		if len(resp.AccountEvents) == 0 {
			blocksSinceLastMessage++
			if blocksSinceLastMessage < heartbeatInterval {
				continue
			}
			blocksSinceLastMessage = 0
		}

		// BlockExecutionData contains CCF encoded events, and the Access API returns JSON-CDC events.
		// convert event payload formats.
		// This is a temporary solution until the Access API supports specifying the encoding in the request
		results, err := convertAccountsStatusesResults(request.GetEventEncodingVersion(), resp)
		if err != nil {
			return err
		}

		err = stream.Send(&executiondata.SubscribeAccountStatusesResponse{
			BlockId:      convert.IdentifierToMessage(resp.BlockID),
			BlockHeight:  resp.Height,
			Results:      results,
			MessageIndex: resp.MessageIndex,
		})
		if err != nil {
			return rpc.ConvertError(err, "could not send response", codes.Internal)
		}
	}
}

func (h *Handler) SubscribeAccountStatusesFromStartHeight(request *executiondata.SubscribeAccountStatusesFromStartHeightRequest, stream executiondata.ExecutionDataAPI_SubscribeAccountStatusesFromStartHeightServer) error {
	// check if the maximum number of streams is reached
	if h.streamCount.Load() >= h.maxStreams {
		return status.Errorf(codes.ResourceExhausted, "maximum number of streams reached")
	}

	h.streamCount.Add(1)
	defer h.streamCount.Add(-1)

	filter, err := h.createAccountStatusesFilter(request.GetFilter())
	if err != nil {
		return err
	}

	sub := h.api.SubscribeAccountStatuses(stream.Context(), flow.ZeroID, request.GetStartBlockHeight(), filter)

	heartbeatInterval := request.HeartbeatInterval
	if heartbeatInterval == 0 {
		heartbeatInterval = h.defaultHeartbeatInterval
	}

	blocksSinceLastMessage := uint64(0)
	for {
		v, ok := <-sub.Channel()
		if !ok {
			if sub.Err() != nil {
				return rpc.ConvertError(sub.Err(), "stream encountered an error", codes.Internal)
			}
			return nil
		}

		resp, ok := v.(*AccountStatusesResponse)
		if !ok {
			return status.Errorf(codes.Internal, "unexpected response type: %T", v)
		}

		// check if there are any events in the response. if not, do not send a message unless the last
		// response was more than HeartbeatInterval blocks ago
		if len(resp.AccountEvents) == 0 {
			blocksSinceLastMessage++
			if blocksSinceLastMessage < heartbeatInterval {
				continue
			}
			blocksSinceLastMessage = 0
		}

		// BlockExecutionData contains CCF encoded events, and the Access API returns JSON-CDC events.
		// convert event payload formats.
		// This is a temporary solution until the Access API supports specifying the encoding in the request
		results, err := convertAccountsStatusesResults(request.GetEventEncodingVersion(), resp)
		if err != nil {
			return err
		}

		err = stream.Send(&executiondata.SubscribeAccountStatusesResponse{
			BlockId:      convert.IdentifierToMessage(resp.BlockID),
			BlockHeight:  resp.Height,
			Results:      results,
			MessageIndex: resp.MessageIndex,
		})
		if err != nil {
			return rpc.ConvertError(err, "could not send response", codes.Internal)
		}
	}
}

func (h *Handler) SubscribeAccountStatusesFromLatestBlock(request *executiondata.SubscribeAccountStatusesFromLatestBlockRequest, stream executiondata.ExecutionDataAPI_SubscribeAccountStatusesFromLatestBlockServer) error {
	// check if the maximum number of streams is reached
	if h.streamCount.Load() >= h.maxStreams {
		return status.Errorf(codes.ResourceExhausted, "maximum number of streams reached")
	}

	h.streamCount.Add(1)
	defer h.streamCount.Add(-1)

	filter, err := h.createAccountStatusesFilter(request.GetFilter())
	if err != nil {
		return err
	}

	sub := h.api.SubscribeAccountStatuses(stream.Context(), flow.ZeroID, 0, filter)

	heartbeatInterval := request.HeartbeatInterval
	if heartbeatInterval == 0 {
		heartbeatInterval = h.defaultHeartbeatInterval
	}

	blocksSinceLastMessage := uint64(0)
	for {
		v, ok := <-sub.Channel()
		if !ok {
			if sub.Err() != nil {
				return rpc.ConvertError(sub.Err(), "stream encountered an error", codes.Internal)
			}
			return nil
		}

		resp, ok := v.(*AccountStatusesResponse)
		if !ok {
			return status.Errorf(codes.Internal, "unexpected response type: %T", v)
		}

		// check if there are any events in the response. if not, do not send a message unless the last
		// response was more than HeartbeatInterval blocks ago
		if len(resp.AccountEvents) == 0 {
			blocksSinceLastMessage++
			if blocksSinceLastMessage < heartbeatInterval {
				continue
			}
			blocksSinceLastMessage = 0
		}

		// BlockExecutionData contains CCF encoded events, and the Access API returns JSON-CDC events.
		// convert event payload formats.
		// This is a temporary solution until the Access API supports specifying the encoding in the request
		results, err := convertAccountsStatusesResults(request.GetEventEncodingVersion(), resp)
		if err != nil {
			return err
		}

		err = stream.Send(&executiondata.SubscribeAccountStatusesResponse{
			BlockId:      convert.IdentifierToMessage(resp.BlockID),
			BlockHeight:  resp.Height,
			Results:      results,
			MessageIndex: resp.MessageIndex,
		})
		if err != nil {
			return rpc.ConvertError(err, "could not send response", codes.Internal)
		}
	}
}
