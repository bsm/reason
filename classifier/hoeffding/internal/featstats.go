package internal

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

func (s *FeatureStats) ObserveExample(target, predictor *core.Feature, x core.Example, weight float64) {
	switch target.Kind {
	case core.Feature_CATEGORICAL: // classification
		tcat := target.Category(x)
		if !core.IsCat(tcat) {
			return
		}

		switch predictor.Kind {
		case core.Feature_CATEGORICAL:
			pcat := predictor.Category(x)
			if !core.IsCat(pcat) {
				return
			}

			acc := s.GetCC()
			if acc == nil {
				acc = new(FeatureStats_ClassificationCategorical)
				s.Kind = &FeatureStats_CC{CC: acc}
			}
			acc.Stats.Add(int(pcat), int(tcat), weight)
		case core.Feature_NUMERICAL:
			pnum := predictor.Number(x)
			if !core.IsNum(pnum) {
				return
			}

			acc := s.GetCN()
			if acc == nil {
				acc = new(FeatureStats_ClassificationNumerical)
				s.Kind = &FeatureStats_CN{CN: acc}
			}
			acc.Stats.ObserveWeight(int(tcat), pnum, weight)
		}
	case core.Feature_NUMERICAL:
		tnum := target.Number(x)
		if !core.IsNum(tnum) {
			return
		}

		switch predictor.Kind {
		case core.Feature_CATEGORICAL:
			pcat := predictor.Category(x)
			if !core.IsCat(pcat) {
				return
			}

			acc := s.GetRC()
			if acc == nil {
				acc = new(FeatureStats_RegressionCategorical)
				s.Kind = &FeatureStats_RC{RC: acc}
			}
			acc.Stats.ObserveWeight(int(pcat), tnum, weight)
		case core.Feature_NUMERICAL:
			panic("TODO")
		}
	}
}

// --------------------------------------------------------------------

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
