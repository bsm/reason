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
func (s *FeatureStats_Categorical) NumCategories() int {
	return util.WrapNumStreams(&s.Matrix).NumCategories()
}

// PostSplit calculates a post-split distribution from previous observations.
func (s *FeatureStats_Categorical) PostSplit() *util.Matrix {
	return &s.Matrix
}

// ObserveWeight adds an observation
func (s *FeatureStats_Categorical) ObserveWeight(featCat core.Category, targetVal, weight float64) {
	util.WrapNumStreams(&s.Matrix).ObserveWeight(int(featCat), targetVal, weight)
}

// --------------------------------------------------------------------

// ObserveWeight adds an observation.
func (s *FeatureStats_Numerical) ObserveWeight(featVal, targetVal, weight float64) {
	if len(s.Observations) == 0 || featVal < s.Min {
		s.Min = featVal
	}
	if len(s.Observations) == 0 || featVal > s.Max {
		s.Max = featVal
	}

	s.Observations = append(s.Observations, FeatureStats_Numerical_Observation{
		FeatureValue: featVal,
		TargetValue:  targetVal,
		Weight:       weight,
	})
}

// PivotPoints determines the optimum split points for the range of values.
func (s *FeatureStats_Numerical) PivotPoints() []float64 {
	return hoeffding.PivotPoints(s.Min, s.Max)
}

// PostSplit calculates a post-split distribution from previous observations
func (s *FeatureStats_Numerical) PostSplit(pivot float64) *util.Matrix {
	post := util.NewMatrix()
	wrap := util.WrapNumStreams(post)
	for _, o := range s.Observations {
		if o.FeatureValue <= pivot {
			wrap.ObserveWeight(0, o.TargetValue, o.Weight)
		} else {
			wrap.ObserveWeight(1, o.TargetValue, o.Weight)
		}
	}
	return post
}
