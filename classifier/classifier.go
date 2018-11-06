package classifier

import (
	"github.com/bsm/reason/core"
)

// SupervisedLearner supports supervised training.
type SupervisedLearner interface {
	// TrainWeight presents the classifier with an example and a weight.
	TrainWeight(x core.Example, weight float64)
}

// Binary supports binary classification.
type Binary interface {
	// Predict returns the probability of the primary outcome.
	Predict(core.Example) float64
}

// MultiCategory supports category classification.
type MultiCategory interface {
	// PredictMC returns the classification.
	PredictMC(core.Example) MultiCategoryClassification
}

// Regressor supports simple regression.
type Regressor interface {
	// PredictNum returns the predicted regression.
	PredictNum(core.Example) Regression
}

// --------------------------------------------------------------------

// MultiCategoryClassification results allow access to individual
// probabilities of all possible categories.
type MultiCategoryClassification interface {
	// Category returns the most probable category.
	Category() core.Category
	// Prob returns the probability of the given category.
	Prob(core.Category) float64
	// Weight is the weight of the observations that have contributed to
	// this result.
	Weight() float64
}

// Regression is a regression prediction result.
type Regression interface {
	// Number returns the predicted regression value.
	Number() float64
	// MSE the  mean squared error of the prediction based on
	// the observations made.
	MSE() float64
	// Weight is the weight of the observations that have contributed to
	// this result.
	Weight() float64
}
