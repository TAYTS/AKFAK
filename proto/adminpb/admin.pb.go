// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/adminpb/admin.proto

package adminpb

import (
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

type Response int32

const (
	Response_SUCCESS Response = 0
	Response_FAIL    Response = 1
)

var Response_name = map[int32]string{
	0: "SUCCESS",
	1: "FAIL",
}

var Response_value = map[string]int32{
	"SUCCESS": 0,
	"FAIL":    1,
}

func (x Response) String() string {
	return proto.EnumName(Response_name, int32(x))
}

func (Response) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b89f257325c17cfa, []int{0}
}

type AdminClientNewTopicRequest struct {
	Topic                string   `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	NumPartitions        int32    `protobuf:"varint,2,opt,name=numPartitions,proto3" json:"numPartitions,omitempty"`
	ReplicationFactor    int32    `protobuf:"varint,3,opt,name=replicationFactor,proto3" json:"replicationFactor,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AdminClientNewTopicRequest) Reset()         { *m = AdminClientNewTopicRequest{} }
func (m *AdminClientNewTopicRequest) String() string { return proto.CompactTextString(m) }
func (*AdminClientNewTopicRequest) ProtoMessage()    {}
func (*AdminClientNewTopicRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b89f257325c17cfa, []int{0}
}

func (m *AdminClientNewTopicRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdminClientNewTopicRequest.Unmarshal(m, b)
}
func (m *AdminClientNewTopicRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdminClientNewTopicRequest.Marshal(b, m, deterministic)
}
func (m *AdminClientNewTopicRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdminClientNewTopicRequest.Merge(m, src)
}
func (m *AdminClientNewTopicRequest) XXX_Size() int {
	return xxx_messageInfo_AdminClientNewTopicRequest.Size(m)
}
func (m *AdminClientNewTopicRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AdminClientNewTopicRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AdminClientNewTopicRequest proto.InternalMessageInfo

func (m *AdminClientNewTopicRequest) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *AdminClientNewTopicRequest) GetNumPartitions() int32 {
	if m != nil {
		return m.NumPartitions
	}
	return 0
}

func (m *AdminClientNewTopicRequest) GetReplicationFactor() int32 {
	if m != nil {
		return m.ReplicationFactor
	}
	return 0
}

type AdminClientNewTopicResponse struct {
	Response             Response `protobuf:"varint,1,opt,name=response,proto3,enum=proto.adminpb.Response" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AdminClientNewTopicResponse) Reset()         { *m = AdminClientNewTopicResponse{} }
func (m *AdminClientNewTopicResponse) String() string { return proto.CompactTextString(m) }
func (*AdminClientNewTopicResponse) ProtoMessage()    {}
func (*AdminClientNewTopicResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b89f257325c17cfa, []int{1}
}

func (m *AdminClientNewTopicResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdminClientNewTopicResponse.Unmarshal(m, b)
}
func (m *AdminClientNewTopicResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdminClientNewTopicResponse.Marshal(b, m, deterministic)
}
func (m *AdminClientNewTopicResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdminClientNewTopicResponse.Merge(m, src)
}
func (m *AdminClientNewTopicResponse) XXX_Size() int {
	return xxx_messageInfo_AdminClientNewTopicResponse.Size(m)
}
func (m *AdminClientNewTopicResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AdminClientNewTopicResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AdminClientNewTopicResponse proto.InternalMessageInfo

func (m *AdminClientNewTopicResponse) GetResponse() Response {
	if m != nil {
		return m.Response
	}
	return Response_SUCCESS
}

type AdminClientNewPartitionRequest struct {
	Topic                string   `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	PartitionID          []int32  `protobuf:"varint,2,rep,packed,name=partitionID,proto3" json:"partitionID,omitempty"`
	ReplicaID            int32    `protobuf:"varint,3,opt,name=replicaID,proto3" json:"replicaID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AdminClientNewPartitionRequest) Reset()         { *m = AdminClientNewPartitionRequest{} }
func (m *AdminClientNewPartitionRequest) String() string { return proto.CompactTextString(m) }
func (*AdminClientNewPartitionRequest) ProtoMessage()    {}
func (*AdminClientNewPartitionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b89f257325c17cfa, []int{2}
}

func (m *AdminClientNewPartitionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdminClientNewPartitionRequest.Unmarshal(m, b)
}
func (m *AdminClientNewPartitionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdminClientNewPartitionRequest.Marshal(b, m, deterministic)
}
func (m *AdminClientNewPartitionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdminClientNewPartitionRequest.Merge(m, src)
}
func (m *AdminClientNewPartitionRequest) XXX_Size() int {
	return xxx_messageInfo_AdminClientNewPartitionRequest.Size(m)
}
func (m *AdminClientNewPartitionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AdminClientNewPartitionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AdminClientNewPartitionRequest proto.InternalMessageInfo

func (m *AdminClientNewPartitionRequest) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *AdminClientNewPartitionRequest) GetPartitionID() []int32 {
	if m != nil {
		return m.PartitionID
	}
	return nil
}

func (m *AdminClientNewPartitionRequest) GetReplicaID() int32 {
	if m != nil {
		return m.ReplicaID
	}
	return 0
}

type AdminClientNewPartitionResponse struct {
	Response             Response `protobuf:"varint,1,opt,name=response,proto3,enum=proto.adminpb.Response" json:"response,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AdminClientNewPartitionResponse) Reset()         { *m = AdminClientNewPartitionResponse{} }
func (m *AdminClientNewPartitionResponse) String() string { return proto.CompactTextString(m) }
func (*AdminClientNewPartitionResponse) ProtoMessage()    {}
func (*AdminClientNewPartitionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b89f257325c17cfa, []int{3}
}

func (m *AdminClientNewPartitionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdminClientNewPartitionResponse.Unmarshal(m, b)
}
func (m *AdminClientNewPartitionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdminClientNewPartitionResponse.Marshal(b, m, deterministic)
}
func (m *AdminClientNewPartitionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdminClientNewPartitionResponse.Merge(m, src)
}
func (m *AdminClientNewPartitionResponse) XXX_Size() int {
	return xxx_messageInfo_AdminClientNewPartitionResponse.Size(m)
}
func (m *AdminClientNewPartitionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AdminClientNewPartitionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AdminClientNewPartitionResponse proto.InternalMessageInfo

func (m *AdminClientNewPartitionResponse) GetResponse() Response {
	if m != nil {
		return m.Response
	}
	return Response_SUCCESS
}

func init() {
	proto.RegisterEnum("proto.adminpb.Response", Response_name, Response_value)
	proto.RegisterType((*AdminClientNewTopicRequest)(nil), "proto.adminpb.AdminClientNewTopicRequest")
	proto.RegisterType((*AdminClientNewTopicResponse)(nil), "proto.adminpb.AdminClientNewTopicResponse")
	proto.RegisterType((*AdminClientNewPartitionRequest)(nil), "proto.adminpb.AdminClientNewPartitionRequest")
	proto.RegisterType((*AdminClientNewPartitionResponse)(nil), "proto.adminpb.AdminClientNewPartitionResponse")
}

func init() {
	proto.RegisterFile("proto/adminpb/admin.proto", fileDescriptor_b89f257325c17cfa)
}

var fileDescriptor_b89f257325c17cfa = []byte{
	// 329 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2c, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x4c, 0xc9, 0xcd, 0xcc, 0x2b, 0x48, 0x82, 0xd0, 0x7a, 0x60, 0x31, 0x21, 0x5e,
	0x30, 0xa5, 0x07, 0x95, 0x52, 0x6a, 0x63, 0xe4, 0x92, 0x72, 0x04, 0xb1, 0x9d, 0x73, 0x32, 0x53,
	0xf3, 0x4a, 0xfc, 0x52, 0xcb, 0x43, 0xf2, 0x0b, 0x32, 0x93, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b,
	0x4b, 0x84, 0x44, 0xb8, 0x58, 0x4b, 0x40, 0x7c, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x08,
	0x47, 0x48, 0x85, 0x8b, 0x37, 0xaf, 0x34, 0x37, 0x20, 0xb1, 0xa8, 0x24, 0xb3, 0x24, 0x33, 0x3f,
	0xaf, 0x58, 0x82, 0x49, 0x81, 0x51, 0x83, 0x35, 0x08, 0x55, 0x50, 0x48, 0x87, 0x4b, 0xb0, 0x28,
	0xb5, 0x20, 0x27, 0x33, 0x39, 0x11, 0xc4, 0x77, 0x4b, 0x4c, 0x2e, 0xc9, 0x2f, 0x92, 0x60, 0x06,
	0xab, 0xc4, 0x94, 0x50, 0x0a, 0xe2, 0x92, 0xc6, 0xea, 0x8e, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54,
	0x21, 0x63, 0x2e, 0x8e, 0x22, 0x28, 0x1b, 0xec, 0x16, 0x3e, 0x23, 0x71, 0x3d, 0x14, 0x9f, 0xe8,
	0xc1, 0x94, 0x06, 0xc1, 0x15, 0x2a, 0x95, 0x71, 0xc9, 0xa1, 0x9a, 0x09, 0x77, 0x1d, 0x7e, 0xff,
	0x29, 0x70, 0x71, 0x17, 0xc0, 0x54, 0x7a, 0xba, 0x48, 0x30, 0x29, 0x30, 0x6b, 0xb0, 0x06, 0x21,
	0x0b, 0x09, 0xc9, 0x70, 0x71, 0x42, 0xbd, 0xe0, 0xe9, 0x02, 0xf5, 0x13, 0x42, 0x40, 0x29, 0x8c,
	0x4b, 0x1e, 0xa7, 0xbd, 0x14, 0xf8, 0x47, 0x4b, 0x91, 0x8b, 0x03, 0x6e, 0x00, 0x37, 0x17, 0x7b,
	0x70, 0xa8, 0xb3, 0xb3, 0x6b, 0x70, 0xb0, 0x00, 0x83, 0x10, 0x07, 0x17, 0x8b, 0x9b, 0xa3, 0xa7,
	0x8f, 0x00, 0xa3, 0xd1, 0x2f, 0x46, 0x2e, 0x1e, 0xb0, 0xdd, 0xc1, 0xa9, 0x45, 0x65, 0x99, 0xc9,
	0xa9, 0x42, 0x79, 0x5c, 0xc2, 0x58, 0xc2, 0x55, 0x48, 0x13, 0xcd, 0x36, 0xdc, 0x69, 0x40, 0x4a,
	0x8b, 0x18, 0xa5, 0xd0, 0x10, 0x67, 0x10, 0xaa, 0xe2, 0x12, 0xc7, 0xe1, 0x77, 0x21, 0x5d, 0xbc,
	0x06, 0xa1, 0xc7, 0x8d, 0x94, 0x1e, 0xb1, 0xca, 0x61, 0x76, 0x3b, 0x89, 0x46, 0x09, 0x3b, 0x7a,
	0xbb, 0x39, 0x7a, 0xeb, 0xa3, 0x24, 0xff, 0x24, 0x36, 0x30, 0xd7, 0x18, 0x10, 0x00, 0x00, 0xff,
	0xff, 0xc5, 0xed, 0x64, 0x98, 0x16, 0x03, 0x00, 0x00,
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
	AdminClientNewTopic(ctx context.Context, in *AdminClientNewTopicRequest, opts ...grpc.CallOption) (*AdminClientNewTopicResponse, error)
	AdminClientNewPartition(ctx context.Context, in *AdminClientNewPartitionRequest, opts ...grpc.CallOption) (*AdminClientNewPartitionResponse, error)
}

type adminServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminServiceClient(cc grpc.ClientConnInterface) AdminServiceClient {
	return &adminServiceClient{cc}
}

func (c *adminServiceClient) AdminClientNewTopic(ctx context.Context, in *AdminClientNewTopicRequest, opts ...grpc.CallOption) (*AdminClientNewTopicResponse, error) {
	out := new(AdminClientNewTopicResponse)
	err := c.cc.Invoke(ctx, "/proto.adminpb.AdminService/AdminClientNewTopic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminServiceClient) AdminClientNewPartition(ctx context.Context, in *AdminClientNewPartitionRequest, opts ...grpc.CallOption) (*AdminClientNewPartitionResponse, error) {
	out := new(AdminClientNewPartitionResponse)
	err := c.cc.Invoke(ctx, "/proto.adminpb.AdminService/AdminClientNewPartition", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServiceServer is the server API for AdminService service.
type AdminServiceServer interface {
	AdminClientNewTopic(context.Context, *AdminClientNewTopicRequest) (*AdminClientNewTopicResponse, error)
	AdminClientNewPartition(context.Context, *AdminClientNewPartitionRequest) (*AdminClientNewPartitionResponse, error)
}

// UnimplementedAdminServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAdminServiceServer struct {
}

func (*UnimplementedAdminServiceServer) AdminClientNewTopic(ctx context.Context, req *AdminClientNewTopicRequest) (*AdminClientNewTopicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AdminClientNewTopic not implemented")
}
func (*UnimplementedAdminServiceServer) AdminClientNewPartition(ctx context.Context, req *AdminClientNewPartitionRequest) (*AdminClientNewPartitionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AdminClientNewPartition not implemented")
}

func RegisterAdminServiceServer(s *grpc.Server, srv AdminServiceServer) {
	s.RegisterService(&_AdminService_serviceDesc, srv)
}

func _AdminService_AdminClientNewTopic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdminClientNewTopicRequest)
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
		return srv.(AdminServiceServer).AdminClientNewTopic(ctx, req.(*AdminClientNewTopicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminService_AdminClientNewPartition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdminClientNewPartitionRequest)
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
		return srv.(AdminServiceServer).AdminClientNewPartition(ctx, req.(*AdminClientNewPartitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AdminService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.adminpb.AdminService",
	HandlerType: (*AdminServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AdminClientNewTopic",
			Handler:    _AdminService_AdminClientNewTopic_Handler,
		},
		{
			MethodName: "AdminClientNewPartition",
			Handler:    _AdminService_AdminClientNewPartition_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/adminpb/admin.proto",
}
