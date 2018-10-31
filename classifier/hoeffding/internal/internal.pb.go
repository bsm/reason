// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: classifier/hoeffding/internal/internal.proto

package internal // import "github.com/bsm/reason/classifier/hoeffding/internal"

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
	return fileDescriptor_internal_b99d5c79526dcd59, []int{0}
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

// Node is a tree node
type Node struct {
	// Nodes can be leaf or split nodes.
	//
	// Types that are valid to be assigned to Kind:
	//	*Node_Leaf
	//	*Node_Split
	Kind isNode_Kind `protobuf_oneof:"kind"`
	// Observation stats for the node
	//
	// Types that are valid to be assigned to Stats:
	//	*Node_Classification
	//	*Node_Regression
	Stats                isNode_Stats `protobuf_oneof:"stats"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_b99d5c79526dcd59, []int{1}
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
type isNode_Stats interface {
	isNode_Stats()
}

type Node_Leaf struct {
	Leaf *LeafNode `protobuf:"bytes,1,opt,name=leaf,oneof"`
}
type Node_Split struct {
	Split *SplitNode `protobuf:"bytes,2,opt,name=split,oneof"`
}
type Node_Classification struct {
	Classification *Node_ClassificationStats `protobuf:"bytes,3,opt,name=classification,oneof"`
}
type Node_Regression struct {
	Regression *Node_RegressionStats `protobuf:"bytes,4,opt,name=regression,oneof"`
}

func (*Node_Leaf) isNode_Kind()            {}
func (*Node_Split) isNode_Kind()           {}
func (*Node_Classification) isNode_Stats() {}
func (*Node_Regression) isNode_Stats()     {}

func (m *Node) GetKind() isNode_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}
func (m *Node) GetStats() isNode_Stats {
	if m != nil {
		return m.Stats
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

func (m *Node) GetClassification() *Node_ClassificationStats {
	if x, ok := m.GetStats().(*Node_Classification); ok {
		return x.Classification
	}
	return nil
}

func (m *Node) GetRegression() *Node_RegressionStats {
	if x, ok := m.GetStats().(*Node_Regression); ok {
		return x.Regression
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Node) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Node_OneofMarshaler, _Node_OneofUnmarshaler, _Node_OneofSizer, []interface{}{
		(*Node_Leaf)(nil),
		(*Node_Split)(nil),
		(*Node_Classification)(nil),
		(*Node_Regression)(nil),
	}
}

func _Node_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Node)
	// kind
	switch x := m.Kind.(type) {
	case *Node_Leaf:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Leaf); err != nil {
			return err
		}
	case *Node_Split:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Split); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Node.Kind has unexpected type %T", x)
	}
	// stats
	switch x := m.Stats.(type) {
	case *Node_Classification:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Classification); err != nil {
			return err
		}
	case *Node_Regression:
		_ = b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Regression); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Node.Stats has unexpected type %T", x)
	}
	return nil
}

func _Node_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Node)
	switch tag {
	case 1: // kind.leaf
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(LeafNode)
		err := b.DecodeMessage(msg)
		m.Kind = &Node_Leaf{msg}
		return true, err
	case 2: // kind.split
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SplitNode)
		err := b.DecodeMessage(msg)
		m.Kind = &Node_Split{msg}
		return true, err
	case 3: // stats.classification
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Node_ClassificationStats)
		err := b.DecodeMessage(msg)
		m.Stats = &Node_Classification{msg}
		return true, err
	case 4: // stats.regression
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Node_RegressionStats)
		err := b.DecodeMessage(msg)
		m.Stats = &Node_Regression{msg}
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
	// stats
	switch x := m.Stats.(type) {
	case *Node_Classification:
		s := proto.Size(x.Classification)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Node_Regression:
		s := proto.Size(x.Regression)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Stats for classifications.
type Node_ClassificationStats struct {
	util.Vector          `protobuf:"bytes,1,opt,name=stats,embedded=stats" json:"stats"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
}

func (m *Node_ClassificationStats) Reset()         { *m = Node_ClassificationStats{} }
func (m *Node_ClassificationStats) String() string { return proto.CompactTextString(m) }
func (*Node_ClassificationStats) ProtoMessage()    {}
func (*Node_ClassificationStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_b99d5c79526dcd59, []int{1, 0}
}
func (m *Node_ClassificationStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node_ClassificationStats.Unmarshal(m, b)
}
func (m *Node_ClassificationStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node_ClassificationStats.Marshal(b, m, deterministic)
}
func (dst *Node_ClassificationStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node_ClassificationStats.Merge(dst, src)
}
func (m *Node_ClassificationStats) XXX_Size() int {
	return xxx_messageInfo_Node_ClassificationStats.Size(m)
}
func (m *Node_ClassificationStats) XXX_DiscardUnknown() {
	xxx_messageInfo_Node_ClassificationStats.DiscardUnknown(m)
}

var xxx_messageInfo_Node_ClassificationStats proto.InternalMessageInfo

// Stats for regressions.
type Node_RegressionStats struct {
	util.NumStream       `protobuf:"bytes,1,opt,name=stats,embedded=stats" json:"stats"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
}

func (m *Node_RegressionStats) Reset()         { *m = Node_RegressionStats{} }
func (m *Node_RegressionStats) String() string { return proto.CompactTextString(m) }
func (*Node_RegressionStats) ProtoMessage()    {}
func (*Node_RegressionStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_b99d5c79526dcd59, []int{1, 1}
}
func (m *Node_RegressionStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node_RegressionStats.Unmarshal(m, b)
}
func (m *Node_RegressionStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node_RegressionStats.Marshal(b, m, deterministic)
}
func (dst *Node_RegressionStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node_RegressionStats.Merge(dst, src)
}
func (m *Node_RegressionStats) XXX_Size() int {
	return xxx_messageInfo_Node_RegressionStats.Size(m)
}
func (m *Node_RegressionStats) XXX_DiscardUnknown() {
	xxx_messageInfo_Node_RegressionStats.DiscardUnknown(m)
}

var xxx_messageInfo_Node_RegressionStats proto.InternalMessageInfo

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
	return fileDescriptor_internal_b99d5c79526dcd59, []int{2}
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
	return fileDescriptor_internal_b99d5c79526dcd59, []int{3}
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

// FeatureStats instances maintain stats based on
// observation of a particular feature.
type FeatureStats struct {
	// Types that are valid to be assigned to Kind:
	//	*FeatureStats_CN
	//	*FeatureStats_CC
	//	*FeatureStats_RN
	//	*FeatureStats_RC
	Kind                 isFeatureStats_Kind `protobuf_oneof:"kind"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
}

func (m *FeatureStats) Reset()         { *m = FeatureStats{} }
func (m *FeatureStats) String() string { return proto.CompactTextString(m) }
func (*FeatureStats) ProtoMessage()    {}
func (*FeatureStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_b99d5c79526dcd59, []int{4}
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

type FeatureStats_CN struct {
	CN *FeatureStats_ClassificationNumerical `protobuf:"bytes,1,opt,name=cn,oneof"`
}
type FeatureStats_CC struct {
	CC *FeatureStats_ClassificationCategorical `protobuf:"bytes,2,opt,name=cc,oneof"`
}
type FeatureStats_RN struct {
	RN *FeatureStats_RegressionNumerical `protobuf:"bytes,3,opt,name=rn,oneof"`
}
type FeatureStats_RC struct {
	RC *FeatureStats_RegressionCategorical `protobuf:"bytes,4,opt,name=rc,oneof"`
}

func (*FeatureStats_CN) isFeatureStats_Kind() {}
func (*FeatureStats_CC) isFeatureStats_Kind() {}
func (*FeatureStats_RN) isFeatureStats_Kind() {}
func (*FeatureStats_RC) isFeatureStats_Kind() {}

func (m *FeatureStats) GetKind() isFeatureStats_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (m *FeatureStats) GetCN() *FeatureStats_ClassificationNumerical {
	if x, ok := m.GetKind().(*FeatureStats_CN); ok {
		return x.CN
	}
	return nil
}

func (m *FeatureStats) GetCC() *FeatureStats_ClassificationCategorical {
	if x, ok := m.GetKind().(*FeatureStats_CC); ok {
		return x.CC
	}
	return nil
}

func (m *FeatureStats) GetRN() *FeatureStats_RegressionNumerical {
	if x, ok := m.GetKind().(*FeatureStats_RN); ok {
		return x.RN
	}
	return nil
}

func (m *FeatureStats) GetRC() *FeatureStats_RegressionCategorical {
	if x, ok := m.GetKind().(*FeatureStats_RC); ok {
		return x.RC
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*FeatureStats) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _FeatureStats_OneofMarshaler, _FeatureStats_OneofUnmarshaler, _FeatureStats_OneofSizer, []interface{}{
		(*FeatureStats_CN)(nil),
		(*FeatureStats_CC)(nil),
		(*FeatureStats_RN)(nil),
		(*FeatureStats_RC)(nil),
	}
}

func _FeatureStats_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*FeatureStats)
	// kind
	switch x := m.Kind.(type) {
	case *FeatureStats_CN:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CN); err != nil {
			return err
		}
	case *FeatureStats_CC:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CC); err != nil {
			return err
		}
	case *FeatureStats_RN:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.RN); err != nil {
			return err
		}
	case *FeatureStats_RC:
		_ = b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.RC); err != nil {
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
	case 1: // kind.cn
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FeatureStats_ClassificationNumerical)
		err := b.DecodeMessage(msg)
		m.Kind = &FeatureStats_CN{msg}
		return true, err
	case 2: // kind.cc
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FeatureStats_ClassificationCategorical)
		err := b.DecodeMessage(msg)
		m.Kind = &FeatureStats_CC{msg}
		return true, err
	case 3: // kind.rn
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FeatureStats_RegressionNumerical)
		err := b.DecodeMessage(msg)
		m.Kind = &FeatureStats_RN{msg}
		return true, err
	case 4: // kind.rc
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FeatureStats_RegressionCategorical)
		err := b.DecodeMessage(msg)
		m.Kind = &FeatureStats_RC{msg}
		return true, err
	default:
		return false, nil
	}
}

func _FeatureStats_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*FeatureStats)
	// kind
	switch x := m.Kind.(type) {
	case *FeatureStats_CN:
		s := proto.Size(x.CN)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *FeatureStats_CC:
		s := proto.Size(x.CC)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *FeatureStats_RN:
		s := proto.Size(x.RN)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *FeatureStats_RC:
		s := proto.Size(x.RC)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type FeatureStats_ClassificationNumerical struct {
	Stats                util.NumStreams `protobuf:"bytes,1,opt,name=stats" json:"stats"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
}

func (m *FeatureStats_ClassificationNumerical) Reset()         { *m = FeatureStats_ClassificationNumerical{} }
func (m *FeatureStats_ClassificationNumerical) String() string { return proto.CompactTextString(m) }
func (*FeatureStats_ClassificationNumerical) ProtoMessage()    {}
func (*FeatureStats_ClassificationNumerical) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_b99d5c79526dcd59, []int{4, 0}
}
func (m *FeatureStats_ClassificationNumerical) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureStats_ClassificationNumerical.Unmarshal(m, b)
}
func (m *FeatureStats_ClassificationNumerical) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureStats_ClassificationNumerical.Marshal(b, m, deterministic)
}
func (dst *FeatureStats_ClassificationNumerical) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureStats_ClassificationNumerical.Merge(dst, src)
}
func (m *FeatureStats_ClassificationNumerical) XXX_Size() int {
	return xxx_messageInfo_FeatureStats_ClassificationNumerical.Size(m)
}
func (m *FeatureStats_ClassificationNumerical) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureStats_ClassificationNumerical.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureStats_ClassificationNumerical proto.InternalMessageInfo

type FeatureStats_ClassificationCategorical struct {
	Stats                util.Matrix `protobuf:"bytes,1,opt,name=stats" json:"stats"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
}

func (m *FeatureStats_ClassificationCategorical) Reset() {
	*m = FeatureStats_ClassificationCategorical{}
}
func (m *FeatureStats_ClassificationCategorical) String() string { return proto.CompactTextString(m) }
func (*FeatureStats_ClassificationCategorical) ProtoMessage()    {}
func (*FeatureStats_ClassificationCategorical) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_b99d5c79526dcd59, []int{4, 1}
}
func (m *FeatureStats_ClassificationCategorical) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureStats_ClassificationCategorical.Unmarshal(m, b)
}
func (m *FeatureStats_ClassificationCategorical) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureStats_ClassificationCategorical.Marshal(b, m, deterministic)
}
func (dst *FeatureStats_ClassificationCategorical) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureStats_ClassificationCategorical.Merge(dst, src)
}
func (m *FeatureStats_ClassificationCategorical) XXX_Size() int {
	return xxx_messageInfo_FeatureStats_ClassificationCategorical.Size(m)
}
func (m *FeatureStats_ClassificationCategorical) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureStats_ClassificationCategorical.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureStats_ClassificationCategorical proto.InternalMessageInfo

type FeatureStats_RegressionNumerical struct {
	Min                  float64                                        `protobuf:"fixed64,1,opt,name=min,proto3" json:"min,omitempty"`
	Max                  float64                                        `protobuf:"fixed64,2,opt,name=max,proto3" json:"max,omitempty"`
	Observations         []FeatureStats_RegressionNumerical_Observation `protobuf:"bytes,3,rep,name=observations" json:"observations"`
	XXX_NoUnkeyedLiteral struct{}                                       `json:"-"`
}

func (m *FeatureStats_RegressionNumerical) Reset()         { *m = FeatureStats_RegressionNumerical{} }
func (m *FeatureStats_RegressionNumerical) String() string { return proto.CompactTextString(m) }
func (*FeatureStats_RegressionNumerical) ProtoMessage()    {}
func (*FeatureStats_RegressionNumerical) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_b99d5c79526dcd59, []int{4, 2}
}
func (m *FeatureStats_RegressionNumerical) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureStats_RegressionNumerical.Unmarshal(m, b)
}
func (m *FeatureStats_RegressionNumerical) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureStats_RegressionNumerical.Marshal(b, m, deterministic)
}
func (dst *FeatureStats_RegressionNumerical) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureStats_RegressionNumerical.Merge(dst, src)
}
func (m *FeatureStats_RegressionNumerical) XXX_Size() int {
	return xxx_messageInfo_FeatureStats_RegressionNumerical.Size(m)
}
func (m *FeatureStats_RegressionNumerical) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureStats_RegressionNumerical.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureStats_RegressionNumerical proto.InternalMessageInfo

type FeatureStats_RegressionNumerical_Observation struct {
	FeatureValue         float64  `protobuf:"fixed64,1,opt,name=feature_value,json=featureValue,proto3" json:"feature_value,omitempty"`
	TargetValue          float64  `protobuf:"fixed64,2,opt,name=target_value,json=targetValue,proto3" json:"target_value,omitempty"`
	Weight               float64  `protobuf:"fixed64,3,opt,name=weight,proto3" json:"weight,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
}

func (m *FeatureStats_RegressionNumerical_Observation) Reset() {
	*m = FeatureStats_RegressionNumerical_Observation{}
}
func (m *FeatureStats_RegressionNumerical_Observation) String() string {
	return proto.CompactTextString(m)
}
func (*FeatureStats_RegressionNumerical_Observation) ProtoMessage() {}
func (*FeatureStats_RegressionNumerical_Observation) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_b99d5c79526dcd59, []int{4, 2, 0}
}
func (m *FeatureStats_RegressionNumerical_Observation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureStats_RegressionNumerical_Observation.Unmarshal(m, b)
}
func (m *FeatureStats_RegressionNumerical_Observation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureStats_RegressionNumerical_Observation.Marshal(b, m, deterministic)
}
func (dst *FeatureStats_RegressionNumerical_Observation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureStats_RegressionNumerical_Observation.Merge(dst, src)
}
func (m *FeatureStats_RegressionNumerical_Observation) XXX_Size() int {
	return xxx_messageInfo_FeatureStats_RegressionNumerical_Observation.Size(m)
}
func (m *FeatureStats_RegressionNumerical_Observation) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureStats_RegressionNumerical_Observation.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureStats_RegressionNumerical_Observation proto.InternalMessageInfo

type FeatureStats_RegressionCategorical struct {
	Stats                util.NumStreams `protobuf:"bytes,1,opt,name=stats" json:"stats"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
}

func (m *FeatureStats_RegressionCategorical) Reset()         { *m = FeatureStats_RegressionCategorical{} }
func (m *FeatureStats_RegressionCategorical) String() string { return proto.CompactTextString(m) }
func (*FeatureStats_RegressionCategorical) ProtoMessage()    {}
func (*FeatureStats_RegressionCategorical) Descriptor() ([]byte, []int) {
	return fileDescriptor_internal_b99d5c79526dcd59, []int{4, 3}
}
func (m *FeatureStats_RegressionCategorical) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeatureStats_RegressionCategorical.Unmarshal(m, b)
}
func (m *FeatureStats_RegressionCategorical) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeatureStats_RegressionCategorical.Marshal(b, m, deterministic)
}
func (dst *FeatureStats_RegressionCategorical) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeatureStats_RegressionCategorical.Merge(dst, src)
}
func (m *FeatureStats_RegressionCategorical) XXX_Size() int {
	return xxx_messageInfo_FeatureStats_RegressionCategorical.Size(m)
}
func (m *FeatureStats_RegressionCategorical) XXX_DiscardUnknown() {
	xxx_messageInfo_FeatureStats_RegressionCategorical.DiscardUnknown(m)
}

var xxx_messageInfo_FeatureStats_RegressionCategorical proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Tree)(nil), "blacksquaremedia.reason.classifier.hoeffding.Tree")
	proto.RegisterType((*Node)(nil), "blacksquaremedia.reason.classifier.hoeffding.Node")
	proto.RegisterType((*Node_ClassificationStats)(nil), "blacksquaremedia.reason.classifier.hoeffding.Node.ClassificationStats")
	proto.RegisterType((*Node_RegressionStats)(nil), "blacksquaremedia.reason.classifier.hoeffding.Node.RegressionStats")
	proto.RegisterType((*SplitNode)(nil), "blacksquaremedia.reason.classifier.hoeffding.SplitNode")
	proto.RegisterType((*LeafNode)(nil), "blacksquaremedia.reason.classifier.hoeffding.LeafNode")
	proto.RegisterMapType((map[string]*FeatureStats)(nil), "blacksquaremedia.reason.classifier.hoeffding.LeafNode.FeatureStatsEntry")
	proto.RegisterType((*FeatureStats)(nil), "blacksquaremedia.reason.classifier.hoeffding.FeatureStats")
	proto.RegisterType((*FeatureStats_ClassificationNumerical)(nil), "blacksquaremedia.reason.classifier.hoeffding.FeatureStats.ClassificationNumerical")
	proto.RegisterType((*FeatureStats_ClassificationCategorical)(nil), "blacksquaremedia.reason.classifier.hoeffding.FeatureStats.ClassificationCategorical")
	proto.RegisterType((*FeatureStats_RegressionNumerical)(nil), "blacksquaremedia.reason.classifier.hoeffding.FeatureStats.RegressionNumerical")
	proto.RegisterType((*FeatureStats_RegressionNumerical_Observation)(nil), "blacksquaremedia.reason.classifier.hoeffding.FeatureStats.RegressionNumerical.Observation")
	proto.RegisterType((*FeatureStats_RegressionCategorical)(nil), "blacksquaremedia.reason.classifier.hoeffding.FeatureStats.RegressionCategorical")
}

func init() {
	proto.RegisterFile("classifier/hoeffding/internal/internal.proto", fileDescriptor_internal_b99d5c79526dcd59)
}

var fileDescriptor_internal_b99d5c79526dcd59 = []byte{
	// 897 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x56, 0xcd, 0x6e, 0x1b, 0x37,
	0x10, 0xd6, 0xae, 0x7e, 0x22, 0x8f, 0xdc, 0xda, 0xa5, 0xdb, 0x54, 0xd5, 0x41, 0x76, 0x95, 0xa2,
	0xd5, 0x21, 0x59, 0x01, 0x0a, 0xd0, 0x9f, 0x9c, 0x5a, 0x39, 0x49, 0x75, 0x70, 0x14, 0x83, 0x0e,
	0x72, 0x70, 0x10, 0x08, 0xd4, 0x2e, 0x25, 0xb1, 0xde, 0x5d, 0xba, 0x24, 0xa5, 0x26, 0x28, 0xd0,
	0x53, 0x1f, 0xa0, 0xf7, 0x3e, 0x44, 0x8f, 0x7d, 0x83, 0x36, 0xc7, 0x1c, 0x7b, 0x32, 0x90, 0xe6,
	0x45, 0x0a, 0xfe, 0xac, 0x24, 0xcb, 0x76, 0x60, 0xb9, 0xbe, 0x08, 0x33, 0xb3, 0xe4, 0xf7, 0xf1,
	0x1b, 0xce, 0x8c, 0x08, 0xb7, 0xc3, 0x98, 0x48, 0xc9, 0x86, 0x8c, 0x8a, 0xd6, 0x98, 0xd3, 0xe1,
	0x30, 0x62, 0xe9, 0xa8, 0xc5, 0x52, 0x45, 0x45, 0x4a, 0xe2, 0x99, 0x11, 0x1c, 0x0b, 0xae, 0x38,
	0xba, 0x3d, 0x88, 0x49, 0x78, 0x24, 0x7f, 0x9c, 0x10, 0x41, 0x13, 0x1a, 0x31, 0x12, 0x08, 0x4a,
	0x24, 0x4f, 0x83, 0x39, 0x4a, 0x30, 0x43, 0xa9, 0x6d, 0x84, 0x5c, 0xd0, 0x96, 0xfe, 0xb1, 0xdb,
	0x6b, 0x1b, 0x13, 0xc5, 0xe2, 0x96, 0xfe, 0x71, 0x81, 0x3b, 0x23, 0xa6, 0xc6, 0x93, 0x41, 0x10,
	0xf2, 0xa4, 0x35, 0xe2, 0x23, 0xde, 0x32, 0xe1, 0xc1, 0x64, 0x68, 0x3c, 0xe3, 0x18, 0xcb, 0x2e,
	0x6f, 0xfc, 0xe9, 0x41, 0xe1, 0x89, 0xa0, 0x14, 0x7d, 0x03, 0xc5, 0x84, 0x47, 0x34, 0xae, 0x7a,
	0x3b, 0x5e, 0xb3, 0xd2, 0xbe, 0x15, 0x5c, 0x78, 0x2e, 0x4d, 0xfe, 0x48, 0x2f, 0xc5, 0x76, 0x07,
	0xba, 0x09, 0x25, 0x45, 0xc4, 0x88, 0xaa, 0xaa, 0xbf, 0xe3, 0x35, 0xd7, 0xb0, 0xf3, 0x10, 0x82,
	0x82, 0xe0, 0x5c, 0x55, 0xf3, 0x3b, 0x5e, 0x33, 0x8f, 0x8d, 0x8d, 0xba, 0x50, 0x4c, 0x79, 0x44,
	0x65, 0xb5, 0xb0, 0x93, 0x6f, 0x56, 0xda, 0xed, 0x60, 0x15, 0xf9, 0x41, 0x8f, 0x47, 0x14, 0x5b,
	0x80, 0xc6, 0xdf, 0x05, 0x28, 0x68, 0x1f, 0xed, 0x41, 0x21, 0xa6, 0x64, 0xe8, 0x0e, 0xfe, 0xe5,
	0x6a, 0x88, 0x7b, 0x94, 0x0c, 0x35, 0x4a, 0x37, 0x87, 0x0d, 0x0a, 0x7a, 0x0c, 0x45, 0x79, 0x1c,
	0x33, 0xab, 0xa5, 0xd2, 0xfe, 0x6a, 0x35, 0xb8, 0x03, 0xbd, 0xd5, 0xe1, 0x59, 0x1c, 0x74, 0x0c,
	0xef, 0x67, 0x4b, 0x43, 0xa2, 0x18, 0x4f, 0x4d, 0x3e, 0x2a, 0xed, 0x87, 0xab, 0x4b, 0x0f, 0x76,
	0x4f, 0x01, 0x1d, 0x28, 0xa2, 0x64, 0xd7, 0xc3, 0x4b, 0xf8, 0x28, 0x02, 0x10, 0x74, 0x24, 0xa8,
	0x94, 0x9a, 0xad, 0x60, 0xd8, 0x3a, 0x57, 0x60, 0xc3, 0x33, 0x90, 0x8c, 0x69, 0x01, 0xb7, 0xf6,
	0x0c, 0xb6, 0xce, 0x39, 0x0e, 0xba, 0x0f, 0x45, 0xa9, 0x0d, 0x77, 0x1d, 0x9f, 0x5d, 0xc8, 0x6b,
	0x6a, 0xf6, 0x29, 0x0d, 0x15, 0x17, 0x9d, 0xf2, 0xab, 0x93, 0xed, 0xdc, 0xeb, 0x93, 0x6d, 0x0f,
	0xdb, 0xcd, 0xb5, 0x43, 0xd8, 0x58, 0x62, 0x47, 0xdf, 0x9f, 0x06, 0xfe, 0xe2, 0xdd, 0xc0, 0xbd,
	0x49, 0x72, 0xa0, 0x04, 0x25, 0xc9, 0x19, 0xec, 0x4e, 0x09, 0x0a, 0x47, 0x2c, 0x8d, 0x3a, 0x37,
	0x1c, 0x60, 0xe3, 0x19, 0xac, 0xcd, 0xee, 0x0d, 0x55, 0xe1, 0xc6, 0x90, 0x12, 0x35, 0x11, 0xd4,
	0x10, 0xad, 0xe1, 0xcc, 0x45, 0x1f, 0x42, 0xf1, 0x98, 0x4d, 0xb9, 0xad, 0x0c, 0x0f, 0x5b, 0x07,
	0xd5, 0xa1, 0x1c, 0x8e, 0x59, 0x1c, 0x09, 0xaa, 0x2f, 0x36, 0xdf, 0xcc, 0x77, 0xfc, 0x4d, 0x0f,
	0xcf, 0x62, 0x8d, 0xbf, 0x7c, 0x28, 0x67, 0x45, 0x86, 0x12, 0x78, 0xcf, 0xa1, 0xf5, 0x33, 0x2d,
	0xba, 0x0b, 0xba, 0x57, 0xab, 0xd9, 0xe0, 0xa1, 0xc5, 0x32, 0xf9, 0x79, 0x90, 0x2a, 0xf1, 0x12,
	0xaf, 0x0f, 0x17, 0x42, 0xe8, 0x0e, 0x6c, 0xfd, 0x44, 0xd9, 0x68, 0xac, 0xfa, 0x44, 0xf5, 0x63,
	0x22, 0x55, 0x9f, 0x4e, 0x49, 0xec, 0xce, 0xbf, 0x69, 0x3f, 0x7d, 0xa7, 0xf6, 0x88, 0x54, 0x0f,
	0xa6, 0x24, 0x46, 0xdb, 0x50, 0x61, 0xb2, 0x1f, 0x31, 0x49, 0x06, 0x31, 0x8d, 0x4c, 0x99, 0x96,
	0x31, 0x30, 0x79, 0xdf, 0x45, 0x6a, 0x3f, 0xc3, 0x07, 0x67, 0x28, 0xd1, 0x26, 0xe4, 0x8f, 0xe8,
	0x4b, 0x97, 0x2c, 0x6d, 0xa2, 0x7d, 0x28, 0x4e, 0x49, 0x3c, 0xa1, 0xae, 0x85, 0xee, 0xad, 0xa6,
	0x6e, 0x91, 0x01, 0x5b, 0xa0, 0x7b, 0xfe, 0xd7, 0x5e, 0xe3, 0xf7, 0x32, 0xac, 0x2f, 0x7e, 0x43,
	0x31, 0xf8, 0x61, 0xea, 0xaa, 0x01, 0x5f, 0x9d, 0x63, 0xa9, 0xa9, 0x7a, 0x93, 0x84, 0x0a, 0x16,
	0x92, 0xb8, 0x53, 0xfa, 0xf7, 0x64, 0xdb, 0xdf, 0xed, 0x75, 0x73, 0xd8, 0x0f, 0x53, 0x94, 0x82,
	0x1f, 0x86, 0x4e, 0xd1, 0x93, 0x6b, 0x63, 0xdb, 0x25, 0x8a, 0x8e, 0xf8, 0x22, 0xdf, 0xae, 0xe1,
	0x0b, 0xd1, 0x18, 0x7c, 0x91, 0x8d, 0x8a, 0xde, 0xff, 0xe0, 0x9b, 0xb7, 0xd1, 0x92, 0x32, 0x6c,
	0x94, 0x89, 0x14, 0xfd, 0x00, 0xbe, 0x08, 0xdd, 0x98, 0xd8, 0xbf, 0x16, 0xa6, 0x33, 0xaa, 0xb0,
	0x51, 0x25, 0xc2, 0x5a, 0x1f, 0x3e, 0xbe, 0x20, 0xdd, 0xcb, 0x83, 0xa3, 0x79, 0xc9, 0xfe, 0x96,
	0x9d, 0x82, 0x6e, 0xf0, 0x6c, 0x70, 0x3c, 0x87, 0x4f, 0x2e, 0xcc, 0x30, 0xfa, 0x76, 0xa5, 0xd9,
	0xf4, 0x88, 0x28, 0xc1, 0x5e, 0x9c, 0x86, 0xff, 0xc3, 0x87, 0xad, 0x73, 0x32, 0xaa, 0x9b, 0x20,
	0x61, 0xb6, 0x18, 0x3d, 0xac, 0x4d, 0x13, 0x21, 0x2f, 0x5c, 0xaf, 0x69, 0x13, 0xfd, 0xea, 0xc1,
	0x3a, 0x1f, 0x48, 0x2a, 0xa6, 0xe6, 0x60, 0xd2, 0x8c, 0x8b, 0x4a, 0xfb, 0xf0, 0x7a, 0x2f, 0x37,
	0x78, 0x3c, 0xa7, 0x70, 0x67, 0x3f, 0xc5, 0x5a, 0x4b, 0xa0, 0xb2, 0xb0, 0x04, 0xdd, 0x9a, 0x8f,
	0x24, 0xdb, 0xb4, 0x56, 0x43, 0x36, 0x48, 0x9e, 0xea, 0x18, 0xfa, 0x14, 0xd6, 0xed, 0x7f, 0x7a,
	0x7f, 0xde, 0xd8, 0x1e, 0xae, 0xd8, 0x98, 0x5d, 0x72, 0x13, 0x4a, 0x76, 0xa0, 0x98, 0x9a, 0xf5,
	0xb0, 0xf3, 0x6a, 0xcf, 0xe1, 0xa3, 0x73, 0x0b, 0xe3, 0x7a, 0xee, 0x7b, 0x36, 0xcc, 0x7f, 0x79,
	0xf5, 0xa6, 0x9e, 0xfb, 0xe7, 0x4d, 0xdd, 0xfb, 0xed, 0x6d, 0x3d, 0xf7, 0xfa, 0x6d, 0x3d, 0x07,
	0x9f, 0x87, 0x3c, 0xb9, 0x44, 0x6a, 0x3b, 0x1b, 0xdd, 0x2c, 0xb7, 0xfb, 0xfa, 0x45, 0x24, 0x0f,
	0xef, 0x2e, 0xbc, 0xa0, 0x06, 0x32, 0x69, 0xd9, 0x2d, 0xad, 0x77, 0xbe, 0xea, 0x06, 0x25, 0xf3,
	0x9c, 0xba, 0xfb, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4c, 0x6b, 0xda, 0xe2, 0xfd, 0x09, 0x00,
	0x00,
}
