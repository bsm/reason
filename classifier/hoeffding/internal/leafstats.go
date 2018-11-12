package internal

import (
	"github.com/bsm/reason"
	"github.com/bsm/reason/common/observer"
)

// Update updates stats by observing an example.
func (s *LeafNode_Stats) Update(target, predictor *reason.Feature, x reason.Example, weight float64) {
	switch target.Kind {
	case reason.Feature_CATEGORICAL:
		if tcat := target.Category(x); reason.IsCat(tcat) {
			switch predictor.Kind {
			case reason.Feature_CATEGORICAL:
				if pcat := predictor.Category(x); reason.IsCat(pcat) {
					s.updateCC(tcat, pcat, weight)
				}
			case reason.Feature_NUMERICAL:
				if pval := predictor.Number(x); reason.IsNum(pval) {
					s.updateCN(tcat, pval, weight)
				}
			}
		}
	case reason.Feature_NUMERICAL:
		if tval := target.Number(x); reason.IsNum(tval) {
			switch predictor.Kind {
			case reason.Feature_CATEGORICAL:
				if pcat := predictor.Category(x); reason.IsCat(pcat) {
					s.updateRC(tval, pcat, weight)
				}
			case reason.Feature_NUMERICAL:
				if pval := predictor.Number(x); reason.IsNum(pval) {
					s.updateRN(tval, pval, weight)
				}
			}
		}
	}
}

func (s *LeafNode_Stats) updateCC(tcat, pcat reason.Category, weight float64) {
	acc := s.GetCC()
	if acc == nil {
		acc = observer.NewClassificationCategorical()
		s.Kind = &LeafNode_Stats_CC{CC: acc}
	}
	acc.ObserveWeight(pcat, tcat, weight)
}

func (s *LeafNode_Stats) updateCN(tcat reason.Category, pval, weight float64) {
	acc := s.GetCN()
	if acc == nil {
		acc = observer.NewClassificationNumerical(0)
		s.Kind = &LeafNode_Stats_CN{CN: acc}
	}
	acc.ObserveWeight(pval, tcat, weight)
}

func (s *LeafNode_Stats) updateRC(tval float64, pcat reason.Category, weight float64) {
	acc := s.GetRC()
	if acc == nil {
		acc = observer.NewRegressionCategorical()
		s.Kind = &LeafNode_Stats_RC{RC: acc}
	}
	acc.ObserveWeight(pcat, tval, weight)
}

func (s *LeafNode_Stats) updateRN(tval, pval, weight float64) {
	acc := s.GetRN()
	if acc == nil {
		acc = observer.NewRegressionNumerical(0)
		s.Kind = &LeafNode_Stats_RN{RN: acc}
	}
	acc.ObserveWeight(pval, tval, weight)
}
