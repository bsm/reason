package regression_test

import (
	"github.com/bsm/reason/regression"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SplitCriterion", func() {
	var pre regression.Stats
	var post, post2 regression.StatsDistribution

	BeforeEach(func() {
		pre = regression.WrapStats(nil)
		post = regression.WrapStatsDistribution(nil)
		post2 = regression.WrapStatsDistribution(nil)

		for _, v := range []float64{1.1, 1.2, 1.3, 1.4, 1.5} {
			pre.Observe(v)
			post.Observe(0, v)
			for i := 0; i < 100; i++ {
				post2.Observe(i, v)
			}
		}
		for _, v := range []float64{6.6, 6.7, 6.8} {
			pre.Observe(v)
			post.Observe(1, v)
			for i := 100; i < 200; i++ {
				post2.Observe(i, v)
			}
		}
	})

	Describe("VarianceReduction", func() {
		var subject = regression.VarianceReduction{MinWeight: 1.0}
		var _ regression.SplitCriterion = subject

		It("should have range", func() {
			Expect(subject.Range(nil)).To(Equal(1.0))
			Expect(subject.Range(pre.Vector)).To(Equal(1.0))
		})

		It("should evaluate split", func() {
			Expect(subject.Merit(nil, nil)).To(Equal(0.0))
			Expect(subject.Merit(pre.Vector, post.Matrix)).To(BeNumerically("~", 7.81, 0.01))

			c := regression.VarianceReduction{MinWeight: 4.0}
			Expect(c.Merit(pre.Vector, post.Matrix)).To(Equal(0.0))
		})
	})

	Describe("GainRatio", func() {
		var base = regression.VarianceReduction{MinWeight: 1.0}
		var subject = regression.GainRatio{SplitCriterion: base}
		var _ regression.SplitCriterion = subject

		It("should reduce merit of 'super-attributes'", func() {
			Expect(base.Merit(pre.Vector, post.Matrix)).To(BeNumerically("~", 7.81, 0.01))
			Expect(base.Merit(pre.Vector, post2.Matrix)).To(BeNumerically("~", 7.81, 0.01))

			Expect(subject.Merit(pre.Vector, post.Matrix)).To(BeNumerically("~", 8.18, 0.01))
			Expect(subject.Merit(pre.Vector, post2.Matrix)).To(BeNumerically("~", 1.03, 0.01))
		})
	})
})
