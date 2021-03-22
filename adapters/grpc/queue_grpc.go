package grpc

import (
	"context"

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// GRPC Interface for registered Nitric Storage Plugins
type QueueServer struct {
	pb.UnimplementedQueueServer
	plugin sdk.QueueService
}

func (s *QueueServer) checkPluginRegistered() (bool, error) {
	if s.plugin == nil {
		return false, status.Errorf(codes.Unimplemented, "Queue plugin not registered")
	}

	return true, nil
}

func (s *QueueServer) Send(ctx context.Context, req *pb.QueueSendRequest) (*pb.QueueSendResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		evt := req.GetEvent()

		nitricEvt := sdk.NitricEvent{
			RequestId:   evt.GetRequestId(),
			PayloadType: evt.GetPayloadType(),
			Payload:     evt.GetPayload().AsMap(),
		}

		if err := s.plugin.Send(req.GetQueue(), nitricEvt); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// Success
	return &pb.QueueSendResponse{}, nil
}

func (s *QueueServer) SendBatch(ctx context.Context, req *pb.QueueSendBatchRequest) (*pb.QueueSendBatchResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		// Translate events
		evts := make([]sdk.NitricEvent, len(req.GetEvents()))
		for i, evt := range req.GetEvents() {
			evts[i] = sdk.NitricEvent{
				RequestId:   evt.GetRequestId(),
				PayloadType: evt.GetPayloadType(),
				Payload:     evt.GetPayload().AsMap(),
			}
		}

		if resp, err := s.plugin.SendBatch(req.GetQueue(), evts); err == nil {
			failedEvents := make([]*pb.FailedEvent, len(resp.FailedMessages))
			for i, fmsg := range resp.FailedMessages {
				st, _ := structpb.NewStruct(fmsg.Event.Payload)
				failedEvents[i] = &pb.FailedEvent{
					Message: fmsg.Message,
					Event: &pb.NitricEvent{
						RequestId:   fmsg.Event.RequestId,
						PayloadType: fmsg.Event.PayloadType,
						Payload:     st,
					},
				}
			}
			return &pb.QueueSendBatchResponse{
				FailedEvents: failedEvents,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (s *QueueServer) Receive(ctx context.Context, req *pb.QueueReceiveRequest) (*pb.QueueReceiveResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		// Convert gRPC request to plugin params
		depth := uint32(req.GetDepth())
		popOptions := sdk.ReceiveOptions{
			QueueName: req.GetQueue(),
			Depth:     &depth,
		}

		// Perform the Queue Pop operation
		queueItems, err := s.plugin.Receive(popOptions)
		if err != nil {
			return nil, err
		}

		// Convert the NitricEvents to the gRPC type
		grpcQueueItems := []*pb.NitricQueueItem{}
		for _, queueItem := range queueItems {
			st, _ := structpb.NewStruct(queueItem.Event.Payload)
			grpcQueueItems = append(grpcQueueItems, &pb.NitricQueueItem{
				Event: &pb.NitricEvent{
					RequestId:   queueItem.Event.RequestId,
					PayloadType: queueItem.Event.PayloadType,
					Payload:     st,
				},
				LeaseId: queueItem.LeaseId,
			})
		}

		// Return the queue items
		res := pb.QueueReceiveResponse{
			Items: grpcQueueItems,
		}
		return &res, nil
	} else {
		return nil, err
	}
}

func (s *QueueServer) Complete(ctx context.Context, req *pb.QueueCompleteRequest) (*pb.QueueCompleteResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		// Convert gRPC request to plugin params
		queueName := req.GetQueue()
		leaseId := req.GetLeaseId()

		// Perform the Queue Complete operation
		err := s.plugin.Complete(queueName, leaseId)
		if err != nil {
			return nil, err
		}

		// Return a successful response
		return &pb.QueueCompleteResponse{}, nil
	} else {
		return nil, err
	}
}

func NewQueueServer(plugin sdk.QueueService) pb.QueueServer {
	return &QueueServer{
		plugin: plugin,
	}
}
