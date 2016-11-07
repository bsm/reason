package hoeffding

import (
	"bufio"
	"io"
	"math"
	"sync"

	"github.com/bsm/reason/classifiers/internal/helpers"
	"github.com/bsm/reason/core"
)

// TreeInfo contains tree information/stats
type TreeInfo struct {
	NumNodes  int
	NumLeaves int
	MaxDepth  int
}

// Tree is an implementation of a HoeffdingTree
type Tree struct {
	conf *Config
	root treeNode

	model      *core.Model
	regression bool

	mu sync.RWMutex
}

func New(model *core.Model, conf *Config) *Tree {
	regression := model.IsRegression()
	if conf == nil {
		conf = new(Config)
	}
	conf.norm(regression)

	return &Tree{
		conf:       conf,
		model:      model,
		regression: regression,
		root:       newLeafNode(helpers.NewObservationStats(regression)),
	}
}

// Info returns information about the tree
func (t *Tree) Info() *TreeInfo {
	t.mu.RLock()
	root := t.root
	t.mu.RUnlock()

	var info TreeInfo
	info.NumNodes, info.NumLeaves, info.MaxDepth = root.Info()
	return &info
}

// WriteGraph write a graph in dot notation to a writer
func (t *Tree) WriteGraph(w io.Writer) error {
	buf := bufio.NewWriter(w)
	defer buf.Flush()

	if _, err := buf.WriteString("digraph ht {\n  edge [arrowsize=0.6, fontsize=10];\n"); err != nil {
		return err
	}

	t.mu.RLock()
	root := t.root
	t.mu.RUnlock()

	if err := root.AppendToGraph(buf, "N"); err != nil {
		return err
	}
	if _, err := buf.WriteString("}\n"); err != nil {
		return err
	}

	return nil
}

// Train passes an instance to the tree for training purposes
func (t *Tree) Train(inst core.Instance) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node, parent, parentIndex := t.root.Filter(inst, nil, -1)
	if node == nil {
		node = newLeafNode(helpers.NewObservationStats(t.regression))
		parent.children[parentIndex] = node
	}

	if leaf, ok := node.(*leafNode); ok {
		leaf.Learn(inst, t)

		weight := leaf.stats.TotalWeight()
		if int(weight-leaf.WeightOnLastEval()) < t.conf.GracePeriod {
			return
		}

		if split := t.attemptSplit(leaf, weight); split != nil {
			if parent == nil {
				t.root = split
			} else {
				parent.SetChild(parentIndex, split)
			}
		}
		leaf.SetWeightOnLastEval(weight)
	}
}

// Predict returns the raw votes by target index
func (t *Tree) Predict(inst core.Instance) core.Prediction {
	var res core.Prediction

	t.mu.Lock()
	node, parent, _ := t.root.Filter(inst, nil, -1)
	if node == nil {
		node = parent
	}
	res = node.Predict()
	t.mu.Unlock()
	return res
}

func (t *Tree) attemptSplit(leaf *leafNode, weight float64) *splitNode {
	if !leaf.stats.IsSufficient() {
		return nil
	}

	// Calculate best splits
	splits := leaf.BestSplits(t)
	bestSplit := splits[0]

	// Calculate the gain between merits of the best and the second-best split
	meritGain := bestSplit.Merit()
	if len(splits) > 1 {
		meritGain -= splits[1].Merit()
	}

	// Don't split if there is no merit gain
	if meritGain <= 0 {
		return nil
	}

	// Calculate hoeffding bound, evaluate split
	srange := bestSplit.Range()
	hbound := math.Sqrt(srange * srange * math.Log(1.0/t.conf.SplitConfidence) / (2.0 * weight))

	if meritGain > hbound || hbound < t.conf.TieThreshold {
		return newSplitNode(
			bestSplit.Condition(),
			bestSplit.PreStats(),
			bestSplit.PostStats(),
		)
	}
	return nil
}
