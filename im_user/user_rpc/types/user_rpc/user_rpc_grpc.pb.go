// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v3.19.4
// source: user_rpc.proto

package user_rpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Users_UserCreate_FullMethodName     = "/user_rpc.Users/UserCreate"
	Users_UserInfo_FullMethodName       = "/user_rpc.Users/UserInfo"
	Users_UserBaseInfo_FullMethodName   = "/user_rpc.Users/UserBaseInfo"
	Users_UserListInfo_FullMethodName   = "/user_rpc.Users/UserListInfo"
	Users_IsFriend_FullMethodName       = "/user_rpc.Users/IsFriend"
	Users_FriendList_FullMethodName     = "/user_rpc.Users/FriendList"
	Users_UserOnlineList_FullMethodName = "/user_rpc.Users/UserOnlineList"
)

// UsersClient is the client API for Users service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersClient interface {
	UserCreate(ctx context.Context, in *UserCreateRequest, opts ...grpc.CallOption) (*UserCreateResponse, error)
	UserInfo(ctx context.Context, in *UserInfoRequest, opts ...grpc.CallOption) (*UserInfoResponse, error)
	UserBaseInfo(ctx context.Context, in *UserBaseInfoRequest, opts ...grpc.CallOption) (*UserBaseInfoResponse, error)
	UserListInfo(ctx context.Context, in *UserListInfoRequest, opts ...grpc.CallOption) (*UserListInfoResponse, error)
	IsFriend(ctx context.Context, in *IsFriendRequest, opts ...grpc.CallOption) (*IsFriendResponse, error)
	FriendList(ctx context.Context, in *FriendListRequest, opts ...grpc.CallOption) (*FriendListResponse, error)
	UserOnlineList(ctx context.Context, in *UserOnlineRequest, opts ...grpc.CallOption) (*UserOnlineResponse, error)
}

type usersClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersClient(cc grpc.ClientConnInterface) UsersClient {
	return &usersClient{cc}
}

func (c *usersClient) UserCreate(ctx context.Context, in *UserCreateRequest, opts ...grpc.CallOption) (*UserCreateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserCreateResponse)
	err := c.cc.Invoke(ctx, Users_UserCreate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UserInfo(ctx context.Context, in *UserInfoRequest, opts ...grpc.CallOption) (*UserInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserInfoResponse)
	err := c.cc.Invoke(ctx, Users_UserInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UserBaseInfo(ctx context.Context, in *UserBaseInfoRequest, opts ...grpc.CallOption) (*UserBaseInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserBaseInfoResponse)
	err := c.cc.Invoke(ctx, Users_UserBaseInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UserListInfo(ctx context.Context, in *UserListInfoRequest, opts ...grpc.CallOption) (*UserListInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserListInfoResponse)
	err := c.cc.Invoke(ctx, Users_UserListInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) IsFriend(ctx context.Context, in *IsFriendRequest, opts ...grpc.CallOption) (*IsFriendResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IsFriendResponse)
	err := c.cc.Invoke(ctx, Users_IsFriend_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) FriendList(ctx context.Context, in *FriendListRequest, opts ...grpc.CallOption) (*FriendListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FriendListResponse)
	err := c.cc.Invoke(ctx, Users_FriendList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UserOnlineList(ctx context.Context, in *UserOnlineRequest, opts ...grpc.CallOption) (*UserOnlineResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserOnlineResponse)
	err := c.cc.Invoke(ctx, Users_UserOnlineList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServer is the server API for Users service.
// All implementations must embed UnimplementedUsersServer
// for forward compatibility
type UsersServer interface {
	UserCreate(context.Context, *UserCreateRequest) (*UserCreateResponse, error)
	UserInfo(context.Context, *UserInfoRequest) (*UserInfoResponse, error)
	UserBaseInfo(context.Context, *UserBaseInfoRequest) (*UserBaseInfoResponse, error)
	UserListInfo(context.Context, *UserListInfoRequest) (*UserListInfoResponse, error)
	IsFriend(context.Context, *IsFriendRequest) (*IsFriendResponse, error)
	FriendList(context.Context, *FriendListRequest) (*FriendListResponse, error)
	UserOnlineList(context.Context, *UserOnlineRequest) (*UserOnlineResponse, error)
	mustEmbedUnimplementedUsersServer()
}

// UnimplementedUsersServer must be embedded to have forward compatible implementations.
type UnimplementedUsersServer struct {
}

func (UnimplementedUsersServer) UserCreate(context.Context, *UserCreateRequest) (*UserCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserCreate not implemented")
}
func (UnimplementedUsersServer) UserInfo(context.Context, *UserInfoRequest) (*UserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserInfo not implemented")
}
func (UnimplementedUsersServer) UserBaseInfo(context.Context, *UserBaseInfoRequest) (*UserBaseInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserBaseInfo not implemented")
}
func (UnimplementedUsersServer) UserListInfo(context.Context, *UserListInfoRequest) (*UserListInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserListInfo not implemented")
}
func (UnimplementedUsersServer) IsFriend(context.Context, *IsFriendRequest) (*IsFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsFriend not implemented")
}
func (UnimplementedUsersServer) FriendList(context.Context, *FriendListRequest) (*FriendListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FriendList not implemented")
}
func (UnimplementedUsersServer) UserOnlineList(context.Context, *UserOnlineRequest) (*UserOnlineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserOnlineList not implemented")
}
func (UnimplementedUsersServer) mustEmbedUnimplementedUsersServer() {}

// UnsafeUsersServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersServer will
// result in compilation errors.
type UnsafeUsersServer interface {
	mustEmbedUnimplementedUsersServer()
}

func RegisterUsersServer(s grpc.ServiceRegistrar, srv UsersServer) {
	s.RegisterService(&Users_ServiceDesc, srv)
}

func _Users_UserCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UserCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Users_UserCreate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UserCreate(ctx, req.(*UserCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Users_UserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UserInfo(ctx, req.(*UserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UserBaseInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserBaseInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UserBaseInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Users_UserBaseInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UserBaseInfo(ctx, req.(*UserBaseInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UserListInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserListInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UserListInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Users_UserListInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UserListInfo(ctx, req.(*UserListInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_IsFriend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).IsFriend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Users_IsFriend_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).IsFriend(ctx, req.(*IsFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_FriendList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FriendListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).FriendList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Users_FriendList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).FriendList(ctx, req.(*FriendListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UserOnlineList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserOnlineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UserOnlineList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Users_UserOnlineList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UserOnlineList(ctx, req.(*UserOnlineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Users_ServiceDesc is the grpc.ServiceDesc for Users service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Users_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user_rpc.Users",
	HandlerType: (*UsersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserCreate",
			Handler:    _Users_UserCreate_Handler,
		},
		{
			MethodName: "UserInfo",
			Handler:    _Users_UserInfo_Handler,
		},
		{
			MethodName: "UserBaseInfo",
			Handler:    _Users_UserBaseInfo_Handler,
		},
		{
			MethodName: "UserListInfo",
			Handler:    _Users_UserListInfo_Handler,
		},
		{
			MethodName: "IsFriend",
			Handler:    _Users_IsFriend_Handler,
		},
		{
			MethodName: "FriendList",
			Handler:    _Users_FriendList_Handler,
		},
		{
			MethodName: "UserOnlineList",
			Handler:    _Users_UserOnlineList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user_rpc.proto",
}
