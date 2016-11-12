package classifiers

import (
	"github.com/bsm/reason/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CSplitCriterion", func() {
	pre := []float64{9.0, 6.0}
	post := [][]float64{
		{3.0, 2.0},
		{4.0, 0.0},
		{2.0, 4.0},
	}

	It("should create default", func() {
		_, ok := DefaultSplitCriterion(false).(CSplitCriterion)
		Expect(ok).To(BeTrue())
	})

	Describe("GiniSplitCriterion", func() {
		var subject = GiniSplitCriterion{}

		It("should have range", func() {
			Expect(subject.Range(nil)).To(Equal(1.0))
			Expect(subject.Range(pre)).To(Equal(1.0))
		})

		It("should evaluate split", func() {
			Expect(subject.Merit(nil, nil)).To(Equal(0.0))
			Expect(subject.Merit(pre, post)).To(BeNumerically("~", 0.338, 0.001))
		})
	})

	Describe("InfoGainSplitCriterion", func() {
		var subject = InfoGainSplitCriterion{MinBranchFrac: 0.1}

		It("should have range", func() {
			Expect(subject.Range(nil)).To(Equal(1.0))
			Expect(subject.Range(pre)).To(Equal(1.0))
			Expect(subject.Range([]float64{1, 2, 3})).To(BeNumerically("~", 1.58, 0.01))
		})

		It("should evaluate split", func() {
			Expect(subject.Merit(nil, nil)).To(Equal(0.0))
			Expect(subject.Merit(pre, [][]float64{pre})).To(Equal(0.0))
			Expect(subject.Merit(pre, post)).To(BeNumerically("~", 0.280, 0.001))
		})

		It("should calculate merit with fraction limit", func() {
			x := InfoGainSplitCriterion{MinBranchFrac: 0.3}
			Expect(x.Merit(pre, post)).To(BeNumerically("~", 0.280, 0.001))

			x = InfoGainSplitCriterion{MinBranchFrac: 0.35}
			Expect(x.Merit(pre, post)).To(Equal(0.0))
		})
	})

})

var _ = Describe("RSplitCriterion", func() {
	var pre *core.NumSeries
	var post []core.NumSeries

	BeforeEach(func() {
		pre = new(core.NumSeries)
		post = make([]core.NumSeries, 2)

		for _, v := range []float64{1.1, 1.2, 1.3, 1.4, 1.5} {
			pre.Append(v, 1)
			post[0].Append(v, 1)
		}
		for _, v := range []float64{6.6, 6.7, 6.8} {
			pre.Append(v, 1)
			post[1].Append(v, 1)
		}
	})

	It("should create default", func() {
		_, ok := DefaultSplitCriterion(true).(RSplitCriterion)
		Expect(ok).To(BeTrue())
	})

	Describe("VarReductionSplitCriterion", func() {
		var subject = VarReductionSplitCriterion{}

		It("should have range", func() {
			Expect(subject.Range(nil)).To(Equal(1.0))
			Expect(subject.Range(pre)).To(Equal(1.0))
		})

		It("should evaluate split", func() {
			Expect(subject.Merit(nil, nil)).To(Equal(0.0))
			Expect(subject.Merit(pre, post)).To(BeNumerically("~", 6.83, 0.01))
		})
	})

})

var _ = Describe("GainRatioSplitCriterion", func() {
	pre := []float64{9.0, 6.0}
	post1 := [][]float64{
		{3.0, 2.0},
		{4.0, 0.0},
		{2.0, 4.0},
	}
	post2 := [][]float64{
		{1.0, 0.0}, {2.0, 0.0}, {1.0, 0.0}, {1.0, 0.0}, {1.0, 0.0},
		{2.0, 1.0}, {1.0, 1.0}, {0.0, 2.0}, {0.0, 1.0}, {0.0, 1.0},
	}

	base := DefaultSplitCriterion(false).(CSplitCriterion)
	subject := GainRatioSplitCriterion(base).(CSplitCriterion)

	It("should reduce merit of 'super-attributes'", func() {
		Expect(base.Merit(pre, post1)).To(BeNumerically("~", 0.28, 0.01))
		Expect(base.Merit(pre, post2)).To(BeNumerically("~", 0.65, 0.01))

		Expect(subject.Merit(pre, post1)).To(BeNumerically("~", 0.18, 0.01))
		Expect(subject.Merit(pre, post2)).To(BeNumerically("~", 0.21, 0.01))
	})

})
