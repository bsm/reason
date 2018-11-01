package util_test

import (
	"math"
	"math/rand"
	"sort"

	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Histogram", func() {
	var blank, subject *util.Histogram

	BeforeEach(func() {
		blank = util.NewHistogram(0)

		subject = util.NewHistogram(4)
		for _, v := range []float64{39, 15, 43, 7, 43, 36, 47, 6, 40, 49, 41} {
			subject.Observe(v)
		}
	})

	DescribeTable("Quantile",
		func(q float64, x float64) {
			Expect(subject.Quantile(q)).To(BeNumerically("~", x, 0.1))
		},

		Entry("0%", 0.0, 6.0),
		Entry("25%", 0.25, 19.6),
		Entry("50%", 0.5, 39.8),
		Entry("75%", 0.75, 44.3),
		Entry("95%", 0.95, 47.2),
		Entry("99%", 0.99, 48.2),
		Entry("100%", 1.0, 49.0),
	)

	It("should observe", func() {
		Expect(blank.Weight).To(Equal(0.0))
		Expect(subject.Weight).To(Equal(11.0))

		Expect(subject.Min).To(Equal(6.0))
		Expect(subject.Max).To(Equal(49.0))
		Expect(subject.Bins).To(Equal([]util.Histogram_Bin{
			{Value: 6.5, Weight: -2},
			{Value: 15, Weight: 1},
			{Value: 39, Weight: -4},
			{Value: 45.5, Weight: -4},
		}))
	})

	It("should observe with weight", func() {
		subject.ObserveWeight(6.5, 2.0)
		subject.ObserveWeight(15, 3.0)
		Expect(subject.Weight).To(Equal(16.0))
		Expect(subject.Sum()).To(Equal(424.0))

		Expect(subject.Bins).To(Equal([]util.Histogram_Bin{
			{Value: 6.5, Weight: -4},
			{Value: 15, Weight: 4},
			{Value: 39, Weight: -4},
			{Value: 45.5, Weight: -4},
		}))
	})

	// inspired by https://github.com/aaw/histosketch/commit/d8284aa#diff-11101c92fbb1d58ccf30ca49764bf202R180
	// released into the public domain
	It("should accurately predict quantile", func() {
		sample := 24000
		rand := rand.New(rand.NewSource(100))
		hist := util.NewHistogram(16)
		exact := make([]float64, 0, sample)

		for i := 0; i < sample; i++ {
			num := rand.NormFloat64()
			hist.Observe(num)
			exact = append(exact, num)
		}
		sort.Float64s(exact)

		for _, q := range []float64{0.0001, 0.001, 0.01, 0.1, 0.25, 0.35, 0.45, 0.55, 0.65, 0.75, 0.9, 0.99, 0.999, 0.9999} {
			hQ := hist.Quantile(q)
			xQ := exact[int(float64(sample)*q)]
			pc := math.Abs((hQ - xQ) / xQ)
			Expect(pc).To(BeNumerically("<", 0.05),
				"s.Quantile(%v) (got %.2f, want %.2f, delta %.1f%%)", q, hQ, xQ, pc*100,
			)
		}
	})

	It("should reject bad quantile inputs", func() {
		Expect(math.IsNaN(blank.Quantile(0.5))).To(BeTrue())
		Expect(math.IsNaN(subject.Quantile(-0.1))).To(BeTrue())
		Expect(math.IsNaN(subject.Quantile(1.1))).To(BeTrue())
	})

	It("should calc sum", func() {
		Expect(math.IsNaN(blank.Sum())).To(BeTrue())
		Expect(subject.Sum()).To(Equal(366.0))
	})

	It("should calc mean", func() {
		Expect(math.IsNaN(blank.Mean())).To(BeTrue())
		Expect(subject.Mean()).To(BeNumerically("~", 33.27, 0.01))
	})
})

var _ = Describe("Histograms", func() {
	var subject *util.Histograms

	BeforeEach(func() {
		subject = util.NewHistograms(4)
		for _, v := range []float64{1.1, 2.2, 3.3, 4.4} {
			subject.Observe(0, v)
		}
		for _, v := range []float64{5.5, 6.6, 7.7, 8.8, 9.9} {
			subject.Observe(2, v)
		}
	})

	It("should return histograms by category", func() {
		Expect(subject.At(-1)).To(BeNil())
		Expect(subject.At(1)).To(BeNil())
		Expect(subject.At(3)).To(BeNil())

		Expect(subject.At(0).Weight).To(Equal(4.0))
		Expect(subject.At(2).Weight).To(Equal(5.0))
	})

	It("should calculate weight sum", func() {
		Expect(subject.WeightSum()).To(Equal(9.0))
	})

	It("should count the number of rows", func() {
		Expect(subject.NumRows()).To(Equal(3))
		subject.Observe(7, 3.3)
		Expect(subject.NumRows()).To(Equal(8))
	})

	It("should count the number of categories", func() {
		Expect(subject.NumCategories()).To(Equal(2))
		subject.Observe(7, 3.3)
		Expect(subject.NumCategories()).To(Equal(3))
	})
})
