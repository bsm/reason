package util

import (
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SparseVector", func() {
	var subject SparseVector

	BeforeEach(func() {
		subject = NewSparseVector()
		subject.Set(0, 2.0)
		subject.Set(3, 7.0)
		subject.Set(4, 9.0)
	})

	It("should set", func() {
		Expect(subject).To(Equal(SparseVector{0: 2, 3: 7, 4: 9}))
		subject.Set(6, 8.0)
		subject.Set(-2, 4.0)
		Expect(subject).To(Equal(SparseVector{0: 2, 3: 7, 4: 9, 6: 8}))
	})

	It("should get", func() {
		Expect(subject.Get(0)).To(Equal(2.0))
		Expect(subject.Get(1)).To(Equal(0.0))
		Expect(subject.Get(-1)).To(Equal(0.0))
	})

	It("should incr", func() {
		subject.Incr(0, 1.0)
		subject.Incr(1, 7.0)
		subject.Incr(-1, 6.0)
		Expect(subject).To(Equal(SparseVector{0: 3, 1: 7, 3: 7, 4: 9}))
	})

	It("should normalize", func() {
		subject.Normalize()
		Expect(subject).To(Equal(SparseVector{0: (1.0 / 9.0), 3: (7.0 / 18.0), 4: 0.5}))
	})

	It("should count", func() {
		Expect(subject.Count()).To(Equal(3))
		Expect(NewSparseVector().Count()).To(Equal(0))
	})

	It("should clear", func() {
		subject.Clear()
		Expect(subject.Count()).To(Equal(0))
	})

	It("should calculate sum", func() {
		Expect(subject.Sum()).To(Equal(18.0))
		Expect(NewSparseVector().Sum()).To(Equal(0.0))
	})

	It("should calculate min", func() {
		Expect(subject.Min()).To(Equal(2.0))
		Expect(math.IsNaN(NewSparseVector().Min())).To(BeTrue())
	})

	It("should calculate max", func() {
		Expect(subject.Max()).To(Equal(9.0))
		Expect(math.IsNaN(NewSparseVector().Max())).To(BeTrue())
	})

	It("should calculate mean", func() {
		Expect(subject.Mean()).To(Equal(6.0))
		Expect(NewSparseVector().Mean()).To(Equal(0.0))
	})

	It("should calculate variances", func() {
		Expect(subject.Variance()).To(BeNumerically("~", 8.67, 0.01))
		Expect(subject.SampleVariance()).To(BeNumerically("~", 13.0, 0.01))
	})

	It("should calculate standard deviations", func() {
		Expect(subject.StdDev()).To(BeNumerically("~", 2.94, 0.01))
		Expect(subject.SampleStdDev()).To(BeNumerically("~", 3.61, 0.01))
	})

	It("should calculate entropy", func() {
		Expect(subject.Entropy()).To(BeNumerically("~", 1.38, 0.01))
	})

})
