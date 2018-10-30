package internal

import (
	util "github.com/bsm/reason/util"
)

// SplitCandidate is a candidate for a split decision
type SplitCandidate struct {
	Feature string  // the feature name
	Merit   float64 // the split merit
	Range   float64 // the split range
	Pivot   float64 // the split pivot, for binary splits

	// Pre-split stats
	PreSplit *util.Vector
	// Post-split stats
	PostSplit *util.Matrix
}

// SplitCandidates are a sortable collection of split candidates
type SplitCandidates []SplitCandidate

func (p SplitCandidates) Len() int           { return len(p) }
func (p SplitCandidates) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p SplitCandidates) Less(i, j int) bool { return p[i].Merit < p[j].Merit }

// --------------------------------------------------------------------

const numPivotBuckets = 11

func pivotPoints(min, max float64) []float64 {
	inc := (max - min) / float64(numPivotBuckets+1)
	if inc <= 0 {
		return nil
	}

	pp := make([]float64, numPivotBuckets)
	for i := 0; i < numPivotBuckets; i++ {
		pp[i] = min + inc*float64(i+1)
	}
	return pp
}

/*

func wrapVector(p classifier.Problem, vv *util.Vector) {

}

type clsVector struct{ vv *util.Vector }


func problemWeight(p classifier.Problem, stats *util.Vector) float64 {

}

func hasStats(p classifier.Problem, row []float64) bool {
	switch p {
	case classifier.Classification:
		for _, w := range row {
			if w > 0 {
				return true
			}
		}
	case classifier.Regression:
		return len(row) == 3 && row[0] > 0
	}
	return false
}
*/
