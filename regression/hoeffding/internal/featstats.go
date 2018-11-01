package internal

import (
	"github.com/bsm/reason/core"
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

// PostSplit calculates a post-split distribution from previous observations.
func (s *FeatureStats_Categorical) PostSplit() *util.NumStreams {
	return &s.NumStreams
}

// ObserveWeight adds an observation
func (s *FeatureStats_Categorical) ObserveWeight(featCat core.Category, targetVal, weight float64) {
	s.NumStreams.ObserveWeight(int(featCat), targetVal, weight)
}

// --------------------------------------------------------------------

// // ObserveWeight adds an observation.
func (s *FeatureStats_Numerical) ObserveWeight(featVal, weight float64) {
	if s.Histogram.Cap == 0 {
		s.Histogram.Cap = 10
	}
	s.Histogram.ObserveWeight(featVal, weight)
}

// PostSplit calculates a post-split distribution from previous observations
func (s *FeatureStats_Numerical) PostSplit(pivot float64) *util.NumStreams {
	post := util.NewNumStreams()
	for _, bin := range s.Bins {
		if bin.Value <= pivot {
			post.ObserveWeight(0, bin.Value, bin.Weight)
		} else {
			post.ObserveWeight(1, bin.Value, bin.Weight)
		}
	}
	return post
}
