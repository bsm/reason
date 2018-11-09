package bayes

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

type prediction struct {
	cat core.Category
	vv  *util.Vector
}

func (p prediction) Category() core.Category        { return p.cat }
func (p prediction) Prob(cat core.Category) float64 { return p.vv.At(int(p.cat)) }
