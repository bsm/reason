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
		subject.ObserveWeight(1.2, 2.2, 1.0)
		subject.ObserveWeight(4.2, 2.2, 1.0)
		subject.ObserveWeight(8.4, 2.2, 1.0)
	})

	It("should observe", func() {
		Expect(subject).To(Equal(&internal.FeatureStats_Numerical{
			Min: 1.2,
			Max: 8.4,
			Observations: []internal.FeatureStats_Numerical_Observation{
				{FeatureValue: 1.2, TargetValue: 2.2, Weight: 1},
				{FeatureValue: 4.2, TargetValue: 2.2, Weight: 1},
				{FeatureValue: 8.4, TargetValue: 2.2, Weight: 1},
			},
		}))
	})

	It("should calculate pivot points", func() {
		pp := subject.PivotPoints()
		Expect(pp).To(HaveLen(11))
		Expect(pp[0]).To(BeNumerically("~", 1.8, 0.01))
		Expect(pp[1]).To(BeNumerically("~", 2.4, 0.01))
		Expect(pp[9]).To(BeNumerically("~", 7.2, 0.01))
		Expect(pp[10]).To(BeNumerically("~", 7.8, 0.01))
		Expect(new(internal.FeatureStats_Numerical).PivotPoints()).To(BeEmpty())
	})

	It("should calculate post-splits", func() {
		s1 := subject.PostSplit(2.4)
		Expect(s1.NumCategories()).To(Equal(2))
		Expect(s1.At(0).Sum).To(Equal(2.2))
		Expect(s1.At(1).Sum).To(Equal(4.4))

		s2 := subject.PostSplit(4.8)
		Expect(s2.NumCategories()).To(Equal(2))
		Expect(s2.At(0).Sum).To(Equal(4.4))
		Expect(s2.At(1).Sum).To(Equal(2.2))
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
