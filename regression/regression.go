// Package regression contains tools for solving regression problems
package regression

import (
	"github.com/bsm/reason/util"
)

// Predictions is a slice of predictions.
type Predictions []*util.NumStream

// Best returns the most accurate prediction.
func (pp Predictions) Best() *util.NumStream {
	if n := len(pp); n != 0 {
		s := pp[n-1]
		return s
	}
	return nil
}
