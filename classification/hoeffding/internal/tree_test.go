package internal_test

import (
	"github.com/bsm/reason/classification/hoeffding/internal"
	"github.com/bsm/reason/common/hoeffding"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {
	var subject *internal.Tree

	model := testdata.ClassificationModel()
	pre := &util.Vector{Sparse: map[int64]float64{0: 9.0, 1: 5.0}}
	post := &util.VectorDistribution{
		Sparse: map[int64]*util.Vector{
			0: &util.Vector{Sparse: map[int64]float64{0: 2, 1: 3}},
			1: &util.Vector{Sparse: map[int64]float64{0: 4}},
			2: &util.Vector{Sparse: map[int64]float64{0: 3, 1: 2}},
		},
	}

	BeforeEach(func() {
		subject = internal.NewTree(model, "play")
	})

	It("should init", func() {
		Expect(subject.Len()).To(Equal(1))
	})

	It("should add (leaf) nodes", func() {
		ref := subject.Add(nil)
		Expect(ref).To(Equal(int64(2)))
		Expect(subject.Len()).To(Equal(2))

		leaf := subject.Get(ref).GetLeaf()
		Expect(leaf).NotTo(BeNil())
	})

	It("should split nodes", func() {
		subject.Split(1, "outlook", pre, post, 0)
		Expect(subject.Len()).To(Equal(4))

		split := subject.Get(1).GetSplit()
		Expect(split).NotTo(BeNil())
		Expect(split.Children.Len()).To(Equal(3))
	})

	It("should traverse", func() {
		subject.Split(1, "outlook", pre, post, 0)
		root := subject.Get(1)

		node, nodeRef, parent, parentIndex := subject.Traverse(core.MapExample{"outlook": "overcast"}, 1, nil, -1, nil)
		Expect(node.Stats.Sparse).To(Equal(map[int64]float64{0: 4}))
		Expect(parent).To(Equal(root))
		Expect(parentIndex).To(Equal(1))

		var traversed []*internal.Node
		subject.Traverse(core.MapExample{"outlook": "overcast"}, 1, nil, -1, func(n *internal.Node) {
			traversed = append(traversed, n)
		})
		Expect(traversed).To(HaveLen(2))
		Expect(traversed[0].Weight()).To(Equal(14.0))
		Expect(traversed[1].Weight()).To(Equal(4.0))

		node, nodeRef, parent, parentIndex = subject.Traverse(core.MapExample{"outlook": "overcast"}, 99, nil, -1, nil)
		Expect(node).To(BeNil())
		Expect(nodeRef).To(Equal(int64(99)))
		Expect(parent).To(BeNil())
		Expect(parentIndex).To(Equal(-1))
	})

	It("should filter leaves", func() {
		subject.Split(1, "outlook", pre, post, 0)
		Expect(subject.FilterLeaves(nil)).To(HaveLen(3))
	})

	It("should accumulate info", func() {
		subject.Split(1, "outlook", pre, post, 0)

		info := new(hoeffding.TreeInfo)
		subject.Accumulate(1, 1, info)
		Expect(info).To(Equal(&hoeffding.TreeInfo{
			NumNodes:    4,
			NumLearning: 3,
			NumDisabled: 0,
			MaxDepth:    2,
		}))
	})

})
