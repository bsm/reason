package internal

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

// ObserveWeight adds an observation.
func (s *FeatureStats_ClassificationCategorical) ObserveWeight(featCat, targetCat core.Category, weight float64) {
	s.Stats.Add(int(featCat), int(targetCat), weight)
}

// --------------------------------------------------------------------

// ObserveWeight adds an observation.
func (s *FeatureStats_ClassificationNumerical) ObserveWeight(featVal float64, targetCat core.Category, weight float64) {
	s.Stats.ObserveWeight(int(targetCat), featVal, weight)
}

// PostSplit calculates a post-split distribution from previous observations
func (s *FeatureStats_ClassificationNumerical) PostSplit(pivot float64) *util.Matrix {
	post := util.NewMatrix()
	rows := s.Stats.NumRows()
	for i := 0; i < rows; i++ {
		t := s.Stats.At(i)
		if t == nil {
			continue
		}

		if t.Min > 0 && pivot < t.Min {
			post.Add(1, i, t.Weight)
		} else if t.Max > 0 && pivot >= t.Max {
			post.Add(0, i, t.Weight)
		} else {
			lt, eq, gt := t.Estimate(pivot)
			post.Add(0, i, lt+eq)
			post.Add(1, i, gt)
		}
	}
	return post
}

// PivotPoints determines the optimum split points for the range of values.
func (s *FeatureStats_ClassificationNumerical) PivotPoints() []float64 {
	var min, max float64

	rows := s.Stats.NumRows()
	for i := 0; i < rows; i++ {
		if t := s.Stats.At(i); t != nil {
			if min == 0 || t.Min < min {
				min = t.Min
			}
			if max == 0 || t.Max > max {
				max = t.Max
			}
		}
	}
	return pivotPoints(min, max)
}
