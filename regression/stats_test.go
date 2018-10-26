package regression_test

import (
	"math"

	"github.com/bsm/reason/regression"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stats", func() {
	var subject, weight1, blank regression.Stats

	BeforeEach(func() {
		subject = regression.WrapStats(nil)
		for _, v := range []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9} {
			subject.Observe(v)
		}

		weight1 = regression.WrapStats(nil)
		weight1.Observe(5.4)

		blank = regression.WrapStats(nil)
	})

	It("should return total weight", func() {
		Expect(subject.TotalWeight()).To(Equal(9.0))
		subject.ObserveWeight(2.2, 2.0)
		Expect(subject.TotalWeight()).To(Equal(11.0))
		Expect(blank.TotalWeight()).To(Equal(0.0))
	})

	It("should return value sum", func() {
		Expect(subject.Sum()).To(Equal(49.5))
		Expect(blank.Sum()).To(Equal(0.0))
	})

	It("should calc mean", func() {
		Expect(subject.Mean()).To(Equal(5.5))
		subject.ObserveWeight(8.8, 8)
		Expect(subject.Mean()).To(BeNumerically("~", 7.05, 0.01))
		Expect(weight1.Mean()).To(Equal(5.4))
		Expect(math.IsNaN(blank.Mean())).To(BeTrue())
	})

	It("should calc variance", func() {
		Expect(subject.Variance()).To(BeNumerically("~", 9.07, 0.01))
		Expect(math.IsNaN(weight1.Variance())).To(BeTrue())
		Expect(math.IsNaN(blank.Variance())).To(BeTrue())
	})

	It("should calc std-dev", func() {
		Expect(subject.StdDev()).To(BeNumerically("~", 3.01, 0.01))
		Expect(math.IsNaN(weight1.StdDev())).To(BeTrue())
		Expect(math.IsNaN(blank.StdDev())).To(BeTrue())
	})

	It("should calculate probability density", func() {
		Expect(subject.Prob(1.2)).To(BeNumerically("~", 0.048, 0.001))
		Expect(subject.Prob(5.5)).To(BeNumerically("~", 0.132, 0.001))
		Expect(subject.Prob(13.3)).To(BeNumerically("~", 0.005, 0.001))
		Expect(subject.Prob(24.6)).To(BeNumerically("~", 0.000, 0.001))
		Expect(math.IsNaN(weight1.Prob(5.5))).To(BeTrue())
		Expect(math.IsNaN(blank.Prob(5.5))).To(BeTrue())
	})

	DescribeTable("should estimate",
		func(v, xlt, xeq, xgt float64) {
			lt, eq, gt := subject.Estimate(v)
			Expect(lt).To(BeNumerically("~", xlt, 0.01))
			Expect(eq).To(BeNumerically("~", xeq, 0.01))
			Expect(gt).To(BeNumerically("~", xgt, 0.01))
		},
		Entry("lower end", 1.2, 0.26, 0.43, 8.31),
		Entry("close to mean", 5.4, 3.19, 1.19, 4.62),
		Entry("top end", 9.1, 7.37, 0.58, 1.04),
	)

	It("should fail to estimate on insufficient weight", func() {
		lt, eq, gt := blank.Estimate(1.2)
		Expect(math.IsNaN(lt)).To(BeTrue())
		Expect(math.IsNaN(eq)).To(BeTrue())
		Expect(math.IsNaN(gt)).To(BeTrue())

		lt, eq, gt = weight1.Estimate(1.2)
		Expect(math.IsNaN(lt)).To(BeTrue())
		Expect(math.IsNaN(eq)).To(BeTrue())
		Expect(math.IsNaN(gt)).To(BeTrue())
	})
})

var _ = Describe("StatsDistribution", func() {
	var subject regression.StatsDistribution

	BeforeEach(func() {
		subject = regression.WrapStatsDistribution(nil)
		for _, v := range []float64{1.1, 2.2, 3.3, 4.4} {
			subject.Observe(0, v)
		}
		for _, v := range []float64{5.5, 6.6, 7.7, 8.8, 9.9} {
			subject.Observe(1, v)
		}
	})

	It("should return total weight", func() {
		Expect(subject.TotalWeight(-1)).To(Equal(0.0))
		Expect(subject.TotalWeight(0)).To(Equal(4.0))
		Expect(subject.TotalWeight(1)).To(Equal(5.0))
		Expect(subject.TotalWeight(2)).To(Equal(0.0))
	})

	It("should return value sum", func() {
		Expect(subject.Sum(-1)).To(Equal(0.0))
		Expect(subject.Sum(0)).To(Equal(11.0))
		Expect(subject.Sum(1)).To(Equal(38.5))
		Expect(subject.Sum(2)).To(Equal(0.0))
	})

	It("should calc mean", func() {
		Expect(math.IsNaN(subject.Mean(-1))).To(BeTrue())
		Expect(subject.Mean(0)).To(Equal(2.75))
		Expect(subject.Mean(1)).To(Equal(7.7))
		Expect(math.IsNaN(subject.Mean(2))).To(BeTrue())
	})

	It("should calc variance", func() {
		Expect(math.IsNaN(subject.Variance(-1))).To(BeTrue())
		Expect(subject.Variance(0)).To(BeNumerically("~", 2.017, 0.001))
		Expect(subject.Variance(1)).To(BeNumerically("~", 3.025, 0.001))
		Expect(math.IsNaN(subject.Variance(2))).To(BeTrue())
	})
})
