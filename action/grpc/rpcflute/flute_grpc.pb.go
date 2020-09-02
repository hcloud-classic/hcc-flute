// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package rpcflute

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	rpcmsgType "hcc/flute/action/grpc/rpcmsgType"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// FluteClient is the client API for Flute service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FluteClient interface {
	// Node
	CreateNode(ctx context.Context, in *ReqCreateNode, opts ...grpc.CallOption) (*ResCreateNode, error)
	GetNode(ctx context.Context, in *ReqGetNode, opts ...grpc.CallOption) (*ResGetNode, error)
	GetNodeList(ctx context.Context, in *ReqGetNodeList, opts ...grpc.CallOption) (*ResGetNodeList, error)
	GetNodeNum(ctx context.Context, in *rpcmsgType.Empty, opts ...grpc.CallOption) (*ResGetNodeNum, error)
	UpdateNode(ctx context.Context, in *ReqUpdateNode, opts ...grpc.CallOption) (*ResUpdateNode, error)
	DeleteNode(ctx context.Context, in *ReqDeleteNode, opts ...grpc.CallOption) (*ResDeleteNode, error)
	// NodeDetail
	CreateNodeDetail(ctx context.Context, in *ReqCreateNodeDetail, opts ...grpc.CallOption) (*ResCreateNodeDetail, error)
	GetNodeDetail(ctx context.Context, in *ReqGetNodeDetail, opts ...grpc.CallOption) (*ResGetNodeDetail, error)
	DeleteNodeDetail(ctx context.Context, in *ReqDeleteNodeDetail, opts ...grpc.CallOption) (*ResDeleteNodeDetail, error)
	// IPMI
	NodePowerControl(ctx context.Context, in *ReqNodePowerControl, opts ...grpc.CallOption) (*ResNodePowerControl, error)
	GetNodePowerState(ctx context.Context, in *ReqNodePowerState, opts ...grpc.CallOption) (*ResNodePowerState, error)
}

type fluteClient struct {
	cc grpc.ClientConnInterface
}

func NewFluteClient(cc grpc.ClientConnInterface) FluteClient {
	return &fluteClient{cc}
}

func (c *fluteClient) CreateNode(ctx context.Context, in *ReqCreateNode, opts ...grpc.CallOption) (*ResCreateNode, error) {
	out := new(ResCreateNode)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/CreateNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fluteClient) GetNode(ctx context.Context, in *ReqGetNode, opts ...grpc.CallOption) (*ResGetNode, error) {
	out := new(ResGetNode)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/GetNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fluteClient) GetNodeList(ctx context.Context, in *ReqGetNodeList, opts ...grpc.CallOption) (*ResGetNodeList, error) {
	out := new(ResGetNodeList)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/GetNodeList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fluteClient) GetNodeNum(ctx context.Context, in *rpcmsgType.Empty, opts ...grpc.CallOption) (*ResGetNodeNum, error) {
	out := new(ResGetNodeNum)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/GetNodeNum", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fluteClient) UpdateNode(ctx context.Context, in *ReqUpdateNode, opts ...grpc.CallOption) (*ResUpdateNode, error) {
	out := new(ResUpdateNode)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/UpdateNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fluteClient) DeleteNode(ctx context.Context, in *ReqDeleteNode, opts ...grpc.CallOption) (*ResDeleteNode, error) {
	out := new(ResDeleteNode)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/DeleteNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fluteClient) CreateNodeDetail(ctx context.Context, in *ReqCreateNodeDetail, opts ...grpc.CallOption) (*ResCreateNodeDetail, error) {
	out := new(ResCreateNodeDetail)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/CreateNodeDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fluteClient) GetNodeDetail(ctx context.Context, in *ReqGetNodeDetail, opts ...grpc.CallOption) (*ResGetNodeDetail, error) {
	out := new(ResGetNodeDetail)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/GetNodeDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fluteClient) DeleteNodeDetail(ctx context.Context, in *ReqDeleteNodeDetail, opts ...grpc.CallOption) (*ResDeleteNodeDetail, error) {
	out := new(ResDeleteNodeDetail)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/DeleteNodeDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fluteClient) NodePowerControl(ctx context.Context, in *ReqNodePowerControl, opts ...grpc.CallOption) (*ResNodePowerControl, error) {
	out := new(ResNodePowerControl)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/NodePowerControl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fluteClient) GetNodePowerState(ctx context.Context, in *ReqNodePowerState, opts ...grpc.CallOption) (*ResNodePowerState, error) {
	out := new(ResNodePowerState)
	err := c.cc.Invoke(ctx, "/RpcFlute.Flute/GetNodePowerState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FluteServer is the server API for Flute service.
// All implementations must embed UnimplementedFluteServer
// for forward compatibility
type FluteServer interface {
	// Node
	CreateNode(context.Context, *ReqCreateNode) (*ResCreateNode, error)
	GetNode(context.Context, *ReqGetNode) (*ResGetNode, error)
	GetNodeList(context.Context, *ReqGetNodeList) (*ResGetNodeList, error)
	GetNodeNum(context.Context, *rpcmsgType.Empty) (*ResGetNodeNum, error)
	UpdateNode(context.Context, *ReqUpdateNode) (*ResUpdateNode, error)
	DeleteNode(context.Context, *ReqDeleteNode) (*ResDeleteNode, error)
	// NodeDetail
	CreateNodeDetail(context.Context, *ReqCreateNodeDetail) (*ResCreateNodeDetail, error)
	GetNodeDetail(context.Context, *ReqGetNodeDetail) (*ResGetNodeDetail, error)
	DeleteNodeDetail(context.Context, *ReqDeleteNodeDetail) (*ResDeleteNodeDetail, error)
	// IPMI
	NodePowerControl(context.Context, *ReqNodePowerControl) (*ResNodePowerControl, error)
	GetNodePowerState(context.Context, *ReqNodePowerState) (*ResNodePowerState, error)
	mustEmbedUnimplementedFluteServer()
}

// UnimplementedFluteServer must be embedded to have forward compatible implementations.
type UnimplementedFluteServer struct {
}

func (*UnimplementedFluteServer) CreateNode(context.Context, *ReqCreateNode) (*ResCreateNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNode not implemented")
}
func (*UnimplementedFluteServer) GetNode(context.Context, *ReqGetNode) (*ResGetNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNode not implemented")
}
func (*UnimplementedFluteServer) GetNodeList(context.Context, *ReqGetNodeList) (*ResGetNodeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodeList not implemented")
}
func (*UnimplementedFluteServer) GetNodeNum(context.Context, *rpcmsgType.Empty) (*ResGetNodeNum, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodeNum not implemented")
}
func (*UnimplementedFluteServer) UpdateNode(context.Context, *ReqUpdateNode) (*ResUpdateNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNode not implemented")
}
func (*UnimplementedFluteServer) DeleteNode(context.Context, *ReqDeleteNode) (*ResDeleteNode, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNode not implemented")
}
func (*UnimplementedFluteServer) CreateNodeDetail(context.Context, *ReqCreateNodeDetail) (*ResCreateNodeDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNodeDetail not implemented")
}
func (*UnimplementedFluteServer) GetNodeDetail(context.Context, *ReqGetNodeDetail) (*ResGetNodeDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodeDetail not implemented")
}
func (*UnimplementedFluteServer) DeleteNodeDetail(context.Context, *ReqDeleteNodeDetail) (*ResDeleteNodeDetail, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNodeDetail not implemented")
}
func (*UnimplementedFluteServer) NodePowerControl(context.Context, *ReqNodePowerControl) (*ResNodePowerControl, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NodePowerControl not implemented")
}
func (*UnimplementedFluteServer) GetNodePowerState(context.Context, *ReqNodePowerState) (*ResNodePowerState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodePowerState not implemented")
}
func (*UnimplementedFluteServer) mustEmbedUnimplementedFluteServer() {}

func RegisterFluteServer(s *grpc.Server, srv FluteServer) {
	s.RegisterService(&_Flute_serviceDesc, srv)
}

func _Flute_CreateNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqCreateNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).CreateNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/CreateNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).CreateNode(ctx, req.(*ReqCreateNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flute_GetNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqGetNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).GetNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/GetNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).GetNode(ctx, req.(*ReqGetNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flute_GetNodeList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqGetNodeList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).GetNodeList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/GetNodeList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).GetNodeList(ctx, req.(*ReqGetNodeList))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flute_GetNodeNum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(rpcmsgType.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).GetNodeNum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/GetNodeNum",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).GetNodeNum(ctx, req.(*rpcmsgType.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flute_UpdateNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqUpdateNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).UpdateNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/UpdateNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).UpdateNode(ctx, req.(*ReqUpdateNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flute_DeleteNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqDeleteNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).DeleteNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/DeleteNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).DeleteNode(ctx, req.(*ReqDeleteNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flute_CreateNodeDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqCreateNodeDetail)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).CreateNodeDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/CreateNodeDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).CreateNodeDetail(ctx, req.(*ReqCreateNodeDetail))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flute_GetNodeDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqGetNodeDetail)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).GetNodeDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/GetNodeDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).GetNodeDetail(ctx, req.(*ReqGetNodeDetail))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flute_DeleteNodeDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqDeleteNodeDetail)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).DeleteNodeDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/DeleteNodeDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).DeleteNodeDetail(ctx, req.(*ReqDeleteNodeDetail))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flute_NodePowerControl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqNodePowerControl)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).NodePowerControl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/NodePowerControl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).NodePowerControl(ctx, req.(*ReqNodePowerControl))
	}
	return interceptor(ctx, in, info, handler)
}

func _Flute_GetNodePowerState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqNodePowerState)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FluteServer).GetNodePowerState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RpcFlute.Flute/GetNodePowerState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FluteServer).GetNodePowerState(ctx, req.(*ReqNodePowerState))
	}
	return interceptor(ctx, in, info, handler)
}

var _Flute_serviceDesc = grpc.ServiceDesc{
	ServiceName: "RpcFlute.Flute",
	HandlerType: (*FluteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateNode",
			Handler:    _Flute_CreateNode_Handler,
		},
		{
			MethodName: "GetNode",
			Handler:    _Flute_GetNode_Handler,
		},
		{
			MethodName: "GetNodeList",
			Handler:    _Flute_GetNodeList_Handler,
		},
		{
			MethodName: "GetNodeNum",
			Handler:    _Flute_GetNodeNum_Handler,
		},
		{
			MethodName: "UpdateNode",
			Handler:    _Flute_UpdateNode_Handler,
		},
		{
			MethodName: "DeleteNode",
			Handler:    _Flute_DeleteNode_Handler,
		},
		{
			MethodName: "CreateNodeDetail",
			Handler:    _Flute_CreateNodeDetail_Handler,
		},
		{
			MethodName: "GetNodeDetail",
			Handler:    _Flute_GetNodeDetail_Handler,
		},
		{
			MethodName: "DeleteNodeDetail",
			Handler:    _Flute_DeleteNodeDetail_Handler,
		},
		{
			MethodName: "NodePowerControl",
			Handler:    _Flute_NodePowerControl_Handler,
		},
		{
			MethodName: "GetNodePowerState",
			Handler:    _Flute_GetNodePowerState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "flute.proto",
}
