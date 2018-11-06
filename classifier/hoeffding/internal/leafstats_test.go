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
		for _, x := range testdata.SimpleDataSet {
			subject.Update(play, outlook, x, 1.0)
		}
		stats := subject.GetCC()
		Expect(stats.Dist.NumRows()).To(Equal(3))
	})

	It("should observe (classification, numeric)", func() {
		for _, x := range testdata.SimpleDataSet {
			subject.Update(play, humidex, x, 1.0)
		}
		stats := subject.GetCN()
		Expect(stats.Dist.NumRows()).To(Equal(2))
	})

	It("should observe (regression, categorical)", func() {
		for _, x := range testdata.SimpleDataSet {
			subject.Update(hours, outlook, x, 1.0)
		}
		stats := subject.GetRC()
		Expect(stats.Dist.NumRows()).To(Equal(3))
	})

	It("should observe (regression, numeric)", func() {
		for _, x := range testdata.SimpleDataSet {
			subject.Update(hours, humidex, x, 1.0)
		}
		stats := subject.GetRN()
		Expect(stats.MaxBuckets).To(Equal(uint32(12)))
		Expect(stats.Dist).To(HaveLen(11))
	})
})
