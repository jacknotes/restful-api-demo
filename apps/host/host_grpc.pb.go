// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: apps/host/pb/host.proto

package host

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

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	// 录入主机信息
	CreateHost(ctx context.Context, in *Host, opts ...grpc.CallOption) (*Host, error)
	// 查询主机列表信息
	QueryHost(ctx context.Context, in *QueryHostRequest, opts ...grpc.CallOption) (*Set, error)
	// 查询主机详情
	DescribeHost(ctx context.Context, in *DescribeHostRequest, opts ...grpc.CallOption) (*Host, error)
	// 修改主机信息
	UpdateHost(ctx context.Context, in *UpdateHostRequest, opts ...grpc.CallOption) (*Host, error)
	// 删除主机， 为了兼容GRPC和 delete event需要返回Host
	DeleteHost(ctx context.Context, in *DeleteHostRequest, opts ...grpc.CallOption) (*Host, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) CreateHost(ctx context.Context, in *Host, opts ...grpc.CallOption) (*Host, error) {
	out := new(Host)
	err := c.cc.Invoke(ctx, "/demo.Service/CreateHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) QueryHost(ctx context.Context, in *QueryHostRequest, opts ...grpc.CallOption) (*Set, error) {
	out := new(Set)
	err := c.cc.Invoke(ctx, "/demo.Service/QueryHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DescribeHost(ctx context.Context, in *DescribeHostRequest, opts ...grpc.CallOption) (*Host, error) {
	out := new(Host)
	err := c.cc.Invoke(ctx, "/demo.Service/DescribeHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) UpdateHost(ctx context.Context, in *UpdateHostRequest, opts ...grpc.CallOption) (*Host, error) {
	out := new(Host)
	err := c.cc.Invoke(ctx, "/demo.Service/UpdateHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) DeleteHost(ctx context.Context, in *DeleteHostRequest, opts ...grpc.CallOption) (*Host, error) {
	out := new(Host)
	err := c.cc.Invoke(ctx, "/demo.Service/DeleteHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	// 录入主机信息
	CreateHost(context.Context, *Host) (*Host, error)
	// 查询主机列表信息
	QueryHost(context.Context, *QueryHostRequest) (*Set, error)
	// 查询主机详情
	DescribeHost(context.Context, *DescribeHostRequest) (*Host, error)
	// 修改主机信息
	UpdateHost(context.Context, *UpdateHostRequest) (*Host, error)
	// 删除主机， 为了兼容GRPC和 delete event需要返回Host
	DeleteHost(context.Context, *DeleteHostRequest) (*Host, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) CreateHost(context.Context, *Host) (*Host, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateHost not implemented")
}
func (UnimplementedServiceServer) QueryHost(context.Context, *QueryHostRequest) (*Set, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryHost not implemented")
}
func (UnimplementedServiceServer) DescribeHost(context.Context, *DescribeHostRequest) (*Host, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeHost not implemented")
}
func (UnimplementedServiceServer) UpdateHost(context.Context, *UpdateHostRequest) (*Host, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHost not implemented")
}
func (UnimplementedServiceServer) DeleteHost(context.Context, *DeleteHostRequest) (*Host, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteHost not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&Service_ServiceDesc, srv)
}

func _Service_CreateHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Host)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).CreateHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.Service/CreateHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).CreateHost(ctx, req.(*Host))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_QueryHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).QueryHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.Service/QueryHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).QueryHost(ctx, req.(*QueryHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DescribeHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DescribeHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.Service/DescribeHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DescribeHost(ctx, req.(*DescribeHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_UpdateHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).UpdateHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.Service/UpdateHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).UpdateHost(ctx, req.(*UpdateHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_DeleteHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteHostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).DeleteHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/demo.Service/DeleteHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).DeleteHost(ctx, req.(*DeleteHostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Service_ServiceDesc is the grpc.ServiceDesc for Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "demo.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateHost",
			Handler:    _Service_CreateHost_Handler,
		},
		{
			MethodName: "QueryHost",
			Handler:    _Service_QueryHost_Handler,
		},
		{
			MethodName: "DescribeHost",
			Handler:    _Service_DescribeHost_Handler,
		},
		{
			MethodName: "UpdateHost",
			Handler:    _Service_UpdateHost_Handler,
		},
		{
			MethodName: "DeleteHost",
			Handler:    _Service_DeleteHost_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "apps/host/pb/host.proto",
}
