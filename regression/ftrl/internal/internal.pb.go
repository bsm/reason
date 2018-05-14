// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: regression/ftrl/internal/internal.proto

/*
Package internal is a generated protocol buffer package.

It is generated from these files:
	regression/ftrl/internal/internal.proto

It has these top-level messages:
	Config
	Optimizer
*/
package internal

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import blacksquaremedia_reason_core "github.com/bsm/reason/core"
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

type Config struct {
	// The number of hash buckets.
	HashBuckets uint32 `protobuf:"varint,1,opt,name=num_buckets,json=numBuckets,proto3" json:"num_buckets,omitempty"`
	// Learn rate alpha.
	Alpha float64 `protobuf:"fixed64,2,opt,name=alpha,proto3" json:"alpha,omitempty"`
	// Learn rate beta.
	Beta float64 `protobuf:"fixed64,3,opt,name=beta,proto3" json:"beta,omitempty"`
	// Regularization parameter 1.
	L1 float64 `protobuf:"fixed64,4,opt,name=l1,proto3" json:"l1,omitempty"`
	// Regularization parameter 2.
	L2 float64 `protobuf:"fixed64,5,opt,name=l2,proto3" json:"l2,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{0} }

// Optimizer wraps the optimizer data.
type Optimizer struct {
	// The underlying model.
	Model *blacksquaremedia_reason_core.Model `protobuf:"bytes,1,opt,name=model" json:"model,omitempty"`
	// The target feature.
	Target string `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	// The config parameters
	Config Config `protobuf:"bytes,3,opt,name=config" json:"config"`
	// The gradient sums.
	Sums []float64 `protobuf:"fixed64,4,rep,packed,name=sums" json:"sums,omitempty"`
	// The weights.
	Weights []float64 `protobuf:"fixed64,5,rep,packed,name=weights" json:"weights,omitempty"`
}

func (m *Optimizer) Reset()                    { *m = Optimizer{} }
func (m *Optimizer) String() string            { return proto.CompactTextString(m) }
func (*Optimizer) ProtoMessage()               {}
func (*Optimizer) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{1} }

func init() {
	proto.RegisterType((*Config)(nil), "blacksquaremedia.reason.regression.ftrl.Config")
	proto.RegisterType((*Optimizer)(nil), "blacksquaremedia.reason.regression.ftrl.Optimizer")
}

func init() { proto.RegisterFile("regression/ftrl/internal/internal.proto", fileDescriptorInternal) }

var fileDescriptorInternal = []byte{
	// 368 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xc1, 0x4a, 0xeb, 0x40,
	0x14, 0x6d, 0xda, 0x24, 0xef, 0x75, 0xca, 0x7b, 0x8b, 0x41, 0x64, 0xe8, 0xa2, 0x2d, 0x15, 0x6d,
	0x37, 0x4e, 0x68, 0x5c, 0xb9, 0x8d, 0xe0, 0x42, 0x2c, 0x4a, 0x70, 0xe5, 0x46, 0x26, 0xe9, 0x34,
	0x1d, 0x9a, 0xc9, 0xd4, 0x99, 0x09, 0x82, 0x5f, 0xe1, 0x67, 0xf5, 0x13, 0x04, 0x41, 0x28, 0xfe,
	0x88, 0xe4, 0xa6, 0xa5, 0x6e, 0x84, 0x6e, 0xc2, 0x39, 0x87, 0x93, 0x33, 0x73, 0xee, 0x1d, 0x34,
	0xd2, 0x3c, 0xd3, 0xdc, 0x18, 0xa1, 0x8a, 0x60, 0x6e, 0x75, 0x1e, 0x88, 0xc2, 0x72, 0x5d, 0xb0,
	0x3d, 0xa0, 0x2b, 0xad, 0xac, 0xc2, 0xa3, 0x24, 0x67, 0xe9, 0xd2, 0x3c, 0x97, 0x4c, 0x73, 0xc9,
	0x67, 0x82, 0x51, 0xcd, 0x99, 0x51, 0x05, 0xdd, 0x07, 0xd0, 0x2a, 0xa0, 0x7b, 0x9a, 0x09, 0xbb,
	0x28, 0x13, 0x9a, 0x2a, 0x19, 0x24, 0x46, 0x06, 0xb5, 0x2d, 0x48, 0x95, 0xe6, 0xf0, 0xa9, 0xf3,
	0xba, 0xe7, 0x3f, 0x6c, 0x99, 0xca, 0x54, 0x00, 0x72, 0x52, 0xce, 0x81, 0x01, 0x01, 0x54, 0xdb,
	0x87, 0x06, 0xf9, 0x57, 0xaa, 0x98, 0x8b, 0x0c, 0xf7, 0x51, 0xa7, 0x28, 0xe5, 0x53, 0x52, 0xa6,
	0x4b, 0x6e, 0x0d, 0x71, 0x06, 0xce, 0xf8, 0x5f, 0x8c, 0x8a, 0x52, 0x46, 0xb5, 0x82, 0x8f, 0x90,
	0xc7, 0xf2, 0xd5, 0x82, 0x91, 0xe6, 0xc0, 0x19, 0x3b, 0x71, 0x4d, 0x30, 0x46, 0x6e, 0xc2, 0x2d,
	0x23, 0x2d, 0x10, 0x01, 0xe3, 0xff, 0xa8, 0x99, 0x4f, 0x88, 0x0b, 0x4a, 0x33, 0x9f, 0x00, 0x0f,
	0x89, 0xb7, 0xe5, 0xe1, 0xf0, 0xc3, 0x41, 0xed, 0xbb, 0x95, 0x15, 0x52, 0xbc, 0x72, 0x8d, 0x2f,
	0x91, 0x27, 0xd5, 0x8c, 0xe7, 0x70, 0x64, 0x27, 0x3c, 0xa1, 0xbf, 0x4d, 0x04, 0x5a, 0x4e, 0x2b,
	0x6b, 0x5c, 0xff, 0x81, 0x8f, 0x91, 0x6f, 0x99, 0xce, 0xb8, 0x85, 0x3b, 0xb5, 0xe3, 0x2d, 0xc3,
	0x53, 0xe4, 0xa7, 0xd0, 0x0a, 0xae, 0xd5, 0x09, 0x03, 0x7a, 0xe0, 0x94, 0x69, 0x3d, 0x8c, 0xc8,
	0x5d, 0x7f, 0xf6, 0x1b, 0xf1, 0x36, 0xa4, 0xea, 0x68, 0x4a, 0x69, 0x88, 0x3b, 0x68, 0x55, 0x1d,
	0x2b, 0x8c, 0x09, 0xfa, 0xf3, 0xc2, 0x45, 0xb6, 0xb0, 0x86, 0x78, 0x20, 0xef, 0x68, 0x74, 0xb3,
	0xde, 0xf4, 0x1a, 0xef, 0x9b, 0x9e, 0xf3, 0xf6, 0xd5, 0x6b, 0xa0, 0xb3, 0x54, 0xc9, 0x03, 0x4e,
	0x8f, 0xd0, 0xf5, 0x43, 0x7c, 0x7b, 0x5f, 0xed, 0xc4, 0x3c, 0xfe, 0xdd, 0xbd, 0x91, 0xc4, 0x87,
	0x2d, 0x5d, 0x7c, 0x07, 0x00, 0x00, 0xff, 0xff, 0x88, 0xcd, 0x71, 0x52, 0x4f, 0x02, 0x00, 0x00,
}
