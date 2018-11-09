// Package hoeffding implements a Hoeffding tree classifier and regressor.
package hoeffding

import (
	"math"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

// Classification prediction implements classifier.Classification and exposes additional stats.
type Classification struct {
	cat    core.Category
	weight float64
	vv     *util.Vector
}

// Weight returns the observed weight that has been used to make this prediction.
func (c Classification) Weight() float64 { return c.weight }

// Category implements classifier.Classification interface.
func (c Classification) Category() core.Category { return c.cat }

// Prob implements classifier.Classification interface.
func (c Classification) Prob(cat core.Category) float64 {
	if c.weight <= 0 || c.vv == nil {
		return 0
	}
	return c.vv.At(int(cat)) / c.weight
}

// Regression prediction implements classifier.Regression and exposes additional stats.
type Regression struct {
	ns *util.NumStream
}

// Weight returns the observed weight that has been used to make this prediction.
func (r Regression) Weight() float64 {
	if r.ns == nil {
		return 0
	}
	return r.ns.Weight
}

// MSE returns the variance of previous observation that have been used to make this prediction.
func (r Regression) MSE() float64 {
	if r.ns == nil {
		return math.NaN()
	}
	return r.ns.Variance()
}

// Number implements classifier.Regression interface.
func (r Regression) Number() float64 {
	if r.ns == nil {
		return math.NaN()
	}
	return r.ns.Mean()
}
