// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package DisysyMandatory2_git

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

// NodeClient is the client API for Node service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeClient interface {
	Request(ctx context.Context, in *RequestMessage, opts ...grpc.CallOption) (*RequestReply, error)
	Accept(ctx context.Context, in *AcceptMessage, opts ...grpc.CallOption) (*AcceptReply, error)
	LeaveCritical(ctx context.Context, in *LeaveMessage, opts ...grpc.CallOption) (*LeaveReply, error)
	EnterCritical(ctx context.Context, in *EnterMessage, opts ...grpc.CallOption) (*EmptyReply, error)
}

type nodeClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeClient(cc grpc.ClientConnInterface) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) Request(ctx context.Context, in *RequestMessage, opts ...grpc.CallOption) (*RequestReply, error) {
	out := new(RequestReply)
	err := c.cc.Invoke(ctx, "/p2p.Node/Request", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Accept(ctx context.Context, in *AcceptMessage, opts ...grpc.CallOption) (*AcceptReply, error) {
	out := new(AcceptReply)
	err := c.cc.Invoke(ctx, "/p2p.Node/Accept", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) LeaveCritical(ctx context.Context, in *LeaveMessage, opts ...grpc.CallOption) (*LeaveReply, error) {
	out := new(LeaveReply)
	err := c.cc.Invoke(ctx, "/p2p.Node/LeaveCritical", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) EnterCritical(ctx context.Context, in *EnterMessage, opts ...grpc.CallOption) (*EmptyReply, error) {
	out := new(EmptyReply)
	err := c.cc.Invoke(ctx, "/p2p.Node/EnterCritical", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServer is the server API for Node service.
// All implementations must embed UnimplementedNodeServer
// for forward compatibility
type NodeServer interface {
	Request(context.Context, *RequestMessage) (*RequestReply, error)
	Accept(context.Context, *AcceptMessage) (*AcceptReply, error)
	LeaveCritical(context.Context, *LeaveMessage) (*LeaveReply, error)
	EnterCritical(context.Context, *EnterMessage) (*EmptyReply, error)
	mustEmbedUnimplementedNodeServer()
}

// UnimplementedNodeServer must be embedded to have forward compatible implementations.
type UnimplementedNodeServer struct {
}

func (UnimplementedNodeServer) Request(context.Context, *RequestMessage) (*RequestReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Request not implemented")
}
func (UnimplementedNodeServer) Accept(context.Context, *AcceptMessage) (*AcceptReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Accept not implemented")
}
func (UnimplementedNodeServer) LeaveCritical(context.Context, *LeaveMessage) (*LeaveReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveCritical not implemented")
}
func (UnimplementedNodeServer) EnterCritical(context.Context, *EnterMessage) (*EmptyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnterCritical not implemented")
}
func (UnimplementedNodeServer) mustEmbedUnimplementedNodeServer() {}

// UnsafeNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServer will
// result in compilation errors.
type UnsafeNodeServer interface {
	mustEmbedUnimplementedNodeServer()
}

func RegisterNodeServer(s grpc.ServiceRegistrar, srv NodeServer) {
	s.RegisterService(&Node_ServiceDesc, srv)
}

func _Node_Request_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Request(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2p.Node/Request",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Request(ctx, req.(*RequestMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Accept_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Accept(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2p.Node/Accept",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Accept(ctx, req.(*AcceptMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_LeaveCritical_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).LeaveCritical(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2p.Node/LeaveCritical",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).LeaveCritical(ctx, req.(*LeaveMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_EnterCritical_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EnterMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).EnterCritical(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2p.Node/EnterCritical",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).EnterCritical(ctx, req.(*EnterMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Node_ServiceDesc is the grpc.ServiceDesc for Node service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Node_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "p2p.Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Request",
			Handler:    _Node_Request_Handler,
		},
		{
			MethodName: "Accept",
			Handler:    _Node_Accept_Handler,
		},
		{
			MethodName: "LeaveCritical",
			Handler:    _Node_LeaveCritical_Handler,
		},
		{
			MethodName: "EnterCritical",
			Handler:    _Node_EnterCritical_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "p2p/p2p.proto",
}
