package internal_test

import (
	"github.com/bsm/reason/regression"
	"github.com/bsm/reason/regression/hoeffding/internal"
	"github.com/bsm/reason/testdata"
	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LeafNode", func() {
	model := testdata.RegressionModel()

	var wrapper *internal.Node
	var subject *internal.LeafNode

	BeforeEach(func() {
		subject = new(internal.LeafNode)
		wrapper = &internal.Node{
			Kind:  &internal.Node_Leaf{Leaf: subject},
			Stats: util.NewNumStream(),
		}

		target := model.Feature("hours")
		for _, x := range testdata.DataSet {
			subject.Observe(model, target, x, 1.0, wrapper)
		}
	})

	It("should observe", func() {
		Expect(wrapper.Weight()).To(Equal(14.0))
		Expect(wrapper.Stats.Mean()).To(BeNumerically("~", 39.8, 0.1))

		Expect(subject.FeatureStats).To(HaveLen(4))
		Expect(subject.FeatureStats).To(HaveKey("temp"))
		Expect(subject.FeatureStats["temp"].GetNumerical()).To(BeNil())
		Expect(subject.FeatureStats["temp"].GetCategorical().NumCategories()).To(Equal(3))
		Expect(subject.WeightAtLastEval).To(Equal(0.0))
	})

	It("should evaluate splits", func() {
		crit := regression.DefaultSplitCriterion()
		Expect(subject.EvaluateSplit("unknown", crit, wrapper)).To(BeNil())

		cat := subject.EvaluateSplit("outlook", crit, wrapper)
		Expect(cat.Feature).To(Equal("outlook"))
		Expect(cat.Merit).To(BeNumerically("~", 9.14, 0.01))
		Expect(cat.Range).To(Equal(1.0))
		Expect(cat.Pivot).To(Equal(0.0))
		Expect(cat.PreSplit.Weight).To(Equal(14.0))
		Expect(cat.PostSplit.NumCategories()).To(Equal(3))

		num := subject.EvaluateSplit("humidex", crit, wrapper)
		Expect(num.Feature).To(Equal("humidex"))
		Expect(num.Merit).To(BeNumerically("~", 93.57, 0.01))
		Expect(num.Range).To(Equal(1.0))
		Expect(num.Pivot).To(BeNumerically("~", 40.0, 0.001))
		Expect(num.PreSplit.Weight).To(Equal(14.0))
		Expect(num.PostSplit.NumCategories()).To(Equal(2))
	})

	It("should allow to disable/enable", func() {
		Expect(subject.FeatureStats).To(HaveLen(4))
		Expect(subject.IsDisabled).To(BeFalse())
		subject.Enable()
		Expect(subject.FeatureStats).To(HaveLen(4))
		Expect(subject.IsDisabled).To(BeFalse())

		subject.Disable()
		Expect(subject.FeatureStats).To(BeNil())
		Expect(subject.IsDisabled).To(BeTrue())

		subject.Enable()
		Expect(subject.FeatureStats).NotTo(BeNil())
		Expect(subject.FeatureStats).To(HaveLen(0))
		Expect(subject.IsDisabled).To(BeFalse())
	})

})
