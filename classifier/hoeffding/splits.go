package hoeffding

import (
	"math"
	"strconv"

	"github.com/bsm/reason/classifier"
	"github.com/bsm/reason/util"
)

// SplitCriterion calculates the merit of an attribute split
// for classifications.
type SplitCriterion interface {
	// Supports returns true if a classifier Goal is supported.
	Supports(classifier.Problem) bool

	// Range returns the range of the split merit
	Range(pre *util.Vector) float64

	// Merit calculates the merit of splitting for a given
	// distribution before and after the split
	Merit(pre *util.Vector, post *util.Matrix) float64
}

// DefaultSplitCriterion returns the default split criterion for the given problem.
func DefaultSplitCriterion(p classifier.Problem) SplitCriterion {
	switch p {
	case classifier.Classification:
		return InformationGain{MinBranchFraction: 0.1}
	case classifier.Regression:
		return VarianceReduction{MinWeight: 4.0}
	default:
		panic("reason: invalid problem " + strconv.FormatUint(uint64(p), 10))
	}
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
func (GiniImpurity) Supports(p classifier.Problem) bool {
	return p == classifier.Classification
}

// Range implements SplitCriterion
func (GiniImpurity) Range(_ *util.Vector) float64 { return 1.0 }

// Merit implements SplitCriterion
func (GiniImpurity) Merit(_ *util.Vector, post *util.Matrix) float64 {
	if post == nil {
		return 0.0
	}

	total := post.Sum()
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

func calcGiniSplit(row []float64, sum float64) float64 {
	res := 1.0
	for _, v := range row {
		sub := v / sum
		res -= sub * sub
	}
	return res
}

// InformationGain determines split merit through information gain.
type InformationGain struct {
	// Requires at least two potential branches to have a sufficient
	// fractional weight in order to qualify for a split.
	MinBranchFraction float64
}

// Supports implements SplitCriterion
func (InformationGain) Supports(p classifier.Problem) bool {
	return p == classifier.Classification
}

// Range implements SplitCriterion
func (InformationGain) Range(pre *util.Vector) float64 {
	if pre != nil {
		if sz := pre.NNZ(); sz > 2 {
			return math.Log2(float64(sz))
		}
	}
	return math.Log2(2.0)
}

// Merit implements SplitCriterion
func (c InformationGain) Merit(pre *util.Vector, post *util.Matrix) float64 {
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
		if w := vv.Weight(); w > 0 {
			e2 += w * vv.Entropy()
		}
	}
	return normSplitMerit(e1 - e2/total)
}

// --------------------------------------------------------------------

// VarianceReduction performs regression splits using variance-reduction.
type VarianceReduction struct {
	// The minimum weight a post-split option requires
	// in order to be considered. Default: 4.0
	MinWeight float64
}

// Supports implements SplitCriterion
func (VarianceReduction) Supports(p classifier.Problem) bool {
	return p == classifier.Regression
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

	return normSplitMerit(preVar - postVar)
}

// --------------------------------------------------------------------

// GainRatio normalises the merits of other split criterions
// by reducing their bias toward attributes that have a large number of
// values over attributes that have a smaller number of values.
type GainRatio struct{ SplitCriterion }

// Merit implements SplitCriterion
func (c GainRatio) Merit(pre *util.Vector, post *util.Matrix) float64 {
	merit := c.SplitCriterion.Merit(pre, post)
	if merit == 0 {
		return merit
	}

	total := 0.0
	if c.Supports(classifier.Classification) {
		total = post.Sum()
	} else if c.Supports(classifier.Regression) {
		total = post.ColSum(0)
	}
	if total == 0 {
		return merit
	}

	rows := post.NumRows()
	frac := 0.0
	for i := 0; i < rows; i++ {
		weight := 0.0
		if c.Supports(classifier.Classification) {
			weight = post.RowSum(i)
		} else if c.Supports(classifier.Regression) {
			weight = post.At(i, 0)
		}

		if weight > 0 {
			rat := weight / total
			frac -= rat * math.Log2(rat)
		}
	}
	return normSplitMerit(merit / frac)
}
