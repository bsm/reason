package hoeffding

import (
	"math"

	"github.com/bsm/reason/core"

	"github.com/bsm/reason/util"
)

// SplitCriterion calculates the merit of an attribute split
// for classifications.
type SplitCriterion interface {
	// Supports returns true if the target feature is supported.
	Supports(target *core.Feature) bool

	// ClassificationMerit applies to classifications only and
	// calculates the merit of splitting distribution pre and post split.
	ClassificationMerit(pre *util.Vector, post *util.Matrix) float64

	// RegressionMerit applies to regresssions only and
	// calculates the merit of splitting distribution pre and post split.
	RegressionMerit(pre *util.NumStream, post *util.NumStreams) float64
}

// DefaultSplitCriterion returns the default split criterion for the given target.
func DefaultSplitCriterion(target *core.Feature) SplitCriterion {
	switch target.Kind {
	case core.Feature_NUMERICAL:
		return VarianceReduction{MinWeight: 4.0}
	}
	return InformationGain{MinBranchFraction: 0.1}
}

func normSplitMerit(m float64) float64 {
	if m < 0.0 || math.IsNaN(m) {
		return 0.0
	}
	return m
}

// --------------------------------------------------------------------

// GiniImpurity determines split merit using Gini impurity
type GiniImpurity struct{}

// Supports implements SplitCriterion
func (GiniImpurity) Supports(target *core.Feature) bool {
	return target != nil && target.Kind == core.Feature_CATEGORICAL
}

// ClassificationMerit implements SplitCriterion
func (GiniImpurity) ClassificationMerit(_ *util.Vector, post *util.Matrix) float64 {
	if post == nil {
		return 0.0
	}

	total := post.WeightSum()
	if total == 0 {
		return 0.0
	}

	merit := 0.0
	rows := post.NumRows()
	for i := 0; i < rows; i++ {
		if sum := post.RowSum(i); sum > 0 {
			merit += sum / total * calcGiniSplit(post.Row(i), sum)
		}
	}
	return normSplitMerit(merit)
}

// RegressionMerit implements SplitCriterion.
func (GiniImpurity) RegressionMerit(_ *util.NumStream, _ *util.NumStreams) float64 {
	return 0.0 // N/A
}

func calcGiniSplit(row []float64, sum float64) float64 {
	res := 1.0
	for _, v := range row {
		sub := v / sum
		res -= sub * sub
	}
	return res
}

// --------------------------------------------------------------------

// InformationGain determines split merit through information gain.
type InformationGain struct {
	// Requires at least two potential branches to have a sufficient
	// fractional weight in order to qualify for a split.
	MinBranchFraction float64
}

// Supports implements SplitCriterion
func (InformationGain) Supports(target *core.Feature) bool {
	return target != nil && target.Kind == core.Feature_CATEGORICAL
}

// ClassificationMerit implements SplitCriterion
func (c InformationGain) ClassificationMerit(pre *util.Vector, post *util.Matrix) float64 {
	if pre == nil || post == nil {
		return 0.0
	}

	total := 0.0
	count := 0
	rows := post.NumRows()
	for i := 0; i < rows; i++ {
		if sum := post.RowSum(i); sum > 0 {
			total += sum
			count++
		}
	}
	if count < 2 || total == 0 {
		return 0.0
	}

	count = 0
	if min := c.MinBranchFraction; min > 0 {
		for i := 0; i < rows; i++ {
			if sum := post.RowSum(i); sum > 0 && sum/total > min {
				if count++; count > 1 {
					break
				}
			}
		}
	}
	if count < 2 {
		return 0.0
	}

	e1, e2 := pre.Entropy(), 0.0
	for i := 0; i < rows; i++ {
		vv := util.NewVectorFromSlice(post.Row(i)...)
		if w := vv.WeightSum(); w > 0 {
			e2 += w * vv.Entropy()
		}
	}
	return normSplitMerit(e1 - e2/total)
}

// RegressionMerit implements SplitCriterion.
func (InformationGain) RegressionMerit(_ *util.NumStream, _ *util.NumStreams) float64 {
	return 0.0 // N/A
}

// --------------------------------------------------------------------

// VarianceReduction performs regression splits using variance-reduction.
type VarianceReduction struct {
	// The minimum weight a post-split option requires
	// in order to be considered. Default: 4.0
	MinWeight float64
}

// Supports implements SplitCriterion
func (VarianceReduction) Supports(target *core.Feature) bool {
	return target != nil && target.Kind == core.Feature_NUMERICAL
}

// ClassificationMerit implements SplitCriterion.
func (VarianceReduction) ClassificationMerit(_ *util.Vector, _ *util.Matrix) float64 {
	return 0.0 // N/A
}

// RegressionMerit implements SplitCriterion.
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

// --------------------------------------------------------------------

// GainRatio normalises the merits of other split criterions
// by reducing their bias toward attributes that have a large number of
// values over attributes that have a smaller number of values.
type GainRatio struct{ SplitCriterion }

// ClassificationMerit implements SplitCriterion.
func (c GainRatio) ClassificationMerit(pre *util.Vector, post *util.Matrix) float64 {
	merit := c.SplitCriterion.ClassificationMerit(pre, post)
	if merit == 0 {
		return 0.0
	}

	weightSum := post.WeightSum()
	if weightSum == 0 {
		return 0.0
	}

	frac := 0.0
	rows := post.NumRows()
	for i := 0; i < rows; i++ {
		if w := post.RowSum(i); w > 0 {
			rat := w / weightSum
			frac -= rat * math.Log2(rat)
		}
	}
	return normSplitMerit(merit / frac)
}

// RegressionMerit implements SplitCriterion.
func (c GainRatio) RegressionMerit(pre *util.NumStream, post *util.NumStreams) float64 {
	merit := c.SplitCriterion.RegressionMerit(pre, post)
	if merit == 0 {
		return 0.0
	}

	weightSum := post.WeightSum()
	if weightSum == 0 {
		return 0.0
	}

	frac := 0.0
	rows := post.NumRows()
	for i := 0; i < rows; i++ {
		if t := post.At(i); t != nil {
			rat := t.Weight / weightSum
			frac -= rat * math.Log2(rat)
		}
	}
	return normSplitMerit(merit / frac)
}
