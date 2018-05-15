package ftrl_test

import (
	"fmt"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/regression/ftrl"
)

func Example() {
	model := core.NewModel(
		core.NewNumericalFeature("probability"),
		core.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"}),
		core.NewCategoricalFeature("temp", []string{"hot", "mild", "cool"}),
		core.NewCategoricalFeature("humidity", []string{"normal", "high"}),
		core.NewCategoricalFeature("windy", []string{"true", "false"}),
	)

	examples := []core.MapExample{
		{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "false", "probability": 0.25},
		{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "true", "probability": 0.30},
		{"outlook": "overcast", "temp": "hot", "humidity": "high", "windy": "false", "probability": 0.46},
		{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "false", "probability": 0.45},
		{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "false", "probability": 0.52},
		{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "true", "probability": 0.23},
		{"outlook": "overcast", "temp": "cool", "humidity": "normal", "windy": "true", "probability": 0.43},
		{"outlook": "rainy", "temp": "mild", "humidity": "high", "windy": "false", "probability": 0.35},
		{"outlook": "rainy", "temp": "cool", "humidity": "normal", "windy": "false", "probability": 0.38},
		{"outlook": "sunny", "temp": "mild", "humidity": "normal", "windy": "false", "probability": 0.46},
		{"outlook": "rainy", "temp": "mild", "humidity": "normal", "windy": "true", "probability": 0.48},
		{"outlook": "overcast", "temp": "mild", "humidity": "high", "windy": "true", "probability": 0.52},
		{"outlook": "overcast", "temp": "hot", "humidity": "normal", "windy": "false", "probability": 0.44},
		{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "true", "probability": 0.30},
	}

	// Init with a model
	opt, err := ftrl.New(model, "probability", nil)
	if err != nil {
		panic(err)
	}

	// Train
	for epoch := 0; epoch < 500; epoch++ {
		for _, x := range examples {
			opt.Train(x, 1.0)
		}
	}

	// Predict
	prediction := opt.Predict(core.MapExample{
		"outlook":  "rainy",
		"temp":     "mild",
		"humidity": "high",
		"windy":    "false",
	})

	// Print prediction
	fmt.Printf("probability: %.2f\n", prediction)

	// Output:
	// probability: 0.42
}
