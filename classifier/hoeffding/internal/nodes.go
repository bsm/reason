package internal

import (
	core "github.com/bsm/reason/core"
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

func (n *SplitNode) childPos(feat *core.Feature, x core.Example) int {
	switch feat.Kind {
	case core.Feature_CATEGORICAL:
		return int(feat.Category(x))
	case core.Feature_NUMERICAL:
		if feat.Number(x) < n.Pivot {
			return 0
		}
		return 1
	default:
		return -1
	}
}

// --------------------------------------------------------------------

// ObserveExample observes an example and updates internal stats.
func (n *LeafNode) ObserveExample(m *core.Model, target *core.Feature, x core.Example, weight float64, me *Node) {
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

// --------------------------------------------------------------------

func (n *Node) incrementStats(target *core.Feature, x core.Example, weight float64) (success bool) {
	switch target.Kind {
	case core.Feature_CATEGORICAL:
		if cat := target.Category(x); core.IsCat(cat) {
			stats := n.GetClassification()
			if stats == nil {
				stats = new(Node_ClassificationStats)
				n.Stats = &Node_Classification{Classification: stats}
			}
			stats.Add(int(cat), weight)
			return true
		}
	case core.Feature_NUMERICAL:
		if num := target.Number(x); core.IsNum(num) {
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
