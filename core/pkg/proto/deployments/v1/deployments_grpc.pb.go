// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: nitric/proto/deployments/v1/deployments.proto

package deploymentspb

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

// DeploymentClient is the client API for Deployment service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeploymentClient interface {
	// Begins a new deployment
	// Server will stream updates back to the connected client
	// on the status of the deployment
	Up(ctx context.Context, in *DeploymentUpRequest, opts ...grpc.CallOption) (Deployment_UpClient, error)
	// Begins a new deployment preview
	// Server will stream updates back to the connected client
	// on the status of the preview request
	Preview(ctx context.Context, in *DeploymentPreviewRequest, opts ...grpc.CallOption) (Deployment_PreviewClient, error)
	// Tears down an existing deployment
	// Server will stream updates back to the connected client
	// on the status of the teardown
	Down(ctx context.Context, in *DeploymentDownRequest, opts ...grpc.CallOption) (Deployment_DownClient, error)
}

type deploymentClient struct {
	cc grpc.ClientConnInterface
}

func NewDeploymentClient(cc grpc.ClientConnInterface) DeploymentClient {
	return &deploymentClient{cc}
}

func (c *deploymentClient) Up(ctx context.Context, in *DeploymentUpRequest, opts ...grpc.CallOption) (Deployment_UpClient, error) {
	stream, err := c.cc.NewStream(ctx, &Deployment_ServiceDesc.Streams[0], "/nitric.proto.deployments.v1.Deployment/Up", opts...)
	if err != nil {
		return nil, err
	}
	x := &deploymentUpClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Deployment_UpClient interface {
	Recv() (*DeploymentUpEvent, error)
	grpc.ClientStream
}

type deploymentUpClient struct {
	grpc.ClientStream
}

func (x *deploymentUpClient) Recv() (*DeploymentUpEvent, error) {
	m := new(DeploymentUpEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *deploymentClient) Preview(ctx context.Context, in *DeploymentPreviewRequest, opts ...grpc.CallOption) (Deployment_PreviewClient, error) {
	stream, err := c.cc.NewStream(ctx, &Deployment_ServiceDesc.Streams[1], "/nitric.proto.deployments.v1.Deployment/Preview", opts...)
	if err != nil {
		return nil, err
	}
	x := &deploymentPreviewClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Deployment_PreviewClient interface {
	Recv() (*DeploymentPreviewEvent, error)
	grpc.ClientStream
}

type deploymentPreviewClient struct {
	grpc.ClientStream
}

func (x *deploymentPreviewClient) Recv() (*DeploymentPreviewEvent, error) {
	m := new(DeploymentPreviewEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *deploymentClient) Down(ctx context.Context, in *DeploymentDownRequest, opts ...grpc.CallOption) (Deployment_DownClient, error) {
	stream, err := c.cc.NewStream(ctx, &Deployment_ServiceDesc.Streams[2], "/nitric.proto.deployments.v1.Deployment/Down", opts...)
	if err != nil {
		return nil, err
	}
	x := &deploymentDownClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Deployment_DownClient interface {
	Recv() (*DeploymentDownEvent, error)
	grpc.ClientStream
}

type deploymentDownClient struct {
	grpc.ClientStream
}

func (x *deploymentDownClient) Recv() (*DeploymentDownEvent, error) {
	m := new(DeploymentDownEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DeploymentServer is the server API for Deployment service.
// All implementations should embed UnimplementedDeploymentServer
// for forward compatibility
type DeploymentServer interface {
	// Begins a new deployment
	// Server will stream updates back to the connected client
	// on the status of the deployment
	Up(*DeploymentUpRequest, Deployment_UpServer) error
	// Begins a new deployment preview
	// Server will stream updates back to the connected client
	// on the status of the preview request
	Preview(*DeploymentPreviewRequest, Deployment_PreviewServer) error
	// Tears down an existing deployment
	// Server will stream updates back to the connected client
	// on the status of the teardown
	Down(*DeploymentDownRequest, Deployment_DownServer) error
}

// UnimplementedDeploymentServer should be embedded to have forward compatible implementations.
type UnimplementedDeploymentServer struct {
}

func (UnimplementedDeploymentServer) Up(*DeploymentUpRequest, Deployment_UpServer) error {
	return status.Errorf(codes.Unimplemented, "method Up not implemented")
}
func (UnimplementedDeploymentServer) Preview(*DeploymentPreviewRequest, Deployment_PreviewServer) error {
	return status.Errorf(codes.Unimplemented, "method Preview not implemented")
}
func (UnimplementedDeploymentServer) Down(*DeploymentDownRequest, Deployment_DownServer) error {
	return status.Errorf(codes.Unimplemented, "method Down not implemented")
}

// UnsafeDeploymentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeploymentServer will
// result in compilation errors.
type UnsafeDeploymentServer interface {
	mustEmbedUnimplementedDeploymentServer()
}

func RegisterDeploymentServer(s grpc.ServiceRegistrar, srv DeploymentServer) {
	s.RegisterService(&Deployment_ServiceDesc, srv)
}

func _Deployment_Up_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DeploymentUpRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DeploymentServer).Up(m, &deploymentUpServer{stream})
}

type Deployment_UpServer interface {
	Send(*DeploymentUpEvent) error
	grpc.ServerStream
}

type deploymentUpServer struct {
	grpc.ServerStream
}

func (x *deploymentUpServer) Send(m *DeploymentUpEvent) error {
	return x.ServerStream.SendMsg(m)
}

func _Deployment_Preview_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DeploymentPreviewRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DeploymentServer).Preview(m, &deploymentPreviewServer{stream})
}

type Deployment_PreviewServer interface {
	Send(*DeploymentPreviewEvent) error
	grpc.ServerStream
}

type deploymentPreviewServer struct {
	grpc.ServerStream
}

func (x *deploymentPreviewServer) Send(m *DeploymentPreviewEvent) error {
	return x.ServerStream.SendMsg(m)
}

func _Deployment_Down_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DeploymentDownRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DeploymentServer).Down(m, &deploymentDownServer{stream})
}

type Deployment_DownServer interface {
	Send(*DeploymentDownEvent) error
	grpc.ServerStream
}

type deploymentDownServer struct {
	grpc.ServerStream
}

func (x *deploymentDownServer) Send(m *DeploymentDownEvent) error {
	return x.ServerStream.SendMsg(m)
}

// Deployment_ServiceDesc is the grpc.ServiceDesc for Deployment service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Deployment_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "nitric.proto.deployments.v1.Deployment",
	HandlerType: (*DeploymentServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Up",
			Handler:       _Deployment_Up_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Preview",
			Handler:       _Deployment_Preview_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Down",
			Handler:       _Deployment_Down_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "nitric/proto/deployments/v1/deployments.proto",
}
