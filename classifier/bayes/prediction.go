package bayes

import (
	"github.com/bsm/reason"
	"github.com/bsm/reason/util"
)

type prediction struct {
	cat reason.Category
	vv  *util.Vector
}

func (p prediction) Category() reason.Category        { return p.cat }
func (p prediction) Prob(cat reason.Category) float64 { return p.vv.At(int(p.cat)) }
