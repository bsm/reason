package helpers

import (
	"github.com/bsm/reason/classifiers"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
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
	// ByteSize returns a required heap-size estimate
	ByteSize() int
	// Promise returns the promise for making predictions
	Promise() float64
	// BestSplit returns a SplitSuggestion
	BestSplit(crit classifiers.SplitCriterion, obs Observer, predictor *core.Attribute) *SplitSuggestion
	// State returns the current state as a prediction
	State() core.Prediction
}

func NewObservationStats(isRegression bool) ObservationStats {
	if isRegression {
		return newObsRStats()
	}
	return newObsCStats()
}

func newCObservationStats(preSplit util.Vector) ObservationStats {
	return &obsCStats{preSplit: preSplit}
}

func newCObservationStatsDist(postSplit util.VectorDistribution) map[int]ObservationStats {
	res := make(map[int]ObservationStats, len(postSplit))
	for i, vv := range postSplit {
		res[i] = &obsCStats{preSplit: vv}
	}
	return res
}

func newRObservationStats(preSplit *util.NumSeries) ObservationStats {
	return &obsRStats{preSplit: *preSplit}
}

func newRObservationStatsDist(postSplit util.NumSeriesDistribution) map[int]ObservationStats {
	res := make(map[int]ObservationStats, len(postSplit))
	for i, vv := range postSplit {
		res[i] = &obsRStats{preSplit: vv}
	}
	return res
}

// --------------------------------------------------------------------

type obsCStats struct {
	preSplit util.Vector
}

func newObsCStats() *obsCStats {
	return &obsCStats{preSplit: util.NewVector()}
}

func (s *obsCStats) ByteSize() int { return 40 + s.preSplit.ByteSize() }

func (s *obsCStats) TotalWeight() float64 { return s.preSplit.Sum() }

func (s *obsCStats) Promise() float64 {
	if w := s.preSplit.Sum(); w != 0 {
		return w - s.preSplit.Max()
	}
	return 0.0
}

func (s *obsCStats) IsSufficient() bool {
	return s.preSplit.Count() > 1
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
	p := make(core.Prediction, 0, s.preSplit.Count())
	s.preSplit.ForEach(func(i int, v float64) {
		p = append(p, core.PredictedValue{
			Value: core.AttributeValue(i),
			Votes: v,
		})
	})
	return p
}

// --------------------------------------------------------------------

type obsRStats struct {
	preSplit util.NumSeries
}

func newObsRStats() *obsRStats {
	return &obsRStats{}
}

func (s *obsRStats) ByteSize() int        { return 40 }
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
