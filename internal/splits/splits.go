package splits

import "math"

// NormMerit normalises merit values
func NormMerit(m float64) float64 {
	if m < 0.0 || math.IsNaN(m) {
		return 0.0
	}
	return m
}

type GainRatioPenalty struct {
	Weight float64
	value  float64
}

func (p *GainRatioPenalty) Update(w float64) {
	if p.Weight != 0 {
		rat := w / p.Weight
		p.value -= rat * math.Log2(rat)
	}
}

func (p *GainRatioPenalty) Value() float64 {
	if p.value <= 0.0 {
		return 1.0
	}
	return p.value
}
