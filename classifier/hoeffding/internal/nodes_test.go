package internal_test

import (
	"github.com/bsm/reason/classifier/hoeffding/internal"
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

	It("should observe examples (classification)", func() {
		model := testdata.ClassificationModel()
		target := model.Feature("play")
		for _, x := range testdata.DataSet {
			subject.ObserveExample(model, target, x, 1.0, node)
		}
		Expect(subject.FeatureStats).To(HaveLen(4))
		Expect(node.GetClassification()).To(Equal(&internal.Node_ClassificationStats{
			Vector: util.Vector{Data: []float64{9, 5}},
		}))
	})

	It("should observe examples (regression)", func() {
		model := testdata.RegressionModel()
		target := model.Feature("hours")
		for _, x := range testdata.DataSet {
			subject.ObserveExample(model, target, x, 1.0, node)
		}
		Expect(subject.FeatureStats).To(HaveLen(4))
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
})
