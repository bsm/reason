package classifier

import "github.com/bsm/reason/core"

// Trainable supports training
type Trainable interface {
	// Train presents the classifier with an example and a weight.
	Train(x core.Example, weight float64)
}
