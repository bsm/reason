package internal

import (
	"github.com/bsm/reason/util"
)

// SplitCandidate is a candidate for a split decision
type SplitCandidate struct {
	Feature   string  // the feature name
	Merit     float64 // the split merit
	Range     float64 // the split range
	Pivot     float64 // the split pivot, for binary splits
	PostSplit PostSplit
}

// SplitCandidates are a sortable collection of split candidates.
type SplitCandidates []SplitCandidate

func (p SplitCandidates) Len() int           { return len(p) }
func (p SplitCandidates) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p SplitCandidates) Less(i, j int) bool { return p[i].Merit < p[j].Merit }

// --------------------------------------------------------------------

// PostSplit contain information about the post-split distribution.
type PostSplit struct {
	Classification *util.Matrix
	Regression     *util.NumStreams
}

func (p PostSplit) forEach(iter func(int, isNode_Stats, float64)) {
	if mat := p.Classification; mat != nil {
		numRows := mat.NumRows()
		for pos := 0; pos < numRows; pos++ {
			if w := mat.RowSum(pos); w > 0 {
				ns := new(Node_ClassificationStats)
				ns.Vector = *util.NewVectorFromSlice(mat.Row(pos)...)
				iter(pos, &Node_Classification{Classification: ns}, w)
			}
		}
	}
	if streams := p.Regression; streams != nil {
		numRows := streams.NumRows()
		for pos := 0; pos < numRows; pos++ {
			if s := streams.At(pos); s != nil {
				ns := new(Node_RegressionStats)
				ns.NumStream = *s
				iter(pos, &Node_Regression{Regression: ns}, s.Weight)
			}
		}
	}
}
