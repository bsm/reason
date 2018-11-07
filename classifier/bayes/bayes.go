// Package bayes implements a Naive Bayes classifier.
package bayes

import (
	"fmt"

	"github.com/bsm/reason/common/observer"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

type featStats struct {
	C *observer.ClassificationCategorical
	N *observer.ClassificationNumerical
}

// Config contains configuration options for the Classifier.
type Config struct{}

// Naive implements a Naive Bayes classifier.
type Naive struct {
	model     *core.Model
	target    *core.Feature
	targets   util.Vector
	featStats map[string]featStats
}

// New inits a new classifer.
func New(model *core.Model, target string, config *Config) (*Naive, error) {
	feat := model.Feature(target)
	if feat == nil {
		return nil, fmt.Errorf("bayes: unknown feature %q", target)
	} else if feat.Kind != core.Feature_CATEGORICAL {
		return nil, fmt.Errorf("bayes: target %q is not suitable for classification", target)
	}

	return &Naive{
		model:     model,
		target:    feat,
		featStats: make(map[string]featStats),
	}, nil
}

// Train trains the optimizer with an example.
func (c *Naive) Train(x core.Example) {
	c.TrainWeight(x, 1.0)
}

// TrainWeight trains the optimizer with an example and a weight.
func (c *Naive) TrainWeight(x core.Example, weight float64) {
	if weight <= 0 {
		return
	}

	target := c.target.Category(x)
	if !core.IsCat(target) {
		return
	}
	c.targets.Incr(int(target), weight)

	for name, feat := range c.model.Features {
		if name == c.target.Name {
			continue
		}

		switch feat.Kind {
		case core.Feature_CATEGORICAL:
			acc := c.featStats[name]
			if acc.C == nil {
				acc = featStats{C: observer.NewClassificationCategorical()}
				c.featStats[name] = acc
			}
			acc.C.ObserveWeight(feat.Category(x), target, weight)
		case core.Feature_NUMERICAL:
			acc := c.featStats[name]
			if acc.N == nil {
				acc = featStats{N: observer.NewClassificationNumerical(0)}
				c.featStats[name] = acc
			}
			acc.N.ObserveWeight(feat.Number(x), target, weight)
		}
	}
}

// Predict implements a TODO.
func (c *Naive) Predict(x core.Example) float64 {
	/*
	   double[] votes = new double[observedClassDistribution.numValues()];
	   double observedClassSum = observedClassDistribution.sumOfValues();
	   for (int classIndex = 0; classIndex < votes.length; classIndex++) {
	       votes[classIndex] = observedClassDistribution.getValue(classIndex)
	               / observedClassSum;
	       for (int attIndex = 0; attIndex < inst.numAttributes() - 1; attIndex++) {
	           int instAttIndex = modelAttIndexToInstanceAttIndex(attIndex,
	                   inst);
	           AttributeClassObserver obs = attributeObservers.get(attIndex);
	           if ((obs != null) && !inst.isMissing(instAttIndex)) {
	               votes[classIndex] *= obs.probabilityOfAttributeValueGivenClass(inst.value(instAttIndex), classIndex);
	           }
	       }
	   }
	   // TODO: need logic to prevent underflow?
	   return votes

	*/

	size := c.targets.NumCols()
	votes := util.NewVector()
	sum := c.targets.WeightSum()
	if sum == 0 {
		return 0
	}

	for target := 0; target < size; target++ {
		votes.Set(target, c.targets.At(target)/sum)
		for name, stats := range c.featStats {
			feat := c.model.Feature(name)
			switch feat.Kind {
			case core.Feature_CATEGORICAL:
				if stats.C != nil {
					if cat := feat.Category(x); core.IsCat(cat) {
						prob := stats.C.Prob(cat, core.Category(target))
						votes.Set(target, votes.At(target)*prob)
					}
				}
			case core.Feature_NUMERICAL:
				if stats.N != nil {
					if val := feat.Number(x); core.IsNum(val) {
						// stats.N.Dist.
						_ = val
					}
				}
			}
		}
	}
	votes.Normalize()
	return votes.At(0)
}
