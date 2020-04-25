// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/consumermetadatapb/consumer_metadata.proto

package consumermetadatapb

import (
	consumepb "AKFAK/proto/consumepb"
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

type ConsumerGroup struct {
	ID                   int32                           `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Assignments          []*consumepb.MetadataAssignment `protobuf:"bytes,2,rep,name=assignments,proto3" json:"assignments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *ConsumerGroup) Reset()         { *m = ConsumerGroup{} }
func (m *ConsumerGroup) String() string { return proto.CompactTextString(m) }
func (*ConsumerGroup) ProtoMessage()    {}
func (*ConsumerGroup) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad10e0c905c64b8e, []int{0}
}

func (m *ConsumerGroup) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsumerGroup.Unmarshal(m, b)
}
func (m *ConsumerGroup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsumerGroup.Marshal(b, m, deterministic)
}
func (m *ConsumerGroup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsumerGroup.Merge(m, src)
}
func (m *ConsumerGroup) XXX_Size() int {
	return xxx_messageInfo_ConsumerGroup.Size(m)
}
func (m *ConsumerGroup) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsumerGroup.DiscardUnknown(m)
}

var xxx_messageInfo_ConsumerGroup proto.InternalMessageInfo

func (m *ConsumerGroup) GetID() int32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *ConsumerGroup) GetAssignments() []*consumepb.MetadataAssignment {
	if m != nil {
		return m.Assignments
	}
	return nil
}

type MetadataConsumerState struct {
	ConsumerGroups       []*ConsumerGroup `protobuf:"bytes,1,rep,name=consumerGroups,proto3" json:"consumerGroups,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *MetadataConsumerState) Reset()         { *m = MetadataConsumerState{} }
func (m *MetadataConsumerState) String() string { return proto.CompactTextString(m) }
func (*MetadataConsumerState) ProtoMessage()    {}
func (*MetadataConsumerState) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad10e0c905c64b8e, []int{1}
}

func (m *MetadataConsumerState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataConsumerState.Unmarshal(m, b)
}
func (m *MetadataConsumerState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataConsumerState.Marshal(b, m, deterministic)
}
func (m *MetadataConsumerState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataConsumerState.Merge(m, src)
}
func (m *MetadataConsumerState) XXX_Size() int {
	return xxx_messageInfo_MetadataConsumerState.Size(m)
}
func (m *MetadataConsumerState) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataConsumerState.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataConsumerState proto.InternalMessageInfo

func (m *MetadataConsumerState) GetConsumerGroups() []*ConsumerGroup {
	if m != nil {
		return m.ConsumerGroups
	}
	return nil
}

func init() {
	proto.RegisterType((*ConsumerGroup)(nil), "proto.consumermetadatapb.ConsumerGroup")
	proto.RegisterType((*MetadataConsumerState)(nil), "proto.consumermetadatapb.MetadataConsumerState")
}

func init() {
	proto.RegisterFile("proto/consumermetadatapb/consumer_metadata.proto", fileDescriptor_ad10e0c905c64b8e)
}

var fileDescriptor_ad10e0c905c64b8e = []byte{
	// 194 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x28, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0xce, 0xcf, 0x2b, 0x2e, 0xcd, 0x4d, 0x2d, 0xca, 0x4d, 0x2d, 0x49, 0x4c, 0x49,
	0x2c, 0x49, 0x2c, 0x48, 0x82, 0x0b, 0xc5, 0xc3, 0xc4, 0xf4, 0xc0, 0x4a, 0x85, 0x24, 0xc0, 0x94,
	0x1e, 0xa6, 0x0e, 0x29, 0x59, 0x14, 0xb3, 0x10, 0x46, 0x40, 0x34, 0x2a, 0xa5, 0x71, 0xf1, 0x3a,
	0x43, 0x35, 0xb9, 0x17, 0xe5, 0x97, 0x16, 0x08, 0xf1, 0x71, 0x31, 0x79, 0xba, 0x48, 0x30, 0x2a,
	0x30, 0x6a, 0xb0, 0x06, 0x31, 0x79, 0xba, 0x08, 0xb9, 0x72, 0x71, 0x27, 0x16, 0x17, 0x67, 0xa6,
	0xe7, 0xe5, 0xa6, 0xe6, 0x95, 0x14, 0x4b, 0x30, 0x29, 0x30, 0x6b, 0x70, 0x1b, 0x29, 0xeb, 0xa1,
	0xd8, 0x57, 0x90, 0xa4, 0xe7, 0x0b, 0xb5, 0xd1, 0x11, 0xae, 0x36, 0x08, 0x59, 0x9f, 0x52, 0x06,
	0x97, 0x28, 0x4c, 0x09, 0xcc, 0xbe, 0xe0, 0x92, 0xc4, 0x92, 0x54, 0x21, 0x7f, 0x2e, 0xbe, 0x64,
	0x64, 0x07, 0x14, 0x4b, 0x30, 0x82, 0xad, 0x50, 0xd7, 0xc3, 0xe5, 0x25, 0x3d, 0x14, 0x07, 0x07,
	0xa1, 0x69, 0x77, 0x52, 0x88, 0x92, 0x73, 0xf4, 0x76, 0x73, 0xf4, 0xd6, 0xc7, 0x15, 0x88, 0x49,
	0x6c, 0x60, 0x19, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x23, 0xf3, 0x45, 0xc4, 0x67, 0x01,
	0x00, 0x00,
}