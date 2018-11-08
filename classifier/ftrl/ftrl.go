package ftrl

import (
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/bsm/reason/classifier"
	"github.com/bsm/reason/classifier/ftrl/internal"
	"github.com/bsm/reason/core"
)

var (
	_ classifier.SupervisedLearner = (*FTRL)(nil)
	_ classifier.Binary            = (*FTRL)(nil)
)

// Optimizer represents an FTRL optimiser. Regressions
// can only predict values between 0 and 1. For correct results,
// please ensure that all your target values are within that range.
type FTRL struct {
	opt        *internal.Optimizer
	target     *core.Feature
	predictors []string
	offsets    []int
	config     Config
	mu         sync.RWMutex
}

// LoadFrom loads an Optimizer from a reader.
func LoadFrom(r io.Reader, config *Config) (*FTRL, error) {
	opt := new(internal.Optimizer)
	if _, err := opt.ReadFrom(r); err != nil {
		return nil, err
	}

	predictors, offsets, _ := parseFeatures(opt.Model.Features, opt.Target)
	return newOptimizer(opt, predictors, offsets, config)
}

// New inits a new Optimizer using a model, a target feature and a config.
func New(model *core.Model, target string, config *Config) (*FTRL, error) {
	predictors, offsets, size := parseFeatures(model.Features, target)
	opt := internal.New(model, target, size)
	return newOptimizer(opt, predictors, offsets, config)
}

func newOptimizer(opt *internal.Optimizer, predictors []string, offsets []int, config *Config) (*FTRL, error) {
	feat := opt.Model.Feature(opt.Target)
	if feat == nil {
		return nil, fmt.Errorf("ftrl: unknown feature %q", opt.Target)
	}
	for _, feat := range opt.Model.Features {
		if feat.Strategy != core.Feature_VOCABULARY {
			return nil, fmt.Errorf("ftrl: feature's %q strategy %q is not supported", feat.Name, feat.Strategy.String())
		}
	}

	var conf Config
	if config != nil {
		conf = *config
	}
	conf.norm()

	return &FTRL{
		opt:        opt,
		target:     feat,
		predictors: predictors,
		offsets:    offsets,
		config:     conf,
	}, nil
}

// Predict performs a prediction.
func (o *FTRL) Predict(x core.Example) classifier.BinaryClassification {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return classifier.BinaryClassification(o.predict(x, nil))
}

// Train trains the optimizer with an example.
func (o *FTRL) Train(x core.Example) {
	o.TrainWeight(x, 1.0)
}

// TrainWeight trains the optimizer with an example and a weight.
func (o *FTRL) TrainWeight(x core.Example, weight float64) {
	if weight <= 0 {
		return
	}

	y := 0.0 // target
	switch o.target.Kind {
	case core.Feature_CATEGORICAL:
		cat := o.target.Category(x)
		if !core.IsCat(cat) {
			return
		}
		if cat > 0 {
			y = 1.0
		}
	case core.Feature_NUMERICAL:
		val := o.target.Number(x)
		if !core.IsNum(val) {
			return
		}
		y = val
	default:
		return
	}

	t := make(map[int]float64, len(o.predictors))

	o.mu.Lock()
	defer o.mu.Unlock()

	delta := o.predict(x, t) - y
	for i, name := range o.predictors {
		feat := o.opt.Model.Features[name]
		bucket, val := featureBV(feat, x, o.offsets[i])
		if bucket < 0 {
			continue
		}

		// calculate gradient
		g := delta * val
		G := g * g

		// calculate sigma
		s := (math.Sqrt(o.opt.Sums[bucket]+G) - math.Sqrt(o.opt.Sums[bucket])) / o.config.Alpha

		// update
		o.opt.Weights[bucket] += g - s*t[bucket]
		o.opt.Sums[bucket] += G
	}
}

// WriteTo implements io.WriterTo
func (o *FTRL) WriteTo(w io.Writer) (int64, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.opt.WriteTo(w)
}

func (o *FTRL) predict(x core.Example, t map[int]float64) float64 {
	var wTx float64

	for i, name := range o.predictors {
		feat := o.opt.Model.Features[name]
		bucket, val := featureBV(feat, x, o.offsets[i])
		if bucket < 0 {
			continue
		}

		sign := 1.0
		if o.opt.Weights[bucket] < 0 {
			sign = -1.0
		}
		fabs := o.opt.Weights[bucket] * sign

		if fabs <= o.config.L1 {
			if t != nil {
				t[bucket] = 0
			}
		} else {
			step := o.config.L2 + (o.config.Beta+math.Sqrt(o.opt.Sums[bucket]))/o.config.Alpha
			factor := sign * (o.config.L1 - fabs) / step
			if t != nil {
				t[bucket] = factor
			}
			wTx += factor * val
		}
	}
	return 1 / (1 + math.Exp(-math.Max(math.Min(wTx, 35), -35)))
}
