// Package regression contains tools for solving regression problems
package regression

// Predictions is a slice of predictions.
type Predictions []Stats

// Best returns the most accurate prediction.
func (pp Predictions) Best() Stats {
	if n := len(pp); n != 0 {
		return pp[n-1]
	}
	return WrapStats(nil)
}
