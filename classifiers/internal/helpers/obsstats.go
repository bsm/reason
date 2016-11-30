package helpers

import (
	"github.com/bsm/reason/classifiers"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/msgpack"
	"github.com/bsm/reason/util"
)

func init() {
	msgpack.Register(7741, (*obsCStats)(nil))
	msgpack.Register(7742, (*obsRStats)(nil))
}

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
	return &obsCStats{PreSplit: preSplit}
}

func newCObservationStatsDist(postSplit util.VectorDistribution) map[int]ObservationStats {
	res := make(map[int]ObservationStats, len(postSplit))
	for i, vv := range postSplit {
		res[i] = &obsCStats{PreSplit: vv}
	}
	return res
}

func newRObservationStats(preSplit *util.NumSeries) ObservationStats {
	return &obsRStats{PreSplit: preSplit}
}

func newRObservationStatsDist(postSplit util.NumSeriesDistribution) map[int]ObservationStats {
	res := make(map[int]ObservationStats, len(postSplit))
	for i, vv := range postSplit {
		res[i] = &obsRStats{PreSplit: vv}
	}
	return res
}

// --------------------------------------------------------------------

type obsCStats struct {
	PreSplit util.Vector
}

func newObsCStats() *obsCStats {
	return &obsCStats{PreSplit: util.NewVector()}
}

func (s *obsCStats) ByteSize() int { return 40 + s.PreSplit.ByteSize() }

func (s *obsCStats) TotalWeight() float64 { return s.PreSplit.Sum() }

func (s *obsCStats) IsSufficient() bool {
	return s.PreSplit.Count() > 1
}

func (s *obsCStats) UpdatePreSplit(tv core.AttributeValue, weight float64) {
	s.PreSplit = s.PreSplit.Incr(tv.Index(), weight)
}

func (s *obsCStats) NewObserver(isNominal bool) Observer {
	if isNominal {
		return NewNominalCObserver()
	}
	return NewNumericCObserver(10)
}

func (s *obsCStats) BestSplit(crit classifiers.SplitCriterion, obs Observer, predictor *core.Attribute) *SplitSuggestion {
	return obs.(CObserver).BestSplit(crit.(classifiers.CSplitCriterion), predictor, s.PreSplit)
}

func (s *obsCStats) State() core.Prediction {
	p := make(core.Prediction, 0, s.PreSplit.Count())
	s.PreSplit.ForEach(func(i int, v float64) {
		p = append(p, core.PredictedValue{
			AttributeValue: core.AttributeValue(i),
			Votes:          v,
		})
	})
	return p
}

func (s *obsCStats) EncodeTo(enc *msgpack.Encoder) error {
	return enc.Encode(s.PreSplit)
}

func (s *obsCStats) DecodeFrom(dec *msgpack.Decoder) error {
	return dec.Decode(&s.PreSplit)
}

// --------------------------------------------------------------------

type obsRStats struct {
	PreSplit *util.NumSeries
}

func newObsRStats() *obsRStats {
	return &obsRStats{PreSplit: new(util.NumSeries)}
}

func (s *obsRStats) ByteSize() int        { return 40 }
func (s *obsRStats) TotalWeight() float64 { return s.PreSplit.TotalWeight() }
func (s *obsRStats) IsSufficient() bool   { return s.PreSplit.SampleVariance() != 0 }

func (s *obsRStats) UpdatePreSplit(tv core.AttributeValue, weight float64) {
	s.PreSplit.Append(tv.Value(), weight)
}

func (s *obsRStats) NewObserver(isNominal bool) Observer {
	if isNominal {
		return NewNominalRObserver()
	}
	return NewNumericRObserver(10)
}

func (s *obsRStats) State() core.Prediction {
	return core.Prediction{{
		AttributeValue: core.AttributeValue(s.PreSplit.Mean()),
		Votes:          s.PreSplit.TotalWeight(),
	}}
}

func (s *obsRStats) BestSplit(crit classifiers.SplitCriterion, obs Observer, predictor *core.Attribute) *SplitSuggestion {
	return obs.(RObserver).BestSplit(crit.(classifiers.RSplitCriterion), predictor, s.PreSplit)
}

func (s *obsRStats) EncodeTo(enc *msgpack.Encoder) error {
	return enc.Encode(s.PreSplit)
}

func (s *obsRStats) DecodeFrom(dec *msgpack.Decoder) error {
	return dec.Decode(&s.PreSplit)
}
