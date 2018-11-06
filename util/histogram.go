package util

import (
	"math"
	"sort"
)

// NewHistogram inits a new histogram.
func NewHistogram(maxBins uint32) *Histogram {
	return &Histogram{Cap: maxBins}
}

// Observe adds a new observation.
func (h *Histogram) Observe(value float64) {
	h.ObserveWeight(value, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (h *Histogram) ObserveWeight(value, weight float64) {
	// update min/max
	if h.Weight == 0 || value < h.Min {
		h.Min = value
	}
	if h.Weight == 0 || value > h.Max {
		h.Max = value
	}

	// insert bin
	h.Weight += weight
	if slot := h.search(value); slot < len(h.Bins) {
		if h.Bins[slot].Value == value {
			h.Bins[slot].Weight += math.Copysign(weight, h.Bins[slot].Weight)
		} else {
			h.Bins = append(h.Bins, Histogram_Bin{})
			copy(h.Bins[slot+1:], h.Bins[slot:])
			h.Bins[slot] = Histogram_Bin{Value: value, Weight: weight}
		}
	} else {
		h.Bins = append(h.Bins, Histogram_Bin{Value: value, Weight: weight})
	}

	// prune bins
	for len(h.Bins) > int(h.Cap) {
		delta := math.MaxFloat64
		slot := 0
		for i := 0; i < len(h.Bins)-1; i++ {
			if x := h.Bins[i+1].Value - h.Bins[i].Value; x < delta {
				slot, delta = i, x
			}
		}

		b1, b2 := h.Bins[slot], h.Bins[slot+1]
		wsum := math.Abs(b1.Weight) + math.Abs(b2.Weight)
		h.Bins[slot+1] = Histogram_Bin{
			Weight: -wsum,
			Value:  (b1.sum() + b2.sum()) / wsum,
		}
		h.Bins = h.Bins[:slot+copy(h.Bins[slot:], h.Bins[slot+1:])]
	}
}

// Sum returns the (approximate) sum of all observed values.
func (h *Histogram) Sum() float64 {
	if h.Weight == 0 {
		return math.NaN()
	}

	sum := 0.0
	for _, bin := range h.Bins {
		sum += bin.sum()
	}
	return sum
}

// Mean returns the (approximate) mean average.
func (h *Histogram) Mean() float64 {
	return h.Sum() / h.Weight
}

// Variance is the (approximate) sample variance of the series.
func (h *Histogram) Variance() float64 {
	if h.Weight <= 1 {
		return math.NaN()
	}

	sls, mean := 0.0, h.Mean()
	for _, bin := range h.Bins {
		delta := mean - bin.Value
		sls += delta * delta * bin.Weight
	}
	return sls / (h.Weight - 1)
}

// StdDev is the (approximate) sample standard deviation of the series.
func (h *Histogram) StdDev() float64 {
	return math.Sqrt(h.Variance())
}

// Quantile returns the (approximate) quantile.
// Accepted values for q are between 0.0 and 1.0.
func (h *Histogram) Quantile(q float64) float64 {
	if h.Weight == 0 || q < 0.0 || q > 1.0 || len(h.Bins) == 0 {
		return math.NaN()
	} else if q == 0.0 {
		return h.Min
	} else if q == 1.0 {
		return h.Max
	}

	delta := q * h.Weight
	slot := 0
	for w0 := 0.0; slot < len(h.Bins); slot++ {
		w1 := math.Abs(h.Bins[slot].Weight) / 2.0
		if delta-w1-w0 < 0 {
			break
		}
		delta -= (w1 + w0)
		w0 = w1
	}

	switch slot {
	case 0: // lower bound
		hi := h.Bins[slot]
		return h.solve(h.Min, 0, hi.Value, hi.Weight, delta)
	case len(h.Bins): // upper bound
		lo := h.Bins[slot-1]
		return h.solve(lo.Value, lo.Weight, h.Max, 0, delta)
	default:
		lo, hi := h.Bins[slot-1], h.Bins[slot]
		return h.solve(lo.Value, lo.Weight, hi.Value, hi.Weight, delta)
	}
}

func (h *Histogram) solve(v1, w1, v2, w2, delta float64) float64 {
	// return if both bins are exact (unmerged)
	if w1 > 0 && w2 > 0 {
		return v2
	}

	// normalize
	w1, w2 = math.Abs(w1), math.Abs(w2)

	// calculate multiplier
	var z float64
	if w1 == w2 {
		z = delta / w1
	} else {
		a := 2 * (w2 - w1)
		b := 2 * w1
		z = (math.Sqrt(b*b+4*a*delta) - b) / a
	}
	return v1 + (v2-v1)*z
}

func (h *Histogram) search(v float64) int {
	return sort.Search(len(h.Bins), func(i int) bool { return h.Bins[i].Value >= v })
}

func (b Histogram_Bin) sum() float64 {
	return math.Abs(b.Weight) * b.Value
}
