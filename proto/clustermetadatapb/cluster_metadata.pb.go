// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/clustermetadatapb/cluster_metadata.proto

package clustermetadatapb

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

type MetadataPartitionState struct {
	TopicName            string   `protobuf:"bytes,1,opt,name=topicName,proto3" json:"topicName,omitempty"`
	PartitionIndex       int32    `protobuf:"varint,2,opt,name=partitionIndex,proto3" json:"partitionIndex,omitempty"`
	Leader               int32    `protobuf:"varint,3,opt,name=leader,proto3" json:"leader,omitempty"`
	Isr                  []int32  `protobuf:"varint,4,rep,packed,name=isr,proto3" json:"isr,omitempty"`
	Replicas             []int32  `protobuf:"varint,5,rep,packed,name=replicas,proto3" json:"replicas,omitempty"`
	OfflineReplicas      []int32  `protobuf:"varint,6,rep,packed,name=offlineReplicas,proto3" json:"offlineReplicas,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MetadataPartitionState) Reset()         { *m = MetadataPartitionState{} }
func (m *MetadataPartitionState) String() string { return proto.CompactTextString(m) }
func (*MetadataPartitionState) ProtoMessage()    {}
func (*MetadataPartitionState) Descriptor() ([]byte, []int) {
	return fileDescriptor_beeae55be1f11968, []int{0}
}

func (m *MetadataPartitionState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataPartitionState.Unmarshal(m, b)
}
func (m *MetadataPartitionState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataPartitionState.Marshal(b, m, deterministic)
}
func (m *MetadataPartitionState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataPartitionState.Merge(m, src)
}
func (m *MetadataPartitionState) XXX_Size() int {
	return xxx_messageInfo_MetadataPartitionState.Size(m)
}
func (m *MetadataPartitionState) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataPartitionState.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataPartitionState proto.InternalMessageInfo

func (m *MetadataPartitionState) GetTopicName() string {
	if m != nil {
		return m.TopicName
	}
	return ""
}

func (m *MetadataPartitionState) GetPartitionIndex() int32 {
	if m != nil {
		return m.PartitionIndex
	}
	return 0
}

func (m *MetadataPartitionState) GetLeader() int32 {
	if m != nil {
		return m.Leader
	}
	return 0
}

func (m *MetadataPartitionState) GetIsr() []int32 {
	if m != nil {
		return m.Isr
	}
	return nil
}

func (m *MetadataPartitionState) GetReplicas() []int32 {
	if m != nil {
		return m.Replicas
	}
	return nil
}

func (m *MetadataPartitionState) GetOfflineReplicas() []int32 {
	if m != nil {
		return m.OfflineReplicas
	}
	return nil
}

type MetadataTopicState struct {
	TopicName            string                    `protobuf:"bytes,1,opt,name=topicName,proto3" json:"topicName,omitempty"`
	PartitionStates      []*MetadataPartitionState `protobuf:"bytes,2,rep,name=partitionStates,proto3" json:"partitionStates,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *MetadataTopicState) Reset()         { *m = MetadataTopicState{} }
func (m *MetadataTopicState) String() string { return proto.CompactTextString(m) }
func (*MetadataTopicState) ProtoMessage()    {}
func (*MetadataTopicState) Descriptor() ([]byte, []int) {
	return fileDescriptor_beeae55be1f11968, []int{1}
}

func (m *MetadataTopicState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataTopicState.Unmarshal(m, b)
}
func (m *MetadataTopicState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataTopicState.Marshal(b, m, deterministic)
}
func (m *MetadataTopicState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataTopicState.Merge(m, src)
}
func (m *MetadataTopicState) XXX_Size() int {
	return xxx_messageInfo_MetadataTopicState.Size(m)
}
func (m *MetadataTopicState) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataTopicState.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataTopicState proto.InternalMessageInfo

func (m *MetadataTopicState) GetTopicName() string {
	if m != nil {
		return m.TopicName
	}
	return ""
}

func (m *MetadataTopicState) GetPartitionStates() []*MetadataPartitionState {
	if m != nil {
		return m.PartitionStates
	}
	return nil
}

type MetadataBroker struct {
	ID                   int32    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Host                 string   `protobuf:"bytes,2,opt,name=host,proto3" json:"host,omitempty"`
	Port                 int32    `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MetadataBroker) Reset()         { *m = MetadataBroker{} }
func (m *MetadataBroker) String() string { return proto.CompactTextString(m) }
func (*MetadataBroker) ProtoMessage()    {}
func (*MetadataBroker) Descriptor() ([]byte, []int) {
	return fileDescriptor_beeae55be1f11968, []int{2}
}

func (m *MetadataBroker) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataBroker.Unmarshal(m, b)
}
func (m *MetadataBroker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataBroker.Marshal(b, m, deterministic)
}
func (m *MetadataBroker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataBroker.Merge(m, src)
}
func (m *MetadataBroker) XXX_Size() int {
	return xxx_messageInfo_MetadataBroker.Size(m)
}
func (m *MetadataBroker) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataBroker.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataBroker proto.InternalMessageInfo

func (m *MetadataBroker) GetID() int32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *MetadataBroker) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *MetadataBroker) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type MetadataCluster struct {
	LiveBrokers          []*MetadataBroker     `protobuf:"bytes,1,rep,name=liveBrokers,proto3" json:"liveBrokers,omitempty"`
	TopicStates          []*MetadataTopicState `protobuf:"bytes,2,rep,name=topicStates,proto3" json:"topicStates,omitempty"`
	Brokers              []*MetadataBroker     `protobuf:"bytes,3,rep,name=brokers,proto3" json:"brokers,omitempty"`
	Controller           *MetadataBroker       `protobuf:"bytes,4,opt,name=controller,proto3" json:"controller,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *MetadataCluster) Reset()         { *m = MetadataCluster{} }
func (m *MetadataCluster) String() string { return proto.CompactTextString(m) }
func (*MetadataCluster) ProtoMessage()    {}
func (*MetadataCluster) Descriptor() ([]byte, []int) {
	return fileDescriptor_beeae55be1f11968, []int{3}
}

func (m *MetadataCluster) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetadataCluster.Unmarshal(m, b)
}
func (m *MetadataCluster) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetadataCluster.Marshal(b, m, deterministic)
}
func (m *MetadataCluster) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetadataCluster.Merge(m, src)
}
func (m *MetadataCluster) XXX_Size() int {
	return xxx_messageInfo_MetadataCluster.Size(m)
}
func (m *MetadataCluster) XXX_DiscardUnknown() {
	xxx_messageInfo_MetadataCluster.DiscardUnknown(m)
}

var xxx_messageInfo_MetadataCluster proto.InternalMessageInfo

func (m *MetadataCluster) GetLiveBrokers() []*MetadataBroker {
	if m != nil {
		return m.LiveBrokers
	}
	return nil
}

func (m *MetadataCluster) GetTopicStates() []*MetadataTopicState {
	if m != nil {
		return m.TopicStates
	}
	return nil
}

func (m *MetadataCluster) GetBrokers() []*MetadataBroker {
	if m != nil {
		return m.Brokers
	}
	return nil
}

func (m *MetadataCluster) GetController() *MetadataBroker {
	if m != nil {
		return m.Controller
	}
	return nil
}

func init() {
	proto.RegisterType((*MetadataPartitionState)(nil), "proto.clustermetadatapb.MetadataPartitionState")
	proto.RegisterType((*MetadataTopicState)(nil), "proto.clustermetadatapb.MetadataTopicState")
	proto.RegisterType((*MetadataBroker)(nil), "proto.clustermetadatapb.MetadataBroker")
	proto.RegisterType((*MetadataCluster)(nil), "proto.clustermetadatapb.MetadataCluster")
}

func init() {
	proto.RegisterFile("proto/clustermetadatapb/cluster_metadata.proto", fileDescriptor_beeae55be1f11968)
}

var fileDescriptor_beeae55be1f11968 = []byte{
	// 371 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0xd1, 0x4a, 0xeb, 0x40,
	0x10, 0x25, 0x49, 0xd3, 0x7b, 0x3b, 0x85, 0x56, 0xe6, 0xa1, 0x06, 0x51, 0x0c, 0x79, 0xd0, 0x80,
	0x90, 0x82, 0x7e, 0x41, 0x6b, 0x51, 0x43, 0xa9, 0xc8, 0xea, 0x8b, 0xbe, 0xc8, 0x36, 0xdd, 0x62,
	0x30, 0xcd, 0x86, 0xcd, 0x2a, 0xfe, 0x84, 0xdf, 0xe3, 0x5f, 0xf8, 0x4d, 0x92, 0x4d, 0xb6, 0xad,
	0xd5, 0x52, 0xfb, 0x94, 0x99, 0x33, 0x67, 0x0e, 0x67, 0x0e, 0x59, 0x08, 0x32, 0xc1, 0x25, 0xef,
	0x46, 0xc9, 0x4b, 0x2e, 0x99, 0x98, 0x31, 0x49, 0x27, 0x54, 0xd2, 0x6c, 0xac, 0x91, 0x47, 0x0d,
	0x95, 0x44, 0xdc, 0x55, 0x9f, 0xe0, 0x07, 0xdf, 0xfb, 0x34, 0xa0, 0x33, 0xaa, 0xda, 0x1b, 0x2a,
	0x64, 0x2c, 0x63, 0x9e, 0xde, 0x4a, 0x2a, 0x19, 0xee, 0x43, 0x43, 0xf2, 0x2c, 0x8e, 0xae, 0xe9,
	0x8c, 0x39, 0x86, 0x6b, 0xf8, 0x0d, 0xb2, 0x00, 0xf0, 0x08, 0x5a, 0x99, 0xe6, 0x87, 0xe9, 0x84,
	0xbd, 0x39, 0xa6, 0x6b, 0xf8, 0x36, 0x59, 0x41, 0xb1, 0x03, 0xf5, 0x84, 0xd1, 0x09, 0x13, 0x8e,
	0xa5, 0xe6, 0x55, 0x87, 0x3b, 0x60, 0xc5, 0xb9, 0x70, 0x6a, 0xae, 0xe5, 0xdb, 0xa4, 0x28, 0x71,
	0x0f, 0xfe, 0x0b, 0x96, 0x25, 0x71, 0x44, 0x73, 0xc7, 0x56, 0xf0, 0xbc, 0x47, 0x1f, 0xda, 0x7c,
	0x3a, 0x4d, 0xe2, 0x94, 0x11, 0x4d, 0xa9, 0x2b, 0xca, 0x2a, 0xec, 0xbd, 0x1b, 0x80, 0xfa, 0xa0,
	0xbb, 0xc2, 0xed, 0x5f, 0x8e, 0xb9, 0x87, 0x76, 0xf6, 0xed, 0xf8, 0xdc, 0x31, 0x5d, 0xcb, 0x6f,
	0x9e, 0x76, 0x83, 0x35, 0xc1, 0x05, 0xbf, 0x87, 0x46, 0x56, 0x75, 0xbc, 0x2b, 0x68, 0x69, 0x6a,
	0x5f, 0xf0, 0x67, 0x26, 0xb0, 0x05, 0x66, 0x38, 0x50, 0x1e, 0x6c, 0x62, 0x86, 0x03, 0x44, 0xa8,
	0x3d, 0xf1, 0x5c, 0xaa, 0xfc, 0x1a, 0x44, 0xd5, 0x05, 0x96, 0x71, 0x21, 0xab, 0xcc, 0x54, 0xed,
	0x7d, 0x98, 0xd0, 0xd6, 0x52, 0xe7, 0xa5, 0x1f, 0x0c, 0xa1, 0x99, 0xc4, 0xaf, 0xac, 0x54, 0xce,
	0x1d, 0x43, 0x99, 0x3e, 0xde, 0x68, 0xba, 0xe4, 0x93, 0xe5, 0x5d, 0x1c, 0x41, 0x53, 0xce, 0xf3,
	0xd2, 0xf7, 0x9f, 0x6c, 0x94, 0x5a, 0x64, 0x4c, 0x96, 0xf7, 0xb1, 0x07, 0xff, 0xc6, 0x95, 0x2b,
	0x6b, 0x3b, 0x57, 0x7a, 0x0f, 0x2f, 0x01, 0x22, 0x9e, 0x4a, 0xc1, 0x93, 0x84, 0x15, 0x7f, 0x8a,
	0xb1, 0x8d, 0xca, 0xd2, 0x6a, 0xff, 0xf0, 0xe1, 0xa0, 0x37, 0xbc, 0xe8, 0x0d, 0xbb, 0x6b, 0x5e,
	0xcd, 0xb8, 0xae, 0x06, 0x67, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc4, 0xa7, 0x3e, 0x07, 0x57,
	0x03, 0x00, 0x00,
}
