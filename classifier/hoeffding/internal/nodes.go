package internal

import (
	"github.com/bsm/reason"
	"github.com/bsm/reason/common/split"
)

// GetChild retrieves the child nref at pos.
func (n *SplitNode) GetChild(pos int) int64 {
	if pos < len(n.Children) {
		return n.Children[pos]
	}
	return 0
}

// SetChild sets references the child at pos to nref.
func (n *SplitNode) SetChild(pos int, nref int64) {
	if sz := pos + 1; sz > cap(n.Children) {
		children := make([]int64, sz, 2*sz)
		copy(children, n.Children)
		n.Children = children
	} else if sz > len(n.Children) {
		n.Children = n.Children[:sz]
	}
	n.Children[pos] = nref
}

func (n *SplitNode) childPos(feat *reason.Feature, x reason.Example) int {
	switch feat.Kind {
	case reason.Feature_CATEGORICAL:
		return int(feat.Category(x))
	case reason.Feature_NUMERICAL:
		if feat.Number(x) < n.Pivot {
			return 0
		}
		return 1
	default:
		return -1
	}
}

// --------------------------------------------------------------------

// Enable enables the node.
func (n *LeafNode) Enable() {
	if n.IsDisabled {
		n.IsDisabled = false
		n.FeatureStats = make(map[string]*LeafNode_Stats)
	}
}

// Disable disables the node.
func (n *LeafNode) Disable() {
	n.IsDisabled = true
	n.FeatureStats = nil
}

// ObserveExample observes an example and updates internal stats.
func (n *LeafNode) ObserveExample(m *reason.Model, target *reason.Feature, x reason.Example, weight float64, me *Node) {
	// Observe example, update node stats
	if success := me.incrementStats(target, x, weight); !success {
		return
	}

	// Skip the remaining steps if this node is disabled
	if n.IsDisabled {
		return
	}

	// Ensure we have stats
	if n.FeatureStats == nil {
		n.FeatureStats = make(map[string]*LeafNode_Stats, len(m.Features)-1)
	}

	// Update each predictor feature's stats with a target-value, predictor-value
	// and weight tuple
	for name, predictor := range m.Features {
		if name == target.Name {
			continue // skip target, we are only interested in predictors
		}

		stats := n.FeatureStats[predictor.Name]
		if stats == nil {
			stats = new(LeafNode_Stats)
			n.FeatureStats[predictor.Name] = stats
		}
		stats.Update(target, predictor, x, weight)
	}
}

// EvaluateSplit evaluates a split for a fiven feature.
// Returns nil if a split is not possible.
func (n *LeafNode) EvaluateSplit(feature string, crit split.Criterion, node *Node) *SplitCandidate {
	if n.IsDisabled {
		return nil
	}

	if n.FeatureStats == nil {
		return nil
	}

	stats, ok := n.FeatureStats[feature]
	if !ok {
		return nil
	}

	var sc *SplitCandidate
	switch kind := stats.GetKind().(type) {
	case *LeafNode_Stats_CN:
		if nstat := node.GetClassification(); nstat != nil {
			pre := &nstat.Vector
			if merit, pivot, post := kind.CN.EvaluateSplit(crit, pre); merit > 0 {
				sc = &SplitCandidate{
					Feature:   feature,
					Range:     crit.ClassificationRange(pre),
					Merit:     merit,
					Pivot:     pivot,
					PostSplit: PostSplit{Classification: post},
				}
			}
		}
	case *LeafNode_Stats_CC:
		if nstat := node.GetClassification(); nstat != nil {
			pre := &nstat.Vector
			if merit, post := kind.CC.EvaluateSplit(crit, pre); merit > 0 {
				sc = &SplitCandidate{
					Feature:   feature,
					Range:     crit.ClassificationRange(pre),
					Merit:     merit,
					PostSplit: PostSplit{Classification: post},
				}
			}
		}
	case *LeafNode_Stats_RN:
		if nstat := node.GetRegression(); nstat != nil {
			pre := &nstat.NumStream
			if merit, pivot, post := kind.RN.EvaluateSplit(crit, pre); merit > 0 {
				sc = &SplitCandidate{
					Feature:   feature,
					Range:     crit.RegressionRange(pre),
					Merit:     merit,
					Pivot:     pivot,
					PostSplit: PostSplit{Regression: post},
				}
			}
		}
	case *LeafNode_Stats_RC:
		if nstat := node.GetRegression(); nstat != nil {
			pre := &nstat.NumStream
			if merit, post := kind.RC.EvaluateSplit(crit, pre); merit > 0 {
				sc = &SplitCandidate{
					Feature:   feature,
					Range:     crit.RegressionRange(pre),
					Merit:     merit,
					PostSplit: PostSplit{Regression: post},
				}
			}
		}
	}
	return sc
}

// --------------------------------------------------------------------

// Weight returns the weight observed at the node.
func (n *Node) Weight() float64 {
	switch kind := n.GetStats().(type) {
	case *Node_Classification:
		return kind.Classification.WeightSum()
	case *Node_Regression:
		return kind.Regression.Weight
	}
	return 0.0
}

// IsSufficient returns true when a node has sufficient stats.
func (n *Node) IsSufficient() bool {
	switch kind := n.GetStats().(type) {
	case *Node_Classification:
		return kind.Classification.NNZ() > 1
	case *Node_Regression:
		return kind.Regression.Weight > 1
	}
	return false
}

func (n *Node) incrementStats(target *reason.Feature, x reason.Example, weight float64) (success bool) {
	switch target.Kind {
	case reason.Feature_CATEGORICAL:
		if cat := target.Category(x); reason.IsCat(cat) {
			stats := n.GetClassification()
			if stats == nil {
				stats = new(Node_ClassificationStats)
				n.Stats = &Node_Classification{Classification: stats}
			}
			stats.Incr(int(cat), weight)
			return true
		}
	case reason.Feature_NUMERICAL:
		if num := target.Number(x); reason.IsNum(num) {
			stats := n.GetRegression()
			if stats == nil {
				stats = new(Node_RegressionStats)
				n.Stats = &Node_Regression{Regression: stats}
			}
			stats.ObserveWeight(num, weight)
			return true
		}
	}
	return false
}
