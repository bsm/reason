// Package bayes implements a Naive Bayes classifier.
package bayes

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/util"
)

// Classifier implements a Naive Bayes classifier.
type Classifier struct {
	model      *core.Model
	target     *core.Feature
	categories util.Vector
	attributes util.Matrix
}

// Train trains the optimizer with an example.
func (c *Classifier) Train(x core.Example) {
	c.TrainWeight(x, 1.0)
}

// TrainWeight trains the optimizer with an example and a weight.
func (c *Classifier) TrainWeight(x core.Example, weight float64) {
	if weight <= 0 {
		return
	}

	cat := c.target.Category(x)
	if !core.IsCat(cat) {
		return
	}
	c.categories.Incr(int(cat), weight)

}
