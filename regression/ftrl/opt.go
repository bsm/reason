package ftrl

import (
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/regression/ftrl/internal"
)

// Optimizer represents an FTRL model
type Optimizer struct {
	opt    *internal.Optimizer
	target *core.Feature
	mu     sync.RWMutex
}

// Load loads an Optimizer from a reader.
func Load(r io.Reader) (*Optimizer, error) {
	opt := new(internal.Optimizer)
	if _, err := opt.ReadFrom(r); err != nil {
		return nil, err
	}
	return newOptimizer(opt)
}

// New inits a new Optimizer using a model, a target feature and a config.
func New(model *core.Model, target string, config *Config) (*Optimizer, error) {
	var conf Config
	if config != nil {
		conf = *config
	}
	conf.Norm()

	return newOptimizer(internal.NewOptimizer(model, target, conf.proto()))
}

func newOptimizer(o *internal.Optimizer) (*Optimizer, error) {
	feat := o.Model.Feature(o.Target)
	if feat == nil {
		return nil, fmt.Errorf("ftrl: unknown feature %q", o.Target)
	} else if !feat.Kind.IsNumerical() {
		return nil, fmt.Errorf("ftrl: feature %q is not numerical", o.Target)
	}
	return &Optimizer{opt: o, target: feat}, nil
}

func (o *Optimizer) predict(x core.Example, t map[int]float64) float64 {
	var wTx float64

	for _, feat := range o.opt.Model.Features {
		if feat.Name == o.opt.Target {
			continue
		}

		bucket, val := featureBV(feat, x, o.opt.Config.HashBuckets)
		if bucket < 0 {
			continue
		}

		if math.Abs(o.opt.Weights[bucket]) <= o.opt.Config.L1 {
			t[bucket] = 0
		} else {
			sign := 1.0
			if o.opt.Weights[bucket] < 0 {
				sign = -1.0
			}
			t[bucket] = -(o.opt.Weights[bucket] - sign*o.opt.Config.L1) /
				(o.opt.Config.L2 + (o.opt.Config.Beta+math.Sqrt(o.opt.Sums[bucket]))/o.opt.Config.Alpha)
		}
		wTx += t[bucket] * val
	}
	return 1 / (1 + math.Exp(-math.Max(math.Min(wTx, 35), -35)))
}

// Predict performs prediction
func (o *Optimizer) Predict(x core.Example) float64 {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.predict(x, make(map[int]float64, len(o.opt.Model.Features)-1))
}

// Trains trains the optimizer with an example and a weight.
func (o *Optimizer) Train(x core.Example, weight float64) {
	if weight <= 0 {
		return
	}

	y := o.target.Number(x)
	if math.IsNaN(y) {
		return
	}

	t := make(map[int]float64, len(o.opt.Model.Features)-1)

	o.mu.Lock()
	defer o.mu.Unlock()

	p := o.predict(x, t)

	for _, feat := range o.opt.Model.Features {
		if feat.Name == o.opt.Target {
			continue
		}

		bucket, val := featureBV(feat, x, o.opt.Config.HashBuckets)
		if bucket < 0 {
			continue
		}

		// calculate gradient g
		g := (p - y) * val
		G := g * g

		// calculate sigma
		s := (math.Sqrt(o.opt.Sums[bucket]+G) - math.Sqrt(o.opt.Sums[bucket])) / o.opt.Config.Alpha

		// update
		o.opt.Weights[bucket] += g - s*t[bucket]*weight
		o.opt.Sums[bucket] += G
	}
}

// WriteTo implements io.WriterTo
func (o *Optimizer) WriteTo(w io.Writer) (int64, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.opt.WriteTo(w)
}
