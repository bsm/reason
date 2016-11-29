package classifiers

import (
	"math"

	"github.com/bsm/reason/util"
)

var (
	_ CSplitCriterion = GiniSplitCriterion{}
	_ CSplitCriterion = InfoGainSplitCriterion{}

	_ RSplitCriterion = VarReductionSplitCriterion{}
)

// SplitCriterion calculates merits of attribute splits
type SplitCriterion interface {
	isSplitCriterion()
}

// CSplitCriterion calculates the merit of an attribute split
// for classifications
type CSplitCriterion interface {
	SplitCriterion

	// Range returns the range of the split merit
	Range(pre util.Vector) float64

	// Merit calculates the merit of splitting for a given
	// distribution before and after the split
	Merit(pre util.Vector, post util.VectorDistribution) float64
}

// DefaultSplitCriterion returns InfoGainSplitCriterion
// for classification or VarReductionSplitCriterion for regressions
// (with a MinBranchFrac or 0.1)
func DefaultSplitCriterion(isRegression bool) SplitCriterion {
	if isRegression {
		return VarReductionSplitCriterion{MinWeight: 5.0}
	}
	return InfoGainSplitCriterion{MinBranchFrac: 0.1}
}

// --------------------------------------------------------------------

// RSplitCriterion calculates the merit of an attribute split
// for regressions
type RSplitCriterion interface {
	SplitCriterion

	// Range returns the range of the split merit
	Range(pre *util.NumSeries) float64

	// Merit calculates the merit of splitting for a given
	// distribution before and after the split
	Merit(pre *util.NumSeries, post util.NumSeriesDistribution) float64
}

// --------------------------------------------------------------------

// GiniSplitCriterion determines split merit using Gini Impurity
type GiniSplitCriterion struct{}

func (GiniSplitCriterion) isSplitCriterion() {}

func (GiniSplitCriterion) Range(pre util.Vector) float64 { return 1.0 }
func (GiniSplitCriterion) Merit(pre util.Vector, post util.VectorDistribution) float64 {
	weights := post.Weights()
	total := 0.0
	for _, w := range weights {
		total += w
	}
	if total == 0 {
		return 0.0
	}

	merit := 0.0
	for i, w := range weights {
		merit += w / total * giniSplitCalc(post[i], w)
	}
	return merit
}

func giniSplitCalc(vv util.Vector, sum float64) float64 {
	res := 1.0
	vv.ForEachValue(func(v float64) {
		sub := v / sum
		res -= sub * sub
	})
	return res
}

// InfoGainSplitCriterion determines split merit through information gain
type InfoGainSplitCriterion struct {
	// Requires at least two potential branches to have a sufficient
	// fractional weight in order to qualify for a split.
	MinBranchFrac float64
}

func (InfoGainSplitCriterion) isSplitCriterion() {}

func (InfoGainSplitCriterion) Range(pre util.Vector) float64 {
	if pre != nil {
		if cnt := pre.Count(); cnt > 2 {
			return math.Log2(float64(cnt))
		}
	}
	return math.Log2(2.0)
}

func (c InfoGainSplitCriterion) Merit(pre util.Vector, post util.VectorDistribution) float64 {
	weights := post.Weights()
	total := 0.0
	for _, w := range weights {
		total += w
	}
	if total == 0 {
		return 0.0
	}

	if min := c.MinBranchFrac; min > 0 {
		n := 0
		for _, w := range weights {
			if w/total > min {
				if n++; n > 1 {
					break
				}
			}
		}
		if n < 2 {
			return 0.0
		}
	}

	e1, e2 := pre.Entropy(), 0.0
	for i, vv := range post {
		e2 += weights[i] * vv.Entropy()
	}
	return e1 - e2/total
}

// VarReductionSplitCriterion performs splits using variance-reduction
type VarReductionSplitCriterion struct {
	// The minimum weight a post-split option requires
	// in order to be considered. Default: 5.0
	MinWeight float64
}

func (VarReductionSplitCriterion) isSplitCriterion() {}

func (VarReductionSplitCriterion) Range(_ *util.NumSeries) float64 { return 1.0 }
func (c VarReductionSplitCriterion) Merit(pre *util.NumSeries, post util.NumSeriesDistribution) float64 {
	if pre == nil {
		return 0.0
	}

	count := 0
	total := 0.0
	for _, n := range post {
		if w := n.TotalWeight(); w >= c.MinWeight {
			total += w
			count++
		}
	}
	if count < 2 || total == 0 {
		return 0.0
	}

	postV := 0.0
	for _, n := range post {
		if w := n.TotalWeight(); w >= c.MinWeight {
			postV += w / total * n.Variance()
		}
	}
	return pre.Variance() - postV
}

// --------------------------------------------------------------------

// GainRatioSplitCriterion normalises the merits of other split criterions
// by reducing their bias toward attributes that have a large number of
// values over attributes that have a smaller number of values.
func GainRatioSplitCriterion(c SplitCriterion) SplitCriterion {
	switch x := c.(type) {
	case CSplitCriterion:
		return cGainRatioSplitCriterion{CSplitCriterion: x}
	case RSplitCriterion:
		return rGainRatioSplitCriterion{RSplitCriterion: x}
	}
	return c
}

type cGainRatioSplitCriterion struct{ CSplitCriterion }

func (c cGainRatioSplitCriterion) Merit(pre util.Vector, post util.VectorDistribution) float64 {
	merit := c.CSplitCriterion.Merit(pre, post)
	return merit / gainRatioSplitInfo(post)
}

type rGainRatioSplitCriterion struct{ RSplitCriterion }

func (c rGainRatioSplitCriterion) Merit(pre *util.NumSeries, post util.NumSeriesDistribution) float64 {
	merit := c.RSplitCriterion.Merit(pre, post)
	return merit / gainRatioSplitInfo(post)
}

func gainRatioSplitInfo(post interface {
	Weights() map[int]float64
}) float64 {
	weights := post.Weights()
	total := 0.0
	for _, w := range weights {
		total += w
	}

	pen := 0.0
	for _, w := range weights {
		rat := w / total
		pen -= rat * math.Log2(rat)
	}
	if pen <= 0.0 {
		return 1.0
	}
	return pen
}
