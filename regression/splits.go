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
	Range(pre *util.NumStream) float64

	// Merit calculates the merit of splitting for a given
	// distribution before and after the split
	Merit(pre *util.NumStream, post *util.NumStreams) float64
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
func (VarianceReduction) Range(_ *util.NumStream) float64 { return 1.0 }

// Merit implements SplitCriterion
func (c VarianceReduction) Merit(pre *util.NumStream, post *util.NumStreams) float64 {
	if pre == nil || post == nil {
		return 0.0
	}

	relevantW := 0.0
	relevantN := 0
	rows := post.NumRows()
	for i := 0; i < rows; i++ {
		if s := post.At(i); s != nil && s.Weight >= c.MinWeight {
			relevantW += s.Weight
			relevantN++
		}
	}

	if relevantN < 2 || relevantW <= 0 {
		return 0.0
	}

	preVar := pre.Variance()
	if math.IsNaN(preVar) {
		return 0.0
	}

	postVar := 0.0
	for i := 0; i < rows; i++ {
		if s := post.At(i); s != nil && s.Weight >= c.MinWeight {
			if v := s.Variance(); !math.IsNaN(v) {
				postVar += s.Weight * v / relevantW
			}
		}
	}
	return splits.NormMerit(preVar - postVar)
}

// --------------------------------------------------------------------

// GainRatio wraps a split criterion and normalises the merits
// by reducing their bias toward attributes that have a large number of
// values over attributes that have a smaller number of values.
type GainRatio struct{ SplitCriterion }

func (c GainRatio) Merit(pre *util.NumStream, post *util.NumStreams) float64 {
	penalty := new(splits.GainRatioPenalty)
	penalty.Weight = post.WeightSum()
	if penalty.Weight == 0 {
		return 0.0
	}

	rows := post.NumRows()
	for i := 0; i < rows; i++ {
		if t := post.At(i); t != nil {
			penalty.Update(t.Weight)
		}
	}
	merit := c.SplitCriterion.Merit(pre, post)
	return splits.NormMerit(merit / penalty.Value())
}
