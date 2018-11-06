package observer

import (
	"github.com/bsm/reason/common/split"
	"github.com/bsm/reason/core"
	util "github.com/bsm/reason/util"
)

// NewClassificationCategorical inits a classification observer for categorical features.
func NewClassificationCategorical() *ClassificationCategorical {
	return new(ClassificationCategorical)
}

// Observe adds a new observation.
func (o *ClassificationCategorical) Observe(cat, target core.Category) {
	o.ObserveWeight(cat, target, 1.0)
}

// ObserveWeight updates stats based on a weighted observation.
func (o *ClassificationCategorical) ObserveWeight(cat, target core.Category, weight float64) {
	if core.IsCat(cat) && core.IsCat(target) && weight > 0 {
		o.Dist.Incr(int(cat), int(target), weight)
	}
}

// EvaluateSplit evaluates a split.
func (o *ClassificationCategorical) EvaluateSplit(crit split.Criterion, pre *util.Vector) (merit float64, post *util.Matrix) {
	if o.numCategories() > 1 {
		post = &o.Dist
		merit = crit.ClassificationMerit(pre, post)
	}
	return
}

// numCategories returns the number of categories of the observed feature.
func (o *ClassificationCategorical) numCategories() int {
	rows := o.Dist.NumRows()
	n := 0
	for i := 0; i < rows; i++ {
		if !o.Dist.IsRowZero(i) {
			n++
		}
	}
	return n
}
