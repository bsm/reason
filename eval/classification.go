package eval

import (
	"sync"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/calc"
)

// Classification is a basic classification evaluator
type Classification struct {
	kappa *calc.KappaStat
	model *core.Model

	weight, correct float64

	mu sync.Mutex
}

// NewClassification inits a new evaluator
func NewClassification(model *core.Model) *Classification {
	return &Classification{
		model: model,
		kappa: calc.NewKappaStat(),
	}
}

// Record records a prediction
func (e *Classification) Record(inst core.Instance, prediction core.Prediction) {
	pi := prediction.Index()
	if pi < 0 {
		return
	}

	av := e.model.Target().Value(inst)
	if av.IsMissing() {
		return
	}

	ai := av.Index()
	weight := inst.GetInstanceWeight()

	e.mu.Lock()
	e.weight += weight
	if pi == ai {
		e.correct += weight
	}
	e.kappa.Record(pi, ai, weight)
	e.mu.Unlock()
}

// TotalWeight returns the total weight observed
func (e *Classification) TotalWeight() float64 {
	e.mu.Lock()
	weight := e.weight
	e.mu.Unlock()
	return weight
}

// Correct returns the fraction of correct observations
func (e *Classification) Correct() float64 {
	e.mu.Lock()
	weight, correct := e.weight, e.correct
	e.mu.Unlock()

	if weight == 0.0 {
		return 0.0
	}
	return correct / weight
}

// Kappa returns the kappa value as fraction
func (e *Classification) Kappa() float64 {
	e.mu.Lock()
	kappa := e.kappa.Kappa()
	e.mu.Unlock()
	return kappa
}
