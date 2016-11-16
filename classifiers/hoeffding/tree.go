package hoeffding

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"math"
	"sort"
	"sync"

	"github.com/bsm/reason/classifiers/internal/helpers"
	"github.com/bsm/reason/core"
)

// TreeInfo contains tree information/stats
type TreeInfo struct {
	NumNodes          int
	NumActiveLeaves   int
	NumInactiveLeaves int
	MaxDepth          int
}

// Tree is an implementation of a HoeffdingTree
type Tree struct {
	conf  *Config
	root  treeNode
	model *core.Model

	traces chan *Trace
	leaves leafNodeSlice
	cycles int64

	mu sync.RWMutex
}

type treeSnapshot struct {
	Model *core.Model
	Root  treeNode
}

// New starts a new hoeffding tree from a model
func New(model *core.Model, conf *Config) *Tree {
	root := newLeafNode(helpers.NewObservationStats(model.IsRegression()))
	return initTree(model, root, conf)
}

// Load loads a tree from a readable source with the given config
func Load(r io.Reader, conf *Config) (*Tree, error) {
	var t *Tree
	if err := gob.NewDecoder(r).Decode(&t); err != nil {
		return nil, err
	}
	return t, nil
}

func initTree(model *core.Model, root treeNode, conf *Config) *Tree {
	tree := &Tree{
		model:  model,
		root:   root,
		traces: make(chan *Trace, 3),
	}
	tree.SetConfig(conf)
	return tree
}

// SetConfig updates config on the fly
func (t *Tree) SetConfig(conf *Config) {
	if conf == nil {
		conf = new(Config)
	}
	conf.norm(t.model.IsRegression())

	t.mu.Lock()
	t.conf = conf
	t.mu.Unlock()
}

// ByteSize estimates the memory required to store the tree
func (t *Tree) ByteSize() int {
	t.mu.RLock()
	byteSize := t.root.ByteSize()
	t.mu.RUnlock()

	return byteSize
}

// Info returns information about the tree
func (t *Tree) Info() *TreeInfo {
	info := new(TreeInfo)

	t.mu.RLock()
	t.root.ReadInfo(1, info)
	t.mu.RUnlock()

	return info
}

// Traces allows users to subscribe to debug traces. When enabled,
// events (can be nil) will be emitted via this channel after each
// training cycle.
func (t *Tree) Traces() <-chan *Trace {
	return t.traces
}

// WriteGraph write a graph in dot notation to a writer
func (t *Tree) WriteGraph(w io.Writer) error {
	buf := bufio.NewWriter(w)
	defer buf.Flush()

	if _, err := buf.WriteString("digraph ht {\n  edge [arrowsize=0.6, fontsize=10];\n"); err != nil {
		return err
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	if err := t.root.WriteGraph(buf, "N"); err != nil {
		return err
	}
	if _, err := buf.WriteString("}\n"); err != nil {
		return err
	}

	return nil
}

// WriteText writes text-based tree output to a writer
func (t *Tree) WriteText(w io.Writer) error {
	buf := bufio.NewWriter(w)
	defer buf.Flush()

	if _, err := buf.WriteString("ROOT"); err != nil {
		return err
	}

	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.root.WriteText(buf, "\t")
}

// Train passes an instance to the tree for training purposes
func (t *Tree) Train(inst core.Instance) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var trace *Trace
	if t.conf.EnableTracing {
		defer func() { t.traces <- trace }()
	}

	node, parent, parentIndex := t.root.Filter(inst, nil, -1)
	if node == nil {
		node = newLeafNode(helpers.NewObservationStats(t.model.IsRegression()))
		parent.Children[parentIndex] = node
	}

	if leaf, ok := node.(*leafNode); ok {
		leaf.Learn(inst, t)

		if t.cycles++; t.cycles%int64(t.conf.PrunePeriod) == 0 {
			if byteSize := t.root.ByteSize(); byteSize >= t.conf.MemTarget*2 {
				t.prune(byteSize)
			}
		}

		weight := leaf.Stats.TotalWeight()
		if int(weight-leaf.WeightOnLastEval) < t.conf.GracePeriod {
			return
		}

		var split *splitNode
		if split, trace = t.attemptSplit(leaf, weight, trace); split != nil {
			if parent == nil {
				t.root = split
			} else {
				parent.SetChild(parentIndex, split)
			}
			if trace != nil {
				trace.Split = true
			}
		}
		if weight > leaf.WeightOnLastEval {
			leaf.WeightOnLastEval = weight
		}
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

// WriteTo writes the tree to a writer
func (t *Tree) WriteTo(w io.Writer) error {
	return gob.NewEncoder(w).Encode(t)
}

// GobEncode implements gob.GobEncoder
func (t *Tree) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(treeSnapshot{
		Model: t.model,
		Root:  t.root,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode implements gob.GobDecoder
func (t *Tree) GobDecode(b []byte) error {
	var snap treeSnapshot
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&snap); err != nil {
		return err
	}

	*t = *initTree(snap.Model, snap.Root, nil)
	return nil
}

func (t *Tree) attemptSplit(leaf *leafNode, weight float64, trace *Trace) (*splitNode, *Trace) {
	if !leaf.Stats.IsSufficient() || leaf.IsInactive {
		return nil, trace
	}

	if t.conf.EnableTracing {
		trace = new(Trace)
	}

	// Calculate best splits
	splits := leaf.BestSplits(t)
	bestSplit := splits[0]

	// Calculate the gain between merits of the best and the second-best split
	meritGain := bestSplit.Merit()
	if len(splits) > 1 {
		meritGain -= splits[1].Merit()
	}

	// Update trace
	if trace != nil {
		trace.MeritGain = meritGain
		trace.PossibleSplits = make([]TracePossibleSplit, 0, len(splits))

		for _, split := range splits {
			if cond := split.Condition(); cond != nil {
				trace.PossibleSplits = append(trace.PossibleSplits, TracePossibleSplit{
					Predictor: cond.Predictor().Name,
					Merit:     split.Merit(),
				})
			}
		}
	}

	// Don't split if there is no merit gain
	if meritGain <= 0 {
		return nil, trace
	}

	// Calculate hoeffding bound, evaluate split
	srange := bestSplit.Range()
	hbound := math.Sqrt(srange * srange * math.Log(1.0/t.conf.SplitConfidence) / (2.0 * weight))

	// Update trace
	if trace != nil {
		trace.HoeffdingBound = hbound
	}

	if meritGain > hbound || hbound < t.conf.TieThreshold {
		return newSplitNode(
			bestSplit.Condition(),
			bestSplit.PreStats(),
			bestSplit.PostStats(),
		), trace
	}
	return nil, trace
}

func (t *Tree) prune(byteSize int) {
	t.leaves = t.root.FindLeaves(t.leaves[:0])
	sort.Sort(sort.Reverse(t.leaves))

	target := t.conf.MemTarget
	piv := 0
	for ; piv < len(t.leaves); piv++ {
		if n := t.leaves[piv]; !n.IsInactive {
			byteSize -= n.ByteSize()
			n.Deactivate()

			if byteSize <= target {
				break
			}
		}
	}

	for ; piv < len(t.leaves); piv++ {
		t.leaves[piv].Activate()
	}
}
