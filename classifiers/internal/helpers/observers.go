package helpers

import (
	"github.com/bsm/reason/classifiers"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/util"
)

// Observer instances monitor and collect distribution stats
type Observer interface {
	// Observe records an instance and updates the attribute stats
	Observe(target, predictor core.AttributeValue, weight float64)

	// HeapSize estimates the required heap-size
	HeapSize() int
}

// CObserver instances monitor and collect distribution stats for
// predictor attributes and can suggest best possible splits.
type CObserver interface {
	Observer
	// Probability returns the probability of a given instance
	Probability(target, predictor core.AttributeValue) float64
	// BestSplit returns a suggestion for the best split
	BestSplit(criterion classifiers.CSplitCriterion, predictor *core.Attribute, preSplit []float64) *SplitSuggestion
}

// NewNominalCObserver monitors a nominal predictor attribute
func NewNominalCObserver() CObserver {
	return &nominalCObserver{}
}

type nominalCObserver struct {
	postSplit util.NumMatrix
}

func (o *nominalCObserver) HeapSize() int {
	return 40 + len(o.postSplit)*8
}

// Observe implements CObserver
func (o *nominalCObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	ti := tv.Index()
	row := util.NumVector(o.postSplit.GetRow(ti))
	row = row.Incr(pv.Index(), weight)
	o.postSplit = o.postSplit.SetRow(ti, row)
}

// Probability implements CObserver
func (o *nominalCObserver) Probability(tv, pv core.AttributeValue) float64 {
	obs := o.postSplit.GetRow(tv.Index())
	if obs == nil {
		return 0.0
	}

	vec := util.NumVector(obs)
	return (vec.Get(pv.Index()) + 1) / (vec.Sum() + float64(vec.Count()))
}

// BestSplit implements CObserver
func (o *nominalCObserver) BestSplit(criterion classifiers.CSplitCriterion, predictor *core.Attribute, preSplit []float64) *SplitSuggestion {
	postSplit := o.calcPostSplit()

	return &SplitSuggestion{
		cond:      NewNominalMultiwaySplitCondition(predictor),
		merit:     normMerit(criterion.Merit(preSplit, postSplit)),
		mrange:    criterion.Range(preSplit),
		preStats:  newCObservationStats(preSplit),
		postStats: newCObservationStatsSlice(postSplit),
	}
}

func (o *nominalCObserver) calcPostSplit() util.NumMatrix {
	var size int
	for _, obs := range o.postSplit {
		if n := util.NumVector(obs).Count(); n > size {
			size = n
		}
	}

	m := make(util.NumMatrix, size)
	for ti, obs := range o.postSplit {
		for pi, val := range obs {
			m[pi] = util.NumVector(m[pi]).Incr(ti, val)
		}
	}
	return m
}

// NewNumericCObserver uses gaussian estimators to monitor a numeric predictor attribute
func NewNumericCObserver(numBins int) CObserver {
	if numBins < 1 {
		numBins = 10
	}

	return &gaussianCObserver{
		numBins: numBins,
	}
}

type gaussianCObserver struct {
	numBins   int
	minMax    minMaxRanges
	postSplit []core.NumSeries
}

func (o *gaussianCObserver) HeapSize() int {
	return 96 + o.minMax.Len()*24 + len(o.postSplit)*24
}

// Observe implements CObserver
func (o *gaussianCObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	ti := tv.Index()
	if n := ti + 1; n > len(o.postSplit) {
		postSplit := make([]core.NumSeries, n)
		copy(postSplit, o.postSplit)
		o.postSplit = postSplit
	}

	pval := pv.Value()
	o.minMax.Update(ti, pval)
	o.postSplit[ti].Append(pval, weight)
}

// Probability implements CObserver
func (o *gaussianCObserver) Probability(tv, pv core.AttributeValue) float64 {
	if ti := tv.Index(); ti < len(o.postSplit) {
		if est := o.postSplit[ti]; !est.IsZero() {
			return est.ProbDensity(pv.Value())
		}
	}
	return 0.0
}

// BestSplit implements Observes using a variance reduction
// algorithm
func (o *gaussianCObserver) BestSplit(criterion classifiers.CSplitCriterion, predictor *core.Attribute, preSplit []float64) *SplitSuggestion {
	var best *SplitSuggestion

	for _, splitVal := range o.minMax.Points(o.numBins) {
		postSplit := o.binarySplitOn(splitVal)
		merit := criterion.Merit(preSplit, postSplit)
		if best != nil && merit <= best.merit {
			continue
		}

		best = &SplitSuggestion{
			cond:      &numericBinarySplitCondition{predictor: predictor, splitValue: splitVal},
			merit:     normMerit(merit),
			mrange:    criterion.Range(preSplit),
			preStats:  newCObservationStats(preSplit),
			postStats: newCObservationStatsSlice(postSplit),
		}
	}
	return best
}

func (o *gaussianCObserver) binarySplitOn(splitVal float64) util.NumMatrix {
	var lhs, rhs util.NumVector

	for i, est := range o.postSplit {
		if est.IsZero() {
			continue
		} else if splitVal < o.minMax.Min(i) {
			rhs = rhs.Incr(i, est.TotalWeight())
		} else if splitVal >= o.minMax.Max(i) {
			lhs = lhs.Incr(i, est.TotalWeight())
		} else {
			lt, eq, gt := est.Estimate(splitVal)
			lhs = lhs.Incr(i, lt+eq)
			rhs = rhs.Incr(i, gt)
		}
	}
	return util.NumMatrix{lhs, rhs}
}

// --------------------------------------------------------------------

// RObserver instances monitor stats for predictor attributes in regressions.
type RObserver interface {
	Observer
	// BestSplit returns a suggestion for the best split
	BestSplit(criterion classifiers.RSplitCriterion, predictor *core.Attribute, preSplit *core.NumSeries) *SplitSuggestion
}

// NewNominalRObserver monitors a nominal predictor attribute for a
// numeric regression target.
func NewNominalRObserver() RObserver {
	return &nominalRObserver{}
}

type nominalRObserver struct {
	postSplit []core.NumSeries
}

func (o *nominalRObserver) HeapSize() int {
	return 40 + len(o.postSplit)*24
}

// Observe implements Observer
func (o *nominalRObserver) Observe(tv, pv core.AttributeValue, weight float64) {
	pi := pv.Index()
	if n := pi + 1; n > len(o.postSplit) {
		postSplit := make([]core.NumSeries, n)
		copy(postSplit, o.postSplit)
		o.postSplit = postSplit
	}
	o.postSplit[pi].Append(tv.Value(), weight)
}

// BestSplit implements RegressionObserves using a variance reduction
// algorithm
func (o *nominalRObserver) BestSplit(criterion classifiers.RSplitCriterion, predictor *core.Attribute, preSplit *core.NumSeries) *SplitSuggestion {
	return &SplitSuggestion{
		cond:      NewNominalMultiwaySplitCondition(predictor),
		merit:     normMerit(criterion.Merit(preSplit, o.postSplit)),
		preStats:  newRObservationStats(preSplit),
		postStats: newRObservationStatsSlice(o.postSplit),
	}
}

// NewNumericRObserver uses gaussian estimators to monitor a numeric predictor
// attribute for a numeric regression target
func NewNumericRObserver(numBins int) RObserver {
	if numBins < 1 {
		numBins = 10
	}

	return &gaussianRObserver{
		numBins: numBins,
		minMax:  *util.NewNumRange(),
	}
}

type gaussianRObserver struct {
	numBins int
	minMax  util.NumRange
	tuples  []gaussianRTuple
}

type gaussianRTuple struct{ pval, tval, weight float64 }

func (o *gaussianRObserver) HeapSize() int {
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

func (o *gaussianRObserver) BestSplit(criterion classifiers.RSplitCriterion, predictor *core.Attribute, preSplit *core.NumSeries) *SplitSuggestion {
	var best *SplitSuggestion
	for _, pivot := range o.minMax.SplitPoints(o.numBins) {
		postSplit := o.postSplit(pivot)
		merit := criterion.Merit(preSplit, postSplit)
		if best != nil && merit <= best.merit {
			continue
		}

		best = &SplitSuggestion{
			cond:      &numericBinarySplitCondition{predictor: predictor, splitValue: pivot},
			merit:     normMerit(merit),
			preStats:  newRObservationStats(preSplit),
			postStats: newRObservationStatsSlice(postSplit),
		}
	}
	return best
}

func (o *gaussianRObserver) postSplit(pivot float64) []core.NumSeries {
	res := make([]core.NumSeries, 2)
	for _, t := range o.tuples {
		if t.pval < pivot {
			res[0].Append(t.tval, t.weight)
		} else {
			res[1].Append(t.tval, t.weight)
		}
	}
	return res
}

func normMerit(merit float64) float64 {
	if merit < 0 {
		return 0.0
	}
	return merit
}

// --------------------------------------------------------------------

type minMaxRanges struct {
	min util.NumVector
	max util.NumVector
	set util.BoolVector
}

func (m *minMaxRanges) Len() int            { return len(m.set) }
func (m *minMaxRanges) Min(pos int) float64 { return m.min.Get(pos) }
func (m *minMaxRanges) Max(pos int) float64 { return m.max.Get(pos) }

func (m *minMaxRanges) Update(pos int, val float64) {
	if m.set.Get(pos) {
		if val < m.Min(pos) {
			m.min = m.min.Set(pos, val)
		}
		if val > m.Max(pos) {
			m.max = m.max.Set(pos, val)
		}
	} else {
		m.min = m.min.Set(pos, val)
		m.max = m.max.Set(pos, val)
		m.set = m.set.Set(pos, true)
	}
}

func (m *minMaxRanges) Points(n int) []float64 {
	rng := util.NewNumRange()
	for i, ok := range m.set {
		if ok {
			rng.Update(m.Min(i))
			rng.Update(m.Max(i))
		}
	}
	return rng.SplitPoints(n)
}
