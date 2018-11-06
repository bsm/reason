package split

import (
	"math"

	"github.com/bsm/reason/util"
)

// GainRatio normalises the merits of other split criterions
// by reducing their bias toward attributes that have a large number of
// values over attributes that have a smaller number of values.
type GainRatio struct{ Criterion }

// ClassificationMerit implements Criterion.
func (c GainRatio) ClassificationMerit(pre *util.Vector, post *util.Matrix) float64 {
	merit := c.Criterion.ClassificationMerit(pre, post)
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

// RegressionMerit implements Criterion.
func (c GainRatio) RegressionMerit(pre *util.NumStream, post *util.NumStreams) float64 {
	merit := c.Criterion.RegressionMerit(pre, post)
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
