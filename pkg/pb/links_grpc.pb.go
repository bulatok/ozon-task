// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: links.proto

package pb

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

// LinksClient is the client API for Links service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LinksClient interface {
	ShortLink(ctx context.Context, in *ShortLinkRequest, opts ...grpc.CallOption) (*ShortLinkResponse, error)
	GetOriginalLink(ctx context.Context, in *GetOriginalRequest, opts ...grpc.CallOption) (*GetOriginalResponse, error)
}

type linksClient struct {
	cc grpc.ClientConnInterface
}

func NewLinksClient(cc grpc.ClientConnInterface) LinksClient {
	return &linksClient{cc}
}

func (c *linksClient) ShortLink(ctx context.Context, in *ShortLinkRequest, opts ...grpc.CallOption) (*ShortLinkResponse, error) {
	out := new(ShortLinkResponse)
	err := c.cc.Invoke(ctx, "/links.Links/OriginalLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *linksClient) GetOriginalLink(ctx context.Context, in *GetOriginalRequest, opts ...grpc.CallOption) (*GetOriginalResponse, error) {
	out := new(GetOriginalResponse)
	err := c.cc.Invoke(ctx, "/links.Links/GetOriginalLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LinksServer is the server API for Links service.
// All implementations must embed UnimplementedLinksServer
// for forward compatibility
type LinksServer interface {
	ShortLink(context.Context, *ShortLinkRequest) (*ShortLinkResponse, error)
	GetOriginalLink(context.Context, *GetOriginalRequest) (*GetOriginalResponse, error)
	mustEmbedUnimplementedLinksServer()
}

// UnimplementedLinksServer must be embedded to have forward compatible implementations.
type UnimplementedLinksServer struct {
}

func (UnimplementedLinksServer) ShortLink(context.Context, *ShortLinkRequest) (*ShortLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OriginalLink not implemented")
}
func (UnimplementedLinksServer) GetOriginalLink(context.Context, *GetOriginalRequest) (*GetOriginalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginalLink not implemented")
}
func (UnimplementedLinksServer) mustEmbedUnimplementedLinksServer() {}

// UnsafeLinksServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LinksServer will
// result in compilation errors.
type UnsafeLinksServer interface {
	mustEmbedUnimplementedLinksServer()
}

func RegisterLinksServer(s grpc.ServiceRegistrar, srv LinksServer) {
	s.RegisterService(&Links_ServiceDesc, srv)
}

func _Links_ShortLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinksServer).ShortLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/links.Links/OriginalLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinksServer).ShortLink(ctx, req.(*ShortLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Links_GetOriginalLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOriginalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinksServer).GetOriginalLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/links.Links/GetOriginalLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinksServer).GetOriginalLink(ctx, req.(*GetOriginalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Links_ServiceDesc is the grpc.ServiceDesc for Links service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Links_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "links.Links",
	HandlerType: (*LinksServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "OriginalLink",
			Handler:    _Links_ShortLink_Handler,
		},
		{
			MethodName: "GetOriginalLink",
			Handler:    _Links_GetOriginalLink_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "links.proto",
}