package helpers

import (
	"github.com/bsm/reason/classifiers"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/msgpack"
	"github.com/bsm/reason/util"
)

func init() {
	msgpack.Register(7737, (*nominalCObserver)(nil))
	msgpack.Register(7738, (*gaussianCObserver)(nil))
	msgpack.Register(7739, (*nominalRObserver)(nil))
	msgpack.Register(7740, (*gaussianRObserver)(nil))
}

// Observer instances monitor and collect distribution stats
type Observer interface {
	// Observe records an instance and updates the attribute stats
	Observe(target, predictor core.AttributeValue, weight float64)

	// ByteSize estimates the required heap-size
	ByteSize() int
}

// CObserver instances monitor and collect distribution stats for
// predictor attributes and can suggest best possible splits.
type CObserver interface {
	Observer
	// Probability returns the probability of a given instance
	Probability(target, predictor core.AttributeValue) float64
	// BestSplit returns a suggestion for the best split
	BestSplit(_ classifiers.CSplitCriterion, predictor *core.Attribute, preSplit util.Vector) *SplitSuggestion
}

// NewNominalCObserver monitors a nominal predictor attribute
func NewNominalCObserver() CObserver {
	return &nominalCObserver{
		PostSplit: util.NewVectorDistribution(),
	}
}

type nominalCObserver struct {
	PostSplit util.VectorDistribution
}

func (o *nominalCObserver) ByteSize() int {
	return 40 + o.PostSplit.ByteSize()
}

// Observe implements CObserver
func (o *nominalCObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	o.PostSplit.Incr(tv.Index(), pv.Index(), weight)
}

// Probability implements CObserver
func (o *nominalCObserver) Probability(tv, pv core.AttributeValue) float64 {
	vec := o.PostSplit.Get(tv.Index())
	if vec == nil {
		return 0.0
	}
	cnt := o.PostSplit.NumTargets()
	return (vec.Get(pv.Index()) + 1) / (vec.Sum() + float64(cnt))
}

// BestSplit implements CObserver
func (o *nominalCObserver) BestSplit(crit classifiers.CSplitCriterion, predictor *core.Attribute, preSplit util.Vector) *SplitSuggestion {
	ncols := o.PostSplit.NumTargets()
	if ncols < 2 {
		return nil
	}

	postSplit := o.calcPostSplit(ncols)
	return &SplitSuggestion{
		cond:      NewNominalMultiwaySplitCondition(predictor),
		merit:     normMerit(crit.Merit(preSplit, postSplit)),
		mrange:    crit.Range(preSplit),
		preStats:  newCObservationStats(preSplit),
		postStats: newCObservationStatsDist(postSplit),
	}
}

func (o *nominalCObserver) calcPostSplit(ncols int) util.VectorDistribution {
	m := make(util.VectorDistribution, ncols)
	for ti, obs := range o.PostSplit {
		obs.ForEach(func(pi int, v float64) { m.Incr(pi, ti, v) })
	}
	return m
}

func (o *nominalCObserver) EncodeTo(enc *msgpack.Encoder) error   { return enc.Encode(o.PostSplit) }
func (o *nominalCObserver) DecodeFrom(dec *msgpack.Decoder) error { return dec.Decode(&o.PostSplit) }

// NewNumericCObserver uses gaussian estimators to monitor a numeric predictor attribute
func NewNumericCObserver(numBins int) CObserver {
	if numBins < 1 {
		numBins = 10
	}

	return &gaussianCObserver{
		NumBins:   numBins,
		Range:     NewMinMaxRanges(),
		PostSplit: util.NewNumSeriesDistribution(),
	}
}

type gaussianCObserver struct {
	NumBins   int
	Range     *MinMaxRanges
	PostSplit util.NumSeriesDistribution
}

func (o *gaussianCObserver) ByteSize() int {
	return 24 + o.Range.ByteSize() + o.PostSplit.ByteSize()
}

// Observe implements CObserver
func (o *gaussianCObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	ti := tv.Index()
	pval := pv.Value()
	o.PostSplit.Append(ti, pval, weight)
	o.Range.Update(ti, pval)
}

// Probability implements CObserver
func (o *gaussianCObserver) Probability(tv, pv core.AttributeValue) float64 {
	if est := o.PostSplit.Get(tv.Index()); est != nil {
		return est.ProbDensity(pv.Value())
	}
	return 0.0
}

// BestSplit implements Observes using a variance reduction
// algorithm
func (o *gaussianCObserver) BestSplit(crit classifiers.CSplitCriterion, predictor *core.Attribute, preSplit util.Vector) *SplitSuggestion {
	var best *SplitSuggestion

	for _, splitVal := range o.Range.SplitPoints(o.NumBins) {
		postSplit := o.binarySplitOn(splitVal)
		merit := crit.Merit(preSplit, postSplit)
		if best != nil && merit <= best.merit {
			continue
		}

		best = &SplitSuggestion{
			cond:      NewNumericBinarySplitCondition(predictor, splitVal),
			merit:     normMerit(merit),
			mrange:    crit.Range(preSplit),
			preStats:  newCObservationStats(preSplit),
			postStats: newCObservationStatsDist(postSplit),
		}
	}
	return best
}

func (o *gaussianCObserver) binarySplitOn(splitVal float64) util.VectorDistribution {
	res := util.NewVectorDistribution()
	for i, est := range o.PostSplit {
		if splitVal < o.Range.GetMin(i) {
			res.Incr(1, i, est.TotalWeight())
		} else if splitVal >= o.Range.GetMax(i) {
			res.Incr(0, i, est.TotalWeight())
		} else {
			lt, eq, gt := est.Estimate(splitVal)
			res.Incr(0, i, lt+eq)
			res.Incr(1, i, gt)
		}
	}
	return res
}

func (o *gaussianCObserver) EncodeTo(enc *msgpack.Encoder) error {
	return enc.Encode(o.NumBins, o.Range, o.PostSplit)
}

func (o *gaussianCObserver) DecodeFrom(dec *msgpack.Decoder) error {
	return dec.Decode(&o.NumBins, &o.Range, &o.PostSplit)
}

// --------------------------------------------------------------------

// RObserver instances monitor stats for predictor attributes in regressions.
type RObserver interface {
	Observer
	// BestSplit returns a suggestion for the best split
	BestSplit(crit classifiers.RSplitCriterion, predictor *core.Attribute, preSplit *util.NumSeries) *SplitSuggestion
}

// NewNominalRObserver monitors a nominal predictor attribute for a
// numeric regression target.
func NewNominalRObserver() RObserver {
	return &nominalRObserver{
		PostSplit: util.NewNumSeriesDistribution(),
	}
}

type nominalRObserver struct {
	PostSplit util.NumSeriesDistribution
}

func (o *nominalRObserver) ByteSize() int {
	return 40 + o.PostSplit.ByteSize()
}

// Observe implements Observer
func (o *nominalRObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	pi := pv.Index()
	o.PostSplit.Append(pi, tv.Value(), weight)
}

// BestSplit implements RegressionObserves using a variance reduction
// algorithm
func (o *nominalRObserver) BestSplit(crit classifiers.RSplitCriterion, predictor *core.Attribute, preSplit *util.NumSeries) *SplitSuggestion {
	if !o.isSplitable() {
		return nil
	}

	return &SplitSuggestion{
		cond:      NewNominalMultiwaySplitCondition(predictor),
		merit:     normMerit(crit.Merit(preSplit, o.PostSplit)),
		mrange:    crit.Range(preSplit),
		preStats:  newRObservationStats(preSplit),
		postStats: newRObservationStatsDist(o.PostSplit),
	}
}

func (o *nominalRObserver) isSplitable() bool {
	n := 0
	for _, s := range o.PostSplit {
		if s.TotalWeight() > 0 {
			if n++; n > 1 {
				return true
			}
		}
	}
	return false
}

func (o *nominalRObserver) EncodeTo(enc *msgpack.Encoder) error   { return enc.Encode(o.PostSplit) }
func (o *nominalRObserver) DecodeFrom(dec *msgpack.Decoder) error { return dec.Decode(&o.PostSplit) }

// NewNumericRObserver uses gaussian estimators to monitor a numeric predictor
// attribute for a numeric regression target
func NewNumericRObserver(numBins int) RObserver {
	if numBins < 1 {
		numBins = 10
	}

	return &gaussianRObserver{
		NumBins: numBins,
		Range:   NewMinMaxRange(),
	}
}

type gaussianRObserver struct {
	NumBins      int
	Range        *MinMaxRange
	Observations []Observation
}

func (o *gaussianRObserver) ByteSize() int {
	return 80 + len(o.Observations)*24
}

// Observe implements RObserver
func (o *gaussianRObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	tval, pval := tv.Value(), pv.Value()
	o.Range.Update(pval)
	o.Observations = append(o.Observations, Observation{
		PVal:   pval,
		TVal:   tval,
		Weight: weight,
	})
}

func (o *gaussianRObserver) BestSplit(crit classifiers.RSplitCriterion, predictor *core.Attribute, preSplit *util.NumSeries) *SplitSuggestion {
	var best *SplitSuggestion
	for _, pivot := range o.Range.SplitPoints(o.NumBins) {
		postSplit := o.postSplit(pivot)
		merit := crit.Merit(preSplit, postSplit)
		if best != nil && merit <= best.merit {
			continue
		}

		best = &SplitSuggestion{
			cond:      NewNumericBinarySplitCondition(predictor, pivot),
			merit:     normMerit(merit),
			mrange:    crit.Range(preSplit),
			preStats:  newRObservationStats(preSplit),
			postStats: newRObservationStatsDist(postSplit),
		}
	}
	return best
}

func (o *gaussianRObserver) postSplit(pivot float64) util.NumSeriesDistribution {
	res := util.NewNumSeriesDistribution()
	for _, t := range o.Observations {
		if t.PVal < pivot {
			res.Append(0, t.TVal, t.Weight)
		} else {
			res.Append(1, t.TVal, t.Weight)
		}
	}
	return res
}

func (o *gaussianRObserver) EncodeTo(enc *msgpack.Encoder) error {
	return enc.Encode(o.NumBins, o.Range, o.Observations)
}

func (o *gaussianRObserver) DecodeFrom(dec *msgpack.Decoder) error {
	return dec.Decode(&o.NumBins, &o.Range, &o.Observations)
}

func normMerit(merit float64) float64 {
	if merit > 0 {
		return merit
	}
	return 0.0
}
