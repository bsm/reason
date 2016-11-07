package calc

import "math"

func Sum(s []float64) (Σ float64) {
	for _, v := range s {
		Σ += v
	}
	return
}

func Min(s []float64) float64 {
	if len(s) == 0 {
		return math.NaN()
	}

	m := s[0]
	for _, v := range s[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func Max(s []float64) float64 {
	if len(s) == 0 {
		return math.NaN()
	}

	m := s[0]
	for _, v := range s[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

func Mean(s []float64) (µ float64) {
	n := len(s)
	if n == 0 {
		return
	}

	µ = Sum(s) / float64(n)
	return
}

func Variance(s []float64) (V float64) {
	n := len(s)
	if n == 0 {
		return
	}

	µ := Mean(s)
	for _, v := range s {
		Δ := v - µ
		V += Δ * Δ
	}

	V /= float64(n)
	return
}

func StdDev(s []float64) (σ float64) {
	V := Variance(s)
	if V == 0 {
		return
	}

	σ = math.Sqrt(V)
	return
}

func SampleVariance(s []float64) (sV float64) {
	n := len(s)
	if n < 2 {
		return
	}

	µ := Mean(s)
	for _, v := range s {
		Δ := v - µ
		sV += Δ * Δ
	}

	sV /= float64(n - 1)
	return
}

func SampleStdDev(s []float64) (sσ float64) {
	sV := SampleVariance(s)
	if sV == 0 {
		return
	}

	sσ = math.Sqrt(sV)
	return
}

// MatrixRowSumsPlusTotal returns a vector of matrix row sums and the overall matrix sum
func MatrixRowSumsPlusTotal(m [][]float64) ([]float64, float64) {
	sums := make([]float64, len(m))
	total := 0.0
	for i, vv := range m {
		sum := Sum(vv)
		sums[i] = sum
		total += sum
	}
	return sums, total
}

func Entropy(s []float64) float64 {
	ent := 0.0
	sum := 0.0
	for _, v := range s {
		if v > 0 {
			ent -= v * math.Log2(v)
			sum += v
		}
	}
	if sum > 0 {
		return (ent + sum*math.Log2(sum)) / sum
	}
	return 0.0
}

func MaxIndex(s []float64) int {
	if len(s) == 0 {
		return -1
	}

	n, m := 0, s[0]
	for i, v := range s[1:] {
		if v > m {
			n, m = i+1, v
		}
	}
	return n
}

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
