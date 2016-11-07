package core

// InstanceValue is a raw attribute value returned by an instance
type InstanceValue interface{}

// Instance is an interface representing a model instance
type Instance interface {
	// GetAttributeValue must return a value for an attribute name
	GetAttributeValue(string) InstanceValue
	// GetInstanceWeight returns a numeric weight which should default to 1.0
	GetInstanceWeight() float64
}

// MapInstance represents and instance as key-value pairs
// To specify weight, assign a @weight attribute with a float64 value.
type MapInstance map[string]InstanceValue

func (m MapInstance) GetAttributeValue(key string) InstanceValue { return m[key] }
func (m MapInstance) GetInstanceWeight() float64 {
	if v, ok := m["@weight"]; ok {
		if w, ok := v.(float64); ok {
			return w
		}
	}
	return 1.0
}
