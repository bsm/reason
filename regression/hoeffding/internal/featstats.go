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

// ObserveWeight adds an observation.
func (s *FeatureStats_Numerical) ObserveWeight(featVal, targetVal, weight float64) {
	s.NumStreamBuckets.ObserveWeight(featVal, targetVal, weight)
}

// PostSplit calculates a post-split distribution from previous observations
func (s *FeatureStats_Numerical) PostSplit(pivot float64) *util.NumStreams {
	split := make([]util.NumStream, 2)
	for _, bucket := range s.Buckets {
		pos := 0
		if bucket.Threshold > pivot {
			pos = 1
		}
		split[pos].Merge(&bucket.NumStream)
	}
	return &util.NumStreams{Data: split}
}
