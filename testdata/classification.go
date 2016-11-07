package testdata

import (
	"github.com/bsm/reason/core"
)

func ClassificationModel() *core.Model {
	return core.NewModel(
		&core.Attribute{
			Name:   "play",
			Kind:   core.AttributeKindNominal,
			Values: core.NewAttributeValues("yes", "no"),
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

func ClassificationData() []core.Instance {
	return []core.Instance{
		core.MapInstance{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "false", "play": "no"},
		core.MapInstance{"outlook": "rainy", "temp": "hot", "humidity": "high", "windy": "true", "play": "no"},
		core.MapInstance{"outlook": "overcast", "temp": "hot", "humidity": "high", "windy": "false", "play": "yes"},
		core.MapInstance{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "false", "play": "yes"},
		core.MapInstance{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "false", "play": "yes"},
		core.MapInstance{"outlook": "sunny", "temp": "cool", "humidity": "normal", "windy": "true", "play": "no"},
		core.MapInstance{"outlook": "overcast", "temp": "cool", "humidity": "normal", "windy": "true", "play": "yes"},
		core.MapInstance{"outlook": "rainy", "temp": "mild", "humidity": "high", "windy": "false", "play": "no"},
		core.MapInstance{"outlook": "rainy", "temp": "cool", "humidity": "normal", "windy": "false", "play": "yes"},
		core.MapInstance{"outlook": "sunny", "temp": "mild", "humidity": "normal", "windy": "false", "play": "yes"},
		core.MapInstance{"outlook": "rainy", "temp": "mild", "humidity": "normal", "windy": "true", "play": "yes"},
		core.MapInstance{"outlook": "overcast", "temp": "mild", "humidity": "high", "windy": "true", "play": "yes"},
		core.MapInstance{"outlook": "overcast", "temp": "hot", "humidity": "normal", "windy": "false", "play": "yes"},
		core.MapInstance{"outlook": "sunny", "temp": "mild", "humidity": "high", "windy": "true", "play": "no"},
	}
}
