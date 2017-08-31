package util

import (
	"bytes"
	"math"

	"github.com/bsm/reason/internal/msgpack"
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
		Expect(new(NumSeries).TotalWeight()).To(Equal(0.0))
	})

	It("should return value sum", func() {
		Expect(subject.Sum()).To(Equal(49.5))
		Expect(new(NumSeries).Sum()).To(Equal(0.0))
	})

	It("should calc mean", func() {
		Expect(subject.Mean()).To(Equal(5.5))
		subject.Append(8.8, 8)
		Expect(subject.Mean()).To(BeNumerically("~", 7.05, 0.01))
		Expect(math.IsNaN(new(NumSeries).Mean())).To(BeTrue())
	})

	It("should calc variance", func() {
		Expect(subject.Variance()).To(BeNumerically("~", 8.07, 0.01))
		subject.Append(8.8, 8)
		Expect(subject.Variance()).To(BeNumerically("~", 6.98, 0.01))
		Expect(math.IsNaN(new(NumSeries).Variance())).To(BeTrue())
	})

	It("should calc std-dev", func() {
		Expect(subject.StdDev()).To(BeNumerically("~", 2.84, 0.01))
		subject.Append(8.8, 8)
		Expect(subject.StdDev()).To(BeNumerically("~", 2.64, 0.01))
		Expect(math.IsNaN(new(NumSeries).StdDev())).To(BeTrue())
	})

	It("should calc sample variance", func() {
		Expect(subject.SampleVariance()).To(BeNumerically("~", 9.07, 0.01))
		Expect(math.IsNaN(new(NumSeries).SampleVariance())).To(BeTrue())
	})

	It("should calc sample std-dev", func() {
		Expect(subject.SampleStdDev()).To(BeNumerically("~", 3.01, 0.01))
		Expect(math.IsNaN(new(NumSeries).SampleStdDev())).To(BeTrue())
	})

	It("should calculate probability density", func() {
		Expect(subject.ProbDensity(1.2)).To(BeNumerically("~", 0.048, 0.001))
		Expect(subject.ProbDensity(5.5)).To(BeNumerically("~", 0.132, 0.001))
		Expect(subject.ProbDensity(13.3)).To(BeNumerically("~", 0.005, 0.001))
		Expect(subject.ProbDensity(24.6)).To(BeNumerically("~", 0.000, 0.001))
		Expect(math.IsNaN(new(NumSeries).ProbDensity(10.0))).To(BeTrue())
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

	It("should encode/decode", func() {
		buf := new(bytes.Buffer)
		enc := msgpack.NewEncoder(buf)
		err := enc.Encode(subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		var out *NumSeries
		err = msgpack.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})

})

var _ = Describe("NumSeriesDistribution", func() {
	var subject NumSeriesDistribution

	BeforeEach(func() {
		subject = NewNumSeriesDistribution()
		for _, v := range []float64{1.1, 2.2, 3.3, 4.4} {
			subject.Append(0, v, 1)
		}
		for _, v := range []float64{5.5, 6.6, 7.7, 8.8, 9.9} {
			subject.Append(1, v, 1)
		}
	})

	It("should append", func() {
		subject.Append(7, 12.12, 1)
		Expect(subject).To(HaveKey(7))
	})

	It("should get", func() {
		Expect(subject.Get(0)).To(Equal(&NumSeries{weight: 4, sum: 11, sumSquares: 36.3}))
		Expect(subject.Get(2)).To(BeNil())
		Expect(subject.Get(-1)).To(BeNil())
	})

	It("should return weights", func() {
		Expect(subject.Weights()).To(Equal(map[int]float64{0: 4, 1: 5}))
	})

	It("should return total weight", func() {
		Expect(subject.TotalWeight()).To(Equal(9.0))
	})

	It("should encode/decode", func() {
		buf := new(bytes.Buffer)
		enc := msgpack.NewEncoder(buf)
		err := enc.Encode(subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		var out NumSeriesDistribution
		err = msgpack.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})

})
