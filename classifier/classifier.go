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
	Predict(core.Example) BinaryClassification
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

// BinaryClassification is a value between 0 and 1. Values < 0.5 indicate that
// outcome 0 is more likely while values >= 0.5 suggest that the more likely
// outcome is 1.
type BinaryClassification float64

// Category returns the more likely category.
func (v BinaryClassification) Category() core.Category {
	if v >= 0.5 {
		return 1
	}
	return 0
}

// Prob returns the probability of the given category.
// Only categories 0 and 1 may yield results > 0.
func (v BinaryClassification) Prob(cat core.Category) float64 {
	switch cat {
	case 0:
		return 1 - float64(v)
	case 1:
		return float64(v)
	}
	return 0.0
}

// MultiCategoryClassification results allow access to individual
// probabilities of all possible categories.
type MultiCategoryClassification interface {
	// Category returns the most probable category.
	Category() core.Category
	// Prob returns the probability of the given category.
	Prob(core.Category) float64
	// Weight is the (optional) weight of the observations that have contributed to
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
	// Weight is the (optional) weight of the observations that have contributed to
	// this result.
	Weight() float64
}
