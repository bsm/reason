package internal_test

import (
	"github.com/bsm/reason/classifier"
	"github.com/bsm/reason/classifier/hoeffding/internal"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {
	var subject *internal.Tree

	model := testdata.ClassificationModel()

	BeforeEach(func() {
		subject = internal.NewTree(model, "play")
	})

	It("should init", func() {
		Expect(subject.NumNodes()).To(Equal(1))
	})

	It("should detect classifier problems", func() {
		Expect(subject.DetectProblem()).To(Equal(classifier.Classification))
	})

	PIt("should add leaf nodes", func() {
		ref := subject.AddLeaf(nil)
		Expect(ref).To(Equal(int64(2)))
		Expect(subject.NumNodes()).To(Equal(2))
		Expect(subject.GetLeaf(ref)).To(Equal(&internal.LeafNode{
			WeightAtLastEval: 7,
		}))
	})

	// It("should split nodes", func() {
	// 	subject.Split(1, "outlook", pre, post, 0)
	// 	Expect(subject.Len()).To(Equal(4))

	// 	split := subject.Get(1).GetSplit()
	// 	Expect(split).NotTo(BeNil())
	// 	Expect(split.Children).To(HaveLen(3))
	// })

	// It("should traverse", func() {
	// 	subject.Split(1, "outlook", pre, post, 0)
	// 	root := subject.Get(1)

	// 	node, nodeRef, parent, parentIndex := subject.Traverse(core.MapExample{"outlook": "overcast"}, 1, nil, -1, nil)
	// 	Expect(node.Stats).To(Equal(util.NewVectorFromSlice(4, 0)))
	// 	Expect(parent).To(Equal(root))
	// 	Expect(parentIndex).To(Equal(1))

	// 	var traversed []*internal.Node
	// 	subject.Traverse(core.MapExample{"outlook": "overcast"}, 1, nil, -1, func(n *internal.Node) {
	// 		traversed = append(traversed, n)
	// 	})
	// 	Expect(traversed).To(HaveLen(2))
	// 	Expect(traversed[0].Weight()).To(Equal(14.0))
	// 	Expect(traversed[1].Weight()).To(Equal(4.0))

	// 	node, nodeRef, parent, parentIndex = subject.Traverse(core.MapExample{"outlook": "overcast"}, 99, nil, -1, nil)
	// 	Expect(node).To(BeNil())
	// 	Expect(nodeRef).To(Equal(int64(99)))
	// 	Expect(parent).To(BeNil())
	// 	Expect(parentIndex).To(Equal(-1))
	// })

	// It("should filter leaves", func() {
	// 	subject.Split(1, "outlook", pre, post, 0)
	// 	Expect(subject.FilterLeaves(nil)).To(HaveLen(3))
	// })

	// It("should accumulate info", func() {
	// 	subject.Split(1, "outlook", pre, post, 0)

	// 	info := new(hoeffding.TreeInfo)
	// 	subject.Accumulate(1, 1, info)
	// 	Expect(info).To(Equal(&hoeffding.TreeInfo{
	// 		NumNodes:    4,
	// 		NumLearning: 3,
	// 		NumDisabled: 0,
	// 		MaxDepth:    2,
	// 	}))
	// })

})
