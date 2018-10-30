package internal_test

import (
	"github.com/bsm/reason/classifier/hoeffding/internal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FeatureStats_ClassificationCategorical", func() {
	var subject *internal.FeatureStats_ClassificationCategorical

	BeforeEach(func() {
		subject = new(internal.FeatureStats_ClassificationCategorical)

		// outlook: rainy=0 overcast=1 sunny=2
		// play:    yes=0 no=1
		subject.ObserveWeight(0, 1, 1.0) // rainy -> no
		subject.ObserveWeight(0, 1, 1.0) // rainy -> no
		subject.ObserveWeight(1, 0, 1.0) // overcast -> yes
		subject.ObserveWeight(2, 0, 1.0) // sunny -> yes
		subject.ObserveWeight(2, 0, 1.0) // sunny -> yes
		subject.ObserveWeight(2, 1, 1.0) // sunny -> no
		subject.ObserveWeight(1, 0, 1.0) // overcast -> yes
		subject.ObserveWeight(0, 1, 1.0) // rainy -> no
		subject.ObserveWeight(0, 0, 1.0) // rainy -> yes
		subject.ObserveWeight(2, 0, 1.0) // sunny -> yes
		subject.ObserveWeight(0, 0, 1.0) // rainy -> yes
		subject.ObserveWeight(1, 0, 1.0) // overcast -> yes
		subject.ObserveWeight(1, 0, 1.0) // overcast -> yes
		subject.ObserveWeight(2, 1, 1.0) // sunny -> no
	})

	It("should observe", func() {
		Expect(subject.NumCategories()).To(Equal(3))
	})

	It("should calculate post-splits", func() {
		s := subject.PostSplit()
		Expect(s.NumRows()).To(Equal(3))
		Expect(s.Row(0)).To(Equal([]float64{2, 3}))
		Expect(s.Row(1)).To(Equal([]float64{4, 0}))
		Expect(s.Row(2)).To(Equal([]float64{3, 2}))
	})

})

var _ = Describe("FeatureStats_ClassificationNumerical", func() {
	var subject *internal.FeatureStats_ClassificationNumerical

	BeforeEach(func() {
		subject = new(internal.FeatureStats_ClassificationNumerical)
		subject.ObserveWeight(1.4, 0, 1.0)
		subject.ObserveWeight(1.3, 0, 1.0)
		subject.ObserveWeight(1.5, 0, 1.0)
		subject.ObserveWeight(4.1, 1, 1.0)
		subject.ObserveWeight(3.7, 1, 1.0)
		subject.ObserveWeight(4.9, 1, 1.0)
		subject.ObserveWeight(4.0, 1, 1.0)
		subject.ObserveWeight(3.3, 1, 1.0)
		subject.ObserveWeight(6.3, 2, 1.0)
		subject.ObserveWeight(5.8, 2, 1.0)
		subject.ObserveWeight(5.1, 2, 1.0)
		subject.ObserveWeight(5.3, 2, 1.0)
	})

	It("should observe", func() {
		Expect(subject.Min.Data).To(Equal([]float64{1.3, 3.3, 5.1}))
		Expect(subject.Max.Data).To(Equal([]float64{1.5, 4.9, 6.3}))

		rows, cols := subject.Stats.Dims()
		Expect(rows).To(Equal(3))
		Expect(cols).To(Equal(3))
	})

	It("should calculate post-splits", func() {
		s1 := subject.PostSplit(2.4)
		Expect(s1.NumRows()).To(Equal(2))
		Expect(s1.Row(0)).To(Equal([]float64{3, 0, 0}))
		Expect(s1.Row(1)).To(Equal([]float64{0, 5, 4}))

		s2 := subject.PostSplit(4.8)
		Expect(s2.NumRows()).To(Equal(2))
		Expect(s2.Row(0)).To(Equal([]float64{3, 4.55925906389872, 0}))
		Expect(s2.Row(1)).To(Equal([]float64{0, 0.44074093610127996, 4}))
	})

	It("should calculate pivot points", func() {
		pp := subject.PivotPoints()
		Expect(pp).To(HaveLen(11))
		Expect(pp[0]).To(BeNumerically("~", 1.72, 0.01))
		Expect(pp[1]).To(BeNumerically("~", 2.13, 0.01))
		Expect(pp[9]).To(BeNumerically("~", 5.47, 0.01))
		Expect(pp[10]).To(BeNumerically("~", 5.88, 0.01))

		Expect(new(internal.FeatureStats_ClassificationNumerical).PivotPoints()).To(BeEmpty())
	})
})