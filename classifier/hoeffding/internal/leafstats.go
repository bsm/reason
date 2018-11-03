package internal

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

/*
// ForEach interates over the stats distribution
func (s *LeafNode_Stats) ForEach(iter func(int, isNode_Stats)) {
	switch wrap := s.GetKind().(type) {
	case *LeafNode_Stats_CC:
		stats := wrap.CC.Stats
		nrows := stats.NumRows()
		for pos := 0; pos < nrows; pos++ {
			if stats.RowSum(pos) > 0 {
				ns := new(Node_ClassificationStats)
				ns.Vector.Data = stats.Row(pos)
				iter(pos, &Node_Classification{Classification: ns})
			}
		}
	case *LeafNode_Stats_CN:
		stats := wrap.CN.Stats
		nrows := stats.NumRows()
		for pos := 0; pos < nrows; pos++ {
			if info := stats.At(pos); info != nil {
				ns := new(Node_RegressionStats)
				ns.NumStream
				ns.Vector.Data = stats.Row(pos)
				iter(pos, &Node_Classification{Classification: ns})
			}
		}
	case *LeafNode_Stats_RC:
	case *LeafNode_Stats_RN:
	}
}
*/

// Update updates stats by observing an example.
func (s *LeafNode_Stats) Update(target, predictor *core.Feature, x core.Example, weight float64) {
	switch target.Kind {
	case core.Feature_CATEGORICAL:
		if tcat := target.Category(x); core.IsCat(tcat) {
			switch predictor.Kind {
			case core.Feature_CATEGORICAL:
				if pcat := predictor.Category(x); core.IsCat(pcat) {
					s.observeCC(tcat, pcat, weight)
				}
			case core.Feature_NUMERICAL:
				if pval := predictor.Number(x); core.IsNum(pval) {
					s.observeCN(tcat, pval, weight)
				}
			}
		}
	case core.Feature_NUMERICAL:
		if tval := target.Number(x); core.IsNum(tval) {
			switch predictor.Kind {
			case core.Feature_CATEGORICAL:
				if pcat := predictor.Category(x); core.IsCat(pcat) {
					s.observeRC(tval, pcat, weight)
				}
			case core.Feature_NUMERICAL:
				if pval := predictor.Number(x); core.IsNum(pval) {
					s.observeRN(tval, pval, weight)
				}
			}
		}
	}
}

func (s *LeafNode_Stats) observeCC(tcat, pcat core.Category, weight float64) {
	acc := s.GetCC()
	if acc == nil {
		acc = new(LeafNode_Stats_ClassificationCategorical)
		s.Kind = &LeafNode_Stats_CC{CC: acc}
	}
	acc.Stats.Add(int(pcat), int(tcat), weight)
}

func (s *LeafNode_Stats) observeCN(tcat core.Category, pval, weight float64) {
	acc := s.GetCN()
	if acc == nil {
		acc = new(LeafNode_Stats_ClassificationNumerical)
		s.Kind = &LeafNode_Stats_CN{CN: acc}
	}
	acc.Stats.ObserveWeight(int(tcat), pval, weight)
}

func (s *LeafNode_Stats) observeRC(tval float64, pcat core.Category, weight float64) {
	acc := s.GetRC()
	if acc == nil {
		acc = new(LeafNode_Stats_RegressionCategorical)
		s.Kind = &LeafNode_Stats_RC{RC: acc}
	}
	acc.Stats.ObserveWeight(int(pcat), tval, weight)
}

func (s *LeafNode_Stats) observeRN(tval, pval, weight float64) {
	acc := s.GetRN()
	if acc == nil {
		acc = new(LeafNode_Stats_RegressionNumerical)
		acc.Stats.MaxBuckets = 12
		s.Kind = &LeafNode_Stats_RN{RN: acc}
	}
	acc.Stats.ObserveWeight(tval, pval, weight)
}

// --------------------------------------------------------------------

// PostSplit calculates a post-split distribution from previous observations.
func (s *LeafNode_Stats_ClassificationCategorical) PostSplit(_ float64) PostSplit {
	return PostSplit{Classification: &s.Stats}
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

// PostSplit calculates a post-split distribution from previous observations.
func (s *LeafNode_Stats_RegressionCategorical) PostSplit(_ float64) PostSplit {
	return PostSplit{Regression: &s.Stats}
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
