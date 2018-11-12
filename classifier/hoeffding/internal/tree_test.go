package internal_test

import (
	"bytes"

	"github.com/bsm/reason"
	"github.com/bsm/reason/classifier/hoeffding/internal"
	"github.com/bsm/reason/testdata"
	"github.com/bsm/reason/util"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {
	var subject *internal.Tree

	model := testdata.SimpleModel
	postSplit := internal.PostSplit{Classification: &util.Matrix{
		Stride: 2,
		Data: []float64{
			2, 3, // rainy: 2 yes, 3 no
			4, 0, // overcast: 4 yes, 0 no
			3, 2, // sunny: 3 yes, 2 no
		},
	}}

	BeforeEach(func() {
		subject = internal.NewTree(model, "play")
	})

	It("should marshal to writer", func() {
		buf := new(bytes.Buffer)
		Expect(subject.WriteTo(buf)).To(Equal(int64(238)))

		t := new(internal.Tree)
		Expect(proto.Unmarshal(buf.Bytes(), t)).To(Succeed())
		Expect(t).To(Equal(subject))
	})

	It("should unmarshal from reader", func() {
		data, err := proto.Marshal(subject)
		Expect(err).NotTo(HaveOccurred())

		t := new(internal.Tree)
		Expect(t.ReadFrom(bytes.NewReader(data))).To(Equal(int64(len(data))))
		Expect(t).To(Equal(subject))
	})

	It("should init", func() {
		Expect(subject.NumNodes()).To(Equal(1))
	})

	It("should add leaf nodes", func() {
		nref := subject.AddLeaf(nil, 7)
		Expect(nref).To(Equal(int64(2)))
		Expect(subject.NumNodes()).To(Equal(2))
		Expect(subject.GetLeaf(nref)).To(Equal(&internal.LeafNode{
			WeightAtLastEval: 7,
		}))
	})

	It("should split nodes", func() {
		Expect(subject.Nodes).To(HaveLen(1))
		subject.SplitNode(1, "outlook", postSplit, 0)
		Expect(subject.Nodes).To(HaveLen(4))

		split := subject.GetSplit(1)
		Expect(split).NotTo(BeNil())
		Expect(split.Children).To(Equal([]int64{2, 3, 4}))

		// rainy
		Expect(subject.GetNode(2).GetClassification()).To(Equal(&internal.Node_ClassificationStats{
			Vector: util.Vector{Data: []float64{2, 3}},
		}))

		// overcast
		Expect(subject.GetNode(3).GetClassification()).To(Equal(&internal.Node_ClassificationStats{
			Vector: util.Vector{Data: []float64{4, 0}},
		}))

		// sunny
		Expect(subject.GetNode(4).GetClassification()).To(Equal(&internal.Node_ClassificationStats{
			Vector: util.Vector{Data: []float64{3, 2}},
		}))
	})

	It("should traverse", func() {
		subject.SplitNode(1, "outlook", postSplit, 0)
		root := subject.GetNode(1)
		example := reason.MapExample{"outlook": "overcast"}

		// valid nref
		node, nref, parent, ppos := subject.Traverse(example, 1, nil, -1)
		Expect(node.GetClassification().WeightSum()).To(Equal(4.0))
		Expect(nref).To(Equal(int64(3)))
		Expect(parent).To(Equal(root))
		Expect(ppos).To(Equal(1)) // 0: rainy, 1: overcast, 2: sunny

		// invalid nref
		node, nref, parent, ppos = subject.Traverse(example, 99, nil, -1)
		Expect(node).To(BeNil())
		Expect(nref).To(Equal(int64(99)))
		Expect(parent).To(BeNil())
		Expect(ppos).To(Equal(-1))
	})
})
