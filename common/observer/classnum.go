package observer

import (
	"github.com/bsm/reason/common/split"
	"github.com/bsm/reason/core"
	util "github.com/bsm/reason/util"
)

// NewClassificationNumerical inits a classification observer for numerical features.
func NewClassificationNumerical(maxBuckets uint32) *ClassificationNumerical {
	if maxBuckets == 0 {
		maxBuckets = 11
	}
	return &ClassificationNumerical{MaxBuckets: maxBuckets}
}

// Observe adds a new observation.
func (o *ClassificationNumerical) Observe(val float64, target core.Category) {
	o.ObserveWeight(val, target, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (o *ClassificationNumerical) ObserveWeight(val float64, target core.Category, weight float64) {
	if !(core.IsCat(target) && core.IsNum(val) && weight > 0) {
		return
	}
	o.Dist.ObserveWeight(int(target), val, weight)
}

// EvaluateSplit evaluates a split.
func (o *ClassificationNumerical) EvaluateSplit(crit split.Criterion, pre *util.Vector) (merit, pivot float64, post *util.Matrix) {
	min, max := o.boundaries()
	inc := (max - min) / float64(o.MaxBuckets+1)
	if inc <= 0 {
		return
	}

	for i := uint32(0); i < o.MaxBuckets; i++ {
		pv := min + inc*float64(i+1)
		pc := o.postSplit(pv)
		mc := crit.ClassificationMerit(pre, pc)

		if post == nil || mc > merit {
			merit, pivot, post = mc, pv, pc
		}
	}
	return
}

func (o *ClassificationNumerical) postSplit(pivot float64) *util.Matrix {
	mat := util.NewMatrix()
	for i, n := 0, o.Dist.NumRows(); i < n; i++ {
		t := o.Dist.At(i)
		if t == nil {
			continue
		}

		if t.Min > 0 && pivot < t.Min {
			mat.Incr(1, i, t.Weight)
		} else if t.Max > 0 && pivot >= t.Max {
			mat.Incr(0, i, t.Weight)
		} else {
			lt, eq, gt := t.Estimate(pivot)
			mat.Incr(0, i, lt+eq)
			mat.Incr(1, i, gt)
		}
	}
	return mat
}

func (o *ClassificationNumerical) boundaries() (min float64, max float64) {
	for i, n := 0, o.Dist.NumRows(); i < n; i++ {
		if t := o.Dist.At(i); t != nil {
			if min == 0 || t.Min < min {
				min = t.Min
			}
			if max == 0 || t.Max > max {
				max = t.Max
			}
		}
	}
	return
}
