// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/adminclientpb/heartbeats.proto

package adminclientpb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type HeartbeatsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatsRequest) Reset()         { *m = HeartbeatsRequest{} }
func (m *HeartbeatsRequest) String() string { return proto.CompactTextString(m) }
func (*HeartbeatsRequest) ProtoMessage()    {}
func (*HeartbeatsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_77be30deda3eb675, []int{0}
}

func (m *HeartbeatsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatsRequest.Unmarshal(m, b)
}
func (m *HeartbeatsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatsRequest.Marshal(b, m, deterministic)
}
func (m *HeartbeatsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatsRequest.Merge(m, src)
}
func (m *HeartbeatsRequest) XXX_Size() int {
	return xxx_messageInfo_HeartbeatsRequest.Size(m)
}
func (m *HeartbeatsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatsRequest proto.InternalMessageInfo

type HeartbeatsResponse struct {
	BrokerID             int32    `protobuf:"varint,1,opt,name=brokerID,proto3" json:"brokerID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatsResponse) Reset()         { *m = HeartbeatsResponse{} }
func (m *HeartbeatsResponse) String() string { return proto.CompactTextString(m) }
func (*HeartbeatsResponse) ProtoMessage()    {}
func (*HeartbeatsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_77be30deda3eb675, []int{1}
}

func (m *HeartbeatsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatsResponse.Unmarshal(m, b)
}
func (m *HeartbeatsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatsResponse.Marshal(b, m, deterministic)
}
func (m *HeartbeatsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatsResponse.Merge(m, src)
}
func (m *HeartbeatsResponse) XXX_Size() int {
	return xxx_messageInfo_HeartbeatsResponse.Size(m)
}
func (m *HeartbeatsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatsResponse proto.InternalMessageInfo

func (m *HeartbeatsResponse) GetBrokerID() int32 {
	if m != nil {
		return m.BrokerID
	}
	return 0
}

func init() {
	proto.RegisterType((*HeartbeatsRequest)(nil), "proto.adminclientpb.HeartbeatsRequest")
	proto.RegisterType((*HeartbeatsResponse)(nil), "proto.adminclientpb.HeartbeatsResponse")
}

func init() {
	proto.RegisterFile("proto/adminclientpb/heartbeats.proto", fileDescriptor_77be30deda3eb675)
}

var fileDescriptor_77be30deda3eb675 = []byte{
	// 130 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x29, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0x4c, 0xc9, 0xcd, 0xcc, 0x4b, 0xce, 0xc9, 0x4c, 0xcd, 0x2b, 0x29, 0x48, 0xd2,
	0xcf, 0x48, 0x4d, 0x2c, 0x2a, 0x49, 0x4a, 0x4d, 0x2c, 0x29, 0xd6, 0x03, 0x4b, 0x0b, 0x09, 0x83,
	0x29, 0x3d, 0x14, 0x55, 0x4a, 0xc2, 0x5c, 0x82, 0x1e, 0x70, 0x85, 0x41, 0xa9, 0x85, 0xa5, 0xa9,
	0xc5, 0x25, 0x4a, 0x06, 0x5c, 0x42, 0xc8, 0x82, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0x52,
	0x5c, 0x1c, 0x49, 0x45, 0xf9, 0xd9, 0xa9, 0x45, 0x9e, 0x2e, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xac,
	0x41, 0x70, 0xbe, 0x93, 0x74, 0x94, 0xa4, 0xa3, 0xb7, 0x9b, 0xa3, 0xb7, 0x3e, 0x16, 0x97, 0x24,
	0xb1, 0x81, 0x05, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xd6, 0x60, 0x38, 0x87, 0xa7, 0x00,
	0x00, 0x00,
}
