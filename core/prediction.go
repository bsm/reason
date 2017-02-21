package core

import (
	"sort"
	"sync"
)

// PredictedValue represents a predicted attribute value
type PredictedValue struct {
	AttributeValue

	// Votes represents the number of votes for this prediction
	Votes float64

	// Variance returns the variance (regressions only)
	Variance float64
}

var preductionsPool sync.Pool

// Prediction is a slice of predicted values
type Prediction []PredictedValue

// NewPrediction allocated a new prediction with zero length and a given minCap.
// It will try to recycle previously released predictions.
func NewPrediction(minCap int) Prediction {
	if v := preductionsPool.Get(); v != nil {
		if p := v.(Prediction); minCap <= cap(p) {
			return p[:0]
		}
	}
	return make(Prediction, 0, minCap)
}

// Rank sorts the predicted values by votes,
// heighest first
func (p Prediction) Rank() {
	sort.Sort(sort.Reverse(p))
}

// Index is a shortcut for Top().Index()
func (p Prediction) Index() int { return p.Top().Index() }

// Value is a shortcut for Top().Value()
func (p Prediction) Value() float64 { return p.Top().Value() }

// Top returns the predicted value with the highest votes
func (p Prediction) Top() PredictedValue {
	if len(p) == 0 {
		return PredictedValue{AttributeValue: MissingValue()}
	}

	if !sort.IsSorted(sort.Reverse(p)) {
		p.Rank()
	}
	return p[0]
}

func (p Prediction) Len() int           { return len(p) }
func (p Prediction) Less(i, j int) bool { return p[i].Votes < p[j].Votes }
func (p Prediction) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Release returns to prediction to a pool. Once called the
// prediction must not be used again. Use this method with care!
func (p Prediction) Release() {
	if cap(p) != 0 {
		preductionsPool.Put(p)
	}
}
