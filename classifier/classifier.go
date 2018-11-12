package classifier

import (
	"github.com/bsm/reason/core"
)

// SupervisedLearner supports supervised training.
type SupervisedLearner interface {
	// TrainWeight presents the classifier with an example and a weight.
	TrainWeight(x core.Example, weight float64)
}

// Classifier supports category classification.
type Classifier interface {
	// Predict returns the classification.
	Predict(core.Example) Classification
}

// Regressor supports simple regression.
type Regressor interface {
	// PredictNum returns the predicted regression.
	PredictNum(core.Example) Regression
}

// Classification results are a type of prediction that allow
// access to probabilities of categories.
type Classification interface {
	// Category returns the most probable category.
	Category() core.Category
	// Prob returns the probability of the given category.
	Prob(core.Category) float64
}

// Regression is a regression prediction result.
type Regression interface {
	// Number returns the predicted regression value.
	Number() float64
}
