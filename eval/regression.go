package eval

import (
	"math"
	"sync"

	"github.com/bsm/reason/core"
)

// Regression is a basic regression evaluator
type Regression struct {
	model *core.Model

	weight float64 // total weight observed
	sum    float64 // sum of all values

	resSum  float64 // residual sum
	resSum2 float64 // residual sum of squares
	totSum2 float64 // total sum of squares

	mu sync.Mutex
}

// NewRegression inits a new evaluator
func NewRegression(model *core.Model) *Regression {
	return &Regression{model: model}
}

// Record records a prediction
func (e *Regression) Record(inst core.Instance, prediction core.Prediction) {
	pv := prediction.Top()
	if pv.IsMissing() {
		return
	}

	av := e.model.Target().Value(inst)
	if av.IsMissing() {
		return
	}

	weight := inst.GetInstanceWeight()
	actual := av.Value()
	residual := actual - pv.Value()

	e.mu.Lock()
	if e.weight != 0 {
		delta := actual - e.sum/e.weight
		e.totSum2 += delta * delta * weight
	}

	e.resSum += math.Abs(residual) * weight
	e.resSum2 += residual * residual * weight

	e.sum += actual * weight
	e.weight += weight
	e.mu.Unlock()
}

// TotalWeight returns the total weight observed
func (e *Regression) TotalWeight() float64 {
	e.mu.Lock()
	weight := e.weight
	e.mu.Unlock()
	return weight
}

// Mean returns the mean value observed
func (e *Regression) Mean() float64 {
	e.mu.Lock()
	weight := e.weight
	sum := e.sum
	e.mu.Unlock()

	if weight > 0 {
		return sum / weight
	}
	return 0.0
}

// MAE returns the mean absolute error
func (e *Regression) MAE() float64 {
	e.mu.Lock()
	weight := e.weight
	resSum := e.resSum
	e.mu.Unlock()

	if weight > 0 {
		return resSum / weight
	}
	return 0.0
}

// MSE returns the mean square error
func (e *Regression) MSE() float64 {
	e.mu.Lock()
	weight := e.weight
	resSum2 := e.resSum2
	e.mu.Unlock()

	if weight > 0 {
		return resSum2 / weight
	}
	return 0.0
}

// RMSE returns the root mean square error
func (e *Regression) RMSE() float64 {
	return math.Sqrt(e.MSE())
}

// R2 returns the R² coefficient of determination
func (e *Regression) R2() float64 {
	e.mu.Lock()
	resSum2 := e.resSum2
	totSum2 := e.totSum2
	e.mu.Unlock()

	if totSum2 > 0 {
		return 1 - resSum2/totSum2
	}
	return 0.0
}
