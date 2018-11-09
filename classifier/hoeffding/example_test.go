package hoeffding_test

import (
	"fmt"

	"github.com/bsm/reason/classifier/hoeffding"
	"github.com/bsm/reason/core"
)

func Example_classification() {
	target := core.NewCategoricalFeature("play", []string{"yes", "no"})
	model := core.NewModel(
		target,
		core.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"}),
		core.NewCategoricalFeature("temp", []string{"hot", "mild", "cool"}),
		core.NewCategoricalFeature("humidity", []string{"normal", "high"}),
		core.NewCategoricalFeature("windy", []string{"true", "false"}),
	)

	examples := []core.MapExample{
		{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "false", "play": "no"},
		{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "true", "play": "no"},
		{"outlook": "overcast", "temp": "hot", "humidity": "high", "windy": "false", "play": "yes"},
		{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "false", "play": "yes"},
		{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "false", "play": "yes"},
		{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "true", "play": "no"},
		{"outlook": "overcast", "temp": "cool", "humidity": "normal", "windy": "true", "play": "yes"},
		{"outlook": "rainy", "temp": "mild", "humidity": "high", "windy": "false", "play": "no"},
		{"outlook": "rainy", "temp": "cool", "humidity": "normal", "windy": "false", "play": "yes"},
		{"outlook": "sunny", "temp": "mild", "humidity": "normal", "windy": "false", "play": "yes"},
		{"outlook": "rainy", "temp": "mild", "humidity": "normal", "windy": "true", "play": "yes"},
		{"outlook": "overcast", "temp": "mild", "humidity": "high", "windy": "true", "play": "yes"},
		{"outlook": "overcast", "temp": "hot", "humidity": "normal", "windy": "false", "play": "yes"},
		{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "true", "play": "no"},
	}

	// Init with a model
	tree, err := hoeffding.New(model, "play", &hoeffding.Config{
		GracePeriod:     2,
		SplitConfidence: 0.1,
	})
	if err != nil {
		panic(err)
	}

	// Train
	for _, x := range examples {
		tree.Train(x)
	}

	// Predict
	prediction := tree.Predict(core.MapExample{
		"outlook":  "rainy",
		"temp":     "mild",
		"humidity": "high",
		"windy":    "false",
	})

	// Print categories with probabilities
	fmt.Printf("yes: %.2f\n", prediction.Prob(target.CategoryOf("yes")))
	fmt.Printf(" no: %.2f\n", prediction.Prob(target.CategoryOf("no")))

	// Output:
	// yes: 0.40
	//  no: 0.60
}

func Example_regression() {
	model := core.NewModel(
		core.NewNumericalFeature("hours"),
		core.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"}),
		core.NewCategoricalFeature("temp", []string{"hot", "mild", "cool"}),
		core.NewCategoricalFeature("humidity", []string{"normal", "high"}),
		core.NewCategoricalFeature("windy", []string{"true", "false"}),
	)

	examples := []core.MapExample{
		{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "false", "hours": 25},
		{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "true", "hours": 30},
		{"outlook": "overcast", "temp": "hot", "humidity": "high", "windy": "false", "hours": 46},
		{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "false", "hours": 45},
		{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "false", "hours": 52},
		{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "true", "hours": 23},
		{"outlook": "overcast", "temp": "cool", "humidity": "normal", "windy": "true", "hours": 43},
		{"outlook": "rainy", "temp": "mild", "humidity": "high", "windy": "false", "hours": 35},
		{"outlook": "rainy", "temp": "cool", "humidity": "normal", "windy": "false", "hours": 38},
		{"outlook": "sunny", "temp": "mild", "humidity": "normal", "windy": "false", "hours": 46},
		{"outlook": "rainy", "temp": "mild", "humidity": "normal", "windy": "true", "hours": 48},
		{"outlook": "overcast", "temp": "mild", "humidity": "high", "windy": "true", "hours": 52},
		{"outlook": "overcast", "temp": "hot", "humidity": "normal", "windy": "false", "hours": 44},
		{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "true", "hours": 30},
	}

	// Init with a model
	tree, err := hoeffding.New(model, "hours", &hoeffding.Config{
		GracePeriod:     2,
		SplitConfidence: 0.1,
	})
	if err != nil {
		panic(err)
	}

	// Train
	for _, x := range examples {
		tree.Train(x)
	}

	// Predict
	prediction := tree.PredictNumExtra(core.MapExample{
		"outlook":  "rainy",
		"temp":     "mild",
		"humidity": "high",
		"windy":    "false",
	})

	// Print value with weight
	fmt.Printf("value: %.2f, weight: %.0f\n",
		prediction.Number(),
		prediction.Weight())

	// Output:
	// value: 42.67, weight: 6
}
