package hoeffding

import (
	"fmt"
	"strings"
)

type Trace struct {
	Split          bool
	MeritGain      float64
	HoeffdingBound float64
	PossibleSplits []TracePossibleSplit
}

type TracePossibleSplit struct {
	Predictor string
	Merit     float64
}

func (t *Trace) String() string {
	splits := make([]string, 0, 3)
	for i, split := range t.PossibleSplits {
		if i == 3 {
			break
		}
		splits = append(splits, fmt.Sprintf("%s: %.2f", split.Predictor, split.Merit))
	}
	return fmt.Sprintf("Split: %v, MeritGain: %.2f, HBound: %.2f, BestSplits: [%s]", t.Split, t.MeritGain, t.HoeffdingBound, strings.Join(splits, ", "))
}
