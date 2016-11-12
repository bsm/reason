package classifiers

import (
	"math"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/calc"
)

var (
	_ CSplitCriterion = GiniSplitCriterion{}
	_ CSplitCriterion = InfoGainSplitCriterion{}

	_ RSplitCriterion = VarReductionSplitCriterion{}
)

// SplitCriterion calculates merits of attribute splits
type SplitCriterion interface {
	isSplitCriterion()
}

// CSplitCriterion calculates the merit of an attribute split
// for classifications
type CSplitCriterion interface {
	SplitCriterion

	// Range returns the range of the split merit
	Range(pre []float64) float64

	// Merit calculates the merit of splitting for a given
	// distribution before and after the split
	Merit(pre []float64, post [][]float64) float64
}

// DefaultSplitCriterion returns InfoGainSplitCriterion
// for classification or VarReductionSplitCriterion for regressions
// (with a MinBranchFrac or 0.1)
func DefaultSplitCriterion(isRegression bool) SplitCriterion {
	if isRegression {
		return VarReductionSplitCriterion{}
	}
	return InfoGainSplitCriterion{MinBranchFrac: 0.1}
}

// --------------------------------------------------------------------

// RSplitCriterion calculates the merit of an attribute split
// for regressions
type RSplitCriterion interface {
	SplitCriterion

	// Range returns the range of the split merit
	Range(pre *core.NumSeries) float64

	// Merit calculates the merit of splitting for a given
	// distribution before and after the split
	Merit(pre *core.NumSeries, post []core.NumSeries) float64
}

// --------------------------------------------------------------------

// GiniSplitCriterion determines split merit using Gini Impurity
type GiniSplitCriterion struct{}

func (GiniSplitCriterion) isSplitCriterion()           {}
func (GiniSplitCriterion) Range(pre []float64) float64 { return 1.0 }
func (GiniSplitCriterion) Merit(pre []float64, post [][]float64) float64 {
	merit := 0.0
	sums, total := calc.MatrixRowSumsPlusTotal(post)
	for i, w := range sums {
		merit += w / total * giniSplitCalc(post[i], w)
	}
	return merit
}

func giniSplitCalc(vv []float64, sum float64) float64 {
	res := 1.0
	for _, v := range vv {
		sub := v / sum
		res -= sub * sub
	}
	return res
}

// InfoGainSplitCriterion determines split merit through information gain
type InfoGainSplitCriterion struct {
	// Requires at least two potential branches to have a sufficient
	// fractional weight in order to qualify for a split.
	MinBranchFrac float64
}

func (InfoGainSplitCriterion) isSplitCriterion() {}

func (InfoGainSplitCriterion) Range(pre []float64) float64 {
	if size := len(pre); size > 2 {
		return math.Log2(float64(size))
	}
	return math.Log2(2.0)
}

func (c InfoGainSplitCriterion) Merit(pre []float64, post [][]float64) float64 {
	sums, total := calc.MatrixRowSumsPlusTotal(post)
	if total == 0 {
		return 0.0
	}

	if min := c.MinBranchFrac; min > 0 {
		n := 0
		for _, sum := range sums {
			if sum/total > min {
				if n++; n > 1 {
					break
				}
			}
		}
		if n < 2 {
			return 0.0
		}
	}

	e1, e2 := calc.Entropy(pre), 0.0
	for i, vv := range post {
		e2 += sums[i] * calc.Entropy(vv)
	}
	return e1 - e2/total
}

// VarReductionSplitCriterion performs splits using variance-reduction
type VarReductionSplitCriterion struct{}

func (VarReductionSplitCriterion) isSplitCriterion() {}

func (VarReductionSplitCriterion) Range(_ *core.NumSeries) float64 { return 1.0 }
func (VarReductionSplitCriterion) Merit(pre *core.NumSeries, post []core.NumSeries) float64 {
	if pre == nil {
		return 0.0
	}

	total := pre.TotalWeight()
	if total == 0 {
		return 0.0
	}

	merit := pre.Variance()
	for _, n := range post {
		ratio := n.TotalWeight() / total
		merit -= n.Variance() * ratio
	}
	return merit
}

// --------------------------------------------------------------------

// GainRatioSplitCriterion normalises the merits of other split criterions
// by reducing their bias toward attributes that have a large number of
// values over attributes that have a smaller number of values.
func GainRatioSplitCriterion(c SplitCriterion) SplitCriterion {
	switch x := c.(type) {
	case CSplitCriterion:
		return cGainRatioSplitCriterion{CSplitCriterion: x}
	case RSplitCriterion:
		return rGainRatioSplitCriterion{RSplitCriterion: x}
	}
	return c
}

type cGainRatioSplitCriterion struct{ CSplitCriterion }

func (c cGainRatioSplitCriterion) Merit(pre []float64, post [][]float64) float64 {
	merit := c.CSplitCriterion.Merit(pre, post)
	total, sums := calcSumsAndTotal(post)
	return merit / gainRatioPenalty(total, sums)
}

type rGainRatioSplitCriterion struct{ RSplitCriterion }

func (c rGainRatioSplitCriterion) Merit(pre *core.NumSeries, post []core.NumSeries) float64 {
	merit := c.RSplitCriterion.Merit(pre, post)
	total := 0.0

	sums := make([]float64, len(post))
	for i, vv := range post {
		sum := vv.TotalWeight()
		total += sum
		sums[i] = sum
	}
	return merit / gainRatioPenalty(total, sums)
}

func gainRatioPenalty(total float64, sums []float64) float64 {
	pen := 0.0
	for _, sum := range sums {
		pen -= sum / total * math.Log2(sum/total)
	}
	if pen <= 0.0 {
		return 1.0
	}
	return pen
}

// --------------------------------------------------------------------

func calcSumsAndTotal(vvv [][]float64) (float64, []float64) {
	total := 0.0
	sums := make([]float64, len(vvv))

	for i, vv := range vvv {
		sum := calc.Sum(vv)
		total += sum
		sums[i] = sum
	}
	return total, sums
}
