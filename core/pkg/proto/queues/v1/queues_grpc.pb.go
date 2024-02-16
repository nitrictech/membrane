// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: nitric/proto/queues/v1/queues.proto

package queuespb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// QueuesClient is the client API for Queues service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueuesClient interface {
	// Send messages to a queue
	Enqueue(ctx context.Context, in *QueueEnqueueRequestBatch, opts ...grpc.CallOption) (*QueueEnqueueResponse, error)
	// Receive message(s) from a queue
	Dequeue(ctx context.Context, in *QueueDequeueRequest, opts ...grpc.CallOption) (*QueueDequeueResponse, error)
	// Complete an item previously popped from a queue
	Complete(ctx context.Context, in *QueueCompleteRequest, opts ...grpc.CallOption) (*QueueCompleteResponse, error)
}

type queuesClient struct {
	cc grpc.ClientConnInterface
}

func NewQueuesClient(cc grpc.ClientConnInterface) QueuesClient {
	return &queuesClient{cc}
}

func (c *queuesClient) Enqueue(ctx context.Context, in *QueueEnqueueRequestBatch, opts ...grpc.CallOption) (*QueueEnqueueResponse, error) {
	out := new(QueueEnqueueResponse)
	err := c.cc.Invoke(ctx, "/nitric.proto.queues.v1.Queues/Enqueue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesClient) Dequeue(ctx context.Context, in *QueueDequeueRequest, opts ...grpc.CallOption) (*QueueDequeueResponse, error) {
	out := new(QueueDequeueResponse)
	err := c.cc.Invoke(ctx, "/nitric.proto.queues.v1.Queues/Dequeue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesClient) Complete(ctx context.Context, in *QueueCompleteRequest, opts ...grpc.CallOption) (*QueueCompleteResponse, error) {
	out := new(QueueCompleteResponse)
	err := c.cc.Invoke(ctx, "/nitric.proto.queues.v1.Queues/Complete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueuesServer is the server API for Queues service.
// All implementations should embed UnimplementedQueuesServer
// for forward compatibility
type QueuesServer interface {
	// Send messages to a queue
	Enqueue(context.Context, *QueueEnqueueRequestBatch) (*QueueEnqueueResponse, error)
	// Receive message(s) from a queue
	Dequeue(context.Context, *QueueDequeueRequest) (*QueueDequeueResponse, error)
	// Complete an item previously popped from a queue
	Complete(context.Context, *QueueCompleteRequest) (*QueueCompleteResponse, error)
}

// UnimplementedQueuesServer should be embedded to have forward compatible implementations.
type UnimplementedQueuesServer struct {
}

func (UnimplementedQueuesServer) Enqueue(context.Context, *QueueEnqueueRequestBatch) (*QueueEnqueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Enqueue not implemented")
}
func (UnimplementedQueuesServer) Dequeue(context.Context, *QueueDequeueRequest) (*QueueDequeueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Dequeue not implemented")
}
func (UnimplementedQueuesServer) Complete(context.Context, *QueueCompleteRequest) (*QueueCompleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Complete not implemented")
}

// UnsafeQueuesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueuesServer will
// result in compilation errors.
type UnsafeQueuesServer interface {
	mustEmbedUnimplementedQueuesServer()
}

func RegisterQueuesServer(s grpc.ServiceRegistrar, srv QueuesServer) {
	s.RegisterService(&Queues_ServiceDesc, srv)
}

func _Queues_Enqueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueEnqueueRequestBatch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServer).Enqueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.proto.queues.v1.Queues/Enqueue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServer).Enqueue(ctx, req.(*QueueEnqueueRequestBatch))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queues_Dequeue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueDequeueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServer).Dequeue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.proto.queues.v1.Queues/Dequeue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServer).Dequeue(ctx, req.(*QueueDequeueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Queues_Complete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueCompleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServer).Complete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/nitric.proto.queues.v1.Queues/Complete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServer).Complete(ctx, req.(*QueueCompleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Queues_ServiceDesc is the grpc.ServiceDesc for Queues service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Queues_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nitric.proto.queues.v1.Queues",
	HandlerType: (*QueuesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Enqueue",
			Handler:    _Queues_Enqueue_Handler,
		},
		{
			MethodName: "Dequeue",
			Handler:    _Queues_Dequeue_Handler,
		},
		{
			MethodName: "Complete",
			Handler:    _Queues_Complete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nitric/proto/queues/v1/queues.proto",
}
