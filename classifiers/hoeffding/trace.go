package hoeffding

import (
	"fmt"
	"strings"
)

type Trace struct {
	SplitPredictor string
	MeritGain      float64
	HoeffdingBound float64
	PossibleSplits []TracePossibleSplit
}

type TracePossibleSplit struct {
	Predictor string
	Merit     float64
}

func (t Trace) String() string {
	splits := make([]string, 0, len(t.PossibleSplits))
	for _, split := range t.PossibleSplits {
		splits = append(splits, fmt.Sprintf("%s: %f", split.Predictor, split.Merit))
	}
	return fmt.Sprintf("Predictor: %s, MeritGain: %f, HBound: %f, Splits: [%s]", t.SplitPredictor, t.MeritGain, t.HoeffdingBound, strings.Join(splits, ", "))
}
