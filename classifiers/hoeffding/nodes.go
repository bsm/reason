package hoeffding

import (
	"bufio"
	"fmt"
	"math"

	"github.com/bsm/reason/classifiers/internal/helpers"
	"github.com/bsm/reason/core"
)

var (
	_ treeNode = (*leafNode)(nil)
	_ treeNode = (*splitNode)(nil)
)

type treeNode interface {
	Filter(inst core.Instance, parent *splitNode, parentIndex int) (treeNode, *splitNode, int)
	WriteGraph(*bufio.Writer, string) error
	WriteText(*bufio.Writer, string) error
	HeapSize() int
	ReadInfo(int, *TreeInfo)
	Predict() core.Prediction
	FindLeaves(leafNodeSlice) leafNodeSlice
}

// --------------------------------------------------------------------

type leafNodeSlice []*leafNode

func (p leafNodeSlice) Len() int      { return len(p) }
func (p leafNodeSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p leafNodeSlice) Less(i, j int) bool {
	a, b := p[i].Promise(), p[j].Promise()
	return a < b || math.IsNaN(a) && !math.IsNaN(b)
}

type leafNode struct {
	stats     helpers.ObservationStats
	observers []helpers.Observer

	weightOnLastEval float64
}

func newLeafNode(stats helpers.ObservationStats) *leafNode {
	return &leafNode{
		stats:            stats,
		weightOnLastEval: stats.TotalWeight(),
		observers:        []helpers.Observer{},
	}
}

func (n *leafNode) IsActive() bool { return n.observers != nil }

func (n *leafNode) Promise() float64 { return n.stats.Promise() }

func (n *leafNode) Predict() core.Prediction { return n.stats.State() }

func (n *leafNode) ReadInfo(depth int, info *TreeInfo) {
	info.NumNodes++

	if n.IsActive() {
		info.NumActiveLeaves++
	} else {
		info.NumInactiveLeaves++
	}
	if depth > info.MaxDepth {
		info.MaxDepth = depth
	}
}

func (n *leafNode) Filter(_ core.Instance, parent *splitNode, parentIndex int) (treeNode, *splitNode, int) {
	return n, parent, parentIndex
}

func (n *leafNode) WriteGraph(w *bufio.Writer, nodeName string) error {
	_, err := fmt.Fprintf(w, "  %s [label=\"%.0f\", fontsize=10, shape=circle];\n", nodeName, n.stats.TotalWeight())
	return err
}

func (n *leafNode) WriteText(w *bufio.Writer, _ string) error {
	_, err := fmt.Fprintf(w, " -> %.2f (%.0f)\n",
		n.Predict().Top().Value.Value(),
		n.stats.TotalWeight(),
	)
	return err
}

func (n *leafNode) WeightOnLastEval() float64 {
	return n.weightOnLastEval
}

func (n *leafNode) SetWeightOnLastEval(w float64) {
	if w > n.weightOnLastEval {
		n.weightOnLastEval = w
	}
}

func (n *leafNode) HeapSize() int {
	size := n.stats.HeapSize()
	if n.observers != nil {
		size += 24
	}
	for _, obs := range n.observers {
		size += obs.HeapSize()
	}
	return size
}

func (n *leafNode) Deactivate() {
	n.observers = nil
}

func (n *leafNode) Activate() {
	if !n.IsActive() {
		n.observers = []helpers.Observer{}
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
	n.stats.UpdatePreSplit(tv, weight)

	// Skip the remaining steps if this node is inactive
	if !n.IsActive() {
		return
	}

	// Update each predictor's observer with a target-value, predictor-value
	// and weight tuple
	predictors := tree.model.Predictors()
	if len(n.observers) == 0 {
		n.observers = make([]helpers.Observer, len(predictors))
	}
	for i, predictor := range predictors {
		pv := predictor.Value(inst)
		if pv.IsMissing() {
			continue
		}

		obs := n.observers[i]
		if obs == nil {
			obs = n.stats.NewObserver(predictor.IsNominal())
			n.observers[i] = obs
		}
		obs.Observe(tv, pv, weight)
	}
}

func (n *leafNode) BestSplits(tree *Tree) helpers.SplitSuggestions {
	if !n.IsActive() {
		return nil
	}

	// Init split-suggestions, including a null suggestion
	suggestions := make(helpers.SplitSuggestions, 1, len(n.observers)+1)

	// Calculate a split suggestion for each of the observed predictors
	predictors := tree.model.Predictors()
	for i, obs := range n.observers {
		if obs != nil {
			split := n.stats.BestSplit(tree.conf.SplitCriterion, obs, predictors[i])
			suggestions = append(suggestions, split)
		}
	}

	// Rank the suggestions by merit
	return suggestions.Rank()
}

func (n *leafNode) FindLeaves(acc leafNodeSlice) leafNodeSlice { return append(acc, n) }

// --------------------------------------------------------------------

type splitNode struct {
	stats helpers.ObservationStats

	condition helpers.SplitCondition
	children  map[int]treeNode
}

func newSplitNode(condition helpers.SplitCondition, preSplit helpers.ObservationStats, postSplit map[int]helpers.ObservationStats) *splitNode {
	children := make(map[int]treeNode, len(postSplit))
	for i, stats := range postSplit {
		children[i] = newLeafNode(stats)
	}

	return &splitNode{
		stats:     preSplit,
		condition: condition,
		children:  children,
	}
}

func (n *splitNode) HeapSize() int {
	size := n.stats.HeapSize() + 64
	for _, c := range n.children {
		size += 8
		size += c.HeapSize()
	}
	return size
}

func (n *splitNode) Predict() core.Prediction { return n.stats.State() }

func (n *splitNode) Filter(inst core.Instance, parent *splitNode, parentIndex int) (treeNode, *splitNode, int) {
	if branch := n.condition.Branch(inst); branch > -1 {
		if child, ok := n.children[branch]; ok {
			return child.Filter(inst, n, branch)
		}
		return nil, n, branch
	}
	return n, parent, parentIndex
}

func (n *splitNode) ReadInfo(depth int, info *TreeInfo) {
	info.NumNodes++
	for _, child := range n.children {
		child.ReadInfo(depth+1, info)
	}
}

func (n *splitNode) WriteGraph(w *bufio.Writer, nodeName string) error {
	if _, err := fmt.Fprintf(w, "  %s [label=%q shape=box];\n", nodeName, n.condition.Predictor().Name); err != nil {
		return err
	}

	for i, child := range n.children {
		subName := fmt.Sprintf("%s_%d", nodeName, i)

		if _, err := fmt.Fprintf(w, "  %s -> %s [label=%q];\n", nodeName, subName, n.condition.Describe(i)); err != nil {
			return err
		}
		if err := child.WriteGraph(w, subName); err != nil {
			return err
		}
	}
	return nil
}

func (n *splitNode) WriteText(w *bufio.Writer, indent string) error {
	if _, err := fmt.Fprintf(w, " -> %.2f (%.0f)\n", n.Predict().Top().Value.Value(), n.stats.TotalWeight()); err != nil {
		return err
	}

	name := n.condition.Predictor().Name
	sind := indent + "\t"
	for i, child := range n.children {
		if _, err := fmt.Fprintf(w, "%s%s %q", indent, name, n.condition.Describe(i)); err != nil {
			return err
		}

		if err := child.WriteText(w, sind); err != nil {
			return err
		}
	}
	return nil
}

func (n *splitNode) SetChild(branch int, child treeNode) { n.children[branch] = child }

func (n *splitNode) FindLeaves(acc leafNodeSlice) leafNodeSlice {
	for _, c := range n.children {
		acc = c.FindLeaves(acc)
	}
	return acc
}
