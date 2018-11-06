package testdata

import "github.com/bsm/reason/core"

// SimpleDataSet contains a simple test dataset.
var SimpleDataSet = []core.MapExample{
	{"outlook": "rainy", "temp": "hot", "humidity": "high", "humidex": 61, "windy": "false", "hours": 25, "play": "no"},
	{"outlook": "rainy", "temp": "hot", "humidity": "high", "humidex": 58, "windy": "true", "hours": 30, "play": "no"},
	{"outlook": "overcast", "temp": "hot", "humidity": "high", "humidex": 60, "windy": "false", "hours": 46, "play": "yes"},
	{"outlook": "sunny", "temp": "mild", "humidity": "high", "humidex": 61, "windy": "false", "hours": 45, "play": "yes"},
	{"outlook": "sunny", "temp": "cool", "humidity": "normal", "humidex": 35, "windy": "false", "hours": 52, "play": "yes"},
	{"outlook": "sunny", "temp": "cool", "humidity": "normal", "humidex": 43, "windy": "true", "hours": 23, "play": "no"},
	{"outlook": "overcast", "temp": "cool", "humidity": "normal", "humidex": 41, "windy": "true", "hours": 43, "play": "yes"},
	{"outlook": "rainy", "temp": "mild", "humidity": "high", "humidex": 61, "windy": "false", "hours": 35, "play": "no"},
	{"outlook": "rainy", "temp": "cool", "humidity": "normal", "humidex": 41, "windy": "false", "hours": 38, "play": "yes"},
	{"outlook": "sunny", "temp": "mild", "humidity": "normal", "humidex": 40, "windy": "false", "hours": 46, "play": "yes"},
	{"outlook": "rainy", "temp": "mild", "humidity": "normal", "humidex": 37, "windy": "true", "hours": 48, "play": "yes"},
	{"outlook": "overcast", "temp": "mild", "humidity": "high", "humidex": 62, "windy": "true", "hours": 52, "play": "yes"},
	{"outlook": "overcast", "temp": "hot", "humidity": "normal", "humidex": 39, "windy": "false", "hours": 44, "play": "yes"},
	{"outlook": "sunny", "temp": "mild", "humidity": "high", "humidex": 64, "windy": "true", "hours": 30, "play": "no"},
}

// SimpleModel is a simple test model.
var SimpleModel = core.NewModel(
	core.NewCategoricalFeature("play", []string{"yes", "no"}),
	core.NewNumericalFeature("hours"),
	core.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"}),
	core.NewCategoricalFeature("temp", []string{"hot", "mild", "cool"}),
	core.NewCategoricalFeature("humidity", []string{"high", "normal"}),
	core.NewNumericalFeature("humidex"),
	core.NewCategoricalFeature("windy", []string{"true", "false"}),
)

// ClassificationScore is used to compare results.
type ClassificationScore struct{ Accuracy, Kappa, LogLoss float64 }

// RegressionScore is used to compare results.
type RegressionScore struct{ R2, RMSE float64 }
