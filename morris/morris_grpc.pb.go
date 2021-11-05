// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package morris

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

// MorrisClient is the client API for Morris service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MorrisClient interface {
	GetBoardStream(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Morris_GetBoardStreamClient, error)
	MakeMove(ctx context.Context, in *Move, opts ...grpc.CallOption) (*BoardState, error)
}

type morrisClient struct {
	cc grpc.ClientConnInterface
}

func NewMorrisClient(cc grpc.ClientConnInterface) MorrisClient {
	return &morrisClient{cc}
}

func (c *morrisClient) GetBoardStream(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Morris_GetBoardStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Morris_ServiceDesc.Streams[0], "/Morris/GetBoardStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &morrisGetBoardStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Morris_GetBoardStreamClient interface {
	Recv() (*BoardState, error)
	grpc.ClientStream
}

type morrisGetBoardStreamClient struct {
	grpc.ClientStream
}

func (x *morrisGetBoardStreamClient) Recv() (*BoardState, error) {
	m := new(BoardState)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *morrisClient) MakeMove(ctx context.Context, in *Move, opts ...grpc.CallOption) (*BoardState, error) {
	out := new(BoardState)
	err := c.cc.Invoke(ctx, "/Morris/MakeMove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MorrisServer is the server API for Morris service.
// All implementations must embed UnimplementedMorrisServer
// for forward compatibility
type MorrisServer interface {
	GetBoardStream(*Empty, Morris_GetBoardStreamServer) error
	MakeMove(context.Context, *Move) (*BoardState, error)
	mustEmbedUnimplementedMorrisServer()
}

// UnimplementedMorrisServer must be embedded to have forward compatible implementations.
type UnimplementedMorrisServer struct {
}

func (UnimplementedMorrisServer) GetBoardStream(*Empty, Morris_GetBoardStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GetBoardStream not implemented")
}
func (UnimplementedMorrisServer) MakeMove(context.Context, *Move) (*BoardState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakeMove not implemented")
}
func (UnimplementedMorrisServer) mustEmbedUnimplementedMorrisServer() {}

// UnsafeMorrisServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MorrisServer will
// result in compilation errors.
type UnsafeMorrisServer interface {
	mustEmbedUnimplementedMorrisServer()
}

func RegisterMorrisServer(s grpc.ServiceRegistrar, srv MorrisServer) {
	s.RegisterService(&Morris_ServiceDesc, srv)
}

func _Morris_GetBoardStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MorrisServer).GetBoardStream(m, &morrisGetBoardStreamServer{stream})
}

type Morris_GetBoardStreamServer interface {
	Send(*BoardState) error
	grpc.ServerStream
}

type morrisGetBoardStreamServer struct {
	grpc.ServerStream
}

func (x *morrisGetBoardStreamServer) Send(m *BoardState) error {
	return x.ServerStream.SendMsg(m)
}

func _Morris_MakeMove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Move)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MorrisServer).MakeMove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Morris/MakeMove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MorrisServer).MakeMove(ctx, req.(*Move))
	}
	return interceptor(ctx, in, info, handler)
}

// Morris_ServiceDesc is the grpc.ServiceDesc for Morris service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Morris_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Morris",
	HandlerType: (*MorrisServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MakeMove",
			Handler:    _Morris_MakeMove_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetBoardStream",
			Handler:       _Morris_GetBoardStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "morris.proto",
}
