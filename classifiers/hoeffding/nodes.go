package hoeffding

import (
	"bufio"
	"fmt"
	"math"

	"github.com/bsm/reason/classifiers/internal/helpers"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/msgpack"
)

func init() {
	msgpack.Register(7748, (*leafNode)(nil))
	msgpack.Register(7749, (*splitNode)(nil))
}

var (
	_ treeNode = (*leafNode)(nil)
	_ treeNode = (*splitNode)(nil)
)

type treeNode interface {
	Filter(inst core.Instance, parent *splitNode, parentIndex int) (treeNode, *splitNode, int)
	Prune(isObsolete PruneEval, parent *splitNode, parentIndex int)
	WriteGraph(*bufio.Writer, string) error
	WriteText(*bufio.Writer, string) error
	ByteSize() int
	TotalWeight() float64
	ReadInfo(int, *TreeInfo)
	FindLeaves(leafNodeSlice) leafNodeSlice
	Predict() core.Prediction
}

var (
	_ Node = (*leafNode)(nil)
	_ Node = (*splitNode)(nil)
)

// Node contains several useful details about the node
type Node interface {
	// IsLeaf returns true for leaves
	IsLeaf() bool
	// TotalWeight returns the total weight seen on this node so far
	TotalWeight() float64
	// Predict returns the current prediction for this node
	Predict() core.Prediction
}

// --------------------------------------------------------------------

type leafNodeSlice []*leafNode

func (p leafNodeSlice) Len() int      { return len(p) }
func (p leafNodeSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p leafNodeSlice) Less(i, j int) bool {
	a, b := p[i].TotalWeight(), p[j].TotalWeight()
	return a < b || math.IsNaN(a) && !math.IsNaN(b)
}

type leafNode struct {
	Stats     helpers.ObservationStats
	Observers []helpers.Observer

	IsInactive       bool
	WeightOnLastEval float64
}

func newLeafNode(stats helpers.ObservationStats) *leafNode {
	return &leafNode{
		Stats:            stats,
		WeightOnLastEval: stats.TotalWeight(),
	}
}

func (n *leafNode) TotalWeight() float64     { return n.Stats.TotalWeight() }
func (n *leafNode) IsLeaf() bool             { return true }
func (n *leafNode) Predict() core.Prediction { return n.Stats.State() }

func (n *leafNode) ReadInfo(depth int, info *TreeInfo) {
	info.NumNodes++

	if n.IsInactive {
		info.NumInactiveLeaves++
	} else {
		info.NumActiveLeaves++
	}
	if depth > info.MaxDepth {
		info.MaxDepth = depth
	}
}

func (n *leafNode) Filter(_ core.Instance, parent *splitNode, parentIndex int) (treeNode, *splitNode, int) {
	return n, parent, parentIndex
}

func (n *leafNode) Prune(isObsolete PruneEval, parent *splitNode, parentIndex int) {
	if parent != nil && isObsolete(n, parent) {
		delete(parent.Children, parentIndex)
	}
}

func (n *leafNode) WriteGraph(w *bufio.Writer, nodeName string) error {
	_, err := fmt.Fprintf(w, "  %s [label=\"%.0f\", fontsize=10, shape=circle];\n", nodeName, n.TotalWeight())
	return err
}

func (n *leafNode) WriteText(w *bufio.Writer, _ string) error {
	_, err := fmt.Fprintf(w, " -> %.2f (%.0f)\n",
		n.Predict().Value(),
		n.TotalWeight(),
	)
	return err
}

func (n *leafNode) ByteSize() int {
	size := 40 + n.Stats.ByteSize()
	if n.Observers != nil {
		size += 24
	}
	for _, obs := range n.Observers {
		size += obs.ByteSize()
	}
	return size
}

func (n *leafNode) Deactivate() {
	n.IsInactive = true
	n.Observers = nil
}

func (n *leafNode) Activate() {
	if n.IsInactive {
		n.IsInactive = false
	}
}

func (n *leafNode) Learn(inst core.Instance, tree *Tree) {
	// Get the target value, skip this instance if missing
	tv := tree.model.Target().Value(inst)
	if tv.IsMissing() {
		return
	}

	// Get instance weight and update pre-split distribution stats
	weight := inst.GetInstanceWeight()
	n.Stats.UpdatePreSplit(tv, weight)

	// Skip the remaining steps if this node is inactive
	if n.IsInactive {
		return
	}

	// Update each predictor's observer with a target-value, predictor-value
	// and weight tuple
	predictors := tree.model.Predictors()
	if len(n.Observers) == 0 {
		n.Observers = make([]helpers.Observer, len(predictors))
	}
	for i, predictor := range predictors {
		pv := predictor.Value(inst)
		if pv.IsMissing() {
			continue
		}

		obs := n.Observers[i]
		if obs == nil {
			obs = n.Stats.NewObserver(predictor.IsNominal())
			n.Observers[i] = obs
		}
		obs.Observe(tv, pv, weight)
	}
}

func (n *leafNode) BestSplits(tree *Tree) helpers.SplitSuggestions {
	if n.IsInactive {
		return nil
	}

	// Init split-suggestions, including a null suggestion
	suggestions := make(helpers.SplitSuggestions, 1, len(n.Observers)+1)

	// Calculate a split suggestion for each of the observed predictors
	predictors := tree.model.Predictors()
	for i, obs := range n.Observers {
		if obs != nil {
			split := n.Stats.BestSplit(tree.conf.SplitCriterion, obs, predictors[i])
			suggestions = append(suggestions, split)
		}
	}

	// Rank the suggestions by merit
	return suggestions.Rank()
}

func (n *leafNode) FindLeaves(acc leafNodeSlice) leafNodeSlice { return append(acc, n) }

func (n *leafNode) EncodeTo(enc *msgpack.Encoder) error {
	return enc.Encode(n.Stats, n.Observers, n.WeightOnLastEval, n.IsInactive)
}

func (n *leafNode) DecodeFrom(dec *msgpack.Decoder) error {
	return dec.Decode(&n.Stats, &n.Observers, &n.WeightOnLastEval, &n.IsInactive)
}

// --------------------------------------------------------------------

type splitNode struct {
	Stats     helpers.ObservationStats
	Condition helpers.SplitCondition
	Children  map[int]treeNode
}

func newSplitNode(condition helpers.SplitCondition, preSplit helpers.ObservationStats, postSplit map[int]helpers.ObservationStats) *splitNode {
	children := make(map[int]treeNode, len(postSplit))
	for i, stats := range postSplit {
		children[i] = newLeafNode(stats)
	}

	return &splitNode{
		Stats:     preSplit,
		Condition: condition,
		Children:  children,
	}
}

func (n *splitNode) ByteSize() int {
	size := n.Stats.ByteSize() + 64
	for _, c := range n.Children {
		size += 8
		size += c.ByteSize()
	}
	return size
}

func (n *splitNode) TotalWeight() float64     { return n.Stats.TotalWeight() }
func (n *splitNode) IsLeaf() bool             { return false }
func (n *splitNode) Predict() core.Prediction { return n.Stats.State() }

func (n *splitNode) Filter(inst core.Instance, parent *splitNode, parentIndex int) (treeNode, *splitNode, int) {
	if branch := n.Condition.Branch(inst); branch > -1 {
		if child, ok := n.Children[branch]; ok {
			return child.Filter(inst, n, branch)
		}
		return nil, n, branch
	}
	return n, parent, parentIndex
}

func (n *splitNode) Prune(isObsolete PruneEval, parent *splitNode, parentIndex int) {
	if parent != nil && isObsolete(n, parent) {
		delete(parent.Children, parentIndex)
		return
	}

	for i, child := range n.Children {
		child.Prune(isObsolete, n, i)
	}
}

func (n *splitNode) ReadInfo(depth int, info *TreeInfo) {
	info.NumNodes++
	for _, child := range n.Children {
		child.ReadInfo(depth+1, info)
	}
}

func (n *splitNode) WriteGraph(w *bufio.Writer, nodeName string) error {
	if _, err := fmt.Fprintf(w, "  %s [label=%q shape=box];\n", nodeName, n.Condition.Predictor()); err != nil {
		return err
	}

	for i, child := range n.Children {
		subName := fmt.Sprintf("%s_%d", nodeName, i)

		if _, err := fmt.Fprintf(w, "  %s -> %s [label=%q];\n", nodeName, subName, n.Condition.Describe(i)); err != nil {
			return err
		}
		if err := child.WriteGraph(w, subName); err != nil {
			return err
		}
	}
	return nil
}

func (n *splitNode) WriteText(w *bufio.Writer, indent string) error {
	if _, err := fmt.Fprintf(w, " -> %.2f (%.0f)\n", n.Predict().Value(), n.TotalWeight()); err != nil {
		return err
	}

	name := n.Condition.Predictor()
	sind := indent + "\t"
	for i, child := range n.Children {
		if _, err := fmt.Fprintf(w, "%s%s %q", indent, name, n.Condition.Describe(i)); err != nil {
			return err
		}

		if err := child.WriteText(w, sind); err != nil {
			return err
		}
	}
	return nil
}

func (n *splitNode) SetChild(branch int, child treeNode) { n.Children[branch] = child }

func (n *splitNode) FindLeaves(acc leafNodeSlice) leafNodeSlice {
	for _, c := range n.Children {
		acc = c.FindLeaves(acc)
	}
	return acc
}

func (n *splitNode) EncodeTo(enc *msgpack.Encoder) error {
	return enc.Encode(n.Stats, n.Condition, n.Children)
}

func (n *splitNode) DecodeFrom(dec *msgpack.Decoder) error {
	return dec.Decode(&n.Stats, &n.Condition, &n.Children)
}
