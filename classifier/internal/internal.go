package internal

import (
	"math"

	"github.com/bsm/reason/classifier"
	"github.com/bsm/reason/core"
)

var (
	_ classifier.MultiCategoryClassification = NoResult{}
	_ classifier.Regression                  = NoResult{}
)

// NoResult is a wrapper for no-result outcomes.
type NoResult struct{}

func (NoResult) Category() core.Category    { return core.NoCategory }
func (NoResult) Prob(core.Category) float64 { return 0.0 }
func (NoResult) Number() float64            { return math.NaN() }
func (NoResult) MSE() float64               { return 0 }
func (NoResult) Weight() float64            { return 0 }
