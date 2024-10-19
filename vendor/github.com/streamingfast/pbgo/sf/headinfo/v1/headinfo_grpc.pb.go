// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: sf/headinfo/v1/headinfo.proto

package pbheadinfo

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	HeadInfo_GetHeadInfo_FullMethodName = "/sf.headinfo.v1.HeadInfo/GetHeadInfo"
)

// HeadInfoClient is the client API for HeadInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HeadInfoClient interface {
	GetHeadInfo(ctx context.Context, in *HeadInfoRequest, opts ...grpc.CallOption) (*HeadInfoResponse, error)
}

type headInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewHeadInfoClient(cc grpc.ClientConnInterface) HeadInfoClient {
	return &headInfoClient{cc}
}

func (c *headInfoClient) GetHeadInfo(ctx context.Context, in *HeadInfoRequest, opts ...grpc.CallOption) (*HeadInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HeadInfoResponse)
	err := c.cc.Invoke(ctx, HeadInfo_GetHeadInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HeadInfoServer is the server API for HeadInfo service.
// All implementations should embed UnimplementedHeadInfoServer
// for forward compatibility.
type HeadInfoServer interface {
	GetHeadInfo(context.Context, *HeadInfoRequest) (*HeadInfoResponse, error)
}

// UnimplementedHeadInfoServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedHeadInfoServer struct{}

func (UnimplementedHeadInfoServer) GetHeadInfo(context.Context, *HeadInfoRequest) (*HeadInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHeadInfo not implemented")
}
func (UnimplementedHeadInfoServer) testEmbeddedByValue() {}

// UnsafeHeadInfoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HeadInfoServer will
// result in compilation errors.
type UnsafeHeadInfoServer interface {
	mustEmbedUnimplementedHeadInfoServer()
}

func RegisterHeadInfoServer(s grpc.ServiceRegistrar, srv HeadInfoServer) {
	// If the following call pancis, it indicates UnimplementedHeadInfoServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&HeadInfo_ServiceDesc, srv)
}

func _HeadInfo_GetHeadInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeadInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HeadInfoServer).GetHeadInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HeadInfo_GetHeadInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HeadInfoServer).GetHeadInfo(ctx, req.(*HeadInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// HeadInfo_ServiceDesc is the grpc.ServiceDesc for HeadInfo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HeadInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sf.headinfo.v1.HeadInfo",
	HandlerType: (*HeadInfoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetHeadInfo",
			Handler:    _HeadInfo_GetHeadInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sf/headinfo/v1/headinfo.proto",
}

const (
	StreamingHeadInfo_StreamHeadInfo_FullMethodName = "/sf.headinfo.v1.StreamingHeadInfo/StreamHeadInfo"
)

// StreamingHeadInfoClient is the client API for StreamingHeadInfo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamingHeadInfoClient interface {
	StreamHeadInfo(ctx context.Context, in *HeadInfoRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[HeadInfoResponse], error)
}

type streamingHeadInfoClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamingHeadInfoClient(cc grpc.ClientConnInterface) StreamingHeadInfoClient {
	return &streamingHeadInfoClient{cc}
}

func (c *streamingHeadInfoClient) StreamHeadInfo(ctx context.Context, in *HeadInfoRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[HeadInfoResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StreamingHeadInfo_ServiceDesc.Streams[0], StreamingHeadInfo_StreamHeadInfo_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[HeadInfoRequest, HeadInfoResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamingHeadInfo_StreamHeadInfoClient = grpc.ServerStreamingClient[HeadInfoResponse]

// StreamingHeadInfoServer is the server API for StreamingHeadInfo service.
// All implementations should embed UnimplementedStreamingHeadInfoServer
// for forward compatibility.
type StreamingHeadInfoServer interface {
	StreamHeadInfo(*HeadInfoRequest, grpc.ServerStreamingServer[HeadInfoResponse]) error
}

// UnimplementedStreamingHeadInfoServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStreamingHeadInfoServer struct{}

func (UnimplementedStreamingHeadInfoServer) StreamHeadInfo(*HeadInfoRequest, grpc.ServerStreamingServer[HeadInfoResponse]) error {
	return status.Errorf(codes.Unimplemented, "method StreamHeadInfo not implemented")
}
func (UnimplementedStreamingHeadInfoServer) testEmbeddedByValue() {}

// UnsafeStreamingHeadInfoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamingHeadInfoServer will
// result in compilation errors.
type UnsafeStreamingHeadInfoServer interface {
	mustEmbedUnimplementedStreamingHeadInfoServer()
}

func RegisterStreamingHeadInfoServer(s grpc.ServiceRegistrar, srv StreamingHeadInfoServer) {
	// If the following call pancis, it indicates UnimplementedStreamingHeadInfoServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StreamingHeadInfo_ServiceDesc, srv)
}

func _StreamingHeadInfo_StreamHeadInfo_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(HeadInfoRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamingHeadInfoServer).StreamHeadInfo(m, &grpc.GenericServerStream[HeadInfoRequest, HeadInfoResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StreamingHeadInfo_StreamHeadInfoServer = grpc.ServerStreamingServer[HeadInfoResponse]

// StreamingHeadInfo_ServiceDesc is the grpc.ServiceDesc for StreamingHeadInfo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StreamingHeadInfo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sf.headinfo.v1.StreamingHeadInfo",
	HandlerType: (*StreamingHeadInfoServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamHeadInfo",
			Handler:       _StreamingHeadInfo_StreamHeadInfo_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "sf/headinfo/v1/headinfo.proto",
}
