// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/adminpb/admin.proto

package adminpb

import (
	adminclientpb "AKFAK/proto/adminclientpb"
	heartbeatspb "AKFAK/proto/heartbeatspb"
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

func init() {
	proto.RegisterFile("proto/adminpb/admin.proto", fileDescriptor_b89f257325c17cfa)
}

var fileDescriptor_b89f257325c17cfa = []byte{
	// 385 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0xcf, 0x4e, 0x2a, 0x31,
	0x18, 0xc5, 0x2f, 0x1b, 0x16, 0x0d, 0xdc, 0x45, 0xc9, 0xcd, 0x8d, 0x2c, 0x35, 0x2a, 0x60, 0x1c,
	0x88, 0xf0, 0x02, 0x48, 0xfc, 0x93, 0x10, 0x8c, 0x11, 0xdd, 0xb8, 0x99, 0x74, 0x86, 0x2f, 0x30,
	0x71, 0x98, 0xd6, 0xe9, 0x87, 0x04, 0x1f, 0xc0, 0xf7, 0xf1, 0x0d, 0x0d, 0x43, 0x5b, 0x3a, 0xb1,
	0x89, 0x75, 0xd5, 0x74, 0xfa, 0x3b, 0xe7, 0x7c, 0xd3, 0x39, 0x19, 0x72, 0x20, 0x72, 0x8e, 0xbc,
	0xcb, 0x66, 0xcb, 0x24, 0x13, 0xd1, 0x6e, 0x0d, 0x8a, 0x67, 0xb4, 0x5e, 0x2c, 0x81, 0x3a, 0x6a,
	0x0e, 0x2c, 0x32, 0x4e, 0x13, 0xc8, 0x50, 0xf3, 0xe1, 0x6e, 0x1b, 0x66, 0xb0, 0x0e, 0x05, 0xcb,
	0x31, 0xc1, 0x84, 0x2b, 0x93, 0x66, 0xcf, 0x4b, 0x85, 0x5c, 0x24, 0xb1, 0x52, 0x9c, 0xbb, 0x14,
	0x31, 0xcf, 0x30, 0xe7, 0x69, 0x0a, 0x79, 0x08, 0x29, 0xc4, 0x56, 0x40, 0xdb, 0x85, 0xaf, 0xc4,
	0x8c, 0x21, 0x84, 0x4b, 0x40, 0x36, 0x63, 0xc8, 0x14, 0xda, 0x72, 0xa1, 0x73, 0xc0, 0x70, 0xef,
	0xae, 0xc8, 0xa3, 0x1d, 0xb9, 0x00, 0x96, 0x63, 0x04, 0x0c, 0xa5, 0x88, 0xac, 0x8d, 0x82, 0x4e,
	0x5d, 0x76, 0x72, 0x93, 0xc5, 0xe1, 0x12, 0xa4, 0x64, 0x73, 0x50, 0xe0, 0xc5, 0x67, 0x95, 0xd4,
	0x86, 0x5b, 0x6a, 0x0a, 0xf9, 0x5b, 0x12, 0x03, 0x5d, 0x13, 0x3a, 0x32, 0x91, 0x57, 0xea, 0x7d,
	0x68, 0x10, 0x58, 0x17, 0xae, 0x0d, 0x83, 0xef, 0xe0, 0x03, 0xbc, 0xae, 0x40, 0x62, 0xb3, 0xeb,
	0xcd, 0x4b, 0xc1, 0x33, 0x09, 0x87, 0x7f, 0xe8, 0x3b, 0x69, 0x14, 0x83, 0x8c, 0x0a, 0xfa, 0x0e,
	0xd6, 0x8f, 0xdb, 0x8b, 0xa7, 0x6e, 0x27, 0x07, 0xa9, 0xa3, 0x7b, 0xfe, 0x02, 0x93, 0xfd, 0x51,
	0x21, 0xff, 0xcb, 0xc4, 0xbd, 0xee, 0x0a, 0xed, 0x7b, 0xf8, 0x19, 0x5a, 0x0f, 0x31, 0xf8, 0x9d,
	0xc8, 0x0c, 0xf2, 0x42, 0xfe, 0x3e, 0x15, 0xfd, 0x98, 0xa8, 0x7a, 0xd0, 0x8e, 0xd3, 0xa9, 0x0c,
	0xe9, 0xd4, 0x33, 0x2f, 0xd6, 0x84, 0x2d, 0x48, 0xfd, 0x06, 0x70, 0xff, 0x51, 0x68, 0xdb, 0xa9,
	0x2f, 0x31, 0x3a, 0xaa, 0xe3, 0x83, 0x9a, 0x24, 0x46, 0xc8, 0xad, 0xa9, 0x28, 0x3d, 0x56, 0x5a,
	0xbb, 0xc2, 0xc1, 0xfe, 0x5c, 0x47, 0x9c, 0xfc, 0x84, 0x69, 0xfb, 0x56, 0xa5, 0x57, 0xa1, 0x09,
	0xa9, 0x4d, 0x37, 0x59, 0x3c, 0x51, 0xf5, 0xa6, 0x2d, 0xe7, 0x80, 0x36, 0xa2, 0x73, 0xda, 0x1e,
	0xa4, 0x1d, 0x75, 0xf9, 0xef, 0xb9, 0x31, 0x1c, 0x5f, 0x0f, 0xc7, 0xdd, 0xd2, 0xff, 0x29, 0xaa,
	0x16, 0xdb, 0xfe, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0xaf, 0x47, 0xa7, 0x55, 0xb7, 0x04, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AdminServiceClient is the client API for AdminService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AdminServiceClient interface {
	ControllerElection(ctx context.Context, in *adminclientpb.ControllerElectionRequest, opts ...grpc.CallOption) (*adminclientpb.ControllerElectionResponse, error)
	AdminClientNewTopic(ctx context.Context, in *adminclientpb.AdminClientNewTopicRequest, opts ...grpc.CallOption) (*adminclientpb.AdminClientNewTopicResponse, error)
	AdminClientNewPartition(ctx context.Context, in *adminclientpb.AdminClientNewPartitionRequest, opts ...grpc.CallOption) (*adminclientpb.AdminClientNewPartitionResponse, error)
	UpdateMetadata(ctx context.Context, in *adminclientpb.UpdateMetadataRequest, opts ...grpc.CallOption) (*adminclientpb.UpdateMetadataResponse, error)
	GetController(ctx context.Context, in *adminclientpb.GetControllerRequest, opts ...grpc.CallOption) (*adminclientpb.GetControllerResponse, error)
	Heartbeats(ctx context.Context, opts ...grpc.CallOption) (AdminService_HeartbeatsClient, error)
	SyncMessages(ctx context.Context, opts ...grpc.CallOption) (AdminService_SyncMessagesClient, error)
}

type adminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminServiceClient(cc grpc.ClientConnInterface) AdminServiceClient {
	return &adminServiceClient{cc}
}

func (c *adminServiceClient) ControllerElection(ctx context.Context, in *adminclientpb.ControllerElectionRequest, opts ...grpc.CallOption) (*adminclientpb.ControllerElectionResponse, error) {
	out := new(adminclientpb.ControllerElectionResponse)
	err := c.cc.Invoke(ctx, "/proto.adminpb.AdminService/ControllerElection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) AdminClientNewTopic(ctx context.Context, in *adminclientpb.AdminClientNewTopicRequest, opts ...grpc.CallOption) (*adminclientpb.AdminClientNewTopicResponse, error) {
	out := new(adminclientpb.AdminClientNewTopicResponse)
	err := c.cc.Invoke(ctx, "/proto.adminpb.AdminService/AdminClientNewTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) AdminClientNewPartition(ctx context.Context, in *adminclientpb.AdminClientNewPartitionRequest, opts ...grpc.CallOption) (*adminclientpb.AdminClientNewPartitionResponse, error) {
	out := new(adminclientpb.AdminClientNewPartitionResponse)
	err := c.cc.Invoke(ctx, "/proto.adminpb.AdminService/AdminClientNewPartition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) UpdateMetadata(ctx context.Context, in *adminclientpb.UpdateMetadataRequest, opts ...grpc.CallOption) (*adminclientpb.UpdateMetadataResponse, error) {
	out := new(adminclientpb.UpdateMetadataResponse)
	err := c.cc.Invoke(ctx, "/proto.adminpb.AdminService/UpdateMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) GetController(ctx context.Context, in *adminclientpb.GetControllerRequest, opts ...grpc.CallOption) (*adminclientpb.GetControllerResponse, error) {
	out := new(adminclientpb.GetControllerResponse)
	err := c.cc.Invoke(ctx, "/proto.adminpb.AdminService/GetController", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) Heartbeats(ctx context.Context, opts ...grpc.CallOption) (AdminService_HeartbeatsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AdminService_serviceDesc.Streams[0], "/proto.adminpb.AdminService/Heartbeats", opts...)
	if err != nil {
		return nil, err
	}
	x := &adminServiceHeartbeatsClient{stream}
	return x, nil
}

type AdminService_HeartbeatsClient interface {
	Send(*heartbeatspb.HeartbeatsRequest) error
	Recv() (*heartbeatspb.HeartbeatsResponse, error)
	grpc.ClientStream
}

type adminServiceHeartbeatsClient struct {
	grpc.ClientStream
}

func (x *adminServiceHeartbeatsClient) Send(m *heartbeatspb.HeartbeatsRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *adminServiceHeartbeatsClient) Recv() (*heartbeatspb.HeartbeatsResponse, error) {
	m := new(heartbeatspb.HeartbeatsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *adminServiceClient) SyncMessages(ctx context.Context, opts ...grpc.CallOption) (AdminService_SyncMessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AdminService_serviceDesc.Streams[1], "/proto.adminpb.AdminService/SyncMessages", opts...)
	if err != nil {
		return nil, err
	}
	x := &adminServiceSyncMessagesClient{stream}
	return x, nil
}

type AdminService_SyncMessagesClient interface {
	Send(*adminclientpb.SyncMessagesRequest) error
	Recv() (*adminclientpb.SyncMessagesResponse, error)
	grpc.ClientStream
}

type adminServiceSyncMessagesClient struct {
	grpc.ClientStream
}

func (x *adminServiceSyncMessagesClient) Send(m *adminclientpb.SyncMessagesRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *adminServiceSyncMessagesClient) Recv() (*adminclientpb.SyncMessagesResponse, error) {
	m := new(adminclientpb.SyncMessagesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AdminServiceServer is the server API for AdminService service.
type AdminServiceServer interface {
	ControllerElection(context.Context, *adminclientpb.ControllerElectionRequest) (*adminclientpb.ControllerElectionResponse, error)
	AdminClientNewTopic(context.Context, *adminclientpb.AdminClientNewTopicRequest) (*adminclientpb.AdminClientNewTopicResponse, error)
	AdminClientNewPartition(context.Context, *adminclientpb.AdminClientNewPartitionRequest) (*adminclientpb.AdminClientNewPartitionResponse, error)
	UpdateMetadata(context.Context, *adminclientpb.UpdateMetadataRequest) (*adminclientpb.UpdateMetadataResponse, error)
	GetController(context.Context, *adminclientpb.GetControllerRequest) (*adminclientpb.GetControllerResponse, error)
	Heartbeats(AdminService_HeartbeatsServer) error
	SyncMessages(AdminService_SyncMessagesServer) error
}

// UnimplementedAdminServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAdminServiceServer struct {
}

func (*UnimplementedAdminServiceServer) ControllerElection(ctx context.Context, req *adminclientpb.ControllerElectionRequest) (*adminclientpb.ControllerElectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ControllerElection not implemented")
}
func (*UnimplementedAdminServiceServer) AdminClientNewTopic(ctx context.Context, req *adminclientpb.AdminClientNewTopicRequest) (*adminclientpb.AdminClientNewTopicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AdminClientNewTopic not implemented")
}
func (*UnimplementedAdminServiceServer) AdminClientNewPartition(ctx context.Context, req *adminclientpb.AdminClientNewPartitionRequest) (*adminclientpb.AdminClientNewPartitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AdminClientNewPartition not implemented")
}
func (*UnimplementedAdminServiceServer) UpdateMetadata(ctx context.Context, req *adminclientpb.UpdateMetadataRequest) (*adminclientpb.UpdateMetadataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetadata not implemented")
}
func (*UnimplementedAdminServiceServer) GetController(ctx context.Context, req *adminclientpb.GetControllerRequest) (*adminclientpb.GetControllerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetController not implemented")
}
func (*UnimplementedAdminServiceServer) Heartbeats(srv AdminService_HeartbeatsServer) error {
	return status.Errorf(codes.Unimplemented, "method Heartbeats not implemented")
}
func (*UnimplementedAdminServiceServer) SyncMessages(srv AdminService_SyncMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method SyncMessages not implemented")
}

func RegisterAdminServiceServer(s *grpc.Server, srv AdminServiceServer) {
	s.RegisterService(&_AdminService_serviceDesc, srv)
}

func _AdminService_ControllerElection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(adminclientpb.ControllerElectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).ControllerElection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.adminpb.AdminService/ControllerElection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).ControllerElection(ctx, req.(*adminclientpb.ControllerElectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_AdminClientNewTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(adminclientpb.AdminClientNewTopicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).AdminClientNewTopic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.adminpb.AdminService/AdminClientNewTopic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).AdminClientNewTopic(ctx, req.(*adminclientpb.AdminClientNewTopicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_AdminClientNewPartition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(adminclientpb.AdminClientNewPartitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).AdminClientNewPartition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.adminpb.AdminService/AdminClientNewPartition",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).AdminClientNewPartition(ctx, req.(*adminclientpb.AdminClientNewPartitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_UpdateMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(adminclientpb.UpdateMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).UpdateMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.adminpb.AdminService/UpdateMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).UpdateMetadata(ctx, req.(*adminclientpb.UpdateMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_GetController_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(adminclientpb.GetControllerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServiceServer).GetController(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.adminpb.AdminService/GetController",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServiceServer).GetController(ctx, req.(*adminclientpb.GetControllerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_Heartbeats_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AdminServiceServer).Heartbeats(&adminServiceHeartbeatsServer{stream})
}

type AdminService_HeartbeatsServer interface {
	Send(*heartbeatspb.HeartbeatsResponse) error
	Recv() (*heartbeatspb.HeartbeatsRequest, error)
	grpc.ServerStream
}

type adminServiceHeartbeatsServer struct {
	grpc.ServerStream
}

func (x *adminServiceHeartbeatsServer) Send(m *heartbeatspb.HeartbeatsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *adminServiceHeartbeatsServer) Recv() (*heartbeatspb.HeartbeatsRequest, error) {
	m := new(heartbeatspb.HeartbeatsRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _AdminService_SyncMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AdminServiceServer).SyncMessages(&adminServiceSyncMessagesServer{stream})
}

type AdminService_SyncMessagesServer interface {
	Send(*adminclientpb.SyncMessagesResponse) error
	Recv() (*adminclientpb.SyncMessagesRequest, error)
	grpc.ServerStream
}

type adminServiceSyncMessagesServer struct {
	grpc.ServerStream
}

func (x *adminServiceSyncMessagesServer) Send(m *adminclientpb.SyncMessagesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *adminServiceSyncMessagesServer) Recv() (*adminclientpb.SyncMessagesRequest, error) {
	m := new(adminclientpb.SyncMessagesRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _AdminService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.adminpb.AdminService",
	HandlerType: (*AdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ControllerElection",
			Handler:    _AdminService_ControllerElection_Handler,
		},
		{
			MethodName: "AdminClientNewTopic",
			Handler:    _AdminService_AdminClientNewTopic_Handler,
		},
		{
			MethodName: "AdminClientNewPartition",
			Handler:    _AdminService_AdminClientNewPartition_Handler,
		},
		{
			MethodName: "UpdateMetadata",
			Handler:    _AdminService_UpdateMetadata_Handler,
		},
		{
			MethodName: "GetController",
			Handler:    _AdminService_GetController_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Heartbeats",
			Handler:       _AdminService_Heartbeats_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "SyncMessages",
			Handler:       _AdminService_SyncMessages_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/adminpb/admin.proto",
}
