package core

import "sort"

// PredictedValue represents a predicted attribute value
type PredictedValue struct {
	Value AttributeValue
	Votes float64
}

// Prediction is a slice of predicted values
type Prediction []PredictedValue

// Rank sorts the predicted values by votes,
// heighest first
func (p Prediction) Rank() {
	sort.Sort(sort.Reverse(p))
}

// Top returns the predicted value with the highest votes
func (p Prediction) Top() PredictedValue {
	if len(p) == 0 {
		return PredictedValue{Value: MissingValue()}
	}

	if !sort.IsSorted(sort.Reverse(p)) {
		p.Rank()
	}
	return p[0]
}

func (p Prediction) Len() int           { return len(p) }
func (p Prediction) Less(i, j int) bool { return p[i].Votes < p[j].Votes }
func (p Prediction) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
