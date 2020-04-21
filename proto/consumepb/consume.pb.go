// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/consumepb/consume.proto

package consumepb

import (
	clustermetadatapb "AKFAK/proto/clustermetadatapb"
	recordpb "AKFAK/proto/recordpb"
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

type ConsumeRequest struct {
	GroupID              int32    `protobuf:"varint,1,opt,name=groupID,proto3" json:"groupID,omitempty"`
	ConsumerID           int32    `protobuf:"varint,2,opt,name=consumerID,proto3" json:"consumerID,omitempty"`
	Partition            int32    `protobuf:"varint,3,opt,name=partition,proto3" json:"partition,omitempty"`
	TopicName            string   `protobuf:"bytes,4,opt,name=topicName,proto3" json:"topicName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsumeRequest) Reset()         { *m = ConsumeRequest{} }
func (m *ConsumeRequest) String() string { return proto.CompactTextString(m) }
func (*ConsumeRequest) ProtoMessage()    {}
func (*ConsumeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0a0fb7c1a5499aef, []int{0}
}

func (m *ConsumeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsumeRequest.Unmarshal(m, b)
}
func (m *ConsumeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsumeRequest.Marshal(b, m, deterministic)
}
func (m *ConsumeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsumeRequest.Merge(m, src)
}
func (m *ConsumeRequest) XXX_Size() int {
	return xxx_messageInfo_ConsumeRequest.Size(m)
}
func (m *ConsumeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsumeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConsumeRequest proto.InternalMessageInfo

func (m *ConsumeRequest) GetGroupID() int32 {
	if m != nil {
		return m.GroupID
	}
	return 0
}

func (m *ConsumeRequest) GetConsumerID() int32 {
	if m != nil {
		return m.ConsumerID
	}
	return 0
}

func (m *ConsumeRequest) GetPartition() int32 {
	if m != nil {
		return m.Partition
	}
	return 0
}

func (m *ConsumeRequest) GetTopicName() string {
	if m != nil {
		return m.TopicName
	}
	return ""
}

type ConsumeResponse struct {
	TopicName            string                `protobuf:"bytes,1,opt,name=topicName,proto3" json:"topicName,omitempty"`
	Partition            int32                 `protobuf:"varint,2,opt,name=partition,proto3" json:"partition,omitempty"`
	RecordSet            *recordpb.RecordBatch `protobuf:"bytes,3,opt,name=recordSet,proto3" json:"recordSet,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ConsumeResponse) Reset()         { *m = ConsumeResponse{} }
func (m *ConsumeResponse) String() string { return proto.CompactTextString(m) }
func (*ConsumeResponse) ProtoMessage()    {}
func (*ConsumeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0a0fb7c1a5499aef, []int{1}
}

func (m *ConsumeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsumeResponse.Unmarshal(m, b)
}
func (m *ConsumeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsumeResponse.Marshal(b, m, deterministic)
}
func (m *ConsumeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsumeResponse.Merge(m, src)
}
func (m *ConsumeResponse) XXX_Size() int {
	return xxx_messageInfo_ConsumeResponse.Size(m)
}
func (m *ConsumeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsumeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ConsumeResponse proto.InternalMessageInfo

func (m *ConsumeResponse) GetTopicName() string {
	if m != nil {
		return m.TopicName
	}
	return ""
}

func (m *ConsumeResponse) GetPartition() int32 {
	if m != nil {
		return m.Partition
	}
	return 0
}

func (m *ConsumeResponse) GetRecordSet() *recordpb.RecordBatch {
	if m != nil {
		return m.RecordSet
	}
	return nil
}

type GetAssignmentRequest struct {
	GroupID              int32    `protobuf:"varint,1,opt,name=groupID,proto3" json:"groupID,omitempty"`
	TopicName            string   `protobuf:"bytes,2,opt,name=topicName,proto3" json:"topicName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAssignmentRequest) Reset()         { *m = GetAssignmentRequest{} }
func (m *GetAssignmentRequest) String() string { return proto.CompactTextString(m) }
func (*GetAssignmentRequest) ProtoMessage()    {}
func (*GetAssignmentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0a0fb7c1a5499aef, []int{2}
}

func (m *GetAssignmentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAssignmentRequest.Unmarshal(m, b)
}
func (m *GetAssignmentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAssignmentRequest.Marshal(b, m, deterministic)
}
func (m *GetAssignmentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAssignmentRequest.Merge(m, src)
}
func (m *GetAssignmentRequest) XXX_Size() int {
	return xxx_messageInfo_GetAssignmentRequest.Size(m)
}
func (m *GetAssignmentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAssignmentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAssignmentRequest proto.InternalMessageInfo

func (m *GetAssignmentRequest) GetGroupID() int32 {
	if m != nil {
		return m.GroupID
	}
	return 0
}

func (m *GetAssignmentRequest) GetTopicName() string {
	if m != nil {
		return m.TopicName
	}
	return ""
}

type MetadataAssignment struct {
	TopicName      string `protobuf:"bytes,1,opt,name=topicName,proto3" json:"topicName,omitempty"`
	PartitionIndex int32  `protobuf:"varint,2,opt,name=partitionIndex,proto3" json:"partitionIndex,omitempty"`
	Offset         int32  `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
	// brokerid for broker holding replica that is assigned to consumer
	// for that partition
	Broker int32 `protobuf:"varint,4,opt,name=broker,proto3" json:"broker,omitempty"`
	// list of brokers holding isr for partition (excludes broker with assigned replica)
	// contacted when an assigned replica fails to respond
	IsrBrokers           []*clustermetadatapb.MetadataBroker `protobuf:"bytes,5,rep,name=isrBrokers,proto3" json:"isrBrokers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                            `json:"-"`
	XXX_unrecognized     []byte                              `json:"-"`
	XXX_sizecache        int32                               `json:"-"`
}

func (m *MetadataAssignment) Reset()         { *m = MetadataAssignment{} }
func (m *MetadataAssignment) String() string { return proto.CompactTextString(m) }
func (*MetadataAssignment) ProtoMessage()    {}
func (*MetadataAssignment) Descriptor() ([]byte, []int) {
	return fileDescriptor_0a0fb7c1a5499aef, []int{3}
}

func (m *MetadataAssignment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataAssignment.Unmarshal(m, b)
}
func (m *MetadataAssignment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataAssignment.Marshal(b, m, deterministic)
}
func (m *MetadataAssignment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataAssignment.Merge(m, src)
}
func (m *MetadataAssignment) XXX_Size() int {
	return xxx_messageInfo_MetadataAssignment.Size(m)
}
func (m *MetadataAssignment) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataAssignment.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataAssignment proto.InternalMessageInfo

func (m *MetadataAssignment) GetTopicName() string {
	if m != nil {
		return m.TopicName
	}
	return ""
}

func (m *MetadataAssignment) GetPartitionIndex() int32 {
	if m != nil {
		return m.PartitionIndex
	}
	return 0
}

func (m *MetadataAssignment) GetOffset() int32 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *MetadataAssignment) GetBroker() int32 {
	if m != nil {
		return m.Broker
	}
	return 0
}

func (m *MetadataAssignment) GetIsrBrokers() []*clustermetadatapb.MetadataBroker {
	if m != nil {
		return m.IsrBrokers
	}
	return nil
}

type GetAssignmentResponse struct {
	// assigned Assignment. len(Assignment) == len(partitions)
	Assignments          []*MetadataAssignment `protobuf:"bytes,1,rep,name=assignments,proto3" json:"assignments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *GetAssignmentResponse) Reset()         { *m = GetAssignmentResponse{} }
func (m *GetAssignmentResponse) String() string { return proto.CompactTextString(m) }
func (*GetAssignmentResponse) ProtoMessage()    {}
func (*GetAssignmentResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0a0fb7c1a5499aef, []int{4}
}

func (m *GetAssignmentResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAssignmentResponse.Unmarshal(m, b)
}
func (m *GetAssignmentResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAssignmentResponse.Marshal(b, m, deterministic)
}
func (m *GetAssignmentResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAssignmentResponse.Merge(m, src)
}
func (m *GetAssignmentResponse) XXX_Size() int {
	return xxx_messageInfo_GetAssignmentResponse.Size(m)
}
func (m *GetAssignmentResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAssignmentResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetAssignmentResponse proto.InternalMessageInfo

func (m *GetAssignmentResponse) GetAssignments() []*MetadataAssignment {
	if m != nil {
		return m.Assignments
	}
	return nil
}

func init() {
	proto.RegisterType((*ConsumeRequest)(nil), "proto.consumepb.ConsumeRequest")
	proto.RegisterType((*ConsumeResponse)(nil), "proto.consumepb.ConsumeResponse")
	proto.RegisterType((*GetAssignmentRequest)(nil), "proto.consumepb.GetAssignmentRequest")
	proto.RegisterType((*MetadataAssignment)(nil), "proto.consumepb.MetadataAssignment")
	proto.RegisterType((*GetAssignmentResponse)(nil), "proto.consumepb.GetAssignmentResponse")
}

func init() {
	proto.RegisterFile("proto/consumepb/consume.proto", fileDescriptor_0a0fb7c1a5499aef)
}

var fileDescriptor_0a0fb7c1a5499aef = []byte{
	// 379 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0xcf, 0x4f, 0xc2, 0x30,
	0x14, 0x4e, 0xc1, 0x61, 0x78, 0x24, 0x90, 0x2c, 0xa2, 0x0b, 0xa8, 0x21, 0x33, 0x51, 0x4e, 0x25,
	0xc1, 0x93, 0x47, 0x10, 0x25, 0x84, 0xc8, 0xa1, 0xde, 0x3c, 0x68, 0xb6, 0x51, 0x70, 0xd1, 0xad,
	0xb3, 0xed, 0x12, 0xff, 0x02, 0xe3, 0xff, 0xe7, 0x3f, 0x64, 0x68, 0xbb, 0x0d, 0xc6, 0x81, 0x13,
	0xef, 0x7d, 0xef, 0xd7, 0xf7, 0x7d, 0x74, 0x70, 0x91, 0x70, 0x26, 0xd9, 0x20, 0x60, 0xb1, 0x48,
	0x23, 0x9a, 0xf8, 0x59, 0x84, 0x15, 0x6e, 0xb7, 0xd4, 0x0f, 0xce, 0xcb, 0x9d, 0xae, 0xee, 0xe7,
	0x34, 0x60, 0x7c, 0x99, 0xf8, 0x26, 0xd0, 0xdd, 0x1d, 0x6c, 0x96, 0x7d, 0xa6, 0x42, 0x52, 0x1e,
	0x51, 0xe9, 0x2d, 0x3d, 0xe9, 0x6d, 0x96, 0x6a, 0xe4, 0x2d, 0x83, 0x74, 0xa3, 0xfb, 0x83, 0xa0,
	0x79, 0xaf, 0x57, 0x13, 0xfa, 0x95, 0x52, 0x21, 0x6d, 0x07, 0x8e, 0xd7, 0x9c, 0xa5, 0xc9, 0x6c,
	0xe2, 0xa0, 0x1e, 0xea, 0x5b, 0x24, 0x4b, 0xed, 0x4b, 0x00, 0x43, 0x83, 0xcf, 0x26, 0x4e, 0x45,
	0x15, 0xb7, 0x10, 0xfb, 0x1c, 0xea, 0x89, 0xc7, 0x65, 0x28, 0x43, 0x16, 0x3b, 0x55, 0x55, 0x2e,
	0x80, 0x4d, 0x55, 0xb2, 0x24, 0x0c, 0x16, 0x5e, 0x44, 0x9d, 0xa3, 0x1e, 0xea, 0xd7, 0x49, 0x01,
	0xb8, 0xbf, 0x08, 0x5a, 0x39, 0x11, 0x91, 0xb0, 0x58, 0xd0, 0xdd, 0x09, 0x54, 0x9a, 0xd8, 0xbd,
	0x56, 0x29, 0x5f, 0xbb, 0x83, 0xba, 0x36, 0xe6, 0x99, 0x4a, 0xc5, 0xa5, 0x31, 0xec, 0x6a, 0xcd,
	0x38, 0x73, 0x0e, 0x13, 0x15, 0x8c, 0x3d, 0x19, 0xbc, 0x93, 0xa2, 0xdb, 0x5d, 0xc0, 0xc9, 0x94,
	0xca, 0x91, 0x10, 0xe1, 0x3a, 0x8e, 0x68, 0x2c, 0x0f, 0x1b, 0xb3, 0x43, 0xb4, 0x52, 0x96, 0xf6,
	0x87, 0xc0, 0x7e, 0x32, 0xb6, 0x17, 0x5b, 0x0f, 0xa8, 0xbb, 0x86, 0x66, 0x2e, 0x66, 0x16, 0x2f,
	0xe9, 0xb7, 0x91, 0x58, 0x42, 0xed, 0x53, 0xa8, 0xb1, 0xd5, 0x4a, 0x18, 0x91, 0x16, 0x31, 0xd9,
	0x06, 0xf7, 0x39, 0xfb, 0xa0, 0x5c, 0x59, 0x6d, 0x11, 0x93, 0xd9, 0x53, 0x80, 0x50, 0xf0, 0xb1,
	0x4a, 0x84, 0x63, 0xf5, 0xaa, 0xfd, 0xc6, 0xf0, 0xc6, 0x18, 0xb3, 0xf7, 0x6a, 0x70, 0x46, 0x5b,
	0xf7, 0x93, 0xad, 0x51, 0xf7, 0x15, 0xda, 0x25, 0x97, 0xcc, 0xbf, 0xf6, 0x00, 0x0d, 0x2f, 0x47,
	0x85, 0x83, 0xd4, 0x89, 0x2b, 0x5c, 0x7a, 0xc6, 0x78, 0xdf, 0x11, 0xb2, 0x3d, 0x37, 0x3e, 0x7b,
	0x69, 0x8f, 0xe6, 0x8f, 0xa3, 0xf9, 0xa0, 0xf4, 0x79, 0xf8, 0x35, 0x05, 0xdc, 0xfe, 0x07, 0x00,
	0x00, 0xff, 0xff, 0xcf, 0xb1, 0xbb, 0xf6, 0x38, 0x03, 0x00, 0x00,
}
