// Package hoeffding implements a Hoeffding tree classifier and regressor.
package hoeffding

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

type classificationResult struct {
	cat    core.Category
	weight float64
	vv     *util.Vector
}

func (c classificationResult) Category() core.Category { return c.cat }
func (c classificationResult) Weight() float64         { return c.weight }
func (c classificationResult) Prob(cat core.Category) float64 {
	if c.weight <= 0 {
		return 0
	}
	return c.vv.At(int(cat)) / c.weight
}

type regressionResult struct {
	ns *util.NumStream
}

func (r regressionResult) Number() float64 { return r.ns.Mean() }
func (r regressionResult) MSE() float64    { return r.ns.Variance() }
func (r regressionResult) Weight() float64 { return r.ns.Weight }
