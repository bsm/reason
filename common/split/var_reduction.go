package split

import (
	"math"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

// VarianceReduction performs regression splits using variance-reduction.
type VarianceReduction struct {
	// The minimum weight a post-split option requires
	// in order to be considered. Default: 4.0
	MinWeight float64
}

// Supports implements Criterion.
func (VarianceReduction) Supports(target *core.Feature) bool {
	return target != nil && target.Kind == core.Feature_NUMERICAL
}

// ClassificationRange implements Criterion.
func (VarianceReduction) ClassificationRange(_ *util.Vector) float64 {
	return 0.0 // N/A
}

// ClassificationMerit implements Criterion.
func (VarianceReduction) ClassificationMerit(_ *util.Vector, _ *util.Matrix) float64 {
	return 0.0 // N/A
}

// RegressionRange implements Criterion.
func (VarianceReduction) RegressionRange(pre *util.NumStream) float64 {
	return 1.0
}

// RegressionMerit implements Criterion.
func (c VarianceReduction) RegressionMerit(pre *util.NumStream, post *util.NumStreams) float64 {
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
	return normSplitMerit(preVar - postVar)
}
