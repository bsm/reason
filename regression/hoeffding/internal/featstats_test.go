package internal_test

import (
	"github.com/bsm/reason/regression/hoeffding/internal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FeatureStats_Numerical", func() {
	var subject *internal.FeatureStats_Numerical

	BeforeEach(func() {
		subject = new(internal.FeatureStats_Numerical)
		subject.ObserveWeight(2.2, 1.2, 1.0)
		subject.ObserveWeight(2.4, 4.2, 1.0)
		subject.ObserveWeight(2.6, 8.4, 1.0)
	})

	It("should observe", func() {
		Expect(subject.WeightSum()).To(Equal(3.0))
		Expect(subject.Buckets).To(HaveLen(3))
	})

	It("should calculate post-splits", func() {
		s1 := subject.PostSplit(2.4)
		Expect(s1.NumCategories()).To(Equal(2))
		Expect(s1.At(0).Mean()).To(BeNumerically("~", 2.2, 0.001))
		Expect(s1.At(1).Mean()).To(BeNumerically("~", 2.5, 0.001))

		s2 := subject.PostSplit(4.8)
		Expect(s2.NumCategories()).To(Equal(2))
		Expect(s2.At(0).Mean()).To(BeNumerically("~", 2.3, 0.001))
		Expect(s2.At(1).Mean()).To(BeNumerically("~", 2.6, 0.001))
	})
})

var _ = Describe("FeatureStats_Categorical", func() {
	var subject *internal.FeatureStats_Categorical

	BeforeEach(func() {
		subject = new(internal.FeatureStats_Categorical)
		subject.ObserveWeight(1, 2.2, 1.0)
		subject.ObserveWeight(4, 2.3, 1.0)
		subject.ObserveWeight(7, 2.4, 1.0)
	})

	It("should observe", func() {
		Expect(subject.NumCategories()).To(Equal(3))
	})

	It("should calculate post-splits", func() {
		s := subject.PostSplit()
		Expect(s.NumCategories()).To(Equal(3))
		Expect(s.At(1).Sum).To(Equal(2.2))
		Expect(s.At(4).Sum).To(Equal(2.3))
		Expect(s.At(7).Sum).To(Equal(2.4))
	})
})
