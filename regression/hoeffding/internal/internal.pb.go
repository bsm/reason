// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: regression/hoeffding/internal/internal.proto

package internal // import "github.com/bsm/reason/regression/hoeffding/internal"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import core "github.com/bsm/reason/core"
import util "github.com/bsm/reason/util"
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

// Tree wraps the tree data.
type Tree struct {
	// The underlying model.
	Model *core.Model `protobuf:"bytes,1,opt,name=model" json:"model,omitempty"`
	// The target feature.
	Target string `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	// The root nodeRef.
	Root int64 `protobuf:"varint,3,opt,name=root,proto3" json:"root,omitempty"`
	// The node registry.
	Nodes                []*Node  `protobuf:"bytes,4,rep,name=nodes" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
}

func (m *Tree) Reset()         { *m = Tree{} }
func (m *Tree) String() string { return proto.CompactTextString(m) }
func (*Tree) ProtoMessage()    {}
func (*Tree) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_8203064dbbb004fd, []int{0}
}
func (m *Tree) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tree.Unmarshal(m, b)
}
func (m *Tree) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tree.Marshal(b, m, deterministic)
}
func (dst *Tree) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tree.Merge(dst, src)
}
func (m *Tree) XXX_Size() int {
	return xxx_messageInfo_Tree.Size(m)
}
func (m *Tree) XXX_DiscardUnknown() {
	xxx_messageInfo_Tree.DiscardUnknown(m)
}

var xxx_messageInfo_Tree proto.InternalMessageInfo

// FeatureStats instances maintain stats based on
// observation of a particular feature.
type FeatureStats struct {
	// Types that are valid to be assigned to Kind:
	//	*FeatureStats_Numerical_
	//	*FeatureStats_Categorical_
	Kind                 isFeatureStats_Kind `protobuf_oneof:"kind"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
}

func (m *FeatureStats) Reset()         { *m = FeatureStats{} }
func (m *FeatureStats) String() string { return proto.CompactTextString(m) }
func (*FeatureStats) ProtoMessage()    {}
func (*FeatureStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_8203064dbbb004fd, []int{1}
}
func (m *FeatureStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureStats.Unmarshal(m, b)
}
func (m *FeatureStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureStats.Marshal(b, m, deterministic)
}
func (dst *FeatureStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureStats.Merge(dst, src)
}
func (m *FeatureStats) XXX_Size() int {
	return xxx_messageInfo_FeatureStats.Size(m)
}
func (m *FeatureStats) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureStats.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureStats proto.InternalMessageInfo

type isFeatureStats_Kind interface {
	isFeatureStats_Kind()
}

type FeatureStats_Numerical_ struct {
	Numerical *FeatureStats_Numerical `protobuf:"bytes,1,opt,name=numerical,oneof"`
}
type FeatureStats_Categorical_ struct {
	Categorical *FeatureStats_Categorical `protobuf:"bytes,2,opt,name=categorical,oneof"`
}

func (*FeatureStats_Numerical_) isFeatureStats_Kind()   {}
func (*FeatureStats_Categorical_) isFeatureStats_Kind() {}

func (m *FeatureStats) GetKind() isFeatureStats_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (m *FeatureStats) GetNumerical() *FeatureStats_Numerical {
	if x, ok := m.GetKind().(*FeatureStats_Numerical_); ok {
		return x.Numerical
	}
	return nil
}

func (m *FeatureStats) GetCategorical() *FeatureStats_Categorical {
	if x, ok := m.GetKind().(*FeatureStats_Categorical_); ok {
		return x.Categorical
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*FeatureStats) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _FeatureStats_OneofMarshaler, _FeatureStats_OneofUnmarshaler, _FeatureStats_OneofSizer, []interface{}{
		(*FeatureStats_Numerical_)(nil),
		(*FeatureStats_Categorical_)(nil),
	}
}

func _FeatureStats_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*FeatureStats)
	// kind
	switch x := m.Kind.(type) {
	case *FeatureStats_Numerical_:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Numerical); err != nil {
			return err
		}
	case *FeatureStats_Categorical_:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Categorical); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("FeatureStats.Kind has unexpected type %T", x)
	}
	return nil
}

func _FeatureStats_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*FeatureStats)
	switch tag {
	case 1: // kind.numerical
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FeatureStats_Numerical)
		err := b.DecodeMessage(msg)
		m.Kind = &FeatureStats_Numerical_{msg}
		return true, err
	case 2: // kind.categorical
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FeatureStats_Categorical)
		err := b.DecodeMessage(msg)
		m.Kind = &FeatureStats_Categorical_{msg}
		return true, err
	default:
		return false, nil
	}
}

func _FeatureStats_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*FeatureStats)
	// kind
	switch x := m.Kind.(type) {
	case *FeatureStats_Numerical_:
		s := proto.Size(x.Numerical)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *FeatureStats_Categorical_:
		s := proto.Size(x.Categorical)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type FeatureStats_Numerical struct {
	util.Histogram       `protobuf:"bytes,1,opt,name=stats,embedded=stats" json:"stats"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
}

func (m *FeatureStats_Numerical) Reset()         { *m = FeatureStats_Numerical{} }
func (m *FeatureStats_Numerical) String() string { return proto.CompactTextString(m) }
func (*FeatureStats_Numerical) ProtoMessage()    {}
func (*FeatureStats_Numerical) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_8203064dbbb004fd, []int{1, 0}
}
func (m *FeatureStats_Numerical) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureStats_Numerical.Unmarshal(m, b)
}
func (m *FeatureStats_Numerical) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureStats_Numerical.Marshal(b, m, deterministic)
}
func (dst *FeatureStats_Numerical) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureStats_Numerical.Merge(dst, src)
}
func (m *FeatureStats_Numerical) XXX_Size() int {
	return xxx_messageInfo_FeatureStats_Numerical.Size(m)
}
func (m *FeatureStats_Numerical) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureStats_Numerical.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureStats_Numerical proto.InternalMessageInfo

type FeatureStats_Categorical struct {
	util.NumStreams      `protobuf:"bytes,1,opt,name=stats,embedded=stats" json:"stats"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
}

func (m *FeatureStats_Categorical) Reset()         { *m = FeatureStats_Categorical{} }
func (m *FeatureStats_Categorical) String() string { return proto.CompactTextString(m) }
func (*FeatureStats_Categorical) ProtoMessage()    {}
func (*FeatureStats_Categorical) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_8203064dbbb004fd, []int{1, 1}
}
func (m *FeatureStats_Categorical) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureStats_Categorical.Unmarshal(m, b)
}
func (m *FeatureStats_Categorical) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureStats_Categorical.Marshal(b, m, deterministic)
}
func (dst *FeatureStats_Categorical) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureStats_Categorical.Merge(dst, src)
}
func (m *FeatureStats_Categorical) XXX_Size() int {
	return xxx_messageInfo_FeatureStats_Categorical.Size(m)
}
func (m *FeatureStats_Categorical) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureStats_Categorical.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureStats_Categorical proto.InternalMessageInfo

// Node is a tree node
type Node struct {
	// Observation stats for the node
	Stats *util.NumStream `protobuf:"bytes,1,opt,name=stats" json:"stats,omitempty"`
	// Nodes can be leaf or split nodes.
	//
	// Types that are valid to be assigned to Kind:
	//	*Node_Leaf
	//	*Node_Split
	Kind                 isNode_Kind `protobuf_oneof:"kind"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_8203064dbbb004fd, []int{2}
}
func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (dst *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(dst, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

type isNode_Kind interface {
	isNode_Kind()
}

type Node_Leaf struct {
	Leaf *LeafNode `protobuf:"bytes,2,opt,name=leaf,oneof"`
}
type Node_Split struct {
	Split *SplitNode `protobuf:"bytes,3,opt,name=split,oneof"`
}

func (*Node_Leaf) isNode_Kind()  {}
func (*Node_Split) isNode_Kind() {}

func (m *Node) GetKind() isNode_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (m *Node) GetLeaf() *LeafNode {
	if x, ok := m.GetKind().(*Node_Leaf); ok {
		return x.Leaf
	}
	return nil
}

func (m *Node) GetSplit() *SplitNode {
	if x, ok := m.GetKind().(*Node_Split); ok {
		return x.Split
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Node) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Node_OneofMarshaler, _Node_OneofUnmarshaler, _Node_OneofSizer, []interface{}{
		(*Node_Leaf)(nil),
		(*Node_Split)(nil),
	}
}

func _Node_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Node)
	// kind
	switch x := m.Kind.(type) {
	case *Node_Leaf:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Leaf); err != nil {
			return err
		}
	case *Node_Split:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Split); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Node.Kind has unexpected type %T", x)
	}
	return nil
}

func _Node_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Node)
	switch tag {
	case 2: // kind.leaf
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(LeafNode)
		err := b.DecodeMessage(msg)
		m.Kind = &Node_Leaf{msg}
		return true, err
	case 3: // kind.split
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SplitNode)
		err := b.DecodeMessage(msg)
		m.Kind = &Node_Split{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Node_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Node)
	// kind
	switch x := m.Kind.(type) {
	case *Node_Leaf:
		s := proto.Size(x.Leaf)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Node_Split:
		s := proto.Size(x.Split)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// SplitNode instances are intermediate nodes within the tree.
type SplitNode struct {
	// The feature name (predictor).
	Feature string `protobuf:"bytes,1,opt,name=feature,proto3" json:"feature,omitempty"`
	// The pivot value for binary splits (numerical predictors).
	Pivot float64 `protobuf:"fixed64,2,opt,name=pivot,proto3" json:"pivot,omitempty"`
	// The child references.
	Children             []int64  `protobuf:"varint,3,rep,packed,name=children" json:"children,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
}

func (m *SplitNode) Reset()         { *m = SplitNode{} }
func (m *SplitNode) String() string { return proto.CompactTextString(m) }
func (*SplitNode) ProtoMessage()    {}
func (*SplitNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_8203064dbbb004fd, []int{3}
}
func (m *SplitNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SplitNode.Unmarshal(m, b)
}
func (m *SplitNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SplitNode.Marshal(b, m, deterministic)
}
func (dst *SplitNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SplitNode.Merge(dst, src)
}
func (m *SplitNode) XXX_Size() int {
	return xxx_messageInfo_SplitNode.Size(m)
}
func (m *SplitNode) XXX_DiscardUnknown() {
	xxx_messageInfo_SplitNode.DiscardUnknown(m)
}

var xxx_messageInfo_SplitNode proto.InternalMessageInfo

// LeafNode instances are the leaves within the tree.
type LeafNode struct {
	// Observation stats, but feature.
	FeatureStats map[string]*FeatureStats `protobuf:"bytes,1,rep,name=feature_stats,json=featureStats" json:"feature_stats,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	// Weight at the time of the last split evaluation.
	WeightAtLastEval float64 `protobuf:"fixed64,2,opt,name=weight_at_last_eval,json=weightAtLastEval,proto3" json:"weight_at_last_eval,omitempty"`
	// Status indicator.
	IsDisabled           bool     `protobuf:"varint,3,opt,name=is_disabled,json=isDisabled,proto3" json:"is_disabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
}

func (m *LeafNode) Reset()         { *m = LeafNode{} }
func (m *LeafNode) String() string { return proto.CompactTextString(m) }
func (*LeafNode) ProtoMessage()    {}
func (*LeafNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_8203064dbbb004fd, []int{4}
}
func (m *LeafNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeafNode.Unmarshal(m, b)
}
func (m *LeafNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeafNode.Marshal(b, m, deterministic)
}
func (dst *LeafNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeafNode.Merge(dst, src)
}
func (m *LeafNode) XXX_Size() int {
	return xxx_messageInfo_LeafNode.Size(m)
}
func (m *LeafNode) XXX_DiscardUnknown() {
	xxx_messageInfo_LeafNode.DiscardUnknown(m)
}

var xxx_messageInfo_LeafNode proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Tree)(nil), "blacksquaremedia.reason.regression.hoeffding.Tree")
	proto.RegisterType((*FeatureStats)(nil), "blacksquaremedia.reason.regression.hoeffding.FeatureStats")
	proto.RegisterType((*FeatureStats_Numerical)(nil), "blacksquaremedia.reason.regression.hoeffding.FeatureStats.Numerical")
	proto.RegisterType((*FeatureStats_Categorical)(nil), "blacksquaremedia.reason.regression.hoeffding.FeatureStats.Categorical")
	proto.RegisterType((*Node)(nil), "blacksquaremedia.reason.regression.hoeffding.Node")
	proto.RegisterType((*SplitNode)(nil), "blacksquaremedia.reason.regression.hoeffding.SplitNode")
	proto.RegisterType((*LeafNode)(nil), "blacksquaremedia.reason.regression.hoeffding.LeafNode")
	proto.RegisterMapType((map[string]*FeatureStats)(nil), "blacksquaremedia.reason.regression.hoeffding.LeafNode.FeatureStatsEntry")
}

func init() {
	proto.RegisterFile("regression/hoeffding/internal/internal.proto", fileDescriptor_internal_8203064dbbb004fd)
}

var fileDescriptor_internal_8203064dbbb004fd = []byte{
	// 666 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xcd, 0x6e, 0xd3, 0x4e,
	0x10, 0xc0, 0xe3, 0x38, 0xe9, 0x3f, 0x99, 0xf4, 0xaf, 0x96, 0x05, 0xa1, 0x28, 0x87, 0x34, 0x0a,
	0x12, 0xe4, 0xd0, 0xda, 0x52, 0x2a, 0xf1, 0x51, 0x89, 0x03, 0xa1, 0x2d, 0x3e, 0x94, 0x52, 0x6d,
	0x2b, 0x21, 0xc1, 0x21, 0xda, 0xd8, 0x6b, 0x67, 0xa9, 0xed, 0x2d, 0xbb, 0xeb, 0xa0, 0x0a, 0x89,
	0x67, 0xe0, 0xc6, 0xa3, 0xf0, 0x06, 0xa8, 0xc7, 0x1e, 0x39, 0x55, 0xaa, 0x7a, 0xe5, 0x21, 0x90,
	0xbd, 0x8e, 0x1b, 0x54, 0xb5, 0x22, 0xe5, 0x62, 0xed, 0xac, 0x77, 0x7e, 0xf3, 0x3d, 0xb0, 0x2a,
	0x68, 0x20, 0xa8, 0x94, 0x8c, 0xc7, 0xf6, 0x98, 0x53, 0xdf, 0xf7, 0x58, 0x1c, 0xd8, 0x2c, 0x56,
	0x54, 0xc4, 0x24, 0x2c, 0x0e, 0xd6, 0x91, 0xe0, 0x8a, 0xa3, 0xd5, 0x51, 0x48, 0xdc, 0x43, 0xf9,
	0x31, 0x21, 0x82, 0x46, 0xd4, 0x63, 0xc4, 0x12, 0x94, 0x48, 0x1e, 0x5b, 0x97, 0x14, 0xab, 0xa0,
	0xb4, 0x96, 0x5c, 0x2e, 0xa8, 0x9d, 0x7e, 0xb4, 0x7a, 0x6b, 0x29, 0x51, 0x2c, 0xb4, 0xd3, 0x4f,
	0x7e, 0xb1, 0x16, 0x30, 0x35, 0x4e, 0x46, 0x96, 0xcb, 0x23, 0x3b, 0xe0, 0x01, 0xb7, 0xb3, 0xeb,
	0x51, 0xe2, 0x67, 0x52, 0x26, 0x64, 0x27, 0xfd, 0xbc, 0xfb, 0xdd, 0x80, 0xca, 0x81, 0xa0, 0x14,
	0x3d, 0x83, 0x6a, 0xc4, 0x3d, 0x1a, 0x36, 0x8d, 0x8e, 0xd1, 0x6b, 0xf4, 0x1f, 0x58, 0xd7, 0xf9,
	0x95, 0x19, 0x7f, 0x9d, 0x3e, 0xc5, 0x5a, 0x03, 0xdd, 0x87, 0x05, 0x45, 0x44, 0x40, 0x55, 0xb3,
	0xdc, 0x31, 0x7a, 0x75, 0x9c, 0x4b, 0x08, 0x41, 0x45, 0x70, 0xae, 0x9a, 0x66, 0xc7, 0xe8, 0x99,
	0x38, 0x3b, 0x23, 0x07, 0xaa, 0x31, 0xf7, 0xa8, 0x6c, 0x56, 0x3a, 0x66, 0xaf, 0xd1, 0xef, 0x5b,
	0xf3, 0x84, 0x6f, 0xed, 0x72, 0x8f, 0x62, 0x0d, 0xe8, 0x7e, 0x33, 0x61, 0x71, 0x9b, 0x12, 0x95,
	0x08, 0xba, 0xaf, 0x88, 0x92, 0xc8, 0x83, 0x7a, 0x9c, 0x44, 0x54, 0x30, 0x97, 0x4c, 0xa3, 0xd8,
	0x9c, 0x0f, 0x3f, 0x8b, 0xb3, 0x76, 0xa7, 0x2c, 0xa7, 0x84, 0x2f, 0xc1, 0xe8, 0x03, 0x34, 0x5c,
	0xa2, 0x68, 0xc0, 0xb5, 0x9d, 0x72, 0x66, 0x67, 0xfb, 0x1f, 0xec, 0xbc, 0xbc, 0xa4, 0x39, 0x25,
	0x3c, 0x0b, 0x6f, 0x1d, 0x40, 0xbd, 0xf0, 0x02, 0xbd, 0x82, 0xaa, 0x4c, 0x15, 0xf2, 0xd0, 0x1e,
	0x5d, 0x6b, 0x32, 0x6b, 0x06, 0x87, 0x49, 0xc5, 0x03, 0x41, 0xa2, 0x41, 0xed, 0xe4, 0x6c, 0xa5,
	0x74, 0x7a, 0xb6, 0x62, 0x60, 0xad, 0xdf, 0x7a, 0x0b, 0x8d, 0x19, 0x9b, 0x69, 0x45, 0x66, 0xb9,
	0xbd, 0x9b, 0xb9, 0xbb, 0x49, 0xb4, 0xaf, 0x04, 0x25, 0x91, 0xbc, 0x02, 0x1e, 0x2c, 0x40, 0xe5,
	0x90, 0xc5, 0x5e, 0xf7, 0x97, 0x01, 0x95, 0xb4, 0x52, 0xe8, 0xf9, 0x5c, 0x2e, 0x17, 0xe8, 0x9c,
	0x87, 0x76, 0xa0, 0x12, 0x52, 0xe2, 0xe7, 0x39, 0x7e, 0x3c, 0x5f, 0x8e, 0x77, 0x28, 0xf1, 0x53,
	0x27, 0x9c, 0x12, 0xce, 0x28, 0xe8, 0x0d, 0x54, 0xe5, 0x51, 0xc8, 0x74, 0x3b, 0x36, 0xfa, 0x4f,
	0xe6, 0xc3, 0xed, 0xa7, 0xaa, 0x39, 0x4f, 0x73, 0x8a, 0x70, 0xdf, 0x43, 0xbd, 0xf8, 0x8b, 0x9a,
	0xf0, 0x9f, 0xaf, 0xab, 0x9b, 0x05, 0x5d, 0xc7, 0x53, 0x11, 0xdd, 0x83, 0xea, 0x11, 0x9b, 0x70,
	0x3d, 0x24, 0x06, 0xd6, 0x02, 0x6a, 0x43, 0xcd, 0x1d, 0xb3, 0xd0, 0x13, 0x34, 0x6e, 0x9a, 0x1d,
	0xb3, 0x67, 0x0e, 0xca, 0xcb, 0x06, 0x2e, 0xee, 0xba, 0x3f, 0xca, 0x50, 0x9b, 0x86, 0x82, 0x22,
	0xf8, 0x3f, 0xa7, 0x0d, 0xa7, 0x79, 0x4d, 0x87, 0xc8, 0xb9, 0x5d, 0x66, 0xfe, 0x68, 0xc3, 0xad,
	0x58, 0x89, 0x63, 0xbc, 0xe8, 0xcf, 0x0e, 0xd4, 0x1a, 0xdc, 0xfd, 0x44, 0x59, 0x30, 0x56, 0x43,
	0xa2, 0x86, 0x21, 0x91, 0x6a, 0x48, 0x27, 0x79, 0xcb, 0x1b, 0x78, 0x59, 0xff, 0x7a, 0xa1, 0x76,
	0x88, 0x54, 0x5b, 0x13, 0x12, 0xa2, 0x15, 0x68, 0x30, 0x39, 0xf4, 0x98, 0x24, 0xa3, 0x90, 0x7a,
	0x59, 0x9a, 0x6b, 0x18, 0x98, 0xdc, 0xcc, 0x6f, 0x5a, 0x9f, 0xe1, 0xce, 0x15, 0x93, 0x68, 0x19,
	0xcc, 0x43, 0x7a, 0x9c, 0x27, 0x2b, 0x3d, 0xa2, 0x3d, 0xa8, 0x4e, 0x48, 0x98, 0xd0, 0xbc, 0xee,
	0x1b, 0xb7, 0x9f, 0x2d, 0xac, 0x41, 0x1b, 0xe5, 0xa7, 0xc6, 0xe0, 0xcb, 0xc9, 0x79, 0xbb, 0xf4,
	0xf3, 0xbc, 0x6d, 0x7c, 0xbd, 0x68, 0x97, 0x4e, 0x2f, 0xda, 0x25, 0x78, 0xe8, 0xf2, 0xe8, 0x2f,
	0xd8, 0x83, 0x25, 0x67, 0x0a, 0xdf, 0x4b, 0x57, 0xa6, 0x7c, 0xb7, 0x3e, 0xb3, 0x62, 0x47, 0x32,
	0xb2, 0xb5, 0x8a, 0x7d, 0xe3, 0xda, 0x1f, 0x2d, 0x64, 0xfb, 0x76, 0xfd, 0x77, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x4c, 0x8e, 0x8f, 0x1b, 0x1e, 0x06, 0x00, 0x00,
}
