// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// MathematicsServiceClient is the client API for MathematicsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MathematicsServiceClient interface {
	// StartLocalEngine starts a Mathematics Engine on the Bhojpur.NET Platform directly.
	// The incoming requests are expected in the following order:
	//   1. metadata
	//   2. all bytes constituting the mathematics/config.yaml
	//   3. all bytes constituting the Engine YAML that will be executed (that the config.yaml points to)
	//   4. all bytes constituting the gzipped Bhojpur.NET Platform application tar stream
	//   5. the Bhojpur.NET Platform application tar stream done marker
	StartLocalEngine(ctx context.Context, opts ...grpc.CallOption) (MathematicsService_StartLocalEngineClient, error)
	// StartFromPreviousEngine starts a new Engine based on a previous one.
	// If the previous Engine does not have the can-replay condition set this call will result in an error.
	StartFromPreviousEngine(ctx context.Context, in *StartFromPreviousEngineRequest, opts ...grpc.CallOption) (*StartEngineResponse, error)
	// StartEngineRequest starts a new Engine based on its specification.
	StartEngine(ctx context.Context, in *StartEngineRequest, opts ...grpc.CallOption) (*StartEngineResponse, error)
	// Searches for Engine(s) known to this Engine
	ListEngines(ctx context.Context, in *ListEnginesRequest, opts ...grpc.CallOption) (*ListEnginesResponse, error)
	// Subscribe listens to new Engine(s) updates
	Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (MathematicsService_SubscribeClient, error)
	// GetEngine retrieves details of a single Engine
	GetEngine(ctx context.Context, in *GetEngineRequest, opts ...grpc.CallOption) (*GetEngineResponse, error)
	// Listen listens to Engine updates and log output of a running Engine
	Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (MathematicsService_ListenClient, error)
	// StopEngine stops a currently running Engine
	StopEngine(ctx context.Context, in *StopEngineRequest, opts ...grpc.CallOption) (*StopEngineResponse, error)
}

type mathematicsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMathematicsServiceClient(cc grpc.ClientConnInterface) MathematicsServiceClient {
	return &mathematicsServiceClient{cc}
}

func (c *mathematicsServiceClient) StartLocalEngine(ctx context.Context, opts ...grpc.CallOption) (MathematicsService_StartLocalEngineClient, error) {
	stream, err := c.cc.NewStream(ctx, &MathematicsService_ServiceDesc.Streams[0], "/v1.MathematicsService/StartLocalEngine", opts...)
	if err != nil {
		return nil, err
	}
	x := &mathematicsServiceStartLocalEngineClient{stream}
	return x, nil
}

type MathematicsService_StartLocalEngineClient interface {
	Send(*StartLocalEngineRequest) error
	CloseAndRecv() (*StartEngineResponse, error)
	grpc.ClientStream
}

type mathematicsServiceStartLocalEngineClient struct {
	grpc.ClientStream
}

func (x *mathematicsServiceStartLocalEngineClient) Send(m *StartLocalEngineRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *mathematicsServiceStartLocalEngineClient) CloseAndRecv() (*StartEngineResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(StartEngineResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mathematicsServiceClient) StartFromPreviousEngine(ctx context.Context, in *StartFromPreviousEngineRequest, opts ...grpc.CallOption) (*StartEngineResponse, error) {
	out := new(StartEngineResponse)
	err := c.cc.Invoke(ctx, "/v1.MathematicsService/StartFromPreviousEngine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mathematicsServiceClient) StartEngine(ctx context.Context, in *StartEngineRequest, opts ...grpc.CallOption) (*StartEngineResponse, error) {
	out := new(StartEngineResponse)
	err := c.cc.Invoke(ctx, "/v1.MathematicsService/StartEngine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mathematicsServiceClient) ListEngines(ctx context.Context, in *ListEnginesRequest, opts ...grpc.CallOption) (*ListEnginesResponse, error) {
	out := new(ListEnginesResponse)
	err := c.cc.Invoke(ctx, "/v1.MathematicsService/ListEngines", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mathematicsServiceClient) Subscribe(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (MathematicsService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &MathematicsService_ServiceDesc.Streams[1], "/v1.MathematicsService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &mathematicsServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MathematicsService_SubscribeClient interface {
	Recv() (*SubscribeResponse, error)
	grpc.ClientStream
}

type mathematicsServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *mathematicsServiceSubscribeClient) Recv() (*SubscribeResponse, error) {
	m := new(SubscribeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mathematicsServiceClient) GetEngine(ctx context.Context, in *GetEngineRequest, opts ...grpc.CallOption) (*GetEngineResponse, error) {
	out := new(GetEngineResponse)
	err := c.cc.Invoke(ctx, "/v1.MathematicsService/GetEngine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mathematicsServiceClient) Listen(ctx context.Context, in *ListenRequest, opts ...grpc.CallOption) (MathematicsService_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &MathematicsService_ServiceDesc.Streams[2], "/v1.MathematicsService/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &mathematicsServiceListenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MathematicsService_ListenClient interface {
	Recv() (*ListenResponse, error)
	grpc.ClientStream
}

type mathematicsServiceListenClient struct {
	grpc.ClientStream
}

func (x *mathematicsServiceListenClient) Recv() (*ListenResponse, error) {
	m := new(ListenResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *mathematicsServiceClient) StopEngine(ctx context.Context, in *StopEngineRequest, opts ...grpc.CallOption) (*StopEngineResponse, error) {
	out := new(StopEngineResponse)
	err := c.cc.Invoke(ctx, "/v1.MathematicsService/StopEngine", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MathematicsServiceServer is the server API for MathematicsService service.
// All implementations must embed UnimplementedMathematicsServiceServer
// for forward compatibility
type MathematicsServiceServer interface {
	// StartLocalEngine starts a Mathematics Engine on the Bhojpur.NET Platform directly.
	// The incoming requests are expected in the following order:
	//   1. metadata
	//   2. all bytes constituting the mathematics/config.yaml
	//   3. all bytes constituting the Engine YAML that will be executed (that the config.yaml points to)
	//   4. all bytes constituting the gzipped Bhojpur.NET Platform application tar stream
	//   5. the Bhojpur.NET Platform application tar stream done marker
	StartLocalEngine(MathematicsService_StartLocalEngineServer) error
	// StartFromPreviousEngine starts a new Engine based on a previous one.
	// If the previous Engine does not have the can-replay condition set this call will result in an error.
	StartFromPreviousEngine(context.Context, *StartFromPreviousEngineRequest) (*StartEngineResponse, error)
	// StartEngineRequest starts a new Engine based on its specification.
	StartEngine(context.Context, *StartEngineRequest) (*StartEngineResponse, error)
	// Searches for Engine(s) known to this Engine
	ListEngines(context.Context, *ListEnginesRequest) (*ListEnginesResponse, error)
	// Subscribe listens to new Engine(s) updates
	Subscribe(*SubscribeRequest, MathematicsService_SubscribeServer) error
	// GetEngine retrieves details of a single Engine
	GetEngine(context.Context, *GetEngineRequest) (*GetEngineResponse, error)
	// Listen listens to Engine updates and log output of a running Engine
	Listen(*ListenRequest, MathematicsService_ListenServer) error
	// StopEngine stops a currently running Engine
	StopEngine(context.Context, *StopEngineRequest) (*StopEngineResponse, error)
	mustEmbedUnimplementedMathematicsServiceServer()
}

// UnimplementedMathematicsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMathematicsServiceServer struct {
}

func (UnimplementedMathematicsServiceServer) StartLocalEngine(MathematicsService_StartLocalEngineServer) error {
	return status.Errorf(codes.Unimplemented, "method StartLocalEngine not implemented")
}
func (UnimplementedMathematicsServiceServer) StartFromPreviousEngine(context.Context, *StartFromPreviousEngineRequest) (*StartEngineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartFromPreviousEngine not implemented")
}
func (UnimplementedMathematicsServiceServer) StartEngine(context.Context, *StartEngineRequest) (*StartEngineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartEngine not implemented")
}
func (UnimplementedMathematicsServiceServer) ListEngines(context.Context, *ListEnginesRequest) (*ListEnginesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEngines not implemented")
}
func (UnimplementedMathematicsServiceServer) Subscribe(*SubscribeRequest, MathematicsService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedMathematicsServiceServer) GetEngine(context.Context, *GetEngineRequest) (*GetEngineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEngine not implemented")
}
func (UnimplementedMathematicsServiceServer) Listen(*ListenRequest, MathematicsService_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (UnimplementedMathematicsServiceServer) StopEngine(context.Context, *StopEngineRequest) (*StopEngineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopEngine not implemented")
}
func (UnimplementedMathematicsServiceServer) mustEmbedUnimplementedMathematicsServiceServer() {}

// UnsafeMathematicsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MathematicsServiceServer will
// result in compilation errors.
type UnsafeMathematicsServiceServer interface {
	mustEmbedUnimplementedMathematicsServiceServer()
}

func RegisterMathematicsServiceServer(s grpc.ServiceRegistrar, srv MathematicsServiceServer) {
	s.RegisterService(&MathematicsService_ServiceDesc, srv)
}

func _MathematicsService_StartLocalEngine_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MathematicsServiceServer).StartLocalEngine(&mathematicsServiceStartLocalEngineServer{stream})
}

type MathematicsService_StartLocalEngineServer interface {
	SendAndClose(*StartEngineResponse) error
	Recv() (*StartLocalEngineRequest, error)
	grpc.ServerStream
}

type mathematicsServiceStartLocalEngineServer struct {
	grpc.ServerStream
}

func (x *mathematicsServiceStartLocalEngineServer) SendAndClose(m *StartEngineResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *mathematicsServiceStartLocalEngineServer) Recv() (*StartLocalEngineRequest, error) {
	m := new(StartLocalEngineRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _MathematicsService_StartFromPreviousEngine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartFromPreviousEngineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MathematicsServiceServer).StartFromPreviousEngine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.MathematicsService/StartFromPreviousEngine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MathematicsServiceServer).StartFromPreviousEngine(ctx, req.(*StartFromPreviousEngineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MathematicsService_StartEngine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartEngineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MathematicsServiceServer).StartEngine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.MathematicsService/StartEngine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MathematicsServiceServer).StartEngine(ctx, req.(*StartEngineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MathematicsService_ListEngines_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListEnginesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MathematicsServiceServer).ListEngines(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.MathematicsService/ListEngines",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MathematicsServiceServer).ListEngines(ctx, req.(*ListEnginesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MathematicsService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MathematicsServiceServer).Subscribe(m, &mathematicsServiceSubscribeServer{stream})
}

type MathematicsService_SubscribeServer interface {
	Send(*SubscribeResponse) error
	grpc.ServerStream
}

type mathematicsServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *mathematicsServiceSubscribeServer) Send(m *SubscribeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _MathematicsService_GetEngine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEngineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MathematicsServiceServer).GetEngine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.MathematicsService/GetEngine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MathematicsServiceServer).GetEngine(ctx, req.(*GetEngineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MathematicsService_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MathematicsServiceServer).Listen(m, &mathematicsServiceListenServer{stream})
}

type MathematicsService_ListenServer interface {
	Send(*ListenResponse) error
	grpc.ServerStream
}

type mathematicsServiceListenServer struct {
	grpc.ServerStream
}

func (x *mathematicsServiceListenServer) Send(m *ListenResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _MathematicsService_StopEngine_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopEngineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MathematicsServiceServer).StopEngine(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.MathematicsService/StopEngine",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MathematicsServiceServer).StopEngine(ctx, req.(*StopEngineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MathematicsService_ServiceDesc is the grpc.ServiceDesc for MathematicsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MathematicsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.MathematicsService",
	HandlerType: (*MathematicsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartFromPreviousEngine",
			Handler:    _MathematicsService_StartFromPreviousEngine_Handler,
		},
		{
			MethodName: "StartEngine",
			Handler:    _MathematicsService_StartEngine_Handler,
		},
		{
			MethodName: "ListEngines",
			Handler:    _MathematicsService_ListEngines_Handler,
		},
		{
			MethodName: "GetEngine",
			Handler:    _MathematicsService_GetEngine_Handler,
		},
		{
			MethodName: "StopEngine",
			Handler:    _MathematicsService_StopEngine_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StartLocalEngine",
			Handler:       _MathematicsService_StartLocalEngine_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _MathematicsService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Listen",
			Handler:       _MathematicsService_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "mathematics.proto",
}