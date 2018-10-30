package internal

import (
	"github.com/bsm/reason/classifier"
	core "github.com/bsm/reason/core"
	util "github.com/bsm/reason/util"
)

// NewTree inits a brand-new tree
func NewTree(model *core.Model, target string) *Tree {
	t := &Tree{
		Model:  model,
		Target: target,
	}
	t.Root = t.AddLeaf(nil) // init root
	return t
}

// DetectProblem detects the classifier problem type.
func (t *Tree) DetectProblem() classifier.Problem {
	if feat := t.Model.Feature(t.Target); feat != nil {
		switch feat.Kind {
		case core.Feature_CATEGORICAL:
			return classifier.Classification
		case core.Feature_NUMERICAL:
			return classifier.Regression
		}
	}
	return 0
}

// AddLeaf adds a new node to registry and returns the reference.
func (t *Tree) AddLeaf(stats *util.Vector) int64 {
	if stats == nil {
		stats = util.NewVector()
	}

	leaf := &LeafNode{WeightAtLastEval: 0} // TODO: how to pass WeightAtLastEval
	kind := &Node_Leaf{Leaf: leaf}
	node := &Node{Kind: kind, Stats: stats}
	t.Nodes = append(t.Nodes, node)
	return int64(len(t.Nodes))
}

// ReplaceNode replaces a node by reference.
func (t *Tree) ReplaceNode(nref int64, node *Node) {
	if pos := int(nref - 1); pos > -1 && pos < len(t.Nodes) {
		t.Nodes[pos] = node
	}
}

// NumNodes returns the number of nodes.
func (t *Tree) NumNodes() int {
	return len(t.Nodes)
}

// GetNode retrieves a node by ref from registry.
func (t *Tree) GetNode(nref int64) *Node {
	if pos := int(nref - 1); pos > -1 && pos < len(t.Nodes) {
		return t.Nodes[pos]
	}
	return nil
}

// GetLeaf retrieves a leaf node by ref from registry.
func (t *Tree) GetLeaf(nref int64) *LeafNode {
	if node := t.GetNode(nref); node != nil {
		return node.GetLeaf()
	}
	return nil
}

// SplitLeaf splits an existing leaf node on feature.
func (t *Tree) SplitLeaf(nref int64, feature string, pre *util.Vector, post *util.Matrix, pivot float64) {
	// skip if node cannot be found
	if leaf := t.GetLeaf(nref); leaf == nil {
		return
	}

	// prepare split node
	split := &SplitNode{
		Feature: feature,
		Pivot:   pivot,
	}

	rows := post.NumRows()
	for i := 0; i < rows; i++ {
		row := post.Row(i)
		// TODO: figure out split
		// if hasStats(p, row) {
		// 	leaf := t.AddLeaf(util.NewVectorFromSlice(row...))
		// 	split.SetChild(i, leaf)
		// }
		_ = row
	}

	t.ReplaceNode(nref, &Node{
		Kind:  &Node_Split{Split: split},
		Stats: pre,
	})
}
