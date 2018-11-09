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

// WeightedPrediction instances expose observation weights to support
// the accuracy of the prediction.
type WeightedPrediction interface {
	// Weight is the weight of the observations that have contributed to
	// this result.
	Weight() float64
}

// VariancePrediction instances expose variance metrics of previous observations
// to support the accuracy of the prediction.
type VariancePrediction interface {
	// MSE the mean squared error of the prediction based on
	// the observations made.
	MSE() float64
}
