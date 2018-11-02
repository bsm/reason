package util_test

import (
	"math"
	"math/rand"

	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("NumStream", func() {
	var subject, weight1, blank *util.NumStream

	BeforeEach(func() {
		subject = util.NewNumStream()
		for _, v := range []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7, 8.8, 9.9} {
			subject.Observe(v)
		}

		weight1 = util.NewNumStream()
		weight1.Observe(5.4)

		blank = util.NewNumStream()
	})

	It("should observe", func() {
		Expect(subject).To(Equal(&util.NumStream{
			Weight:     9,
			Sum:        49.5,
			SumSquares: 344.84999999999997,
			Min:        1.1,
			Max:        9.9,
		}))
	})

	It("should calculate total weight", func() {
		Expect(subject.Weight).To(Equal(9.0))
		subject.ObserveWeight(2.2, 2.0)
		Expect(subject.Weight).To(Equal(11.0))
		Expect(blank.Weight).To(Equal(0.0))
	})

	It("should calculate value sum", func() {
		Expect(subject.Sum).To(Equal(49.5))
		Expect(blank.Sum).To(Equal(0.0))
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

var _ = Describe("NumStreams", func() {
	var subject *util.NumStreams

	BeforeEach(func() {
		subject = util.NewNumStreams()
		for _, v := range []float64{1.1, 2.2, 3.3, 4.4} {
			subject.Observe(0, v)
		}
		for _, v := range []float64{5.5, 6.6, 7.7, 8.8, 9.9} {
			subject.Observe(2, v)
		}
	})

	It("should return stats by category", func() {
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

var _ = Describe("NumStreamBuckets", func() {
	var subject *util.NumStreamBuckets

	BeforeEach(func() {
		subject = util.NewNumStreamBuckets(4)
		rnd := rand.New(rand.NewSource(10))
		for _, tv := range []float64{39, 15, 43, 7, 43, 36, 47, 6, 40, 49, 41} {
			pv := math.Ceil(tv * (rnd.Float64() + 0.5))
			subject.Observe(tv, pv)
		}
	})

	It("should observe", func() {
		Expect(subject.WeightSum()).To(Equal(11.0))

		Expect(subject.Buckets).To(HaveLen(4))
		Expect(subject.Buckets[0].Threshold).To(Equal(6.5))
		Expect(subject.Buckets[0].Weight).To(Equal(2.0))
		Expect(subject.Buckets[0].Sum).To(Equal(15.0))

		Expect(subject.Buckets[1].Threshold).To(Equal(15.0))
		Expect(subject.Buckets[1].Weight).To(Equal(1.0))
		Expect(subject.Buckets[1].Sum).To(Equal(14.0))

		Expect(subject.Buckets[2].Threshold).To(Equal(39.0))
		Expect(subject.Buckets[2].Weight).To(Equal(4.0))
		Expect(subject.Buckets[2].Sum).To(Equal(204.0))

		Expect(subject.Buckets[3].Threshold).To(Equal(45.5))
		Expect(subject.Buckets[3].Weight).To(Equal(4.0))
		Expect(subject.Buckets[3].Sum).To(Equal(172.0))
	})
})
