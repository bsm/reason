package internal_test

import (
	"github.com/bsm/reason/classifier/hoeffding/internal"
	"github.com/bsm/reason/common/split"
	"github.com/bsm/reason/testdata"
	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SplitNode", func() {
	var subject *internal.SplitNode

	BeforeEach(func() {
		subject = new(internal.SplitNode)
	})

	It("should set/get child references", func() {
		Expect(subject.GetChild(0)).To(Equal(int64(0)))
		Expect(subject.GetChild(1)).To(Equal(int64(0)))
		subject.SetChild(1, 33)
		Expect(subject.GetChild(1)).To(Equal(int64(33)))
	})
})

var _ = Describe("LeafNode", func() {
	var subject *internal.LeafNode
	var node *internal.Node

	BeforeEach(func() {
		subject = new(internal.LeafNode)
		node = &internal.Node{Kind: &internal.Node_Leaf{Leaf: subject}}
	})

	Describe("classification", func() {
		crit := split.InformationGain{MinBranchFraction: 0.1}

		BeforeEach(func() {
			model := testdata.SimpleModel
			target := model.Feature("play")
			for _, x := range testdata.SimpleDataSet {
				subject.ObserveExample(model, target, x, 1.0, node)
			}
		})

		It("should observe examples", func() {
			Expect(subject.FeatureStats).To(HaveLen(6))
			Expect(node.GetClassification()).To(Equal(&internal.Node_ClassificationStats{
				Vector: util.Vector{Data: []float64{9, 5}},
			}))
		})

		It("should allow to disable/enable", func() {
			Expect(subject.FeatureStats).To(HaveLen(6))
			Expect(subject.IsDisabled).To(BeFalse())
			subject.Enable()
			Expect(subject.FeatureStats).To(HaveLen(6))
			Expect(subject.IsDisabled).To(BeFalse())
			Expect(subject.EvaluateSplit("outlook", crit, node)).NotTo(BeNil())

			subject.Disable()
			Expect(subject.FeatureStats).To(BeNil())
			Expect(subject.IsDisabled).To(BeTrue())
			Expect(subject.EvaluateSplit("outlook", crit, node)).To(BeNil())

			subject.Enable()
			Expect(subject.FeatureStats).NotTo(BeNil())
			Expect(subject.FeatureStats).To(HaveLen(0))
			Expect(subject.IsDisabled).To(BeFalse())
		})

		It("should evaluate splits", func() {
			Expect(subject.EvaluateSplit("unknown", crit, node)).To(BeNil())

			cc := subject.EvaluateSplit("outlook", crit, node)
			Expect(cc.Feature).To(Equal("outlook"))
			Expect(cc.Merit).To(BeNumerically("~", 0.246, 0.001))
			Expect(cc.Range).To(BeNumerically("~", 1.0, 0.001))
			Expect(cc.PostSplit.Classification).NotTo(BeNil())

			cn := subject.EvaluateSplit("humidex", crit, node)
			Expect(cn.Feature).To(Equal("humidex"))
			Expect(cn.Merit).To(BeNumerically("~", 0.176, 0.001))
			Expect(cn.Range).To(BeNumerically("~", 1.0, 0.001))
			Expect(cn.PostSplit.Classification).NotTo(BeNil())
		})
	})

	Describe("regression", func() {
		crit := split.VarianceReduction{MinWeight: 1.0}

		BeforeEach(func() {
			model := testdata.SimpleModel
			target := model.Feature("hours")
			for _, x := range testdata.SimpleDataSet {
				subject.ObserveExample(model, target, x, 1.0, node)
			}
		})

		It("should observe examples", func() {
			Expect(subject.FeatureStats).To(HaveLen(6))
			Expect(node.GetRegression()).To(Equal(&internal.Node_RegressionStats{
				NumStream: util.NumStream{
					Min:        23,
					Max:        52,
					Weight:     14,
					Sum:        557,
					SumSquares: 23377,
				},
			}))
		})

		It("should evaluate splits", func() {
			Expect(subject.EvaluateSplit("unknown", crit, node)).To(BeNil())

			rc := subject.EvaluateSplit("outlook", crit, node)
			Expect(rc.Feature).To(Equal("outlook"))
			Expect(rc.Merit).To(BeNumerically("~", 9.137, 0.001))
			Expect(rc.Range).To(BeNumerically("~", 1.0, 0.001))
			Expect(rc.PostSplit.Regression).NotTo(BeNil())

			rn := subject.EvaluateSplit("humidex", crit, node)
			Expect(rn.Feature).To(Equal("humidex"))
			Expect(rn.Merit).To(BeNumerically("~", 22.923, 0.001))
			Expect(rn.Range).To(BeNumerically("~", 1.0, 0.001))
			Expect(rn.PostSplit.Regression).NotTo(BeNil())
		})
	})
})

var _ = Describe("Node", func() {
	var subject *internal.Node
	var (
		clsStats *internal.Node_Classification
		regStats *internal.Node_Regression
	)
	BeforeEach(func() {
		subject = new(internal.Node)
		clsStats = &internal.Node_Classification{Classification: &internal.Node_ClassificationStats{
			Vector: *util.NewVectorFromSlice(1.1, 2.2, 3.3, 4.4, 5.5),
		}}
		regStats = &internal.Node_Regression{Regression: &internal.Node_RegressionStats{
			NumStream: util.NumStream{Weight: 5, Sum: 25, SumSquares: 125},
		}}
	})

	It("should calculate weight", func() {
		Expect(subject.Weight()).To(Equal(0.0))

		subject.Stats = clsStats
		Expect(subject.Weight()).To(Equal(16.5))

		subject.Stats = regStats
		Expect(subject.Weight()).To(Equal(5.0))
	})

	It("should evaluate sufficiency", func() {
		Expect(subject.IsSufficient()).To(BeFalse())

		subject.Stats = clsStats
		Expect(subject.IsSufficient()).To(BeTrue())
		clsStats.Classification.Data = clsStats.Classification.Data[:1]
		Expect(subject.IsSufficient()).To(BeFalse())

		subject.Stats = regStats
		Expect(subject.IsSufficient()).To(BeTrue())
		regStats.Regression.Weight = 1.0
		Expect(subject.IsSufficient()).To(BeFalse())
	})
})
