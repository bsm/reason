package calc

import "math"

const sqrth = 7.07106781186547524401E-1

// NormProb returns the area under the Normal (Gaussian) probability
// density function
func NormProb(a float64) float64 {
	var p float64

	x := a * sqrth
	if z := math.Abs(x); z < sqrth {
		p = 0.5 + 0.5*math.Erf(x)
	} else {
		p = 0.5 * math.Erfc(z)
		if x > 0 {
			p = 1.0 - p
		}
	}
	return p
}
