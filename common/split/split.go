// Package split contains implementations of various methods for computing
// splitting criteria with respect to distributions of class values.
package split

import (
	"math"

	"github.com/bsm/reason"
	"github.com/bsm/reason/util"
)

// Criterion calculates the merit of an attribute split
// for classifications.
type Criterion interface {
	// Supports returns true if the target feature is supported.
	Supports(target *reason.Feature) bool

	// ClassificationRange applies to classifications only and
	// determines the split range.
	ClassificationRange(pre *util.Vector) float64

	// ClassificationMerit applies to classifications only and
	// calculates the merit of splitting distribution pre and post split.
	ClassificationMerit(pre *util.Vector, post *util.Matrix) float64

	// RegressionRange applies to regressions only and
	// determines the split range.
	RegressionRange(pre *util.NumStream) float64

	// RegressionMerit applies to regressions only and
	// calculates the merit of splitting distribution pre and post split.
	RegressionMerit(pre *util.NumStream, post *util.NumStreams) float64
}

// DefaultCriterion returns the default split criterion for the given target.
func DefaultCriterion(target *reason.Feature) Criterion {
	switch target.Kind {
	case reason.Feature_NUMERICAL:
		return VarianceReduction{MinWeight: 4.0}
	}
	return InformationGain{MinBranchFraction: 0.1}
}

func normSplitMerit(m float64) float64 {
	if m < 0.0 || math.IsNaN(m) {
		return 0.0
	}
	return m
}
