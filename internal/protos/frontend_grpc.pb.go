// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protos

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

// FrontendClient is the client API for Frontend service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FrontendClient interface {
	RequestServiceJWT(ctx context.Context, in *RequestJWTMessage, opts ...grpc.CallOption) (*JWTResponse, error)
	GetHeroData(ctx context.Context, in *RequestHeroMessage, opts ...grpc.CallOption) (*HeroData, error)
	GetConnectionStatus(ctx context.Context, in *ConnectionStatusMessage, opts ...grpc.CallOption) (*ConnectionStatusResponse, error)
	NewGameJWT(ctx context.Context, in *RequestGameJWTMessage, opts ...grpc.CallOption) (*JWTResponse, error)
	GetGameJWT(ctx context.Context, in *RequestGameJWTMessage, opts ...grpc.CallOption) (*JWTResponse, error)
}

type frontendClient struct {
	cc grpc.ClientConnInterface
}

func NewFrontendClient(cc grpc.ClientConnInterface) FrontendClient {
	return &frontendClient{cc}
}

func (c *frontendClient) RequestServiceJWT(ctx context.Context, in *RequestJWTMessage, opts ...grpc.CallOption) (*JWTResponse, error) {
	out := new(JWTResponse)
	err := c.cc.Invoke(ctx, "/blt.adoptahero.Frontend/RequestServiceJWT", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frontendClient) GetHeroData(ctx context.Context, in *RequestHeroMessage, opts ...grpc.CallOption) (*HeroData, error) {
	out := new(HeroData)
	err := c.cc.Invoke(ctx, "/blt.adoptahero.Frontend/GetHeroData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frontendClient) GetConnectionStatus(ctx context.Context, in *ConnectionStatusMessage, opts ...grpc.CallOption) (*ConnectionStatusResponse, error) {
	out := new(ConnectionStatusResponse)
	err := c.cc.Invoke(ctx, "/blt.adoptahero.Frontend/GetConnectionStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frontendClient) NewGameJWT(ctx context.Context, in *RequestGameJWTMessage, opts ...grpc.CallOption) (*JWTResponse, error) {
	out := new(JWTResponse)
	err := c.cc.Invoke(ctx, "/blt.adoptahero.Frontend/NewGameJWT", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *frontendClient) GetGameJWT(ctx context.Context, in *RequestGameJWTMessage, opts ...grpc.CallOption) (*JWTResponse, error) {
	out := new(JWTResponse)
	err := c.cc.Invoke(ctx, "/blt.adoptahero.Frontend/GetGameJWT", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FrontendServer is the server API for Frontend service.
// All implementations must embed UnimplementedFrontendServer
// for forward compatibility
type FrontendServer interface {
	RequestServiceJWT(context.Context, *RequestJWTMessage) (*JWTResponse, error)
	GetHeroData(context.Context, *RequestHeroMessage) (*HeroData, error)
	GetConnectionStatus(context.Context, *ConnectionStatusMessage) (*ConnectionStatusResponse, error)
	NewGameJWT(context.Context, *RequestGameJWTMessage) (*JWTResponse, error)
	GetGameJWT(context.Context, *RequestGameJWTMessage) (*JWTResponse, error)
	mustEmbedUnimplementedFrontendServer()
}

// UnimplementedFrontendServer must be embedded to have forward compatible implementations.
type UnimplementedFrontendServer struct {
}

func (UnimplementedFrontendServer) RequestServiceJWT(context.Context, *RequestJWTMessage) (*JWTResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestServiceJWT not implemented")
}
func (UnimplementedFrontendServer) GetHeroData(context.Context, *RequestHeroMessage) (*HeroData, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHeroData not implemented")
}
func (UnimplementedFrontendServer) GetConnectionStatus(context.Context, *ConnectionStatusMessage) (*ConnectionStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConnectionStatus not implemented")
}
func (UnimplementedFrontendServer) NewGameJWT(context.Context, *RequestGameJWTMessage) (*JWTResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewGameJWT not implemented")
}
func (UnimplementedFrontendServer) GetGameJWT(context.Context, *RequestGameJWTMessage) (*JWTResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGameJWT not implemented")
}
func (UnimplementedFrontendServer) mustEmbedUnimplementedFrontendServer() {}

// UnsafeFrontendServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FrontendServer will
// result in compilation errors.
type UnsafeFrontendServer interface {
	mustEmbedUnimplementedFrontendServer()
}

func RegisterFrontendServer(s grpc.ServiceRegistrar, srv FrontendServer) {
	s.RegisterService(&Frontend_ServiceDesc, srv)
}

func _Frontend_RequestServiceJWT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestJWTMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontendServer).RequestServiceJWT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blt.adoptahero.Frontend/RequestServiceJWT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontendServer).RequestServiceJWT(ctx, req.(*RequestJWTMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Frontend_GetHeroData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestHeroMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontendServer).GetHeroData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blt.adoptahero.Frontend/GetHeroData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontendServer).GetHeroData(ctx, req.(*RequestHeroMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Frontend_GetConnectionStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectionStatusMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontendServer).GetConnectionStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blt.adoptahero.Frontend/GetConnectionStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontendServer).GetConnectionStatus(ctx, req.(*ConnectionStatusMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Frontend_NewGameJWT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestGameJWTMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontendServer).NewGameJWT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blt.adoptahero.Frontend/NewGameJWT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontendServer).NewGameJWT(ctx, req.(*RequestGameJWTMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Frontend_GetGameJWT_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestGameJWTMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FrontendServer).GetGameJWT(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/blt.adoptahero.Frontend/GetGameJWT",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FrontendServer).GetGameJWT(ctx, req.(*RequestGameJWTMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Frontend_ServiceDesc is the grpc.ServiceDesc for Frontend service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Frontend_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "blt.adoptahero.Frontend",
	HandlerType: (*FrontendServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestServiceJWT",
			Handler:    _Frontend_RequestServiceJWT_Handler,
		},
		{
			MethodName: "GetHeroData",
			Handler:    _Frontend_GetHeroData_Handler,
		},
		{
			MethodName: "GetConnectionStatus",
			Handler:    _Frontend_GetConnectionStatus_Handler,
		},
		{
			MethodName: "NewGameJWT",
			Handler:    _Frontend_NewGameJWT_Handler,
		},
		{
			MethodName: "GetGameJWT",
			Handler:    _Frontend_GetGameJWT_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "frontend.proto",
}
