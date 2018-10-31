package hoeffding

import (
	"fmt"
	"io"
	"sync"

	"github.com/bsm/reason/classifier"
	"github.com/bsm/reason/classifier/hoeffding/internal"
	"github.com/bsm/reason/core"
)

// TreeInfo contains tree information/stats.
type TreeInfo struct {
	NumNodes    int // the total number of nodes
	NumLearning int // the number of learning leaves
	NumDisabled int // the number of disable leaves
	MaxDepth    int // the maximum depth
}

// Tree is an implementation of a Hoeffding tree.
type Tree struct {
	tree    *internal.Tree
	target  *core.Feature
	problem classifier.Problem

	config  Config
	cycles  int
	scratch []*internal.Node

	mu sync.RWMutex
}

// LoadFrom loads a new tree from a reader.
func LoadFrom(r io.Reader, config *Config) (*Tree, error) {
	tt := new(internal.Tree)
	if _, err := tt.ReadFrom(r); err != nil {
		return nil, err
	}
	return newTree(tt, config)
}

// New inits a new Tree using a model, a target feature and a config.
func New(model *core.Model, target string, config *Config) (*Tree, error) {
	return newTree(internal.NewTree(model, target), config)
}

func newTree(t *internal.Tree, c *Config) (*Tree, error) {
	target := t.Model.Feature(t.Target)
	if target == nil {
		return nil, fmt.Errorf("hoeffding: unknown feature %q", t.Target)
	}

	problem := classifier.ProblemFromTarget(target)
	if !problem.IsValid() {
		return nil, fmt.Errorf("hoeffding: unsupported feature %q", t.Target)
	}

	var config Config
	if c != nil {
		config = *c
	}
	config.norm(problem)

	return &Tree{
		tree:    t,
		target:  target,
		problem: problem,
		config:  config,
	}, nil
}

// Info returns information about the tree
func (t *Tree) Info() *TreeInfo {
	acc := new(TreeInfo)

	t.mu.RLock()
	t.recursiveInfo(t.tree.Root, 1, acc)
	t.mu.RUnlock()

	return acc
}

// WriteTo implements io.WriterTo
func (t *Tree) WriteTo(w io.Writer) (int64, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.tree.WriteTo(w)
}

// Train trains the tree with an example x.
func (t *Tree) Train(x core.Example) {
	t.TrainWeight(x, 1.0)
}

// TrainWeight trains the tree with an example x and a custom weight.
func (t *Tree) TrainWeight(x core.Example, weight float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node, nref, parent, ppos := t.tree.Traverse(x, t.tree.Root, nil, -1)
	if node == nil && ppos > -1 {
		if split := parent.GetSplit(); split != nil {
			lref := t.tree.AddLeaf(t.problem, nil)
			node = t.tree.GetNode(lref)
			split.SetChild(ppos, lref)
		}
	}
	if node == nil {
		return
	}

	if leaf := node.GetLeaf(); leaf != nil {
		// Observe an example
		leaf.ObserveExample(t.tree.Model, t.target, x, weight, node)

		// Pre-prune, if enabled
		if t.config.PrunePeriod > 0 {
			if t.cycles++; t.config.MaxLearningNodes > 0 && t.cycles%t.config.PrunePeriod == 0 {
				t.prune(t.config.MaxLearningNodes)
			}
		}

		// Check if a split should be attempted
		nodeWeight := node.Weight()
		if leaf.IsDisabled || int(nodeWeight-leaf.WeightAtLastEval) < t.config.GracePeriod {
			return nil
		}

		// Store new weight
		leaf.WeightAtLastEval = nodeWeight

		// Check if we have sufficient stats to perform the split
		if !node.IsSufficient() {
			return nil
		}

		// Try to split
		info := t.attemptSplit(leaf, node, nodeRef, nodeWeight)
		if info.Success {
			if parent == nil {
				t.tree.Root = nodeRef
			} else if split := parent.GetSplit(); split != nil {
				split.SetChild(parentIndex, nodeRef)
			}
		}
		return info
	}

	return nil
}

func (t *Tree) recursiveInfo(nref int64, depth int, acc *TreeInfo) {
	node := t.tree.GetNode(nref)
	if node == nil {
		return
	}

	acc.NumNodes++
	if depth > acc.MaxDepth {
		acc.MaxDepth = depth
	}

	if split := node.GetSplit(); split != nil {
		for _, cnref := range split.Children {
			if cnref != 0 {
				t.recursiveInfo(cnref, depth+1, acc)
			}
		}
	} else if leaf := node.GetLeaf(); leaf != nil {
		if leaf.IsDisabled {
			acc.NumDisabled++
		} else {
			acc.NumLearning++
		}
	}
}
