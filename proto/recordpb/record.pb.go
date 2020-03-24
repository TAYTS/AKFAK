// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/recordpb/record.proto

package recordpb

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

type Record struct {
	Length     int32 `protobuf:"varint,1,opt,name=length,proto3" json:"length,omitempty"`
	Attributes int32 `protobuf:"varint,2,opt,name=attributes,proto3" json:"attributes,omitempty"`
	// attributes
	//bit 0~7: unused
	TimestampDelta       int32    `protobuf:"varint,3,opt,name=timestampDelta,proto3" json:"timestampDelta,omitempty"`
	OffsetDelta          int32    `protobuf:"varint,4,opt,name=offsetDelta,proto3" json:"offsetDelta,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Record) Reset()         { *m = Record{} }
func (m *Record) String() string { return proto.CompactTextString(m) }
func (*Record) ProtoMessage()    {}
func (*Record) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea920a3f23369c2d, []int{0}
}

func (m *Record) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Record.Unmarshal(m, b)
}
func (m *Record) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Record.Marshal(b, m, deterministic)
}
func (m *Record) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Record.Merge(m, src)
}
func (m *Record) XXX_Size() int {
	return xxx_messageInfo_Record.Size(m)
}
func (m *Record) XXX_DiscardUnknown() {
	xxx_messageInfo_Record.DiscardUnknown(m)
}

var xxx_messageInfo_Record proto.InternalMessageInfo

func (m *Record) GetLength() int32 {
	if m != nil {
		return m.Length
	}
	return 0
}

func (m *Record) GetAttributes() int32 {
	if m != nil {
		return m.Attributes
	}
	return 0
}

func (m *Record) GetTimestampDelta() int32 {
	if m != nil {
		return m.TimestampDelta
	}
	return 0
}

func (m *Record) GetOffsetDelta() int32 {
	if m != nil {
		return m.OffsetDelta
	}
	return 0
}

type RecordBatch struct {
	BaseOffset           int64  `protobuf:"varint,1,opt,name=baseOffset,proto3" json:"baseOffset,omitempty"`
	BatchLength          []byte `protobuf:"bytes,2,opt,name=batchLength,proto3" json:"batchLength,omitempty"`
	PartitionLeaderEpoch int32  `protobuf:"varint,3,opt,name=partitionLeaderEpoch,proto3" json:"partitionLeaderEpoch,omitempty"`
	Magic                int32  `protobuf:"varint,4,opt,name=magic,proto3" json:"magic,omitempty"`
	Crc                  int32  `protobuf:"varint,5,opt,name=crc,proto3" json:"crc,omitempty"`
	Attributes           int32  `protobuf:"varint,6,opt,name=attributes,proto3" json:"attributes,omitempty"`
	// attributes
	//bit 0~2:
	//0: no compression
	//1: gzip
	//2: snappy
	//3: lz4
	//4: zstd
	//bit 3: timestampType
	//bit 4: isTransactional (0 means not transactional)
	//bit 5: isControlBatch (0 means not a control batch)
	//bit 6~15: unused
	LastOffsetDelta      int32     `protobuf:"varint,7,opt,name=lastOffsetDelta,proto3" json:"lastOffsetDelta,omitempty"`
	FirstTimestamp       int64     `protobuf:"varint,8,opt,name=firstTimestamp,proto3" json:"firstTimestamp,omitempty"`
	MaxTimestamp         int64     `protobuf:"varint,9,opt,name=maxTimestamp,proto3" json:"maxTimestamp,omitempty"`
	ProducerId           int64     `protobuf:"varint,10,opt,name=producerId,proto3" json:"producerId,omitempty"`
	ProducerEpoch        int32     `protobuf:"varint,11,opt,name=producerEpoch,proto3" json:"producerEpoch,omitempty"`
	BaseSequence         int32     `protobuf:"varint,12,opt,name=baseSequence,proto3" json:"baseSequence,omitempty"`
	Records              []*Record `protobuf:"bytes,13,rep,name=records,proto3" json:"records,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *RecordBatch) Reset()         { *m = RecordBatch{} }
func (m *RecordBatch) String() string { return proto.CompactTextString(m) }
func (*RecordBatch) ProtoMessage()    {}
func (*RecordBatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_ea920a3f23369c2d, []int{1}
}

func (m *RecordBatch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RecordBatch.Unmarshal(m, b)
}
func (m *RecordBatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RecordBatch.Marshal(b, m, deterministic)
}
func (m *RecordBatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecordBatch.Merge(m, src)
}
func (m *RecordBatch) XXX_Size() int {
	return xxx_messageInfo_RecordBatch.Size(m)
}
func (m *RecordBatch) XXX_DiscardUnknown() {
	xxx_messageInfo_RecordBatch.DiscardUnknown(m)
}

var xxx_messageInfo_RecordBatch proto.InternalMessageInfo

func (m *RecordBatch) GetBaseOffset() int64 {
	if m != nil {
		return m.BaseOffset
	}
	return 0
}

func (m *RecordBatch) GetBatchLength() []byte {
	if m != nil {
		return m.BatchLength
	}
	return nil
}

func (m *RecordBatch) GetPartitionLeaderEpoch() int32 {
	if m != nil {
		return m.PartitionLeaderEpoch
	}
	return 0
}

func (m *RecordBatch) GetMagic() int32 {
	if m != nil {
		return m.Magic
	}
	return 0
}

func (m *RecordBatch) GetCrc() int32 {
	if m != nil {
		return m.Crc
	}
	return 0
}

func (m *RecordBatch) GetAttributes() int32 {
	if m != nil {
		return m.Attributes
	}
	return 0
}

func (m *RecordBatch) GetLastOffsetDelta() int32 {
	if m != nil {
		return m.LastOffsetDelta
	}
	return 0
}

func (m *RecordBatch) GetFirstTimestamp() int64 {
	if m != nil {
		return m.FirstTimestamp
	}
	return 0
}

func (m *RecordBatch) GetMaxTimestamp() int64 {
	if m != nil {
		return m.MaxTimestamp
	}
	return 0
}

func (m *RecordBatch) GetProducerId() int64 {
	if m != nil {
		return m.ProducerId
	}
	return 0
}

func (m *RecordBatch) GetProducerEpoch() int32 {
	if m != nil {
		return m.ProducerEpoch
	}
	return 0
}

func (m *RecordBatch) GetBaseSequence() int32 {
	if m != nil {
		return m.BaseSequence
	}
	return 0
}

func (m *RecordBatch) GetRecords() []*Record {
	if m != nil {
		return m.Records
	}
	return nil
}

func init() {
	proto.RegisterType((*Record)(nil), "proto.recordpb.Record")
	proto.RegisterType((*RecordBatch)(nil), "proto.recordpb.RecordBatch")
}

func init() {
	proto.RegisterFile("proto/recordpb/record.proto", fileDescriptor_ea920a3f23369c2d)
}

var fileDescriptor_ea920a3f23369c2d = []byte{
	// 354 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xd1, 0x6a, 0xc2, 0x30,
	0x14, 0x86, 0x71, 0x9d, 0x75, 0x3b, 0x55, 0x37, 0x82, 0x48, 0x60, 0x30, 0x44, 0xc6, 0xf0, 0xaa,
	0x0e, 0xf7, 0x04, 0xca, 0x36, 0x18, 0x0a, 0x42, 0xb7, 0xab, 0xdd, 0xa5, 0x69, 0xd4, 0x82, 0xb5,
	0x5d, 0x7a, 0x84, 0x3d, 0xc3, 0xde, 0x71, 0xef, 0x32, 0x72, 0xd2, 0x62, 0xdb, 0xed, 0xca, 0xfc,
	0xdf, 0xf9, 0x31, 0x7f, 0xfe, 0x1e, 0xb8, 0xc9, 0x74, 0x8a, 0xe9, 0x54, 0x2b, 0x99, 0xea, 0x28,
	0x0b, 0x8b, 0x83, 0x4f, 0x94, 0xf5, 0xe9, 0xc7, 0x2f, 0x87, 0xe3, 0xef, 0x16, 0xb8, 0x01, 0x09,
	0x36, 0x04, 0x77, 0xaf, 0x0e, 0x5b, 0xdc, 0xf1, 0xd6, 0xa8, 0x35, 0x69, 0x07, 0x85, 0x62, 0xb7,
	0x00, 0x02, 0x51, 0xc7, 0xe1, 0x11, 0x55, 0xce, 0xcf, 0x68, 0x56, 0x21, 0xec, 0x1e, 0xfa, 0x18,
	0x27, 0x2a, 0x47, 0x91, 0x64, 0x4f, 0x6a, 0x8f, 0x82, 0x3b, 0xe4, 0x69, 0x50, 0x36, 0x02, 0x2f,
	0xdd, 0x6c, 0x72, 0x85, 0xd6, 0x74, 0x4e, 0xa6, 0x2a, 0x1a, 0xff, 0x38, 0xe0, 0xd9, 0x30, 0x0b,
	0x81, 0x92, 0x6e, 0x0e, 0x45, 0xae, 0xd6, 0x64, 0xa1, 0x54, 0x4e, 0x50, 0x21, 0xe6, 0x1f, 0x43,
	0x63, 0x5c, 0xd9, 0xd8, 0x26, 0x5a, 0x37, 0xa8, 0x22, 0x36, 0x83, 0x41, 0x26, 0x34, 0xc6, 0x18,
	0xa7, 0x87, 0x95, 0x12, 0x91, 0xd2, 0xcf, 0x59, 0x2a, 0x77, 0x45, 0xc2, 0x7f, 0x67, 0x6c, 0x00,
	0xed, 0x44, 0x6c, 0x63, 0x59, 0x24, 0xb4, 0x82, 0x5d, 0x83, 0x23, 0xb5, 0xe4, 0x6d, 0x62, 0xe6,
	0xd8, 0xe8, 0xc5, 0xfd, 0xd3, 0xcb, 0x04, 0xae, 0xf6, 0x22, 0xc7, 0x75, 0xe5, 0xcd, 0x1d, 0x32,
	0x35, 0xb1, 0x69, 0x70, 0x13, 0xeb, 0x1c, 0xdf, 0xcb, 0xc2, 0xf8, 0x05, 0xbd, 0xb5, 0x41, 0xd9,
	0x18, 0xba, 0x89, 0xf8, 0x3a, 0xb9, 0x2e, 0xc9, 0x55, 0x63, 0x26, 0x55, 0xa6, 0xd3, 0xe8, 0x28,
	0x95, 0x7e, 0x8d, 0x38, 0xd8, 0xce, 0x4e, 0x84, 0xdd, 0x41, 0xaf, 0x54, 0xb6, 0x0a, 0x8f, 0x32,
	0xd5, 0xa1, 0xb9, 0xc9, 0xf4, 0xfc, 0xa6, 0x3e, 0x8f, 0xea, 0x20, 0x15, 0xef, 0x92, 0xa9, 0xc6,
	0xd8, 0x03, 0x74, 0xec, 0x1a, 0xe5, 0xbc, 0x37, 0x72, 0x26, 0xde, 0x6c, 0xe8, 0xd7, 0x97, 0xcb,
	0xb7, 0xdf, 0x32, 0x28, 0x6d, 0x8b, 0xe1, 0xc7, 0x60, 0xbe, 0x7c, 0x99, 0x2f, 0xa7, 0xf5, 0x0d,
	0x0d, 0x5d, 0xd2, 0x8f, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x55, 0x13, 0x10, 0xd3, 0xba, 0x02,
	0x00, 0x00,
}
