package ftrl

import (
	"io"

	"github.com/bsm/reason/classification"
	"github.com/bsm/reason/core"
	regression "github.com/bsm/reason/regression/ftrl"
	"github.com/bsm/reason/util"
)

// Config configures behaviour
type Config struct {
	// Learn rate alpha parameter.
	// Default: 0.1
	Alpha float64
	// Learn rate beta parameter.
	// Default: 1.0
	Beta float64
	// Regularization strength #1.
	// Default: 1.0
	L1 float64
	// Regularization strength #2.
	// Default: 0.1
	L2 float64
}

// Optimizer is a thin wrapper around the regression/ftrl.Optimizer
// with convenience methods for classifications.
type Optimizer struct{ *regression.Optimizer }

// Load loads an Optimizer from a reader.
func Load(r io.Reader, config *Config) (*Optimizer, error) {
	opt, err := regression.Load(r, (*regression.Config)(config))
	if err != nil {
		return nil, err
	}
	return &Optimizer{Optimizer: opt}, nil
}

// New inits a new Optimizer using a model, a target feature and a config.
func New(model *core.Model, target string, config *Config) (*Optimizer, error) {
	opt, err := regression.New(model, target, (*regression.Config)(config))
	if err != nil {
		return nil, err
	}
	return &Optimizer{Optimizer: opt}, nil
}

// Predict returns the prefiction
func (o *Optimizer) Predict(x core.Example) *classification.Prediction {
	p := o.Optimizer.Predict(x)
	return &classification.Prediction{
		Vector: util.Vector{Dense: []float64{1 - p, p}},
	}
}
