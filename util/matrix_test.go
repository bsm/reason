package util_test

import (
	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Matrix", func() {
	var subject *util.Matrix

	BeforeEach(func() {
		subject = util.NewMatrix()
	})

	ItBehavesLikeAMatrix := func() {
		BeforeEach(func() {
			subject.Set(0, 0, 1.0)
			subject.Set(0, 2, 2.0)
			subject.Set(0, 4, 3.0)
			subject.Set(0, 6, 4.0)

			subject.Set(1, 1, 5.0)
			subject.Set(1, 3, 6.0)
			subject.Set(1, 4, 7.0)
			subject.Set(1, 5, 8.0)
		})

		It("should have dims", func() {
			rows, cols := subject.Dims()
			Expect(rows).To(Equal(2))
			Expect(cols).To(Equal(7))
		})

		It("should at", func() {
			Expect(subject.At(0, 0)).To(Equal(1.0))
			Expect(subject.At(0, 1)).To(Equal(0.0))
			Expect(subject.At(0, 4)).To(Equal(3.0))
			Expect(subject.At(1, 1)).To(Equal(5.0))
			Expect(subject.At(1, 2)).To(Equal(0.0))
			Expect(subject.At(1, 5)).To(Equal(8.0))
			Expect(subject.At(2, 0)).To(Equal(0.0))
			Expect(subject.At(-1, 0)).To(Equal(0.0))
			Expect(subject.At(0, -1)).To(Equal(0.0))
		})

		It("should set", func() {
			Expect(subject.At(1, 1)).To(Equal(5.0))
			subject.Set(1, 1, 8.0)
			subject.Set(-1, 0, 4.0)
			Expect(subject.At(1, 1)).To(Equal(8.0))
		})

		It("should add", func() {
			Expect(subject.At(1, 1)).To(Equal(5.0))
			subject.Add(1, 1, 4.0)
			subject.Add(-1, 0, 4.0)
			Expect(subject.At(1, 1)).To(Equal(9.0))
		})

		It("should returns rows", func() {
			Expect(subject.Row(-1)).To(BeNil())
			Expect(subject.Row(0)).To(Equal([]float64{1, 0, 2, 0, 3, 0, 4}))
			Expect(subject.Row(1)).To(Equal([]float64{0, 5, 0, 6, 7, 8, 0}))
			Expect(subject.Row(2)).To(BeNil())
		})

		It("should calculate row sums", func() {
			Expect(subject.RowSum(-1)).To(Equal(0.0))
			Expect(subject.RowSum(0)).To(Equal(10.0))
			Expect(subject.RowSum(1)).To(Equal(26.0))
			Expect(subject.RowSum(2)).To(Equal(0.0))
		})

		It("should calculate sums", func() {
			Expect(subject.Sum()).To(Equal(0.0))
		})
	}

	Describe("Dense", func() {
		ItBehavesLikeAMatrix()
	})
})
