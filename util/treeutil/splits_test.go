package treeutil_test

import (
	"github.com/bsm/reason/util"
	"github.com/bsm/reason/util/treeutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SplitCriterion", func() {
	clspre := util.NewVectorFromSlice(
		9, 0, 6,
	)
	clspost1 := &util.Matrix{Stride: 3, Data: []float64{
		3, 0, 2,
		4, 0, 0,
		2, 0, 4,
	}}
	clspost2 := &util.Matrix{Stride: 3, Data: []float64{
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
	clspost3 := &util.Matrix{Stride: 3, Data: []float64{
		9, 0, 6,
	}}

	regpre := &util.NumStream{
		Weight:     8,
		Sum:        26.6,
		SumSquares: 143.24,
	}
	regpost1 := &util.NumStreams{Data: []util.NumStream{
		{Weight: 5, Sum: 6.5, SumSquares: 8.55},
		{Weight: 0, Sum: 0.0, SumSquares: 0.0},
		{Weight: 3, Sum: 20.1, SumSquares: 134.69},
	}}
	regpost2 := &util.NumStreams{Data: []util.NumStream{
		{Weight: 1, Sum: 1.1, SumSquares: 1.21},
		{Weight: 1, Sum: 1.2, SumSquares: 1.44},
		{Weight: 1, Sum: 1.3, SumSquares: 1.69},
		{Weight: 1, Sum: 1.4, SumSquares: 1.96},
		{Weight: 2, Sum: 8.1, SumSquares: 45.81},
		{Weight: 1, Sum: 6.7, SumSquares: 44.89},
		{Weight: 1, Sum: 6.8, SumSquares: 46.24},
	}}

	Describe("GiniImpurity", func() {
		var subject = treeutil.GiniImpurity{}
		var _ treeutil.SplitCriterion = subject

		It("should evaluate split (classification)", func() {
			Expect(subject.ClassificationRange(nil)).To(Equal(1.0))
			Expect(subject.ClassificationRange(clspre)).To(Equal(1.0))

			Expect(subject.ClassificationMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.338, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.311, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost3)).To(BeNumerically("~", 0.480, 0.001))
		})

		It("should evaluate split (regression)", func() {
			Expect(subject.RegressionRange(nil)).To(Equal(0.0))
			Expect(subject.RegressionRange(regpre)).To(Equal(0.0))

			Expect(subject.RegressionMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.RegressionMerit(regpre, regpost1)).To(Equal(0.0))
			Expect(subject.RegressionMerit(regpre, regpost2)).To(Equal(0.0))
		})
	})

	Describe("InformationGain", func() {
		var subject = treeutil.InformationGain{MinBranchFraction: 0.1}
		var _ treeutil.SplitCriterion = subject

		It("should evaluate split (classification)", func() {
			Expect(subject.ClassificationRange(nil)).To(Equal(1.0))
			Expect(subject.ClassificationRange(clspre)).To(Equal(1.0))
			Expect(subject.ClassificationRange(util.NewVectorFromSlice(1, 2, 3))).To(BeNumerically("~", 1.585, 0.001))

			Expect(subject.ClassificationMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.336, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost3)).To(Equal(0.0))

			x := treeutil.InformationGain{MinBranchFraction: 0.3}
			Expect(x.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))

			x = treeutil.InformationGain{MinBranchFraction: 0.35}
			Expect(x.ClassificationMerit(clspre, clspost1)).To(Equal(0.0))
		})

		It("should evaluate split (regression)", func() {
			Expect(subject.RegressionRange(nil)).To(Equal(0.0))
			Expect(subject.RegressionRange(regpre)).To(Equal(0.0))

			Expect(subject.RegressionMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.RegressionMerit(regpre, regpost1)).To(Equal(0.0))
			Expect(subject.RegressionMerit(regpre, regpost2)).To(Equal(0.0))
		})
	})

	Describe("VarianceReduction", func() {
		var subject = treeutil.VarianceReduction{MinWeight: 1.0}
		var _ treeutil.SplitCriterion = subject

		It("should evaluate split (classification)", func() {
			Expect(subject.ClassificationRange(nil)).To(Equal(0.0))
			Expect(subject.ClassificationRange(clspre)).To(Equal(0.0))

			Expect(subject.ClassificationMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.ClassificationMerit(clspre, clspost1)).To(Equal(0.0))
			Expect(subject.ClassificationMerit(clspre, clspost2)).To(Equal(0.0))
			Expect(subject.ClassificationMerit(clspre, clspost3)).To(Equal(0.0))
		})

		It("should evaluate split (regression)", func() {
			Expect(subject.RegressionRange(nil)).To(Equal(1.0))
			Expect(subject.RegressionRange(regpre)).To(Equal(1.0))

			Expect(subject.RegressionMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.RegressionMerit(regpre, regpost1)).To(BeNumerically("~", 7.808, 0.001))

			c := treeutil.VarianceReduction{MinWeight: 4.0}
			Expect(c.RegressionMerit(regpre, regpost1)).To(Equal(0.0))
		})
	})

	Describe("GainRatio", func() {
		var subject treeutil.GainRatio
		var _ treeutil.SplitCriterion = subject

		It("should reduce merit of 'super-attributes' (classification)", func() {
			parent := treeutil.InformationGain{MinBranchFraction: 0.1}
			subject.SplitCriterion = parent

			Expect(parent.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.179, 0.001))

			Expect(parent.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.337, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.117, 0.001))

			Expect(parent.ClassificationMerit(clspre, clspost3)).To(BeNumerically("~", 0.0, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost3)).To(BeNumerically("~", 0.0, 0.001))
		})

		It("should reduce merit of 'super-attributes' (regression)", func() {
			parent := treeutil.VarianceReduction{MinWeight: 1.0}
			subject.SplitCriterion = parent

			Expect(parent.RegressionMerit(regpre, regpost1)).To(BeNumerically("~", 7.808, 0.001))
			Expect(subject.RegressionMerit(regpre, regpost1)).To(BeNumerically("~", 8.181, 0.001))

			Expect(parent.RegressionMerit(regpre, regpost2)).To(BeNumerically("~", 4.576, 0.001))
			Expect(subject.RegressionMerit(regpre, regpost2)).To(BeNumerically("~", 1.664, 0.001))
		})
	})
})
