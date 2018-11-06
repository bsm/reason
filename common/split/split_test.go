package split_test

import (
	"testing"

	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "common/split")
}

// --------------------------------------------------------------------

var (
	clspre = util.NewVectorFromSlice(
		9, 0, 6,
	)
	clspost1 = &util.Matrix{Stride: 3, Data: []float64{
		3, 0, 2,
		4, 0, 0,
		2, 0, 4,
	}}
	clspost2 = &util.Matrix{Stride: 3, Data: []float64{
		1, 0, 1,
		2, 0, 1,
		0, 0, 0,
		1, 0, 0,
		1, 0, 0,
		2, 0, 0,
		1, 0, 1,
		1, 0, 2,
		0, 0, 1,
	}}
	clspost3 = &util.Matrix{Stride: 3, Data: []float64{
		9, 0, 6,
	}}

	regpre = &util.NumStream{
		Weight:     8,
		Sum:        26.6,
		SumSquares: 143.24,
	}
	regpost1 = &util.NumStreams{Data: []util.NumStream{
		{Weight: 5, Sum: 6.5, SumSquares: 8.55},
		{Weight: 0, Sum: 0.0, SumSquares: 0.0},
		{Weight: 3, Sum: 20.1, SumSquares: 134.69},
	}}
	regpost2 = &util.NumStreams{Data: []util.NumStream{
		{Weight: 1, Sum: 1.1, SumSquares: 1.21},
		{Weight: 1, Sum: 1.2, SumSquares: 1.44},
		{Weight: 1, Sum: 1.3, SumSquares: 1.69},
		{Weight: 1, Sum: 1.4, SumSquares: 1.96},
		{Weight: 2, Sum: 8.1, SumSquares: 45.81},
		{Weight: 1, Sum: 6.7, SumSquares: 44.89},
		{Weight: 1, Sum: 6.8, SumSquares: 46.24},
	}}
)
