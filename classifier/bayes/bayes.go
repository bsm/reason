// Package bayes implements a Naive Bayes classifier.
package bayes

import (
	"fmt"
	"io"
	"sync"

	"github.com/bsm/reason/classifier"
	"github.com/bsm/reason/classifier/bayes/internal"
	cinternal "github.com/bsm/reason/classifier/internal"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

var (
	_ classifier.SupervisedLearner = (*NaiveBayes)(nil)
	_ classifier.Classifier     = (*NaiveBayes)(nil)
)

// Config contains configuration options for the Classifier.
type Config struct{}

// NaiveBayes implements a Naive Bayes classifier.
type NaiveBayes struct {
	nb     *internal.NaiveBayes
	target *core.Feature

	mu sync.RWMutex
}

// LoadFrom loads a classifier from a reader.
func LoadFrom(r io.Reader, config *Config) (*NaiveBayes, error) {
	nb := new(internal.NaiveBayes)
	if _, err := nb.ReadFrom(r); err != nil {
		return nil, err
	}
	return newNaiveBayes(nb, config)
}

// New inits a new classifer.
func New(model *core.Model, target string, config *Config) (*NaiveBayes, error) {
	return newNaiveBayes(&internal.NaiveBayes{
		Model:  model,
		Target: target,
	}, config)
}

func newNaiveBayes(nb *internal.NaiveBayes, config *Config) (*NaiveBayes, error) {
	target := nb.Model.Feature(nb.Target)
	if target == nil {
		return nil, fmt.Errorf("bayes: unknown feature %q", nb.Target)
	} else if target.Kind != core.Feature_CATEGORICAL {
		return nil, fmt.Errorf("bayes: target %q is not suitable for classification", nb.Target)
	}
	return &NaiveBayes{nb: nb, target: target}, nil
}

// WriteTo implements io.WriterTo
func (b *NaiveBayes) WriteTo(w io.Writer) (int64, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.nb.WriteTo(w)
}

// Train trains the optimizer with an example.
func (b *NaiveBayes) Train(x core.Example) {
	b.TrainWeight(x, 1.0)
}

// TrainWeight trains the optimizer with an example and a weight.
func (b *NaiveBayes) TrainWeight(x core.Example, weight float64) {
	if weight <= 0 {
		return
	}

	tcat := b.target.Category(x)
	if !core.IsCat(tcat) {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.nb.ObserveWeight(x, tcat, weight)
}

// PredictMC implements classifier.MultiCategory interface.
func (b *NaiveBayes) Predict(x core.Example) classifier.Classification {
	b.mu.RLock()
	defer b.mu.RUnlock()

	sum := b.nb.TargetStats.WeightSum()
	if sum == 0 {
		return nil
	}

	ncols := b.nb.TargetStats.NumCols()
	votes := util.NewVector()
	for col := 0; col < ncols; col++ {
		votes.Set(col, b.nb.TargetStats.At(col)/sum)

		for name, wrap := range b.nb.FeatureStats {
			feat := b.nb.Model.Feature(name)
			switch feat.Kind {
			case core.Feature_CATEGORICAL:
				if stats := wrap.GetCat(); stats != nil {
					if cat := feat.Category(x); core.IsCat(cat) {
						votes.Set(col, votes.At(col)*stats.Prob(cat, core.Category(col)))
					}
				}
			case core.Feature_NUMERICAL:
				if stats := wrap.GetNum(); stats != nil {
					if val := feat.Number(x); core.IsNum(val) {
						votes.Set(col, votes.At(col)*stats.Prob(val, core.Category(col)))
					}
				}
			}
		}
	}

	pos, _ := votes.Max()
	if pos < 0 {
		return cinternal.NoResult{}
	}

	votes.Normalize()
	return prediction{cat: core.Category(pos), vv: votes}
}
