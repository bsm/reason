package classifier

import (
	"fmt"
	"strings"

	"github.com/bsm/reason/core"
)

const (
	Classification Problem = iota + 1
	Regression
	problemTerminator
)

type Problem uint8

// ParseProblem parses a Problem from a name
func ParseProblem(s string) (Problem, error) {
	switch strings.ToLower(s) {
	case "c", "cls", "classification":
		return Classification, nil
	case "r", "reg", "regression":
		return Regression, nil
	}
	return 0, fmt.Errorf("reason: unable to parse problem %q", s)
}

// IsValid returns true if the problem is valid and known.
func (p Problem) IsValid() bool {
	return p >= Classification && p <= Regression
}

// String returns the Problem name.
func (p Problem) String() string {
	switch p {
	case Classification:
		return "classification"
	case Regression:
		return "regression"
	}
	return "(unknown)"
}

// --------------------------------------------------------------------

// Trainable supports training
type Trainable interface {
	// Train presents the classifier with an example and a weight.
	Train(x core.Example, weight float64)
}
