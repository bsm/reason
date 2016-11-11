package calc

import (
	"math"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func() {
	var null []float64
	N := []float64{470, 600, 170, 430, 300}
	D := []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9}

	It("should calc sum", func() {
		Expect(Sum(null)).To(Equal(0.0))
		Expect(Sum(N)).To(Equal(1970.0))
		Expect(Sum(D)).To(Equal(49.5))
	})

	It("should calc min/max", func() {
		Expect(math.IsNaN(Min(null))).To(BeTrue())
		Expect(math.IsNaN(Max(null))).To(BeTrue())
		Expect(Min(N)).To(Equal(170.0))
		Expect(Max(N)).To(Equal(600.0))
	})

	It("should calc mean", func() {
		Expect(Mean(null)).To(Equal(0.0))
		Expect(Mean(N)).To(Equal(394.0))
		Expect(Mean(D)).To(Equal(5.5))
	})

	It("should calc variance", func() {
		Expect(Variance(null)).To(Equal(0.0))
		Expect(Variance(N)).To(Equal(21704.0))
		Expect(Variance(D)).To(BeNumerically("~", 8.07, 0.01))
	})

	It("should calc stddev", func() {
		Expect(StdDev(null)).To(Equal(0.0))
		Expect(StdDev(N)).To(BeNumerically("~", 147.32, 0.01))
		Expect(StdDev(D)).To(BeNumerically("~", 2.84, 0.01))
	})

	It("should calc sample variance", func() {
		Expect(SampleVariance(null)).To(Equal(0.0))
		Expect(SampleVariance(N)).To(Equal(27130.0))
		Expect(SampleVariance(D)).To(BeNumerically("~", 9.08, 0.01))
	})

	It("should calc sample stddev", func() {
		Expect(SampleStdDev(null)).To(Equal(0.0))
		Expect(SampleStdDev(N)).To(BeNumerically("~", 164.71, 0.01))
		Expect(SampleStdDev(D)).To(BeNumerically("~", 3.01, 0.01))
	})

	It("should sum matrix rows and total", func() {
		rsum, tsum := MatrixRowSumsPlusTotal([][]float64{
			{1.0, 5.0, 3.0},
			{4.0, 6.0, 2.0},
		})
		Expect(rsum).To(Equal([]float64{9, 12}))
		Expect(tsum).To(Equal(21.0))
	})

	It("should calc entropy", func() {
		Expect(Entropy(null)).To(Equal(0.0))
		Expect(Entropy(N)).To(BeNumerically("~", 2.21, 0.01))
		Expect(Entropy(D)).To(BeNumerically("~", 2.95, 0.01))

		x := append(D, D...)
		x = append(x, x...)
		Expect(Entropy(x)).To(BeNumerically("~", 4.96, 0.01))
	})

	It("should calc norm-probability", func() {
		Expect(NormProb(1.23)).To(BeNumerically("~", 0.891, 0.001))
		Expect(NormProb(0.12)).To(BeNumerically("~", 0.548, 0.001))
		Expect(NormProb(-0.76)).To(BeNumerically("~", 0.224, 0.001))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/calc")
}
