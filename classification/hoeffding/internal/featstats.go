package internal

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/hoeffding"
	"github.com/bsm/reason/util"
)

// FetchCategorical fetches categorical stats.
func (s *FeatureStats) FetchCategorical() *FeatureStats_Categorical {
	stats := s.GetCategorical()
	if stats == nil {
		stats = new(FeatureStats_Categorical)
		s.Kind = &FeatureStats_Categorical_{Categorical: stats}
	}
	return stats
}

// FetchNumerical fetches numerical stats.
func (s *FeatureStats) FetchNumerical() *FeatureStats_Numerical {
	stats := s.GetNumerical()
	if stats == nil {
		stats = new(FeatureStats_Numerical)
		s.Kind = &FeatureStats_Numerical_{Numerical: stats}
	}
	return stats
}

// --------------------------------------------------------------------

// NumCategories returns the number of categories.
func (s *FeatureStats_Categorical) NumCategories() (n int) {
	rows := s.NumRows()
	for i := 0; i < rows; i++ {
		if s.RowSum(i) > 0 {
			n++
		}
	}
	return
}

// PostSplit calculates a post-split distribution from previous observations.
func (s *FeatureStats_Categorical) PostSplit() *util.Matrix {
	return &s.Matrix
}

// ObserveWeight adds an observation.
func (s *FeatureStats_Categorical) ObserveWeight(featCat, targetCat core.Category, weight float64) {
	s.Matrix.Add(int(featCat), int(targetCat), weight)
}

// --------------------------------------------------------------------

// ObserveWeight adds an observation.
func (s *FeatureStats_Numerical) ObserveWeight(featVal float64, targetCat core.Category, weight float64) {
	targetPos := int(targetCat)
	if v := s.Min.At(targetPos); v == 0 || featVal < v {
		s.Min.Set(targetPos, featVal)
	}
	if v := s.Max.At(targetPos); v == 0 || featVal > v {
		s.Max.Set(targetPos, featVal)
	}
	streams := util.WrapNumStreams(&s.Stats)
	streams.ObserveWeight(targetPos, featVal, weight)
}

// PivotPoints determines the optimum split points for the range of values.
func (s *FeatureStats_Numerical) PivotPoints() []float64 {
	var tmin, tmax float64
	s.Min.ForEach(func(i int, min float64) bool {
		if tmin == 0 || min < tmin {
			tmin = min
		}
		if max := s.Max.At(i); tmax == 0 || max > tmax {
			tmax = max
		}
		return true
	})
	return hoeffding.PivotPoints(tmin, tmax)
}

// PostSplit calculates a post-split distribution from previous observations
func (s *FeatureStats_Numerical) PostSplit(pivot float64) *util.Matrix {
	post := util.NewMatrix()
	wrap := util.WrapNumStreams(&s.Stats)
	wrap.ForEach(func(cat int) {
		if min := s.Min.At(cat); min > 0 && pivot < min {
			post.Add(1, cat, wrap.TotalWeight(cat))
		} else if max := s.Max.At(cat); max > 0 && pivot >= max {
			post.Add(0, cat, wrap.TotalWeight(cat))
		} else {
			lt, eq, gt := wrap.Estimate(cat, pivot)
			post.Add(0, cat, lt+eq)
			post.Add(1, cat, gt)
		}
	})
	return post
}
