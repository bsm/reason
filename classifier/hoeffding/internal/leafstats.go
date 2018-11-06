package internal

import (
	"github.com/bsm/reason/common/observer"
	"github.com/bsm/reason/core"
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
		acc = observer.NewClassificationCategorical()
		s.Kind = &LeafNode_Stats_CC{CC: acc}
	}
	acc.ObserveWeight(pcat, tcat, weight)
}

func (s *LeafNode_Stats) updateCN(tcat core.Category, pval, weight float64) {
	acc := s.GetCN()
	if acc == nil {
		acc = observer.NewClassificationNumerical()
		s.Kind = &LeafNode_Stats_CN{CN: acc}
	}
	acc.ObserveWeight(pval, tcat, weight)
}

func (s *LeafNode_Stats) updateRC(tval float64, pcat core.Category, weight float64) {
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
		acc = observer.NewRegressionNumerical(12)
		s.Kind = &LeafNode_Stats_RN{RN: acc}
	}
	acc.ObserveWeight(pval, tval, weight)
}
