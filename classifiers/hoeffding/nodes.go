package hoeffding

import (
	"bufio"
	"fmt"

	"github.com/bsm/reason/classifiers/internal/helpers"
	"github.com/bsm/reason/core"
)

var (
	_ treeNode = (*leafNode)(nil)
	_ treeNode = (*splitNode)(nil)
)

type treeNode interface {
	Filter(inst core.Instance, parent *splitNode, parentIndex int) (treeNode, *splitNode, int)
	AppendToGraph(*bufio.Writer, string) error
	Info() (numNodes, numLeaves, maxDepth int)
	Predict() core.Prediction
}

// --------------------------------------------------------------------

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

func (n *leafNode) Promise() float64 { return n.stats.Promise() }

func (n *leafNode) Predict() core.Prediction { return n.stats.State() }

func (n *leafNode) Info() (numNodes, numLeaves, maxDepth int) { return 1, 1, 1 }

func (n *leafNode) Filter(_ core.Instance, parent *splitNode, parentIndex int) (treeNode, *splitNode, int) {
	return n, parent, parentIndex
}

func (n *leafNode) AppendToGraph(w *bufio.Writer, nodeName string) error {
	_, err := fmt.Fprintf(w, "  %s [label=\"%.0f\", fontsize=10, shape=circle];\n", nodeName, n.stats.TotalWeight())
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

func (n *leafNode) Deactivate() {
	n.observers = nil
}

func (n *leafNode) Activate() {
	if n.observers == nil {
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

	// Skip the remaining steps if this node is deactivated
	if n.observers == nil {
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
	if n.observers == nil {
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

// --------------------------------------------------------------------

type splitNode struct {
	stats helpers.ObservationStats

	condition helpers.SplitCondition
	children  map[int]treeNode
}

func newSplitNode(condition helpers.SplitCondition, preSplit helpers.ObservationStats, postSplit []helpers.ObservationStats) *splitNode {
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

func (n *splitNode) Predict() core.Prediction { return n.stats.State() }

func (n *splitNode) Filter(inst core.Instance, parent *splitNode, parentIndex int) (treeNode, *splitNode, int) {
	if childIndex := n.condition.Branch(inst); childIndex > -1 {
		if child, ok := n.children[childIndex]; ok {
			return child.Filter(inst, n, childIndex)
		}
		return nil, n, childIndex
	}
	return n, parent, parentIndex
}

func (n *splitNode) Info() (numNodes, numLeaves, maxDepth int) {
	for _, child := range n.children {
		cNodes, cLeaves, cMaxDepth := child.Info()
		numNodes += cNodes
		numLeaves += cLeaves
		if cMaxDepth > maxDepth {
			maxDepth = cMaxDepth
		}
	}
	return numNodes + 1, numLeaves, maxDepth + 1
}

func (n *splitNode) AppendToGraph(w *bufio.Writer, nodeName string) error {
	if _, err := fmt.Fprintf(w, "  %s [label=%q shape=box];\n", nodeName, n.condition.Predictor().Name); err != nil {
		return err
	}

	for i, child := range n.children {
		subName := fmt.Sprintf("%s_%d", nodeName, i)

		if _, err := fmt.Fprintf(w, "  %s -> %s [label=%q];\n", nodeName, subName, n.condition.Describe(i)); err != nil {
			return err
		}
		if err := child.AppendToGraph(w, subName); err != nil {
			return err
		}
	}
	return nil
}

func (n *splitNode) SetChild(branch int, child treeNode) { n.children[branch] = child }
