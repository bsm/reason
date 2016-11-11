package helpers

import (
	"fmt"
	"sort"

	"github.com/bsm/reason/core"
)

// SplitSuggestion is used for computing attribute split
// suggestions given a split condition.
type SplitSuggestion struct {
	cond      SplitCondition
	merit     float64
	mrange    float64
	preStats  ObservationStats
	postStats []ObservationStats
}

// Condition returns the conditional test
func (s *SplitSuggestion) Condition() SplitCondition {
	if s != nil {
		return s.cond
	}
	return nil
}

// Merit returns the merit and range of a possible split
func (s *SplitSuggestion) Merit() float64 {
	if s != nil {
		return s.merit
	}
	return 0.0
}

// Range returns the merit range of the split
func (s *SplitSuggestion) Range() float64 {
	if s != nil {
		return s.mrange
	}
	return 0.0
}

// PreStats returns the pre-split observation stats
func (s *SplitSuggestion) PreStats() ObservationStats {
	if s != nil {
		return s.preStats
	}
	return nil
}

// PostStats returns the post-split observation stats
func (s *SplitSuggestion) PostStats() []ObservationStats {
	if s != nil {
		return s.postStats
	}
	return nil
}

// SplitSuggestions is a slice if SplitSuggestion options
type SplitSuggestions []*SplitSuggestion

// Rank ranks suggestions, highest merit first
func (p SplitSuggestions) Rank() SplitSuggestions {
	sort.Sort(sort.Reverse(p))
	return p
}

func (p SplitSuggestions) Len() int           { return len(p) }
func (p SplitSuggestions) Less(i, j int) bool { return p[i].Merit() < p[j].Merit() }
func (p SplitSuggestions) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// --------------------------------------------------------------------

var (
	_ SplitCondition = (*nominalMultiwaySplitCondition)(nil)
	_ SplitCondition = (*numericBinarySplitCondition)(nil)
)

type SplitCondition interface {
	// Branch returns the branch index for an instance
	Branch(inst core.Instance) int
	// Predictor returns the predictor attribute
	Predictor() *core.Attribute
	// Describe returns a branch description
	Describe(branch int) string
}

// NewNominalMultiwaySplitCondition inits a new split-condition
func NewNominalMultiwaySplitCondition(predictor *core.Attribute) SplitCondition {
	return &nominalMultiwaySplitCondition{predictor: predictor}
}

type nominalMultiwaySplitCondition struct {
	predictor *core.Attribute
}

func (c *nominalMultiwaySplitCondition) Predictor() *core.Attribute { return c.predictor }
func (c *nominalMultiwaySplitCondition) Branch(inst core.Instance) int {
	return c.predictor.Value(inst).Index()
}
func (c *nominalMultiwaySplitCondition) Describe(branch int) string {
	if branch < 0 {
		return ""
	}
	if vals := c.predictor.Values.Values(); branch < len(vals) {
		return vals[branch]
	}
	return ""
}

type numericBinarySplitCondition struct {
	predictor  *core.Attribute
	splitValue float64
}

func (c *numericBinarySplitCondition) Predictor() *core.Attribute { return c.predictor }
func (c *numericBinarySplitCondition) Branch(inst core.Instance) int {
	v := c.predictor.Value(inst)
	if v.IsMissing() {
		return -1
	}

	if n := v.Value(); n > c.splitValue {
		return 1
	}
	return 0
}
func (c *numericBinarySplitCondition) Describe(branch int) string {
	if branch == 0 {
		return fmt.Sprintf("> %f", c.splitValue)
	} else if branch == 1 {
		return fmt.Sprintf("<= %f", c.splitValue)
	}
	return ""
}
