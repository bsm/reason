package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("NumSeries", func() {
	var subject *NumSeries

	BeforeEach(func() {
		subject = new(NumSeries)
		for _, v := range []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9} {
			subject.Append(v, 1)
		}
	})

	It("should return total weight", func() {
		Expect(subject.TotalWeight()).To(Equal(9.0))
		subject.Append(2.2, 2)
		Expect(subject.TotalWeight()).To(Equal(11.0))
	})

	It("should return value sum", func() {
		Expect(subject.Sum()).To(Equal(49.5))
	})

	It("should calc mean", func() {
		Expect(subject.Mean()).To(Equal(5.5))
		subject.Append(8.8, 8)
		Expect(subject.Mean()).To(BeNumerically("~", 7.05, 0.01))
	})

	It("should calc variance", func() {
		Expect(subject.Variance()).To(BeNumerically("~", 8.07, 0.01))
		subject.Append(8.8, 8)
		Expect(subject.Variance()).To(BeNumerically("~", 6.98, 0.01))
	})

	It("should calc std-dev", func() {
		Expect(subject.StdDev()).To(BeNumerically("~", 2.84, 0.01))
		subject.Append(8.8, 8)
		Expect(subject.StdDev()).To(BeNumerically("~", 2.64, 0.01))
	})

	It("should calc sample variance", func() {
		Expect(subject.SampleVariance()).To(BeNumerically("~", 9.07, 0.01))
	})

	It("should calc sample std-dev", func() {
		Expect(subject.SampleStdDev()).To(BeNumerically("~", 3.01, 0.01))
	})

	It("should calculate probability density", func() {
		Expect(subject.ProbDensity(1.2)).To(BeNumerically("~", 0.048, 0.001))
		Expect(subject.ProbDensity(5.5)).To(BeNumerically("~", 0.132, 0.001))
		Expect(subject.ProbDensity(13.3)).To(BeNumerically("~", 0.005, 0.001))
		Expect(subject.ProbDensity(24.6)).To(BeNumerically("~", 0.000, 0.001))
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

})
