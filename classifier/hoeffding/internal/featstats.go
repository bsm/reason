package internal

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

// ObserveWeight adds an observation.
func (s *FeatureStats_ClassificationCategorical) ObserveWeight(featCat, targetCat core.Category, weight float64) {
	s.Matrix.Add(int(featCat), int(targetCat), weight)
}

// PostSplit calculates a post-split distribution from previous observations.
func (s *FeatureStats_ClassificationCategorical) PostSplit() *util.Matrix {
	return &s.Matrix
}

// NumCategories returns the number of categories.
func (s *FeatureStats_ClassificationCategorical) NumCategories() (n int) {
	for i, rows := 0, s.NumRows(); i < rows; i++ {
		if s.RowSum(i) > 0 {
			n++
		}
	}
	return
}

// --------------------------------------------------------------------

// ObserveWeight adds an observation.
func (s *FeatureStats_ClassificationNumerical) ObserveWeight(featVal float64, targetCat core.Category, weight float64) {
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

// PostSplit calculates a post-split distribution from previous observations
func (s *FeatureStats_ClassificationNumerical) PostSplit(pivot float64) *util.Matrix {
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

// PivotPoints determines the optimum split points for the range of values.
func (s *FeatureStats_ClassificationNumerical) PivotPoints() []float64 {
	_, min := s.Min.Min()
	_, max := s.Max.Max()
	return pivotPoints(min, max)
}
