package internal_test

import (
	"github.com/bsm/reason/classifier/hoeffding/internal"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FeatureStats", func() {
	var subject *internal.FeatureStats

	play := core.NewCategoricalFeature("play", []string{"yes", "no"})
	outlook := core.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"})
	hours := core.NewNumericalFeature("hours")
	humidex := core.NewNumericalFeature("humidex")

	BeforeEach(func() {
		subject = new(internal.FeatureStats)
	})

	It("should observe (classification, categorical)", func() {
		for _, x := range testdata.DataSet {
			subject.ObserveExample(play, outlook, x, 1.0)
		}
		Expect(subject.GetCC().Stats.NumRows()).To(Equal(3))
	})

	It("should observe (classification, numeric)", func() {
		for _, x := range testdata.DataSet {
			subject.ObserveExample(play, humidex, x, 1.0)
		}
		Expect(subject.GetCN().Stats.NumRows()).To(Equal(2))
	})

	It("should observe (regression, categorical)", func() {
		for _, x := range testdata.DataSet {
			subject.ObserveExample(hours, outlook, x, 1.0)
		}
		Expect(subject.GetRC().Stats.NumRows()).To(Equal(3))
	})

	PIt("should observe (regression, numeric)", func() {
		for _, x := range testdata.DataSet {
			subject.ObserveExample(hours, humidex, x, 1.0)
		}
		Expect(subject.GetRN()).To(Equal(3))
	})
})

var _ = Describe("FeatureStats_ClassificationNumerical", func() {
	var subject *internal.FeatureStats_ClassificationNumerical

	BeforeEach(func() {
		subject = new(internal.FeatureStats_ClassificationNumerical)
		subject.Stats.ObserveWeight(0, 1.4, 1.0)
		subject.Stats.ObserveWeight(0, 1.3, 1.0)
		subject.Stats.ObserveWeight(0, 1.5, 1.0)
		subject.Stats.ObserveWeight(1, 4.1, 1.0)
		subject.Stats.ObserveWeight(1, 3.7, 1.0)
		subject.Stats.ObserveWeight(1, 4.9, 1.0)
		subject.Stats.ObserveWeight(1, 4.0, 1.0)
		subject.Stats.ObserveWeight(1, 3.3, 1.0)
		subject.Stats.ObserveWeight(2, 6.3, 1.0)
		subject.Stats.ObserveWeight(2, 5.8, 1.0)
		subject.Stats.ObserveWeight(2, 5.1, 1.0)
		subject.Stats.ObserveWeight(2, 5.3, 1.0)
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
		pts := subject.PivotPoints()
		Expect(pts).To(HaveLen(11))
		Expect(pts[0]).To(BeNumerically("~", 1.72, 0.01))
		Expect(pts[1]).To(BeNumerically("~", 2.13, 0.01))
		Expect(pts[9]).To(BeNumerically("~", 5.47, 0.01))
		Expect(pts[10]).To(BeNumerically("~", 5.88, 0.01))

		Expect(new(internal.FeatureStats_ClassificationNumerical).PivotPoints()).To(BeEmpty())
	})
})
