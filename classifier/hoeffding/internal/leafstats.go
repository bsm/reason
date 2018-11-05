package internal

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
	"github.com/bsm/reason/util/treeutil"
)

// Update updates stats by observing an example.
func (s *LeafNode_Stats) Update(target, predictor *core.Feature, x core.Example, weight float64) {
	switch target.Kind {
	case core.Feature_CATEGORICAL:
		if tcat := target.Category(x); core.IsCat(tcat) {
			switch predictor.Kind {
			case core.Feature_CATEGORICAL:
				if pcat := predictor.Category(x); core.IsCat(pcat) {
					s.updateCC(tcat, pcat, weight)
				}
			case core.Feature_NUMERICAL:
				if pval := predictor.Number(x); core.IsNum(pval) {
					s.updateCN(tcat, pval, weight)
				}
			}
		}
	case core.Feature_NUMERICAL:
		if tval := target.Number(x); core.IsNum(tval) {
			switch predictor.Kind {
			case core.Feature_CATEGORICAL:
				if pcat := predictor.Category(x); core.IsCat(pcat) {
					s.updateRC(tval, pcat, weight)
				}
			case core.Feature_NUMERICAL:
				if pval := predictor.Number(x); core.IsNum(pval) {
					s.updateRN(tval, pval, weight)
				}
			}
		}
	}
}

func (s *LeafNode_Stats) updateCC(tcat, pcat core.Category, weight float64) {
	acc := s.GetCC()
	if acc == nil {
		acc = new(LeafNode_Stats_ClassificationCategorical)
		s.Kind = &LeafNode_Stats_CC{CC: acc}
	}
	acc.Stats.Add(int(pcat), int(tcat), weight)
}

func (s *LeafNode_Stats) updateCN(tcat core.Category, pval, weight float64) {
	acc := s.GetCN()
	if acc == nil {
		acc = new(LeafNode_Stats_ClassificationNumerical)
		s.Kind = &LeafNode_Stats_CN{CN: acc}
	}
	acc.Stats.ObserveWeight(int(tcat), pval, weight)
}

func (s *LeafNode_Stats) updateRC(tval float64, pcat core.Category, weight float64) {
	acc := s.GetRC()
	if acc == nil {
		acc = new(LeafNode_Stats_RegressionCategorical)
		s.Kind = &LeafNode_Stats_RC{RC: acc}
	}
	acc.Stats.ObserveWeight(int(pcat), tval, weight)
}

func (s *LeafNode_Stats) updateRN(tval, pval, weight float64) {
	acc := s.GetRN()
	if acc == nil {
		acc = new(LeafNode_Stats_RegressionNumerical)
		acc.Stats.MaxBuckets = 12
		s.Kind = &LeafNode_Stats_RN{RN: acc}
	}
	acc.Stats.ObserveWeight(pval, tval, weight)
}

// --------------------------------------------------------------------

// EvaluateSplit evaluates a split.
func (s *LeafNode_Stats_ClassificationCategorical) evaluateSplit(crit treeutil.SplitCriterion, pre *util.Vector) *SplitCandidate {
	if s.numCategories() < 2 {
		return nil
	}

	post := PostSplit{Classification: &s.Stats}
	return &SplitCandidate{
		Range:     crit.ClassificationRange(pre),
		Merit:     crit.ClassificationMerit(pre, post.Classification),
		PostSplit: post,
	}
}

func (s *LeafNode_Stats_ClassificationCategorical) numCategories() (n int) {
	rows := s.Stats.NumRows()
	for i := 0; i < rows; i++ {
		if !s.Stats.IsRowZero(i) {
			n++
		}
	}
	return
}

// --------------------------------------------------------------------

// EvaluateSplit evaluates a split.
func (s *LeafNode_Stats_ClassificationNumerical) evaluateSplit(crit treeutil.SplitCriterion, pre *util.Vector) (sc *SplitCandidate) {
	rang := crit.ClassificationRange(pre)

	for _, pivot := range s.PivotPoints() {
		post := s.PostSplit(pivot)
		merit := crit.ClassificationMerit(pre, post.Classification)

		if sc == nil || merit > sc.Merit {
			sc = &SplitCandidate{
				Merit:     merit,
				Range:     rang,
				Pivot:     pivot,
				PostSplit: post,
			}
		}
	}
	return
}

// PostSplit calculates a post-split distribution from previous observations.
func (s *LeafNode_Stats_ClassificationNumerical) PostSplit(pivot float64) PostSplit {
	mat := util.NewMatrix()
	rows := s.Stats.NumRows()
	for i := 0; i < rows; i++ {
		t := s.Stats.At(i)
		if t == nil {
			continue
		}

		if t.Min > 0 && pivot < t.Min {
			mat.Add(1, i, t.Weight)
		} else if t.Max > 0 && pivot >= t.Max {
			mat.Add(0, i, t.Weight)
		} else {
			lt, eq, gt := t.Estimate(pivot)
			mat.Add(0, i, lt+eq)
			mat.Add(1, i, gt)
		}
	}
	return PostSplit{Classification: mat}
}

// PivotPoints determines the optimum split points for the range of values.
func (s *LeafNode_Stats_ClassificationNumerical) PivotPoints() []float64 {
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

// --------------------------------------------------------------------

// EvaluateSplit evaluates a split.
func (s *LeafNode_Stats_RegressionCategorical) evaluateSplit(crit treeutil.SplitCriterion, pre *util.NumStream) *SplitCandidate {
	if n := s.Stats.NumCategories(); n < 2 {
		return nil
	}

	post := PostSplit{Regression: &s.Stats}
	return &SplitCandidate{
		Range:     crit.RegressionRange(pre),
		Merit:     crit.RegressionMerit(pre, post.Regression),
		PostSplit: post,
	}
}

// --------------------------------------------------------------------

// EvaluateSplit evaluates a split.
func (s *LeafNode_Stats_RegressionNumerical) evaluateSplit(crit treeutil.SplitCriterion, pre *util.NumStream) (sc *SplitCandidate) {
	rang := crit.RegressionRange(pre)

	for i := 0; i < len(s.Stats.Buckets)-1; i++ {
		pivot := s.Stats.Buckets[i].Threshold
		post := s.PostSplit(pivot)
		merit := crit.RegressionMerit(pre, post.Regression)

		if sc == nil || merit > sc.Merit {
			sc = &SplitCandidate{
				Range:     rang,
				Merit:     merit,
				Pivot:     pivot,
				PostSplit: post,
			}
		}
	}
	return
}

// PostSplit calculates a post-split distribution from previous observations.
func (s *LeafNode_Stats_RegressionNumerical) PostSplit(pivot float64) PostSplit {
	data := make([]util.NumStream, 2)
	for _, bucket := range s.Stats.Buckets {
		pos := 0
		if bucket.Threshold > pivot {
			pos = 1
		}
		data[pos].Merge(&bucket.NumStream)
	}
	return PostSplit{Regression: &util.NumStreams{Data: data}}
}

// --------------------------------------------------------------------

const numPivotBuckets = 11

func pivotPoints(min, max float64) []float64 {
	inc := (max - min) / float64(numPivotBuckets+1)
	if inc <= 0 {
		return nil
	}

	pp := make([]float64, numPivotBuckets)
	for i := 0; i < numPivotBuckets; i++ {
		pp[i] = min + inc*float64(i+1)
	}
	return pp
}
