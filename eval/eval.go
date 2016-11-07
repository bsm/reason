package eval

import "github.com/bsm/reason/core"

// Evaluator implementation can record actual instances against predictions
type Evaluator interface {
	// Record records an actual instance with a prediction
	Record(core.Instance, core.Prediction)
	// TotalWeight records the total instance weight observed
	TotalWeight() float64
}
