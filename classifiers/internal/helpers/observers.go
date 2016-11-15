package helpers

import (
	"github.com/bsm/reason/classifiers"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

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
		postSplit: util.NewVectorDistribution(),
	}
}

type nominalCObserver struct {
	postSplit util.VectorDistribution
}

func (o *nominalCObserver) ByteSize() int {
	return 40 + o.postSplit.ByteSize()
}

// Observe implements CObserver
func (o *nominalCObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	o.postSplit.Incr(tv.Index(), pv.Index(), weight)
}

// Probability implements CObserver
func (o *nominalCObserver) Probability(tv, pv core.AttributeValue) float64 {
	vec := o.postSplit.Get(tv.Index())
	if vec == nil {
		return 0.0
	}
	cnt := o.postSplit.NumTargets()
	return (vec.Get(pv.Index()) + 1) / (vec.Sum() + float64(cnt))
}

// BestSplit implements CObserver
func (o *nominalCObserver) BestSplit(crit classifiers.CSplitCriterion, predictor *core.Attribute, preSplit util.Vector) *SplitSuggestion {
	ncols := o.postSplit.NumTargets()
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
	for ti, obs := range o.postSplit {
		obs.ForEach(func(pi int, v float64) { m.Incr(pi, ti, v) })
	}
	return m
}

// NewNumericCObserver uses gaussian estimators to monitor a numeric predictor attribute
func NewNumericCObserver(numBins int) CObserver {
	if numBins < 1 {
		numBins = 10
	}

	return &gaussianCObserver{
		numBins:   numBins,
		minMax:    *newMinMaxRanges(),
		postSplit: util.NewNumSeriesDistribution(),
	}
}

type gaussianCObserver struct {
	numBins   int
	minMax    minMaxRanges
	postSplit util.NumSeriesDistribution
}

func (o *gaussianCObserver) ByteSize() int {
	return 24 + o.minMax.ByteSize() + o.postSplit.ByteSize()
}

// Observe implements CObserver
func (o *gaussianCObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	ti := tv.Index()
	pval := pv.Value()
	o.postSplit.Append(ti, pval, weight)
	o.minMax.Update(ti, pval)
}

// Probability implements CObserver
func (o *gaussianCObserver) Probability(tv, pv core.AttributeValue) float64 {
	if est := o.postSplit.Get(tv.Index()); est != nil {
		return est.ProbDensity(pv.Value())
	}
	return 0.0
}

// BestSplit implements Observes using a variance reduction
// algorithm
func (o *gaussianCObserver) BestSplit(crit classifiers.CSplitCriterion, predictor *core.Attribute, preSplit util.Vector) *SplitSuggestion {
	var best *SplitSuggestion

	for _, splitVal := range o.minMax.SplitPoints(o.numBins) {
		postSplit := o.binarySplitOn(splitVal)
		merit := crit.Merit(preSplit, postSplit)
		if best != nil && merit <= best.merit {
			continue
		}

		best = &SplitSuggestion{
			cond:      &numericBinarySplitCondition{predictor: predictor, splitValue: splitVal},
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
	for i, est := range o.postSplit {
		if splitVal < o.minMax.Min(i) {
			res.Incr(1, i, est.TotalWeight())
		} else if splitVal >= o.minMax.Max(i) {
			res.Incr(0, i, est.TotalWeight())
		} else {
			lt, eq, gt := est.Estimate(splitVal)
			res.Incr(0, i, lt+eq)
			res.Incr(1, i, gt)
		}
	}
	return res
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
		postSplit: util.NewNumSeriesDistribution(),
	}
}

type nominalRObserver struct {
	postSplit util.NumSeriesDistribution
}

func (o *nominalRObserver) ByteSize() int {
	return 40 + o.postSplit.ByteSize()
}

// Observe implements Observer
func (o *nominalRObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	pi := pv.Index()
	o.postSplit.Append(pi, tv.Value(), weight)
}

// BestSplit implements RegressionObserves using a variance reduction
// algorithm
func (o *nominalRObserver) BestSplit(crit classifiers.RSplitCriterion, predictor *core.Attribute, preSplit *util.NumSeries) *SplitSuggestion {
	if !o.isSplitable() {
		return nil
	}

	return &SplitSuggestion{
		cond:      NewNominalMultiwaySplitCondition(predictor),
		merit:     normMerit(crit.Merit(preSplit, o.postSplit)),
		mrange:    crit.Range(preSplit),
		preStats:  newRObservationStats(preSplit),
		postStats: newRObservationStatsDist(o.postSplit),
	}
}

func (o *nominalRObserver) isSplitable() bool {
	n := 0
	for _, s := range o.postSplit {
		if s.TotalWeight() > 0 {
			if n++; n > 1 {
				return true
			}
		}
	}
	return false
}

// NewNumericRObserver uses gaussian estimators to monitor a numeric predictor
// attribute for a numeric regression target
func NewNumericRObserver(numBins int) RObserver {
	if numBins < 1 {
		numBins = 10
	}

	return &gaussianRObserver{
		numBins: numBins,
		minMax:  *newMinMaxRange(),
	}
}

type gaussianRObserver struct {
	numBins int
	minMax  minMaxRange
	tuples  []gaussianRTuple
}

type gaussianRTuple struct{ pval, tval, weight float64 }

func (o *gaussianRObserver) ByteSize() int {
	return 80 + len(o.tuples)*24
}

// Observe implements RObserver
func (o *gaussianRObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	tval, pval := tv.Value(), pv.Value()
	o.minMax.Update(pval)
	o.tuples = append(o.tuples, gaussianRTuple{
		pval:   pval,
		tval:   tval,
		weight: weight,
	})
}

func (o *gaussianRObserver) BestSplit(crit classifiers.RSplitCriterion, predictor *core.Attribute, preSplit *util.NumSeries) *SplitSuggestion {
	var best *SplitSuggestion
	for _, pivot := range o.minMax.SplitPoints(o.numBins) {
		postSplit := o.postSplit(pivot)
		merit := crit.Merit(preSplit, postSplit)
		if best != nil && merit <= best.merit {
			continue
		}

		best = &SplitSuggestion{
			cond:      &numericBinarySplitCondition{predictor: predictor, splitValue: pivot},
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
	for _, t := range o.tuples {
		if t.pval < pivot {
			res.Append(0, t.tval, t.weight)
		} else {
			res.Append(1, t.tval, t.weight)
		}
	}
	return res
}

func normMerit(merit float64) float64 {
	if merit > 0 {
		return merit
	}
	return 0.0
}
