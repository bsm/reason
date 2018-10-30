package classification

import (
	"math"

	"github.com/bsm/reason/internal/splits"
	"github.com/bsm/reason/util"
)

// SplitCriterion calculates the merit of an attribute split
// for classifications.
type SplitCriterion interface {
	// Range returns the range of the split merit
	Range(pre *util.Vector) float64

	// Merit calculates the merit of splitting for a given
	// distribution before and after the split
	Merit(pre *util.Vector, post *util.Matrix) float64
}

// DefaultSplitCriterion returns the default split criterion:
//   InformationGain{MinBranchFraction: 0.1}
func DefaultSplitCriterion() SplitCriterion {
	return InformationGain{MinBranchFraction: 0.1}
}

// --------------------------------------------------------------------

// GiniImpurity determines split merit using Gini impurity
type GiniImpurity struct{}

// Range implements SplitCriterion
func (GiniImpurity) Range(_ *util.Vector) float64 { return 1.0 }

// Merit implements SplitCriterion
func (GiniImpurity) Merit(pre *util.Vector, post *util.Matrix) float64 {
	if pre == nil || post == nil {
		return 0.0
	}

	total := post.Sum()
	if total == 0 {
		return 0.0
	}

	merit := 0.0
	rows, _ := post.Dims()
	for i := 0; i < rows; i++ {
		if sum := post.RowSum(i); sum > 0 {
			merit += sum / total * calcGiniSplit(post.Row(i), sum)
		}
	}
	return splits.NormMerit(merit)
}

func calcGiniSplit(row []float64, sum float64) float64 {
	res := 1.0
	for _, v := range row {
		sub := v / sum
		res -= sub * sub
	}
	return res
}

// InformationGain determines split merit through information gain
type InformationGain struct {
	// Requires at least two potential branches to have a sufficient
	// fractional weight in order to qualify for a split.
	MinBranchFraction float64
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
	rows, _ := post.Dims()
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
	return splits.NormMerit(e1 - e2/total)
}

// --------------------------------------------------------------------

// GainRatio normalises the merits of other split criterions
// by reducing their bias toward attributes that have a large number of
// values over attributes that have a smaller number of values.
type GainRatio struct{ SplitCriterion }

// Merit implements SplitCriterion
func (c GainRatio) Merit(pre *util.Vector, post *util.Matrix) float64 {
	rows, _ := post.Dims()
	penalty := new(splits.GainRatioPenalty)
	for i := 0; i < rows; i++ {
		penalty.Weight += post.RowSum(i)
	}
	for i := 0; i < rows; i++ {
		penalty.Update(post.RowSum(i))
	}

	merit := c.SplitCriterion.Merit(pre, post)
	return splits.NormMerit(merit / penalty.Value())
}
