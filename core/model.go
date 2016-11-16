package core

import (
	"bytes"
	"encoding/gob"
)

// Model represents a model of a domain with specific attributes
type Model struct {
	target     *Attribute
	predictors []*Attribute
	lookup     map[string]int
}

type modelSnapshot struct {
	Target     *Attribute
	Predictors []*Attribute
}

// NewModel creates a new model with attributes
func NewModel(target, predictor *Attribute, predictors ...*Attribute) *Model {
	m := &Model{
		target:     target,
		predictors: append([]*Attribute{predictor}, predictors...),
	}
	m.postInit()
	return m
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

// GobEncode implements gob.GobEncoder
func (m *Model) GobEncode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(modelSnapshot{
		Target:     m.target,
		Predictors: m.predictors,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GobDecode implements gob.GobDecoder
func (m *Model) GobDecode(b []byte) error {
	var snap modelSnapshot
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&snap); err != nil {
		return err
	}

	*m = Model{
		target:     snap.Target,
		predictors: snap.Predictors,
	}
	m.postInit()
	return nil
}

func (m *Model) postInit() {
	m.lookup = make(map[string]int, len(m.predictors))
	for i, attr := range m.predictors {
		m.lookup[attr.Name] = i
	}
}
