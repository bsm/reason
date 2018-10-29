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

	fullWeight := 0.0
	catCount := 0
	postStats := WrapStatsDistribution(post)
	postStats.ForEach(func(cat int) {
		if w := postStats.TotalWeight(cat); w >= c.MinWeight {
			fullWeight += w
			catCount++
		}
	})
	if catCount < 2 || fullWeight == 0 {
		return 0.0
	}

	preStats := WrapStats(pre)
	preVar := preStats.Variance()
	if math.IsNaN(preVar) {
		return 0.0
	}

	postVar := 0.0
	postStats.ForEach(func(cat int) {
		if w, v := postStats.TotalWeight(cat), postStats.Variance(cat); w >= c.MinWeight && !math.IsNaN(v) {
			postVar += w * v / fullWeight
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
	merit := c.SplitCriterion.Merit(pre, post)
	postStats := WrapStatsDistribution(post)
	numCats := postStats.NumCategories()

	penalty := new(splits.GainRatioPenalty)
	for i := 0; i < numCats; i++ {
		penalty.Weight += postStats.TotalWeight(i)
	}
	for i := 0; i < numCats; i++ {
		penalty.Update(postStats.TotalWeight(i))
	}
	return splits.NormMerit(merit / penalty.Value())
}
