package hoeffding

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"sort"
	"sync"

	"github.com/bsm/reason/classifier"
	"github.com/bsm/reason/classifier/hoeffding/internal"
	cinternal "github.com/bsm/reason/classifier/internal"
	"github.com/bsm/reason/core"
)

var (
	_ classifier.SupervisedLearner = (*Hoeffding)(nil)
	_ classifier.MultiCategory     = (*Hoeffding)(nil)
	_ classifier.Regressor         = (*Hoeffding)(nil)
)

// TreeInfo contains tree information/stats.
type TreeInfo struct {
	NumNodes    int // the total number of nodes
	NumLearning int // the number of learning leaves
	NumDisabled int // the number of disable leaves
	MaxDepth    int // the maximum depth
}

// Hoeffding is an implementation of a Hoeffding tree.
type Hoeffding struct {
	tree   *internal.Tree
	target *core.Feature

	config  Config
	cycles  int
	scratch []*internal.Node

	mu sync.RWMutex
}

// LoadFrom loads a new tree from a reader.
func LoadFrom(r io.Reader, config *Config) (*Hoeffding, error) {
	tt := new(internal.Tree)
	if _, err := tt.ReadFrom(r); err != nil {
		return nil, err
	}
	return newTree(tt, config)
}

// New inits a new Tree using a model, a target feature and a config.
func New(model *core.Model, target string, config *Config) (*Hoeffding, error) {
	return newTree(internal.NewTree(model, target), config)
}

func newTree(t *internal.Tree, c *Config) (*Hoeffding, error) {
	target := t.Model.Feature(t.Target)
	if target == nil {
		return nil, fmt.Errorf("hoeffding: unknown feature %q", t.Target)
	}

	var config Config
	if c != nil {
		config = *c
	}
	config.norm(target)

	if !config.SplitCriterion.Supports(target) {
		return nil, fmt.Errorf("hoeffding: split criterion is incompatible with target %q", t.Target)
	}

	return &Hoeffding{
		tree:   t,
		target: target,
		config: config,
	}, nil
}

// Info returns information about the tree
func (t *Hoeffding) Info() *TreeInfo {
	acc := new(TreeInfo)

	t.mu.RLock()
	t.recursiveInfo(t.tree.Root, 1, acc)
	t.mu.RUnlock()

	return acc
}

// WriteTo implements io.WriterTo
func (t *Hoeffding) WriteTo(w io.Writer) (int64, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.tree.WriteTo(w)
}

// Train trains the tree with an example x.
func (t *Hoeffding) Train(x core.Example) {
	t.TrainWeight(x, 1.0)
}

// TrainWeight trains the tree with an example x and a custom weight.
func (t *Hoeffding) TrainWeight(x core.Example, weight float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node, nref, parent, ppos := t.tree.Traverse(x, t.tree.Root, nil, -1)
	if node == nil && ppos > -1 {
		if split := parent.GetSplit(); split != nil {
			lref := t.tree.AddLeaf(nil, 0)
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
		weight := node.Weight()
		if leaf.IsDisabled || int(weight-leaf.WeightAtLastEval) < t.config.GracePeriod {
			return
		}

		// Store new weight
		leaf.WeightAtLastEval = weight

		// Check if we have sufficient stats to perform the split
		if !node.IsSufficient() {
			return
		}

		// Try to split
		if success := t.attemptSplit(leaf, node, nref, weight); success {
			if parent == nil {
				t.tree.Root = nref
			} else if split := parent.GetSplit(); split != nil {
				split.SetChild(ppos, nref)
			}
		}
	}
}

// PredictMC performs a classification and returns a prediction.
func (t *Hoeffding) PredictMC(x core.Example) classifier.MultiCategoryClassification {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.target.Kind != core.Feature_CATEGORICAL {
		return cinternal.NoResult{}
	}

	node, _, parent, _ := t.tree.Traverse(x, t.tree.Root, nil, -1)
	if node == nil {
		node = parent
	}

	stats := node.GetClassification()
	if stats == nil {
		return cinternal.NoResult{}
	}

	weight := stats.WeightSum()
	if weight <= 0 {
		return cinternal.NoResult{}
	}

	cat, _ := stats.Max()
	return classificationResult{cat: core.Category(cat), weight: weight, vv: &stats.Vector}
}

// PredictNum performs a regression and returns a prediction.
func (t *Hoeffding) PredictNum(x core.Example) classifier.Regression {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.target.Kind != core.Feature_NUMERICAL {
		return cinternal.NoResult{}
	}

	node, _, parent, _ := t.tree.Traverse(x, t.tree.Root, nil, -1)
	if node == nil {
		node = parent
	}

	stats := node.GetRegression()
	if stats == nil {
		return cinternal.NoResult{}
	}

	return regressionResult{ns: &stats.NumStream}
}

// Prune manually prunes the tree to limit it to maxLearningNodes.
func (t *Hoeffding) Prune(maxLearningNodes int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.prune(maxLearningNodes)
}

// WriteText writes text-based tree output to a writer
func (t *Hoeffding) WriteText(w io.Writer) (int64, error) {
	buf := bufio.NewWriter(w)

	t.mu.RLock()
	defer t.mu.RUnlock()

	nn, err := t.tree.WriteText(w, t.tree.Root, "", "ROOT")
	if err != nil {
		return nn, err
	}
	return nn, buf.Flush()
}

// WriteDOT writes a graph in dot notation to a writer
func (t *Hoeffding) WriteDOT(w io.Writer) (int64, error) {
	buf := bufio.NewWriter(w)
	nw := int64(0)

	n, err := buf.WriteString(`digraph ht {
  edge [fontsize=10];
  node [fontsize=10,shape=box];

`)
	nw += int64(n)
	if err != nil {
		return nw, err
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	nn, err := t.tree.WriteDOT(buf, t.tree.Root, "N", "")
	nw += nn
	if err != nil {
		return nw, err
	}

	n, err = buf.WriteString("}\n")
	nw += int64(n)
	if err != nil {
		return nw, err
	}
	return nw, buf.Flush()
}

func (t *Hoeffding) recursiveInfo(nref int64, depth int, acc *TreeInfo) {
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

func (t *Hoeffding) attemptSplit(leaf *internal.LeafNode, node *internal.Node, nref int64, weight float64) bool {
	// Init candidates, including a null result
	candidates := make(internal.SplitCandidates, 1, len(leaf.FeatureStats)+1)

	// Calculate a split candiate from each of the leaf stats
	for name := range leaf.FeatureStats {
		if c := leaf.EvaluateSplit(name, t.config.SplitCriterion, node); c != nil {
			candidates = append(candidates, *c)
		}
	}

	// Sort candidates by merit, select first
	sort.Stable(sort.Reverse(candidates))
	best := candidates[0]

	// Calculate the gain between merits of the best and the second-best split
	meritGain := best.Merit
	if len(candidates) > 1 {
		meritGain -= candidates[1].Merit
	}

	// Give up if there is no merit gain
	if meritGain <= 0 {
		return false
	}

	// Calculate confidence interval + hoeffding bound
	confiv := math.Log(1.0 / t.config.SplitConfidence)
	hbound := math.Sqrt(best.Range * best.Range * confiv * 0.5 / weight)
	if meritGain <= hbound && hbound >= t.config.TieThreshold {
		return false
	}

	// Split node
	t.tree.SplitNode(nref, best.Feature, best.PostSplit, best.Pivot)
	return true
}

func (t *Hoeffding) prune(maxLearningNodes int) {
	if maxLearningNodes < 0 {
		return
	}

	t.scratch = t.scratch[:0]
	for _, node := range t.tree.Nodes {
		if _, ok := node.GetKind().(*internal.Node_Leaf); ok {
			t.scratch = append(t.scratch, node)
		}
	}
	if len(t.scratch) <= maxLearningNodes {
		return
	}

	// Sort leaves by weight (highest first)
	sort.Slice(t.scratch, func(i, j int) bool {
		return t.scratch[i].Weight() >= t.scratch[j].Weight()
	})

	// Update node status
	for i, node := range t.scratch {
		if leaf := node.GetLeaf(); leaf != nil {
			if i < maxLearningNodes {
				leaf.Enable()
			} else {
				leaf.Disable()
			}
		}
	}
}
