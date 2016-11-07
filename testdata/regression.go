package testdata

import "github.com/bsm/reason/core"

func RegressionModel() *core.Model {
	return core.NewModel(
		&core.Attribute{
			Name: "hours",
			Kind: core.AttributeKindNumeric,
		},
		&core.Attribute{
			Name:   "outlook",
			Kind:   core.AttributeKindNominal,
			Values: core.NewAttributeValues("rainy", "overcast", "sunny"),
		},
		&core.Attribute{
			Name:   "temp",
			Kind:   core.AttributeKindNominal,
			Values: core.NewAttributeValues("hot", "mild", "cool"),
		},
		&core.Attribute{
			Name:   "humidity",
			Kind:   core.AttributeKindNominal,
			Values: core.NewAttributeValues("high", "normal"),
		},
		&core.Attribute{
			Name:   "windy",
			Kind:   core.AttributeKindNominal,
			Values: core.NewAttributeValues("true", "false"),
		},
	)
}

func RegressionData() []core.Instance {
	return []core.Instance{
		core.MapInstance{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "false", "hours": 25},
		core.MapInstance{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "true", "hours": 30},
		core.MapInstance{"outlook": "overcast", "temp": "hot", "humidity": "high", "windy": "false", "hours": 46},
		core.MapInstance{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "false", "hours": 45},
		core.MapInstance{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "false", "hours": 52},
		core.MapInstance{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "true", "hours": 23},
		core.MapInstance{"outlook": "overcast", "temp": "cool", "humidity": "normal", "windy": "true", "hours": 43},
		core.MapInstance{"outlook": "rainy", "temp": "mild", "humidity": "high", "windy": "false", "hours": 35},
		core.MapInstance{"outlook": "rainy", "temp": "cool", "humidity": "normal", "windy": "false", "hours": 38},
		core.MapInstance{"outlook": "sunny", "temp": "mild", "humidity": "normal", "windy": "false", "hours": 46},
		core.MapInstance{"outlook": "rainy", "temp": "mild", "humidity": "normal", "windy": "true", "hours": 48},
		core.MapInstance{"outlook": "overcast", "temp": "mild", "humidity": "high", "windy": "true", "hours": 52},
		core.MapInstance{"outlook": "overcast", "temp": "hot", "humidity": "normal", "windy": "false", "hours": 44},
		core.MapInstance{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "true", "hours": 30},
	}
}
