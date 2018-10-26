package regression

import (
	"math"

	"github.com/bsm/reason/util"
	"gonum.org/v1/gonum/stat/distuv"
)

// Stats accumulate regression stream stats.
type Stats struct{ *util.Vector }

// WrapStats init stats with a vector.
func WrapStats(vv *util.Vector) Stats {
	if vv == nil {
		vv = util.NewVector()
	}
	return Stats{Vector: vv}
}

// Observe adds a new observation.
func (s Stats) Observe(value float64) {
	s.ObserveWeight(value, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (s Stats) ObserveWeight(value, weight float64) {
	if math.IsNaN(value) || math.IsInf(value, 0) || weight <= 0 {
		return
	}

	wv := weight * value
	s.Add(0, weight)
	s.Add(1, wv)
	s.Add(2, wv*value)
}

// TotalWeight returns the total weight observed.
func (s Stats) TotalWeight() float64 {
	return s.At(0)
}

// Sum returns the weighted sum of all observations.
func (s Stats) Sum() float64 {
	return s.At(1)
}

// Mean returns a mean average
func (s Stats) Mean() float64 {
	return s.Sum() / s.TotalWeight()
}

// Variance is the sample variance of the series
func (s Stats) Variance() float64 {
	w := s.TotalWeight()
	if w <= 1 {
		return math.NaN()
	}
	z := s.Sum()
	return (s.At(2) - (z * z / w)) / (w - 1)
}

// StdDev is the sample standard deviation of the series
func (s Stats) StdDev() float64 {
	return math.Sqrt(s.Variance())
}

// Prob calculates the gaussian probability density of a value
func (s Stats) Prob(value float64) float64 {
	if sig := s.StdDev(); !math.IsNaN(sig) {
		return distuv.Normal{Mu: s.Mean(), Sigma: sig}.Prob(value)
	}
	return math.NaN()
}

// Estimate estimates weight boundaries for a given value
func (s Stats) Estimate(value float64) (lessThan float64, equalTo float64, greaterThan float64) {
	total := s.TotalWeight()
	equalTo = s.Prob(value) * total

	if sig := s.StdDev(); !math.IsNaN(sig) {
		lessThan = distuv.Normal{
			Mu:    s.Mean(),
			Sigma: sig,
		}.CDF(value)*total - equalTo
	} else {
		lessThan = math.NaN()
	}

	if greaterThan = total - equalTo - lessThan; greaterThan < 0 {
		greaterThan = 0
	}

	return
}

// --------------------------------------------------------------------

// StatsDistribution contains accumulated stream stats by category.
type StatsDistribution struct{ *util.Matrix }

// WrapStatsDistribution init stats with a vector.
func WrapStatsDistribution(mat *util.Matrix) StatsDistribution {
	if mat == nil {
		mat = util.NewMatrix()
	}
	return StatsDistribution{Matrix: mat}
}

// Observe adds a new observation.
func (s StatsDistribution) Observe(cat int, value float64) {
	s.ObserveWeight(cat, value, 1.0)
}

// ObserveWeight adds a new observation with a weight.
func (s StatsDistribution) ObserveWeight(cat int, value, weight float64) {
	if math.IsNaN(value) || math.IsInf(value, 0) || weight <= 0 {
		return
	}

	wv := weight * value
	s.Add(cat, 0, weight)
	s.Add(cat, 1, wv)
	s.Add(cat, 2, wv*value)
}

// NumCategories returns the number of categories.
func (s StatsDistribution) NumCategories() int {
	rows, _ := s.Dims()
	return rows
}

// TotalWeight returns the total weight observed for cat.
func (s StatsDistribution) TotalWeight(cat int) float64 {
	return s.At(cat, 0)
}

// Sum returns the weighted sum of all observations of cat.
func (s StatsDistribution) Sum(cat int) float64 {
	return s.At(cat, 1)
}

// Mean returns a mean average of ovserved cat values.
func (s StatsDistribution) Mean(cat int) float64 {
	return s.Sum(cat) / s.TotalWeight(cat)
}

// Variance is the sample variance of the cat series.
func (s StatsDistribution) Variance(cat int) float64 {
	w := s.TotalWeight(cat)
	if w <= 1 {
		return math.NaN()
	}
	z := s.Sum(cat)
	return (s.At(cat, 2) - (z * z / w)) / (w - 1)
}

// StdDev is the sample standard deviation of the cat series.
func (s StatsDistribution) StdDev(cat int) float64 {
	return math.Sqrt(s.Variance(cat))
}
