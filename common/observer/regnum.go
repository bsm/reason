package observer

import (
	"math"
	"sort"

	"github.com/bsm/reason/common/split"
	util "github.com/bsm/reason/util"
)

// NewRegressionNumerical inits a regression observer for numerical features.
func NewRegressionNumerical(maxBuckets uint32) *RegressionNumerical {
	if maxBuckets == 0 {
		maxBuckets = 16
	}
	return &RegressionNumerical{MaxBuckets: maxBuckets}
}

// Observe adds a new observation.
func (o *RegressionNumerical) Observe(target, prediction float64) {
	o.ObserveWeight(target, prediction, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (o *RegressionNumerical) ObserveWeight(target, prediction, weight float64) {
	if math.IsNaN(target) || math.IsInf(target, 0) || math.IsNaN(prediction) || math.IsInf(prediction, 0) || weight <= 0 {
		return
	}

	// upsert bucket
	slot := o.findSlot(target)
	if slot < len(o.Dist) {
		if o.Dist[slot].Threshold != target {
			o.Dist = append(o.Dist, RegressionNumerical_Bucket{})
			copy(o.Dist[slot+1:], o.Dist[slot:])
			o.Dist[slot] = RegressionNumerical_Bucket{Threshold: target}
		}
	} else {
		o.Dist = append(o.Dist, RegressionNumerical_Bucket{Threshold: target})
	}
	o.Dist[slot].ObserveWeight(prediction, weight)

	// prune buckets
	for uint32(len(o.Dist)) > o.MaxBuckets {
		delta := math.MaxFloat64
		slot := 0
		for i := 0; i < len(o.Dist)-1; i++ {
			if x := o.Dist[i+1].Threshold - o.Dist[i].Threshold; x < delta {
				slot, delta = i, x
			}
		}

		b1, b2 := o.Dist[slot], o.Dist[slot+1]
		b2.Threshold = (b1.Threshold*b1.Weight + b2.Threshold*b2.Weight) / (b1.Weight + b2.Weight)
		b2.NumStream.Merge(&b1.NumStream)

		o.Dist[slot+1] = b2
		o.Dist = o.Dist[:slot+copy(o.Dist[slot:], o.Dist[slot+1:])]
	}
}

// EvaluateSplit evaluates a split.
func (o *RegressionNumerical) EvaluateSplit(crit split.Criterion, pre *util.NumStream) (merit, pivot float64, post *util.NumStreams) {
	for i := 0; i < len(o.Dist)-1; i++ {
		pv := o.Dist[i].Threshold
		pc := o.postSplit(pv)
		mc := crit.RegressionMerit(pre, pc)

		if post == nil || mc > merit {
			merit, pivot, post = mc, pv, pc
		}
	}
	return
}

func (o *RegressionNumerical) findSlot(v float64) int {
	return sort.Search(len(o.Dist), func(i int) bool { return o.Dist[i].Threshold >= v })
}

func (o *RegressionNumerical) postSplit(pivot float64) *util.NumStreams {
	data := make([]util.NumStream, 2)
	for _, b := range o.Dist {
		pos := 0
		if b.Threshold > pivot {
			pos = 1
		}
		data[pos].Merge(&b.NumStream)
	}
	return &util.NumStreams{Data: data}
}
