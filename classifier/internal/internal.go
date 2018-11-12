package internal

import (
	"math"

	"github.com/bsm/reason"
	"github.com/bsm/reason/classifier"
)

var (
	_ classifier.Classification = NoResult{}
	_ classifier.Regression     = NoResult{}
	_ classifier.Regression     = StdRegression(0)
)

// NoResult is a wrapper for no-result outcomes.
type NoResult struct{}

func (NoResult) Category() reason.Category    { return reason.NoCategory }
func (NoResult) Prob(reason.Category) float64 { return 0.0 }
func (NoResult) Number() float64              { return math.NaN() }

// StdRegression is a wrapper for minimal regression predictions.
type StdRegression float64

// Number implements classifier.Regression.
func (v StdRegression) Number() float64 { return float64(v) }
