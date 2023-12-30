// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: framework/grpc/proto/choice.proto

package choice

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

// ChoiceServiceClient is the client API for ChoiceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChoiceServiceClient interface {
	GetChoiceListById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*ChoiceList, error)
	UpsertChoiceList(ctx context.Context, in *ChoiceList, opts ...grpc.CallOption) (*UpsertResponse, error)
}

type choiceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChoiceServiceClient(cc grpc.ClientConnInterface) ChoiceServiceClient {
	return &choiceServiceClient{cc}
}

func (c *choiceServiceClient) GetChoiceListById(ctx context.Context, in *GetByIdRequest, opts ...grpc.CallOption) (*ChoiceList, error) {
	out := new(ChoiceList)
	err := c.cc.Invoke(ctx, "/choice.ChoiceService/GetChoiceListById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *choiceServiceClient) UpsertChoiceList(ctx context.Context, in *ChoiceList, opts ...grpc.CallOption) (*UpsertResponse, error) {
	out := new(UpsertResponse)
	err := c.cc.Invoke(ctx, "/choice.ChoiceService/UpsertChoiceList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChoiceServiceServer is the server API for ChoiceService service.
// All implementations should embed UnimplementedChoiceServiceServer
// for forward compatibility
type ChoiceServiceServer interface {
	GetChoiceListById(context.Context, *GetByIdRequest) (*ChoiceList, error)
	UpsertChoiceList(context.Context, *ChoiceList) (*UpsertResponse, error)
}

// UnimplementedChoiceServiceServer should be embedded to have forward compatible implementations.
type UnimplementedChoiceServiceServer struct {
}

func (UnimplementedChoiceServiceServer) GetChoiceListById(context.Context, *GetByIdRequest) (*ChoiceList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChoiceListById not implemented")
}
func (UnimplementedChoiceServiceServer) UpsertChoiceList(context.Context, *ChoiceList) (*UpsertResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertChoiceList not implemented")
}

// UnsafeChoiceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChoiceServiceServer will
// result in compilation errors.
type UnsafeChoiceServiceServer interface {
	mustEmbedUnimplementedChoiceServiceServer()
}

func RegisterChoiceServiceServer(s grpc.ServiceRegistrar, srv ChoiceServiceServer) {
	s.RegisterService(&ChoiceService_ServiceDesc, srv)
}

func _ChoiceService_GetChoiceListById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChoiceServiceServer).GetChoiceListById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/choice.ChoiceService/GetChoiceListById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChoiceServiceServer).GetChoiceListById(ctx, req.(*GetByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChoiceService_UpsertChoiceList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChoiceList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChoiceServiceServer).UpsertChoiceList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/choice.ChoiceService/UpsertChoiceList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChoiceServiceServer).UpsertChoiceList(ctx, req.(*ChoiceList))
	}
	return interceptor(ctx, in, info, handler)
}

// ChoiceService_ServiceDesc is the grpc.ServiceDesc for ChoiceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChoiceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "choice.ChoiceService",
	HandlerType: (*ChoiceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetChoiceListById",
			Handler:    _ChoiceService_GetChoiceListById_Handler,
		},
		{
			MethodName: "UpsertChoiceList",
			Handler:    _ChoiceService_UpsertChoiceList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "framework/grpc/proto/choice.proto",
}
