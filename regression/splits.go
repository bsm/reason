package regression

import (
	"math"

	"github.com/bsm/reason/internal/splits"
	"github.com/bsm/reason/util"
)

// SplitCriterion calculates the merit of an attribute split
// for regressions
type SplitCriterion interface {
	// Range returns the range of the split merit.
	Range(pre *util.Vector) float64

	// Merit calculates the merit of splitting for a given
	// distribution before and after the split
	Merit(pre *util.Vector, post *util.Matrix) float64
}

// DefaultSplitCriterion returns the default split criterion:
//   VarReductionSplitCriterion{MinWeight: 4.0}
func DefaultSplitCriterion() SplitCriterion {
	return VarianceReduction{MinWeight: 4.0}
}

// VarianceReduction performs splits using variance-reduction
type VarianceReduction struct {
	// The minimum weight a post-split option requires
	// in order to be considered. Default: 4.0
	MinWeight float64
}

// Range implements SplitCriterion
func (VarianceReduction) Range(_ *util.Vector) float64 { return 1.0 }

// Merit implements SplitCriterion
func (c VarianceReduction) Merit(pre *util.Vector, post *util.Matrix) float64 {
	if pre == nil {
		return 0.0
	}

	postWeight := 0.0
	postCount := 0
	postStreams := util.WrapNumStreams(post)
	postStreams.ForEach(func(cat int) {
		if w := postStreams.TotalWeight(cat); w >= c.MinWeight {
			postWeight += w
			postCount++
		}
	})
	if postCount < 2 || postWeight == 0 {
		return 0.0
	}

	preStream := util.WrapNumStream(pre)
	preVar := preStream.Variance()
	if math.IsNaN(preVar) {
		return 0.0
	}

	postVar := 0.0
	postStreams.ForEach(func(cat int) {
		if w, v := postStreams.TotalWeight(cat), postStreams.Variance(cat); w >= c.MinWeight && !math.IsNaN(v) {
			postVar += w * v / postWeight
		}
	})

	return splits.NormMerit(preVar - postVar)
}

// --------------------------------------------------------------------

// GainRatio wraps a split criterion and normalises the merits
// by reducing their bias toward attributes that have a large number of
// values over attributes that have a smaller number of values.
type GainRatio struct{ SplitCriterion }

func (c GainRatio) Merit(pre *util.Vector, post *util.Matrix) float64 {
	postStreams := util.WrapNumStreams(post)
	penalty := new(splits.GainRatioPenalty)
	postStreams.ForEach(func(cat int) {
		penalty.Weight += postStreams.TotalWeight(cat)
	})
	postStreams.ForEach(func(cat int) {
		penalty.Update(postStreams.TotalWeight(cat))
	})

	merit := c.SplitCriterion.Merit(pre, post)
	return splits.NormMerit(merit / penalty.Value())
}
