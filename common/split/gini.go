package split

import (
	"github.com/bsm/reason"
	"github.com/bsm/reason/util"
)

// GiniImpurity determines split merit using Gini impurity
type GiniImpurity struct{}

// Supports implements Criterion.
func (GiniImpurity) Supports(target *reason.Feature) bool {
	return target != nil && target.Kind == reason.Feature_CATEGORICAL
}

// ClassificationRange implements Criterion.
func (GiniImpurity) ClassificationRange(_ *util.Vector) float64 {
	return 1.0
}

// ClassificationMerit implements Criterion.
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

// RegressionRange implements Criterion.
func (GiniImpurity) RegressionRange(_ *util.NumStream) float64 {
	return 0.0 // N/A
}

// RegressionMerit implements Criterion.
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
