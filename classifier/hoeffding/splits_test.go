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

	regpre := util.NewVectorFromSlice(
		8, 26.6, 143.24,
	)
	regpost1 := &util.Matrix{Stride: 3, Data: []float64{
		5, 6.5, 8.55,
		0, 0, 0,
		3, 20.1, 134.69,
	}}
	regpost2 := &util.Matrix{Stride: 3, Data: []float64{
		1, 1.1, 1.21,
		1, 1.2, 1.44,
		1, 1.3, 1.69,
		1, 1.4, 1.96,
		2, 8.1, 45.81,
		1, 6.7, 44.89,
		1, 6.8, 46.24,
	}}

	Describe("GiniImpurity", func() {
		var subject = hoeffding.GiniImpurity{}
		var _ hoeffding.SplitCriterion = subject

		It("should have range", func() {
			Expect(subject.Range(nil)).To(Equal(1.0))
			Expect(subject.Range(clspre)).To(Equal(1.0))
		})

		It("should evaluate split", func() {
			Expect(subject.Merit(nil, nil)).To(Equal(0.0))
			Expect(subject.Merit(clspre, clspost1)).To(BeNumerically("~", 0.338, 0.001))
			Expect(subject.Merit(clspre, clspost2)).To(BeNumerically("~", 0.311, 0.001))
			Expect(subject.Merit(clspre, clspost3)).To(BeNumerically("~", 0.480, 0.001))
		})
	})

	Describe("InformationGain", func() {
		var subject = hoeffding.InformationGain{MinBranchFraction: 0.1}
		var _ hoeffding.SplitCriterion = subject

		It("should have range", func() {
			Expect(subject.Range(nil)).To(Equal(1.0))
			Expect(subject.Range(clspre)).To(Equal(1.0))
			Expect(subject.Range(util.NewVectorFromSlice(1, 2, 3))).To(BeNumerically("~", 1.585, 0.001))
		})

		It("should evaluate split", func() {
			Expect(subject.Merit(nil, nil)).To(Equal(0.0))
			Expect(subject.Merit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))
			Expect(subject.Merit(clspre, clspost2)).To(BeNumerically("~", 0.336, 0.001))
			Expect(subject.Merit(clspre, clspost3)).To(Equal(0.0))

		})

		It("should calculate merit with fraction limit", func() {
			x := hoeffding.InformationGain{MinBranchFraction: 0.3}
			Expect(x.Merit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))

			x = hoeffding.InformationGain{MinBranchFraction: 0.35}
			Expect(x.Merit(clspre, clspost1)).To(Equal(0.0))
		})
	})

	Describe("VarianceReduction", func() {
		var subject = hoeffding.VarianceReduction{MinWeight: 1.0}
		var _ hoeffding.SplitCriterion = subject

		It("should have range", func() {
			Expect(subject.Range(nil)).To(Equal(1.0))
			Expect(subject.Range(regpre)).To(Equal(1.0))
		})

		It("should evaluate split", func() {
			Expect(subject.Merit(nil, nil)).To(Equal(0.0))
			Expect(subject.Merit(regpre, regpost1)).To(BeNumerically("~", 7.808, 0.001))

			c := hoeffding.VarianceReduction{MinWeight: 4.0}
			Expect(c.Merit(regpre, regpost1)).To(Equal(0.0))
		})
	})

	Describe("GainRatio (classification)", func() {
		var parent = hoeffding.InformationGain{MinBranchFraction: 0.1}
		var subject = hoeffding.GainRatio{SplitCriterion: parent}
		var _ hoeffding.SplitCriterion = subject

		It("should reduce merit of 'super-attributes'", func() {
			Expect(parent.Merit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))
			Expect(subject.Merit(clspre, clspost1)).To(BeNumerically("~", 0.179, 0.001))

			Expect(parent.Merit(clspre, clspost2)).To(BeNumerically("~", 0.337, 0.001))
			Expect(subject.Merit(clspre, clspost2)).To(BeNumerically("~", 0.117, 0.001))

			Expect(parent.Merit(clspre, clspost3)).To(BeNumerically("~", 0.0, 0.001))
			Expect(subject.Merit(clspre, clspost3)).To(BeNumerically("~", 0.0, 0.001))
		})
	})

	Describe("GainRatio (regression)", func() {
		var parent = hoeffding.VarianceReduction{MinWeight: 1.0}
		var subject = hoeffding.GainRatio{SplitCriterion: parent}
		var _ hoeffding.SplitCriterion = subject

		It("should reduce merit of 'super-attributes'", func() {
			Expect(parent.Merit(regpre, regpost1)).To(BeNumerically("~", 7.808, 0.001))
			Expect(subject.Merit(regpre, regpost1)).To(BeNumerically("~", 8.181, 0.001))

			Expect(parent.Merit(regpre, regpost2)).To(BeNumerically("~", 4.576, 0.001))
			Expect(subject.Merit(regpre, regpost2)).To(BeNumerically("~", 1.664, 0.001))
		})
	})
})
