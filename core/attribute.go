package core

import (
	"fmt"
	"math"
	"sync"
)

const (
	// This type of attribute represents a floating-point number.
	AttributeKindNumeric AttributeKind = iota
	// This type of attribute represents a fixed set of nominal values.
	AttributeKindNominal
)

// AttributeKind defines the attribute kind
type AttributeKind uint8

func (k AttributeKind) String() string {
	switch k {
	case AttributeKindNumeric:
		return "numeric"
	case AttributeKindNominal:
		return "nominal"
	}
	return fmt.Sprintf("unknown (%d)", k)
}

// AttributeValue is an attribute value extracted from an instance
type AttributeValue float64

// MissingValue returns a missing value
func MissingValue() AttributeValue { return AttributeValue(math.NaN()) }

// IsMissing retrns true if the attribute value is missing
func (v AttributeValue) IsMissing() bool { return math.IsNaN(float64(v)) }

// Value returns the actual numeric value
func (v AttributeValue) Value() float64 { return float64(v) }

// Index returns the index of the value if the attribute value is nominal
// It will return -1 for numeric attributes or if the value is missing.
func (v AttributeValue) Index() int {
	if v.IsMissing() {
		return -1
	}
	return int(v)
}

// --------------------------------------------------------------------

// Attributes represents a simple model attribute
type Attribute struct {
	// The attribute name
	Name string
	// The attribute kind
	Kind AttributeKind
	// Values contain the attribute's possible values.
	// Only relevant for nominal attributes
	Values *AttributeValues
}

// IsNominal returns true if the attribute is nominal
func (a *Attribute) IsNominal() bool { return a.Kind == AttributeKindNominal }

// IsNumeric returns true if the attribute is numeric
func (a *Attribute) IsNumeric() bool { return !a.IsNominal() }

// Len returns the number of attribute values
func (a *Attribute) Len() int {
	if a.Kind == AttributeKindNumeric {
		return 0
	}
	return a.Values.Len()
}

// Value extracts the attribute value from an instance
func (a *Attribute) Value(inst Instance) AttributeValue {
	return a.ValueOf(inst.GetAttributeValue(a.Name))
}

// ValueOf converts an instance value an attribute value
func (a *Attribute) ValueOf(v InstanceValue) AttributeValue {
	switch a.Kind {
	case AttributeKindNumeric:
		switch n := v.(type) {
		case float64:
			return AttributeValue(n)
		case float32:
			return AttributeValue(n)
		case int:
			return AttributeValue(n)
		case int64:
			return AttributeValue(n)
		case int32:
			return AttributeValue(n)
		case int16:
			return AttributeValue(n)
		case int8:
			return AttributeValue(n)
		case uint:
			return AttributeValue(n)
		case uint64:
			return AttributeValue(n)
		case uint32:
			return AttributeValue(n)
		case uint16:
			return AttributeValue(n)
		case uint8:
			return AttributeValue(n)
		}
	case AttributeKindNominal:
		if a.Values == nil {
			a.Values = NewAttributeValues()
		}

		switch s := v.(type) {
		case string:
			return AttributeValue(a.Values.IndexOf(s))
		case []byte:
			return AttributeValue(a.Values.IndexOf(string(s)))
		}
	}
	return MissingValue()
}

// --------------------------------------------------------------------

// AttributeValues hold a slice of possible values
type AttributeValues struct {
	vi map[string]int
	iv []string
	mu sync.RWMutex
}

// NewAttributeValues inits AttributeValues
func NewAttributeValues(vals ...string) *AttributeValues {
	v := &AttributeValues{
		vi: make(map[string]int, len(vals)),
	}
	for i, val := range vals {
		v.vi[val] = i
	}
	return v
}

// Len returns the number of known values
func (v *AttributeValues) Len() int {
	if v == nil {
		return 0
	}

	v.mu.RLock()
	count := len(v.vi)
	v.mu.RUnlock()
	return count
}

// Values returns the values as a slice
func (v *AttributeValues) Values() []string {
	if v == nil {
		return nil
	}

	v.mu.RLock()
	vals := v.iv
	v.mu.RUnlock()

	if vals != nil {
		return vals
	}

	v.mu.RLock()
	vals = make([]string, len(v.vi))
	for val, i := range v.vi {
		vals[i] = val
	}
	v.mu.RUnlock()

	v.mu.Lock()
	if v.iv == nil {
		v.iv = vals
	} else {
		vals = v.iv
	}
	v.mu.Unlock()
	return vals
}

// ValueAt returns the value at an index
func (v *AttributeValues) ValueAt(i int) (string, bool) {
	if v != nil && i > -1 {
		if vals := v.Values(); i < len(vals) {
			return vals[i], true
		}
	}
	return "", false
}

// IndexOf returns the index of a value. If the value has not
// been seen before, it will be appended and a new index will
// be returned.
func (v *AttributeValues) IndexOf(s string) int {
	v.mu.RLock()
	n, ok := v.vi[s]
	v.mu.RUnlock()
	if ok {
		return n
	}

	v.mu.Lock()
	if n, ok = v.vi[s]; !ok {
		n = len(v.vi)
		v.vi[s] = n
	}
	v.mu.Unlock()
	return n
}
