// Package hoeffding implements a Hoeffding tree classifier and regressor.
package hoeffding

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

type classification struct {
	cat    core.Category
	weight float64
	vv     *util.Vector
}

func (c classification) Category() core.Category { return c.cat }
func (c classification) Weight() float64         { return c.weight }
func (c classification) Prob(cat core.Category) float64 {
	if c.weight <= 0 {
		return 0
	}
	return c.vv.At(int(cat)) / c.weight
}

type regression struct {
	ns *util.NumStream
}

func (r regression) Number() float64 { return r.ns.Mean() }
func (r regression) MSE() float64    { return r.ns.Variance() }
func (r regression) Weight() float64 { return r.ns.Weight }
