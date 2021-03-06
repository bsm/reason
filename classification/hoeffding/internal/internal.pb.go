// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: classification/hoeffding/internal/internal.proto

/*
Package internal is a generated protocol buffer package.

It is generated from these files:
	classification/hoeffding/internal/internal.proto

It has these top-level messages:
	Tree
	FeatureStats
	Node
	SplitNode
	LeafNode
*/
package internal

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import blacksquaremedia_reason_core "github.com/bsm/reason/core"
import blacksquaremedia_reason_util "github.com/bsm/reason/util"
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
	Model *blacksquaremedia_reason_core.Model `protobuf:"bytes,1,opt,name=model" json:"model,omitempty"`
	// The target feature.
	Target string `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	// The root nodeRef.
	Root int64 `protobuf:"varint,3,opt,name=root,proto3" json:"root,omitempty"`
	// The node registry.
	Nodes []*Node `protobuf:"bytes,4,rep,name=nodes" json:"nodes,omitempty"`
}

func (m *Tree) Reset()                    { *m = Tree{} }
func (m *Tree) String() string            { return proto.CompactTextString(m) }
func (*Tree) ProtoMessage()               {}
func (*Tree) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{0} }

// FeatureStats instances maintain stats based on
// observation of a particular feature.
type FeatureStats struct {
	// Types that are valid to be assigned to Kind:
	//	*FeatureStats_Numerical_
	//	*FeatureStats_Categorical_
	Kind isFeatureStats_Kind `protobuf_oneof:"kind"`
}

func (m *FeatureStats) Reset()                    { *m = FeatureStats{} }
func (m *FeatureStats) String() string            { return proto.CompactTextString(m) }
func (*FeatureStats) ProtoMessage()               {}
func (*FeatureStats) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{1} }

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
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *FeatureStats_Categorical_:
		s := proto.Size(x.Categorical)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type FeatureStats_Numerical struct {
	Min   blacksquaremedia_reason_util.Vector                  `protobuf:"bytes,1,opt,name=min" json:"min"`
	Max   blacksquaremedia_reason_util.Vector                  `protobuf:"bytes,2,opt,name=max" json:"max"`
	Stats blacksquaremedia_reason_util.StreamStatsDistribution `protobuf:"bytes,3,opt,name=stats" json:"stats"`
}

func (m *FeatureStats_Numerical) Reset()         { *m = FeatureStats_Numerical{} }
func (m *FeatureStats_Numerical) String() string { return proto.CompactTextString(m) }
func (*FeatureStats_Numerical) ProtoMessage()    {}
func (*FeatureStats_Numerical) Descriptor() ([]byte, []int) {
	return fileDescriptorInternal, []int{1, 0}
}

type FeatureStats_Categorical struct {
	blacksquaremedia_reason_util.VectorDistribution `protobuf:"bytes,1,opt,name=stats,embedded=stats" json:"stats"`
}

func (m *FeatureStats_Categorical) Reset()         { *m = FeatureStats_Categorical{} }
func (m *FeatureStats_Categorical) String() string { return proto.CompactTextString(m) }
func (*FeatureStats_Categorical) ProtoMessage()    {}
func (*FeatureStats_Categorical) Descriptor() ([]byte, []int) {
	return fileDescriptorInternal, []int{1, 1}
}

// Node is a tree node
type Node struct {
	// Observation stats for the node
	Stats *blacksquaremedia_reason_util.Vector `protobuf:"bytes,1,opt,name=stats" json:"stats,omitempty"`
	// Nodes can be leaf or split nodes.
	//
	// Types that are valid to be assigned to Kind:
	//	*Node_Leaf
	//	*Node_Split
	Kind isNode_Kind `protobuf_oneof:"kind"`
}

func (m *Node) Reset()                    { *m = Node{} }
func (m *Node) String() string            { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()               {}
func (*Node) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{2} }

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
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Node_Split:
		s := proto.Size(x.Split)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
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
	Children SplitNode_Children `protobuf:"bytes,3,opt,name=children" json:"children"`
}

func (m *SplitNode) Reset()                    { *m = SplitNode{} }
func (m *SplitNode) String() string            { return proto.CompactTextString(m) }
func (*SplitNode) ProtoMessage()               {}
func (*SplitNode) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{3} }

// Children is a collection of child node references.
type SplitNode_Children struct {
	Dense     []int64         `protobuf:"varint,1,rep,packed,name=dense" json:"dense,omitempty"`
	Sparse    map[int64]int64 `protobuf:"bytes,2,rep,name=sparse" json:"sparse,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	SparseCap int64           `protobuf:"varint,3,opt,name=sparse_cap,json=sparseCap,proto3" json:"sparse_cap,omitempty"`
}

func (m *SplitNode_Children) Reset()                    { *m = SplitNode_Children{} }
func (m *SplitNode_Children) String() string            { return proto.CompactTextString(m) }
func (*SplitNode_Children) ProtoMessage()               {}
func (*SplitNode_Children) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{3, 0} }

// LeafNode instances are the leaves within the tree.
type LeafNode struct {
	// Observation stats, but feature.
	FeatureStats map[string]*FeatureStats `protobuf:"bytes,1,rep,name=feature_stats,json=featureStats" json:"feature_stats,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	// Weight at the time of the last split evaluation.
	WeightAtLastEval float64 `protobuf:"fixed64,2,opt,name=weight_at_last_eval,json=weightAtLastEval,proto3" json:"weight_at_last_eval,omitempty"`
	// Status indicator.
	IsDisabled bool `protobuf:"varint,3,opt,name=is_disabled,json=isDisabled,proto3" json:"is_disabled,omitempty"`
}

func (m *LeafNode) Reset()                    { *m = LeafNode{} }
func (m *LeafNode) String() string            { return proto.CompactTextString(m) }
func (*LeafNode) ProtoMessage()               {}
func (*LeafNode) Descriptor() ([]byte, []int) { return fileDescriptorInternal, []int{4} }

func init() {
	proto.RegisterType((*Tree)(nil), "blacksquaremedia.reason.classification.hoeffding.Tree")
	proto.RegisterType((*FeatureStats)(nil), "blacksquaremedia.reason.classification.hoeffding.FeatureStats")
	proto.RegisterType((*FeatureStats_Numerical)(nil), "blacksquaremedia.reason.classification.hoeffding.FeatureStats.Numerical")
	proto.RegisterType((*FeatureStats_Categorical)(nil), "blacksquaremedia.reason.classification.hoeffding.FeatureStats.Categorical")
	proto.RegisterType((*Node)(nil), "blacksquaremedia.reason.classification.hoeffding.Node")
	proto.RegisterType((*SplitNode)(nil), "blacksquaremedia.reason.classification.hoeffding.SplitNode")
	proto.RegisterType((*SplitNode_Children)(nil), "blacksquaremedia.reason.classification.hoeffding.SplitNode.Children")
	proto.RegisterType((*LeafNode)(nil), "blacksquaremedia.reason.classification.hoeffding.LeafNode")
}

func init() {
	proto.RegisterFile("classification/hoeffding/internal/internal.proto", fileDescriptorInternal)
}

var fileDescriptorInternal = []byte{
	// 785 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x95, 0xdf, 0x8a, 0x23, 0x45,
	0x14, 0xc6, 0xd3, 0xe9, 0xce, 0x98, 0x9c, 0x5e, 0x71, 0x2d, 0x45, 0x42, 0xc0, 0x99, 0x21, 0x2a,
	0x04, 0x61, 0x3b, 0x43, 0x44, 0x71, 0x47, 0x11, 0xcc, 0xec, 0x4a, 0x90, 0xb8, 0xc4, 0xca, 0xe2,
	0x85, 0x37, 0xa1, 0xd2, 0x5d, 0x9d, 0x14, 0xd3, 0xdd, 0x95, 0xad, 0xaa, 0x8e, 0xbb, 0x57, 0xde,
	0x7b, 0xe5, 0x93, 0xf8, 0x02, 0xbe, 0xc0, 0x5e, 0x89, 0x97, 0xe2, 0xc5, 0xc2, 0xa2, 0x2f, 0xe0,
	0x1b, 0x48, 0xfd, 0xe9, 0x4c, 0x07, 0x19, 0x99, 0x38, 0xde, 0x34, 0x75, 0xba, 0xeb, 0xfc, 0xce,
	0x77, 0xbe, 0xaa, 0xae, 0x82, 0xb3, 0x38, 0x23, 0x52, 0xb2, 0x94, 0xc5, 0x44, 0x31, 0x5e, 0x0c,
	0xd7, 0x9c, 0xa6, 0x69, 0xc2, 0x8a, 0xd5, 0x90, 0x15, 0x8a, 0x8a, 0x82, 0x64, 0xbb, 0x41, 0xb4,
	0x11, 0x5c, 0x71, 0x74, 0xb6, 0xcc, 0x48, 0x7c, 0x29, 0x9f, 0x94, 0x44, 0xd0, 0x9c, 0x26, 0x8c,
	0x44, 0x82, 0x12, 0xc9, 0x8b, 0x68, 0x9f, 0x14, 0xed, 0x48, 0xbd, 0xf7, 0x56, 0x4c, 0xad, 0xcb,
	0x65, 0x14, 0xf3, 0x7c, 0xb8, 0x94, 0xf9, 0xd0, 0xce, 0x1f, 0xc6, 0x5c, 0x50, 0xf3, 0xb0, 0xe0,
	0xeb, 0xa6, 0x95, 0x8a, 0x65, 0xe6, 0xe1, 0xa6, 0xdd, 0xab, 0x4d, 0x5b, 0xf1, 0x15, 0x1f, 0x9a,
	0xd7, 0xcb, 0x32, 0x35, 0x91, 0x09, 0xcc, 0xc8, 0x4e, 0xef, 0xff, 0xec, 0x41, 0xf0, 0x58, 0x50,
	0x8a, 0xee, 0x43, 0x2b, 0xe7, 0x09, 0xcd, 0xba, 0xde, 0xa9, 0x37, 0x08, 0x47, 0xef, 0x44, 0xd7,
	0xf6, 0xa1, 0x25, 0x7d, 0xa5, 0xa7, 0x62, 0x9b, 0x81, 0xde, 0x82, 0x23, 0x45, 0xc4, 0x8a, 0xaa,
	0x6e, 0xf3, 0xd4, 0x1b, 0x74, 0xb0, 0x8b, 0x10, 0x82, 0x40, 0x70, 0xae, 0xba, 0xfe, 0xa9, 0x37,
	0xf0, 0xb1, 0x19, 0xa3, 0x29, 0xb4, 0x0a, 0x9e, 0x50, 0xd9, 0x0d, 0x4e, 0xfd, 0x41, 0x38, 0xfa,
	0x28, 0x3a, 0xd4, 0xae, 0xe8, 0x11, 0x4f, 0x28, 0xb6, 0x90, 0xfe, 0x4f, 0x01, 0xdc, 0xf9, 0x82,
	0x12, 0x55, 0x0a, 0x3a, 0x57, 0x44, 0x49, 0xb4, 0x86, 0x4e, 0x51, 0xe6, 0x54, 0xb0, 0x98, 0x54,
	0x9d, 0x4c, 0x0e, 0x2f, 0x51, 0x47, 0x46, 0x8f, 0x2a, 0xde, 0xa4, 0x81, 0xaf, 0xe0, 0xa8, 0x80,
	0x30, 0x26, 0x8a, 0xae, 0xb8, 0xad, 0xd5, 0x34, 0xb5, 0xbe, 0xbc, 0x65, 0xad, 0x8b, 0x2b, 0xe2,
	0xa4, 0x81, 0xeb, 0x05, 0x7a, 0xbf, 0x7b, 0xd0, 0xd9, 0x49, 0x41, 0x9f, 0x82, 0x9f, 0xb3, 0xc2,
	0x75, 0xf8, 0xee, 0xb5, 0x55, 0xcd, 0xbe, 0xf8, 0x86, 0xc6, 0x8a, 0x8b, 0x71, 0xf0, 0xfc, 0xc5,
	0x49, 0x03, 0xeb, 0x34, 0x93, 0x4d, 0x9e, 0x3a, 0xcd, 0x87, 0x65, 0x93, 0xa7, 0xe8, 0x6b, 0x68,
	0x49, 0xad, 0xd6, 0xac, 0x6b, 0x38, 0xfa, 0xf0, 0xdf, 0xf3, 0xe7, 0x4a, 0x50, 0x92, 0x9b, 0xf6,
	0x1e, 0x30, 0xa9, 0x04, 0x5b, 0x96, 0xda, 0x01, 0x07, 0xb4, 0xa4, 0xde, 0x02, 0xc2, 0x5a, 0xeb,
	0x68, 0x56, 0x55, 0xb0, 0xfd, 0x9d, 0xdd, 0x44, 0xe1, 0x1e, 0xbc, 0xad, 0xe1, 0xbf, 0xbe, 0x38,
	0xf1, 0x5c, 0x81, 0xf1, 0x11, 0x04, 0x97, 0xac, 0x48, 0xfa, 0x7f, 0x79, 0x10, 0xe8, 0x0d, 0x84,
	0xce, 0xf7, 0x4b, 0xdc, 0xc8, 0x04, 0x07, 0x43, 0x33, 0x08, 0x32, 0x4a, 0x52, 0xe7, 0xdf, 0xf9,
	0xe1, 0x6b, 0x3e, 0xa5, 0x24, 0xd5, 0x2a, 0x26, 0x0d, 0x6c, 0x48, 0x68, 0x0e, 0x2d, 0xb9, 0xc9,
	0x98, 0x72, 0x96, 0x7e, 0x72, 0x38, 0x72, 0xae, 0xd3, 0x1d, 0xd3, 0xb2, 0x76, 0x3d, 0xff, 0xe0,
	0x43, 0x67, 0xf7, 0x19, 0x75, 0xe1, 0x95, 0xd4, 0x6e, 0x39, 0xd3, 0x7a, 0x07, 0x57, 0x21, 0x7a,
	0x13, 0x5a, 0x1b, 0xb6, 0xe5, 0xf6, 0x2f, 0xf6, 0xb0, 0x0d, 0x50, 0x0a, 0xed, 0x78, 0xcd, 0xb2,
	0x44, 0xd0, 0xc2, 0xa9, 0x7b, 0x70, 0x0b, 0x75, 0xd1, 0x85, 0x63, 0xb9, 0xf5, 0xdf, 0xb1, 0x7b,
	0x7f, 0x7a, 0xd0, 0xae, 0x3e, 0x6a, 0x29, 0x09, 0x2d, 0xa4, 0x96, 0xe8, 0x0f, 0x7c, 0x6c, 0x03,
	0xb4, 0x86, 0x23, 0xb9, 0x21, 0x42, 0xd2, 0x6e, 0xd3, 0x1c, 0x1e, 0xb3, 0xff, 0x43, 0x48, 0x34,
	0x37, 0xc8, 0x87, 0x85, 0x12, 0xcf, 0xb0, 0xe3, 0xa3, 0xb7, 0x01, 0xec, 0x68, 0x11, 0x93, 0x8d,
	0x3b, 0xbf, 0x3a, 0xf6, 0xcd, 0x05, 0xd9, 0xf4, 0xee, 0x43, 0x58, 0xcb, 0x42, 0x77, 0xc1, 0xbf,
	0xa4, 0xcf, 0x8c, 0x9d, 0x3e, 0xd6, 0x43, 0xad, 0x7f, 0x4b, 0xb2, 0x92, 0x1a, 0x2b, 0x7d, 0x6c,
	0x83, 0xf3, 0xe6, 0xc7, 0x5e, 0xff, 0x97, 0x26, 0xb4, 0xab, 0xe5, 0x47, 0x4f, 0xe0, 0x55, 0x67,
	0xfe, 0xa2, 0xda, 0x8c, 0xba, 0xaf, 0xe9, 0x7f, 0xdf, 0x51, 0x7b, 0xc7, 0x89, 0xed, 0xe9, 0x4e,
	0x5a, 0x3f, 0x20, 0xef, 0xc1, 0x1b, 0xdf, 0x51, 0xb6, 0x5a, 0xab, 0x05, 0x51, 0x8b, 0x8c, 0x48,
	0xb5, 0xa0, 0x5b, 0x77, 0x7c, 0x79, 0xf8, 0xae, 0xfd, 0xf4, 0xb9, 0x9a, 0x12, 0xa9, 0x1e, 0x6e,
	0x49, 0x86, 0x4e, 0x20, 0x64, 0x72, 0x91, 0x30, 0x49, 0x96, 0x19, 0x4d, 0x8c, 0x13, 0x6d, 0x0c,
	0x4c, 0xff, 0xca, 0xe6, 0x4d, 0xef, 0x7b, 0x78, 0xfd, 0x1f, 0x25, 0xeb, 0x86, 0x74, 0xac, 0x21,
	0x8f, 0xeb, 0x86, 0x84, 0xa3, 0xcf, 0x6e, 0x77, 0x4e, 0xd6, 0x0c, 0x1d, 0xcf, 0x9f, 0xbf, 0x3c,
	0x6e, 0xfc, 0xf6, 0xf2, 0xd8, 0xfb, 0xf1, 0x8f, 0xe3, 0x06, 0xbc, 0x1f, 0xf3, 0xfc, 0x86, 0xec,
	0xf1, 0x6b, 0x93, 0x0a, 0x3e, 0xd3, 0x57, 0xa1, 0xfc, 0xb6, 0x5d, 0x5d, 0xe5, 0xcb, 0x23, 0x73,
	0x39, 0x7e, 0xf0, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x94, 0xf4, 0x01, 0xcc, 0xff, 0x07, 0x00,
	0x00,
}
