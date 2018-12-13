package ingest

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/bsm/reason"
)

const (
	featureAutoDetectLimitVocabulary = 100
	featureAutoDetectNumBuckets      = 100
)

type featureAutoDetect struct {
	name    string
	kind    reason.Feature_Kind
	numSeen int
	values  map[string]struct{}
}

func newFeatureAutoDetect(name string) *featureAutoDetect {
	return &featureAutoDetect{
		name:   name,
		kind:   reason.Feature_NUMERICAL,
		values: make(map[string]struct{}),
	}
}

func (f *featureAutoDetect) ObserveString(s string) {
	if s == "" {
		return
	}

	f.numSeen++

	if len(f.values) < featureAutoDetectLimitVocabulary {
		f.values[s] = struct{}{}
	}

	if f.kind == reason.Feature_NUMERICAL {
		if _, err := strconv.ParseFloat(s, 64); err != nil {
			f.kind = reason.Feature_CATEGORICAL
		}
	}
}

func (f *featureAutoDetect) Feature() (*reason.Feature, error) {
	if f.numSeen == 0 {
		return nil, fmt.Errorf("reason: unable to detect model feature %q", f.name)
	}

	if f.kind == reason.Feature_NUMERICAL {
		return reason.NewNumericalFeature(f.name), nil
	}

	if len(f.values) < featureAutoDetectLimitVocabulary {
		vals := make([]string, 0, len(f.values))
		for val := range f.values {
			vals = append(vals, val)
		}
		sort.Strings(vals)
		return reason.NewCategoricalFeature(f.name, vals), nil
	}

	return reason.NewCategoricalFeatureHashBuckets(f.name, featureAutoDetectNumBuckets), nil
}
