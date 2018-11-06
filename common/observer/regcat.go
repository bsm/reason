package observer

import (
	"github.com/bsm/reason/common/split"
	"github.com/bsm/reason/core"
	util "github.com/bsm/reason/util"
)

// NewRegressionCategorical inits a regression observer for categorical features.
func NewRegressionCategorical() *RegressionCategorical {
	return new(RegressionCategorical)
}

// Observe adds a new observation.
func (o *RegressionCategorical) Observe(cat core.Category, target float64) {
	o.ObserveWeight(cat, target, 1.0)
}

// ObserveWeight updates stats based on a weighted observation.
func (o *RegressionCategorical) ObserveWeight(cat core.Category, target float64, weight float64) {
	if core.IsCat(cat) && core.IsNum(target) && weight > 0 {
		o.Dist.ObserveWeight(int(cat), target, weight)
	}
}

// EvaluateSplit evaluates a split.
func (o *RegressionCategorical) EvaluateSplit(crit split.Criterion, pre *util.NumStream) (merit float64, post *util.NumStreams) {
	if n := o.Dist.NumCategories(); n > 1 {
		post = &o.Dist
		merit = crit.RegressionMerit(pre, post)
	}
	return
}
