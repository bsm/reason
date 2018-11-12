package ftrl_test

import (
	"fmt"

	"github.com/bsm/reason"
	"github.com/bsm/reason/classifier/ftrl"
)

func Example() {
	target := reason.NewCategoricalFeature("play", []string{"yes", "no"})
	model := reason.NewModel(
		target,
		reason.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"}),
		reason.NewCategoricalFeature("temp", []string{"hot", "mild", "cool"}),
		reason.NewCategoricalFeature("humidity", []string{"normal", "high"}),
		reason.NewCategoricalFeature("windy", []string{"true", "false"}),
	)

	examples := []reason.MapExample{
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
	opt, err := ftrl.New(model, target.Name, nil)
	if err != nil {
		panic(err)
	}

	// Train
	for epoch := 0; epoch < 2000; epoch++ {
		for _, x := range examples {
			opt.Train(x)
		}
	}

	// Predict
	prediction := opt.Predict(reason.MapExample{
		"outlook":  "rainy",
		"temp":     "mild",
		"humidity": "high",
		"windy":    "false",
	})

	// Print categories with probabilities
	fmt.Printf("yes: %.2f\n", prediction.Prob(target.CategoryOf("yes")))
	fmt.Printf(" no: %.2f\n", prediction.Prob(target.CategoryOf("no")))

	// Output:
	// yes: 0.41
	//  no: 0.59
}
