package helpers

import (
	"github.com/bsm/reason/classifiers"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/util"
)

// ObservationStats stats are used to maintain sufficient
// stats across multiple attributes.
type ObservationStats interface {
	// IsSufficient returns true if stats contain a sufficient distribution of data
	IsSufficient() bool
	// UpdatePreSplit updates pre-split stats
	UpdatePreSplit(target core.AttributeValue, weight float64)
	// NewObserver creates a new attribute observer
	NewObserver(isNominal bool) Observer
	// TotalWeight returns the total weight observed
	TotalWeight() float64
	// HeapSize returns an in-memory heap-size estimate
	HeapSize() int
	// Promise returns the promise for making predictions
	Promise() float64
	// BestSplit returns a SplitSuggestion
	BestSplit(crit classifiers.SplitCriterion, obs Observer, predictor *core.Attribute) *SplitSuggestion
	// State returns the current state as a prediction
	State() core.Prediction
}

func NewObservationStats(isRegression bool) ObservationStats {
	if isRegression {
		return new(obsRStats)
	}
	return new(obsCStats)
}

func newCObservationStats(preSplit []float64) ObservationStats {
	return &obsCStats{preSplit: preSplit}
}

func newCObservationStatsSlice(postSplit [][]float64) []ObservationStats {
	slice := make([]ObservationStats, len(postSplit))
	for i, vv := range postSplit {
		slice[i] = &obsCStats{preSplit: vv}
	}
	return slice
}

func newRObservationStats(preSplit *core.NumSeries) ObservationStats {
	return &obsRStats{preSplit: *preSplit}
}

func newRObservationStatsSlice(postSplit []core.NumSeries) []ObservationStats {
	slice := make([]ObservationStats, len(postSplit))
	for i, vv := range postSplit {
		slice[i] = &obsRStats{preSplit: vv}
	}
	return slice
}

// --------------------------------------------------------------------

type obsCStats struct {
	preSplit util.NumVector
}

func (s *obsCStats) HeapSize() int { return 40 + len(s.preSplit)*8 }

func (s *obsCStats) TotalWeight() float64 { return s.preSplit.Sum() }

func (s *obsCStats) Promise() float64 {
	if w := s.preSplit.Sum(); w != 0 {
		return w - s.preSplit.Max()
	}
	return 0.0
}

func (s *obsCStats) IsSufficient() bool {
	m := 0
	for _, w := range s.preSplit {
		if w != 0 {
			if m++; m == 2 {
				return true
			}
		}
	}
	return false
}

func (s *obsCStats) UpdatePreSplit(tv core.AttributeValue, weight float64) {
	s.preSplit = s.preSplit.Incr(tv.Index(), weight)
}

func (s *obsCStats) NewObserver(isNominal bool) Observer {
	if isNominal {
		return NewNominalCObserver()
	}
	return NewNumericCObserver(10)
}

func (s *obsCStats) BestSplit(crit classifiers.SplitCriterion, obs Observer, predictor *core.Attribute) *SplitSuggestion {
	return obs.(CObserver).BestSplit(crit.(classifiers.CSplitCriterion), predictor, s.preSplit)
}

func (s *obsCStats) State() core.Prediction {
	p := make(core.Prediction, len(s.preSplit))
	for i, w := range s.preSplit {
		p[i].Value = core.AttributeValue(i)
		p[i].Votes = w
	}
	return p
}

// --------------------------------------------------------------------

type obsRStats struct {
	preSplit core.NumSeries
}

func (s *obsRStats) HeapSize() int        { return 40 }
func (s *obsRStats) TotalWeight() float64 { return s.preSplit.TotalWeight() }
func (s *obsRStats) Promise() float64     { return s.preSplit.TotalWeight() }
func (s *obsRStats) IsSufficient() bool   { return s.preSplit.SampleVariance() != 0 }

func (s *obsRStats) UpdatePreSplit(tv core.AttributeValue, weight float64) {
	s.preSplit.Append(tv.Value(), weight)
}

func (s *obsRStats) NewObserver(isNominal bool) Observer {
	if isNominal {
		return NewNominalRObserver()
	}
	return NewNumericRObserver(10)
}

func (s *obsRStats) State() core.Prediction {
	return core.Prediction{
		{Value: core.AttributeValue(s.preSplit.Mean()), Votes: s.preSplit.TotalWeight()},
	}
}

func (s *obsRStats) BestSplit(crit classifiers.SplitCriterion, obs Observer, predictor *core.Attribute) *SplitSuggestion {
	return obs.(RObserver).BestSplit(crit.(classifiers.RSplitCriterion), predictor, &s.preSplit)
}
