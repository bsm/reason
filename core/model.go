package core

// Model represents a model of a domain with specific attributes
type Model struct {
	target     *Attribute
	predictors []*Attribute
	lookup     map[string]int
}

// NewModel creates a new model with attributes
func NewModel(target, predictor *Attribute, predictors ...*Attribute) *Model {
	lookup := make(map[string]int, len(predictors)+1)
	lookup[predictor.Name] = 0
	for i, attr := range predictors {
		lookup[attr.Name] = i + 1
	}
	return &Model{
		target:     target,
		predictors: append([]*Attribute{predictor}, predictors...),
		lookup:     lookup,
	}
}

// NumPredictors returns the number of predictors
func (m *Model) NumPredictors() int { return len(m.predictors) }

// Attribute returns an attribute by name. May return nil
func (m *Model) Attribute(name string) *Attribute {
	if m.target.Name == name {
		return m.Target()
	}
	return m.Predictor(name)
}

// Attributes returns all predictor attributes
func (m *Model) Predictors() []*Attribute { return m.predictors }

// Predictor returns a predictor by name, may return nil
func (m *Model) Predictor(name string) *Attribute {
	if index, ok := m.lookup[name]; ok {
		return m.predictors[index]
	}
	return nil
}

// Target returns the target attribute
func (m *Model) Target() *Attribute {
	return m.target
}

// IsClassification returns true if the target is a nominal class
func (m *Model) IsClassification() bool { return m.target.IsNominal() }

// IsRegression returns true if the target is a numeric value
func (m *Model) IsRegression() bool { return m.target.IsNumeric() }
