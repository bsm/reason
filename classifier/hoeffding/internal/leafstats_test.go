package internal_test

import (
	"github.com/bsm/reason/classifier/hoeffding/internal"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LeafNode_Stats", func() {
	var subject *internal.LeafNode_Stats

	play := core.NewCategoricalFeature("play", []string{"yes", "no"})
	outlook := core.NewCategoricalFeature("outlook", []string{"rainy", "overcast", "sunny"})
	hours := core.NewNumericalFeature("hours")
	humidex := core.NewNumericalFeature("humidex")

	BeforeEach(func() {
		subject = new(internal.LeafNode_Stats)
	})

	It("should observe (classification, categorical)", func() {
		for _, x := range testdata.DataSet {
			subject.Update(play, outlook, x, 1.0)
		}
		stats := subject.GetCC().Stats
		Expect(stats.NumRows()).To(Equal(3))
	})

	It("should observe (classification, numeric)", func() {
		for _, x := range testdata.DataSet {
			subject.Update(play, humidex, x, 1.0)
		}
		stats := subject.GetCN().Stats
		Expect(stats.NumRows()).To(Equal(2))
	})

	It("should observe (regression, categorical)", func() {
		for _, x := range testdata.DataSet {
			subject.Update(hours, outlook, x, 1.0)
		}
		stats := subject.GetRC().Stats
		Expect(stats.NumRows()).To(Equal(3))
	})

	It("should observe (regression, numeric)", func() {
		for _, x := range testdata.DataSet {
			subject.Update(hours, humidex, x, 1.0)
		}
		stats := subject.GetRN().Stats
		Expect(stats.MaxBuckets).To(Equal(uint32(12)))
		Expect(stats.Buckets).To(HaveLen(11))
	})
})

var _ = Describe("LeafNode_Stats_ClassificationNumerical", func() {
	var subject *internal.LeafNode_Stats_ClassificationNumerical

	BeforeEach(func() {
		subject = new(internal.LeafNode_Stats_ClassificationNumerical)
		subject.Stats.Observe(0, 1.4)
		subject.Stats.Observe(0, 1.3)
		subject.Stats.Observe(0, 1.5)
		subject.Stats.Observe(1, 4.1)
		subject.Stats.Observe(1, 3.7)
		subject.Stats.Observe(1, 4.9)
		subject.Stats.Observe(1, 4.0)
		subject.Stats.Observe(1, 3.3)
		subject.Stats.Observe(2, 6.3)
		subject.Stats.Observe(2, 5.8)
		subject.Stats.Observe(2, 5.1)
		subject.Stats.Observe(2, 5.3)
	})

	It("should calculate post-splits", func() {
		s1 := subject.PostSplit(2.4).Classification
		Expect(s1.NumRows()).To(Equal(2))
		Expect(s1.Row(0)).To(Equal([]float64{3, 0, 0}))
		Expect(s1.Row(1)).To(Equal([]float64{0, 5, 4}))

		s2 := subject.PostSplit(4.8).Classification
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

		Expect(new(internal.LeafNode_Stats_ClassificationNumerical).PivotPoints()).To(BeEmpty())
	})
})

var _ = Describe("LeafNode_Stats_RegressionNumerical", func() {
	var subject *internal.LeafNode_Stats_RegressionNumerical

	BeforeEach(func() {
		subject = new(internal.LeafNode_Stats_RegressionNumerical)
		subject.Stats.Observe(0.2, 1.4)
		subject.Stats.Observe(0.4, 1.3)
		subject.Stats.Observe(0.3, 1.5)
		subject.Stats.Observe(1.1, 4.1)
		subject.Stats.Observe(1.4, 3.7)
		subject.Stats.Observe(1.0, 4.9)
		subject.Stats.Observe(0.8, 4.0)
		subject.Stats.Observe(0.6, 3.3)
		subject.Stats.Observe(1.6, 6.3)
		subject.Stats.Observe(2.2, 5.8)
		subject.Stats.Observe(2.0, 5.1)
		subject.Stats.Observe(1.7, 5.3)
	})

	It("should calculate post-splits", func() {
		s1 := subject.PostSplit(0.2).Regression
		Expect(s1.NumRows()).To(Equal(2))
		Expect(s1.At(0).Mean()).To(BeNumerically("~", 1.4, 0.1))
		Expect(s1.At(1).Mean()).To(BeNumerically("~", 4.1, 0.1))

		s2 := subject.PostSplit(1.0).Regression
		Expect(s2.NumRows()).To(Equal(2))
		Expect(s2.At(0).Mean()).To(BeNumerically("~", 2.7, 0.1))
		Expect(s2.At(1).Mean()).To(BeNumerically("~", 5.1, 0.1))

		s3 := subject.PostSplit(2.0).Regression
		Expect(s3.NumRows()).To(Equal(2))
		Expect(s3.At(0).Mean()).To(BeNumerically("~", 3.7, 0.1))
		Expect(s3.At(1).Mean()).To(BeNumerically("~", 5.8, 0.1))
	})
})
