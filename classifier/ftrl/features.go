// Package ftrl performs Follow-The-Regularized-Leader adaptive learning and binary classification.
package ftrl

import (
	"math"
	"sort"

	"github.com/bsm/reason"
)

func featureBV(feat *reason.Feature, x reason.Example, offset int) (int, float64) {
	switch feat.Kind {
	case reason.Feature_CATEGORICAL:
		if cat := feat.Category(x); cat != reason.NoCategory {
			return offset + int(cat), 1.0
		}
	case reason.Feature_NUMERICAL:
		if num := feat.Number(x); !math.IsNaN(num) {
			return offset, num
		}
	}
	return -1, 0.0
}

func parseFeatures(features map[string]*reason.Feature, target string) (predictors []string, offsets []int, size int) {
	predictors = make([]string, 0, len(features)-1)
	for _, feat := range features {
		if feat.Name != target {
			predictors = append(predictors, feat.Name)
		}
	}
	sort.Strings(predictors)

	offsets = make([]int, len(predictors))
	for i, name := range predictors {
		offsets[i] = size

		feat := features[name]
		switch feat.Kind {
		case reason.Feature_CATEGORICAL:
			size += feat.NumCategories()
		case reason.Feature_NUMERICAL:
			size += 1
		}
	}
	return
}
