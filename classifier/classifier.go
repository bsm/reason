package classifier

import (
	"github.com/bsm/reason/core"
)

// SupervisedLearner supports supervised training.
type SupervisedLearner interface {
	// TrainWeight presents the classifier with an example and a weight.
	TrainWeight(x core.Example, weight float64)
}

// BinaryClassifier supports binary classification.
type BinaryClassifier interface {
	// PredictProb returns the probability of the primary outcome.
	PredictProb(core.Example) float64
}

// MultiCategoryClassifier supports category classification.
type MultiCategoryClassifier interface {
	// PredictCategory returns the classification
	PredictCategory(core.Example) *MultiCategoryClassification
}

// Regressor supports simple regression.
type Regressor interface {
	// PredictValue returns the predicted value.
	PredictValue(core.Example) float64
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

// Regression is a minimal regression prediction result.
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
