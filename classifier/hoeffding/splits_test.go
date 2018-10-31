package hoeffding_test

import (
	"github.com/bsm/reason/classifier/hoeffding"
	"github.com/bsm/reason/util"
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
		var subject = hoeffding.GiniImpurity{}
		var _ hoeffding.SplitCriterion = subject

		It("should evaluate split (classification)", func() {
			Expect(subject.ClassificationMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.338, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.311, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost3)).To(BeNumerically("~", 0.480, 0.001))
		})

		It("should evaluate split (regression)", func() {
			Expect(subject.RegressionMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.RegressionMerit(regpre, regpost1)).To(Equal(0.0))
			Expect(subject.RegressionMerit(regpre, regpost2)).To(Equal(0.0))
		})
	})

	Describe("InformationGain", func() {
		var subject = hoeffding.InformationGain{MinBranchFraction: 0.1}
		var _ hoeffding.SplitCriterion = subject

		It("should evaluate split (classification)", func() {
			Expect(subject.ClassificationMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.336, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost3)).To(Equal(0.0))

			x := hoeffding.InformationGain{MinBranchFraction: 0.3}
			Expect(x.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))

			x = hoeffding.InformationGain{MinBranchFraction: 0.35}
			Expect(x.ClassificationMerit(clspre, clspost1)).To(Equal(0.0))
		})

		It("should evaluate split (regression)", func() {
			Expect(subject.RegressionMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.RegressionMerit(regpre, regpost1)).To(Equal(0.0))
			Expect(subject.RegressionMerit(regpre, regpost2)).To(Equal(0.0))
		})
	})

	Describe("VarianceReduction", func() {
		var subject = hoeffding.VarianceReduction{MinWeight: 1.0}
		var _ hoeffding.SplitCriterion = subject

		It("should evaluate split (classification)", func() {
			Expect(subject.ClassificationMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.ClassificationMerit(clspre, clspost1)).To(Equal(0.0))
			Expect(subject.ClassificationMerit(clspre, clspost2)).To(Equal(0.0))
			Expect(subject.ClassificationMerit(clspre, clspost3)).To(Equal(0.0))
		})

		It("should evaluate split (regression)", func() {
			Expect(subject.RegressionMerit(nil, nil)).To(Equal(0.0))
			Expect(subject.RegressionMerit(regpre, regpost1)).To(BeNumerically("~", 7.808, 0.001))

			c := hoeffding.VarianceReduction{MinWeight: 4.0}
			Expect(c.RegressionMerit(regpre, regpost1)).To(Equal(0.0))
		})
	})

	Describe("GainRatio", func() {
		var subject hoeffding.GainRatio
		var _ hoeffding.SplitCriterion = subject

		It("should reduce merit of 'super-attributes' (classification)", func() {
			parent := hoeffding.InformationGain{MinBranchFraction: 0.1}
			subject.SplitCriterion = parent

			Expect(parent.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.179, 0.001))

			Expect(parent.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.337, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.117, 0.001))

			Expect(parent.ClassificationMerit(clspre, clspost3)).To(BeNumerically("~", 0.0, 0.001))
			Expect(subject.ClassificationMerit(clspre, clspost3)).To(BeNumerically("~", 0.0, 0.001))
		})

		It("should reduce merit of 'super-attributes' (regression)", func() {
			parent := hoeffding.VarianceReduction{MinWeight: 4.0}
			subject.SplitCriterion = parent

			Expect(parent.RegressionMerit(regpre, regpost1)).To(BeNumerically("~", 7.808, 0.001))
			Expect(subject.RegressionMerit(regpre, regpost1)).To(BeNumerically("~", 8.181, 0.001))

			Expect(parent.RegressionMerit(regpre, regpost2)).To(BeNumerically("~", 4.576, 0.001))
			Expect(subject.RegressionMerit(regpre, regpost2)).To(BeNumerically("~", 1.664, 0.001))
		})
	})
})
