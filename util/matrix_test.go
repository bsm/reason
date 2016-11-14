package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SparseMatrix", func() {
	var subject SparseMatrix

	BeforeEach(func() {
		subject = SparseMatrix{
			0: {1: 2, 2: 3},
			2: {2: 7, 9: 4},
			3: {3: 4},
		}
	})

	It("should get", func() {
		Expect(subject.Get(0)).To(Equal(SparseVector{1: 2, 2: 3}))
		Expect(subject.Get(1)).To(BeNil())
		Expect(subject.Get(2)).To(Equal(SparseVector{2: 7, 9: 4}))
		Expect(subject.Get(4)).To(BeNil())
		Expect(subject.Get(-1)).To(BeNil())
	})

	It("should count rows/cols", func() {
		Expect(subject.NumRows()).To(Equal(3))
		Expect(subject.NumCols()).To(Equal(2))
	})

	It("should return weights", func() {
		weights := subject.Weights()
		Expect(weights).To(Equal(map[int]float64{0: 5, 2: 11, 3: 4}))
	})

	It("should increment", func() {
		m := SparseMatrix{}
		m.Incr(2, 3, 4)
		Expect(m).To(Equal(SparseMatrix{
			2: {3: 4},
		}))
		m.Incr(1, 2, 3)
		m.Incr(1, 2, 4)
		Expect(m).To(Equal(SparseMatrix{
			1: {2: 7},
			2: {3: 4},
		}))
	})

})
