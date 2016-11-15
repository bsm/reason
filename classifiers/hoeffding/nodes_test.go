package hoeffding

import (
	"github.com/bsm/reason/classifiers/internal/helpers"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("leafNode", func() {
	var subject *leafNode

	Describe("classification", func() {
		model := testdata.ClassificationModel()
		tree := New(model, nil)
		instances := testdata.ClassificationData()

		BeforeEach(func() {
			subject = newLeafNode(helpers.NewObservationStats(model.IsRegression()))
			for _, inst := range instances {
				subject.Learn(inst, tree)
			}
		})

		It("should init", func() {
			Expect(subject.weightOnLastEval).To(Equal(0.0))

			leaf := newLeafNode(subject.stats)
			Expect(leaf.weightOnLastEval).To(Equal(14.0))
		})

		It("should learn", func() {
			Expect(subject.stats.State()).To(ConsistOf(core.Prediction{
				{Value: 0, Votes: 9},
				{Value: 1, Votes: 5},
			}))
			Expect(subject.weightOnLastEval).To(Equal(0.0))
			Expect(subject.observers).To(HaveLen(4))
		})

		It("should estimate heap-size", func() {
			Expect(subject.ByteSize()).To(BeNumerically("~", 940, 20))
		})

		It("should calc promise split", func() {
			Expect(subject.Promise()).To(Equal(5.0))
		})

		It("should calc best split", func() {
			splits := subject.BestSplits(tree)
			Expect(splits).To(HaveLen(5))
			Expect(splits[0].Merit()).To(BeNumerically("~", 0.247, 0.001))
			Expect(splits[0].Condition().Predictor().Name).To(Equal("outlook"))
		})
	})

	Describe("regression", func() {
		model := testdata.RegressionModel()
		tree := New(model, nil)
		instances := testdata.RegressionData()

		BeforeEach(func() {
			subject = newLeafNode(helpers.NewObservationStats(model.IsRegression()))
			for _, inst := range instances {
				subject.Learn(inst, tree)
			}
		})

		It("should init", func() {
			Expect(subject.weightOnLastEval).To(Equal(0.0))

			leaf := newLeafNode(subject.stats)
			Expect(leaf.weightOnLastEval).To(Equal(14.0))
		})

		It("should learn", func() {
			state := subject.stats.State()
			Expect(state).To(HaveLen(1))
			Expect(state[0].Votes).To(Equal(14.0))
			Expect(state[0].Value.Value()).To(BeNumerically("~", 39.8, 0.1))
			Expect(subject.weightOnLastEval).To(Equal(0.0))
			Expect(subject.observers).To(HaveLen(4))
		})

		It("should estimate heap-size", func() {
			Expect(subject.ByteSize()).To(BeNumerically("~", 3000, 20))
		})

		It("should calc promise split", func() {
			Expect(subject.Promise()).To(Equal(14.0))
		})

		It("should calc best split", func() {
			splits := subject.BestSplits(tree)
			Expect(splits).To(HaveLen(5))
			Expect(splits[0].Merit()).To(BeNumerically("~", 19.57, 0.01))
			Expect(splits[0].Condition().Predictor().Name).To(Equal("outlook"))
		})
	})

})

var _ = Describe("splitNode", func() {
	var subject *splitNode
	model := testdata.ClassificationModel()

	BeforeEach(func() {
		condition := helpers.NewNominalMultiwaySplitCondition(model.Attribute("outlook"))
		stats := helpers.NewObservationStats(model.IsRegression())
		subject = newSplitNode(condition, stats, map[int]helpers.ObservationStats{
			1: helpers.NewObservationStats(model.IsRegression()),
			3: helpers.NewObservationStats(model.IsRegression()),
		})
	})

	It("should initialize", func() {
		Expect(subject.children).To(HaveLen(2))
	})

	It("should find leaves", func() {
		Expect(subject.FindLeaves(nil)).To(HaveLen(2))
	})

})
