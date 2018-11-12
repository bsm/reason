package split

import (
	"math"

	"github.com/bsm/reason"
	"github.com/bsm/reason/util"
)

// InformationGain determines split merit through information gain.
type InformationGain struct {
	// Requires at least two potential branches to have a sufficient
	// fractional weight in order to qualify for a split.
	MinBranchFraction float64
}

// Supports implements Criterion.
func (InformationGain) Supports(target *reason.Feature) bool {
	return target != nil && target.Kind == reason.Feature_CATEGORICAL
}

// ClassificationRange implements Criterion.
func (InformationGain) ClassificationRange(pre *util.Vector) float64 {
	if pre != nil {
		if sz := pre.NNZ(); sz > 2 {
			return math.Log2(float64(sz))
		}
	}
	return 1.0
}

// ClassificationMerit implements Criterion.
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

// RegressionRange implements Criterion.
func (InformationGain) RegressionRange(_ *util.NumStream) float64 {
	return 0.0 // N/A
}

// RegressionMerit implements Criterion.
func (InformationGain) RegressionMerit(_ *util.NumStream, _ *util.NumStreams) float64 {
	return 0.0 // N/A
}
