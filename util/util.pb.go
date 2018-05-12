// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: util/util.proto

/*
Package util is a generated protocol buffer package.

It is generated from these files:
	util/util.proto

It has these top-level messages:
	StreamStats
	StreamStatsDistribution
	Vector
	VectorDistribution
*/
package util

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// StreamStats is a data structure that can accumulate
// common statical metrics about a stream of numbers
// memory-efficiently, without actually storing the data.
type StreamStats struct {
	// The total weight of the sample.
	Weight float64 `protobuf:"fixed64,1,opt,name=weight,proto3" json:"weight,omitempty"`
	// The weighted sum of all values in the sample.
	Sum float64 `protobuf:"fixed64,2,opt,name=sum,proto3" json:"sum,omitempty"`
	// The weighted sum of of all square values.
	SumSquares float64 `protobuf:"fixed64,3,opt,name=sum_squares,json=sumSquares,proto3" json:"sum_squares,omitempty"`
}

func (m *StreamStats) Reset()                    { *m = StreamStats{} }
func (m *StreamStats) String() string            { return proto.CompactTextString(m) }
func (*StreamStats) ProtoMessage()               {}
func (*StreamStats) Descriptor() ([]byte, []int) { return fileDescriptorUtil, []int{0} }

// StreamStatsDistribution maintains a distribution of
// stream stats by a particular outcome.
type StreamStatsDistribution struct {
	Dense     []StreamStatsDistribution_Dense `protobuf:"bytes,1,rep,name=dense" json:"dense"`
	Sparse    map[int64]*StreamStats          `protobuf:"bytes,2,rep,name=sparse" json:"sparse,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	SparseCap int64                           `protobuf:"varint,3,opt,name=sparse_cap,json=sparseCap,proto3" json:"sparse_cap,omitempty"`
}

func (m *StreamStatsDistribution) Reset()                    { *m = StreamStatsDistribution{} }
func (m *StreamStatsDistribution) String() string            { return proto.CompactTextString(m) }
func (*StreamStatsDistribution) ProtoMessage()               {}
func (*StreamStatsDistribution) Descriptor() ([]byte, []int) { return fileDescriptorUtil, []int{1} }

type StreamStatsDistribution_Dense struct {
	*StreamStats `protobuf:"bytes,1,opt,name=stats,embedded=stats" json:"stats,omitempty"`
}

func (m *StreamStatsDistribution_Dense) Reset()         { *m = StreamStatsDistribution_Dense{} }
func (m *StreamStatsDistribution_Dense) String() string { return proto.CompactTextString(m) }
func (*StreamStatsDistribution_Dense) ProtoMessage()    {}
func (*StreamStatsDistribution_Dense) Descriptor() ([]byte, []int) {
	return fileDescriptorUtil, []int{1, 0}
}

// Vector represents a vector of weights.
// The minimum value of each vector element is 0, which indicates "not set".
type Vector struct {
	Dense     []float64         `protobuf:"fixed64,1,rep,packed,name=dense" json:"dense,omitempty"`
	Sparse    map[int64]float64 `protobuf:"bytes,2,rep,name=sparse" json:"sparse,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
	SparseCap int64             `protobuf:"varint,3,opt,name=sparse_cap,json=sparseCap,proto3" json:"sparse_cap,omitempty"`
}

func (m *Vector) Reset()                    { *m = Vector{} }
func (m *Vector) String() string            { return proto.CompactTextString(m) }
func (*Vector) ProtoMessage()               {}
func (*Vector) Descriptor() ([]byte, []int) { return fileDescriptorUtil, []int{2} }

// Vector maintains a distribution of vectors.
type VectorDistribution struct {
	Dense     []VectorDistribution_Dense `protobuf:"bytes,1,rep,name=dense" json:"dense"`
	Sparse    map[int64]*Vector          `protobuf:"bytes,2,rep,name=sparse" json:"sparse,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	SparseCap int64                      `protobuf:"varint,3,opt,name=sparse_cap,json=sparseCap,proto3" json:"sparse_cap,omitempty"`
}

func (m *VectorDistribution) Reset()                    { *m = VectorDistribution{} }
func (m *VectorDistribution) String() string            { return proto.CompactTextString(m) }
func (*VectorDistribution) ProtoMessage()               {}
func (*VectorDistribution) Descriptor() ([]byte, []int) { return fileDescriptorUtil, []int{3} }

type VectorDistribution_Dense struct {
	*Vector `protobuf:"bytes,1,opt,name=vector,embedded=vector" json:"vector,omitempty"`
}

func (m *VectorDistribution_Dense) Reset()                    { *m = VectorDistribution_Dense{} }
func (m *VectorDistribution_Dense) String() string            { return proto.CompactTextString(m) }
func (*VectorDistribution_Dense) ProtoMessage()               {}
func (*VectorDistribution_Dense) Descriptor() ([]byte, []int) { return fileDescriptorUtil, []int{3, 0} }

func init() {
	proto.RegisterType((*StreamStats)(nil), "blacksquaremedia.reason.util.StreamStats")
	proto.RegisterType((*StreamStatsDistribution)(nil), "blacksquaremedia.reason.util.StreamStatsDistribution")
	proto.RegisterType((*StreamStatsDistribution_Dense)(nil), "blacksquaremedia.reason.util.StreamStatsDistribution.Dense")
	proto.RegisterType((*Vector)(nil), "blacksquaremedia.reason.util.Vector")
	proto.RegisterType((*VectorDistribution)(nil), "blacksquaremedia.reason.util.VectorDistribution")
	proto.RegisterType((*VectorDistribution_Dense)(nil), "blacksquaremedia.reason.util.VectorDistribution.Dense")
}

func init() { proto.RegisterFile("util/util.proto", fileDescriptorUtil) }

var fileDescriptorUtil = []byte{
	// 479 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xdd, 0x6a, 0x13, 0x41,
	0x14, 0xc7, 0x33, 0xd9, 0x64, 0xc1, 0xb3, 0x17, 0x96, 0x41, 0x74, 0x89, 0x9a, 0x94, 0xe0, 0x45,
	0xbd, 0x70, 0x23, 0x15, 0x44, 0x5b, 0x41, 0x5c, 0x5b, 0x10, 0x04, 0x91, 0x89, 0xdf, 0x37, 0x61,
	0x76, 0x33, 0xa6, 0x43, 0xb3, 0x99, 0x38, 0x1f, 0x95, 0x3e, 0x83, 0x37, 0x3e, 0x83, 0x0f, 0x23,
	0xbd, 0xf4, 0xd2, 0xab, 0x42, 0xf1, 0x45, 0x64, 0x66, 0xb6, 0xb0, 0x29, 0xb4, 0xd9, 0xe4, 0x26,
	0x9c, 0x99, 0xec, 0xf9, 0x9d, 0x73, 0xfe, 0xe7, 0xcf, 0xc0, 0x75, 0xa3, 0xf9, 0x74, 0x60, 0x7f,
	0x92, 0xb9, 0x14, 0x5a, 0xe0, 0x3b, 0xd9, 0x94, 0xe6, 0x87, 0xea, 0x9b, 0xa1, 0x92, 0x15, 0x6c,
	0xcc, 0x69, 0x22, 0x19, 0x55, 0x62, 0x96, 0xd8, 0x6f, 0x3a, 0x0f, 0x26, 0x5c, 0x1f, 0x98, 0x2c,
	0xc9, 0x45, 0x31, 0x98, 0x88, 0x89, 0x18, 0xb8, 0xa4, 0xcc, 0x7c, 0x75, 0x27, 0x77, 0x70, 0x91,
	0x87, 0xf5, 0x3f, 0x41, 0x34, 0xd4, 0x92, 0xd1, 0x62, 0xa8, 0xa9, 0x56, 0xf8, 0x26, 0x84, 0xdf,
	0x19, 0x9f, 0x1c, 0xe8, 0x18, 0x6d, 0xa2, 0x2d, 0x44, 0xca, 0x13, 0xde, 0x80, 0x40, 0x99, 0x22,
	0x6e, 0xba, 0x4b, 0x1b, 0xe2, 0x1e, 0x44, 0xca, 0x14, 0x23, 0xdf, 0x86, 0x8a, 0x03, 0xf7, 0x0f,
	0x28, 0x53, 0x0c, 0xfd, 0x4d, 0xff, 0x57, 0x00, 0xb7, 0x2a, 0xe8, 0x3d, 0xae, 0xb4, 0xe4, 0x99,
	0xd1, 0x5c, 0xcc, 0xf0, 0x47, 0x68, 0x8f, 0xd9, 0x4c, 0xb1, 0x18, 0x6d, 0x06, 0x5b, 0xd1, 0xf6,
	0x6e, 0x72, 0xd5, 0x48, 0xc9, 0x25, 0x94, 0x64, 0xcf, 0x22, 0xd2, 0xd6, 0xc9, 0x69, 0xaf, 0x41,
	0x3c, 0x0f, 0x7f, 0x86, 0x50, 0xcd, 0xa9, 0x54, 0x2c, 0x6e, 0x3a, 0xf2, 0x8b, 0xf5, 0xc8, 0x43,
	0xc7, 0xd8, 0x9f, 0x69, 0x79, 0x4c, 0x4a, 0x20, 0xbe, 0x0b, 0xe0, 0xa3, 0x51, 0x4e, 0xe7, 0x6e,
	0xde, 0x80, 0x5c, 0xf3, 0x37, 0x2f, 0xe9, 0xbc, 0xf3, 0x06, 0xda, 0xae, 0x1f, 0xbc, 0x0f, 0x6d,
	0x65, 0x81, 0x4e, 0xc1, 0x68, 0xfb, 0x7e, 0xed, 0x0e, 0xd2, 0xd6, 0x9f, 0xd3, 0x1e, 0x22, 0x3e,
	0xbb, 0x33, 0x86, 0xa8, 0xd2, 0x85, 0x5d, 0xc0, 0x21, 0x3b, 0x76, 0xcc, 0x80, 0xd8, 0x10, 0x3f,
	0x87, 0xf6, 0x11, 0x9d, 0x1a, 0xe6, 0x96, 0xb2, 0x4a, 0x1d, 0xe2, 0xf3, 0x76, 0x9a, 0x4f, 0x50,
	0xff, 0x37, 0x82, 0xf0, 0x03, 0xcb, 0xb5, 0x90, 0x38, 0xae, 0xee, 0x04, 0xa5, 0xcd, 0x0d, 0x74,
	0x2e, 0xea, 0xab, 0x0b, 0xa2, 0x3e, 0xbc, 0xba, 0x94, 0xe7, 0xad, 0xa3, 0xe1, 0xd3, 0x65, 0x33,
	0xdf, 0xa8, 0xce, 0x8c, 0xaa, 0x83, 0xfc, 0x08, 0x00, 0xfb, 0xc2, 0x0b, 0x46, 0x23, 0x8b, 0x46,
	0x7b, 0x5c, 0xa7, 0xf3, 0x65, 0x1e, 0x7b, 0x77, 0x41, 0x8e, 0x67, 0x2b, 0x43, 0xd7, 0x90, 0xe6,
	0xf5, 0xb9, 0xbd, 0x52, 0x08, 0x8f, 0x1c, 0xb1, 0xf4, 0xd7, 0xbd, 0x3a, 0xd5, 0x4b, 0x6b, 0x95,
	0x99, 0x9d, 0xd1, 0x32, 0x9d, 0x77, 0x16, 0xbd, 0x55, 0xab, 0x46, 0x65, 0x1b, 0xe9, 0xee, 0xc9,
	0x59, 0xb7, 0xf1, 0xf7, 0xac, 0x8b, 0x7e, 0xfe, 0xeb, 0x36, 0xe0, 0x76, 0x2e, 0x8a, 0xcb, 0x08,
	0x29, 0xbc, 0xd7, 0x7c, 0xfa, 0xd6, 0xbe, 0x45, 0xea, 0x4b, 0xcb, 0xe2, 0xb2, 0xd0, 0xbd, 0x4c,
	0x8f, 0xfe, 0x07, 0x00, 0x00, 0xff, 0xff, 0x32, 0x81, 0xc8, 0x2d, 0xf9, 0x04, 0x00, 0x00,
}