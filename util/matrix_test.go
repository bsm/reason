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
		subject.Set(0, 0, 1.0)
		subject.Set(0, 2, 2.0)
		subject.Set(0, 4, 3.0)
		subject.Set(0, 6, 4.0)

		subject.Set(1, 0, 2.0)
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
		Expect(subject.Stride).To(Equal(uint32(7)))
		Expect(subject.Data).To(Equal([]float64{
			1, 0, 2, 0, 3, 0, 4,
			2, 5, 0, 6, 7, 8, 0,
		}))

		subject.Set(1, 1, 8.0)
		subject.Set(-1, 0, 4.0)
		Expect(subject.At(1, 1)).To(Equal(8.0))
	})

	It("should incr", func() {
		Expect(subject.At(1, 1)).To(Equal(5.0))
		subject.Incr(1, 1, 4.0)
		subject.Incr(-1, 0, 4.0)
		Expect(subject.At(1, 1)).To(Equal(9.0))
	})

	It("should returns rows", func() {
		Expect(subject.Row(-1)).To(BeNil())
		Expect(subject.Row(0)).To(Equal([]float64{1, 0, 2, 0, 3, 0, 4}))
		Expect(subject.Row(1)).To(Equal([]float64{2, 5, 0, 6, 7, 8, 0}))
		Expect(subject.Row(2)).To(BeNil())
	})

	It("should identify zero-rows", func() {
		Expect(subject.IsRowZero(-1)).To(BeTrue())
		Expect(subject.IsRowZero(0)).To(BeFalse())
		Expect(subject.IsRowZero(1)).To(BeFalse())
		Expect(subject.IsRowZero(2)).To(BeTrue())

		subject.Set(2, 2, 4.0)
		Expect(subject.IsRowZero(2)).To(BeFalse())
		subject.Set(2, 2, 0.0)
		Expect(subject.IsRowZero(2)).To(BeTrue())
	})

	It("should calculate row sums", func() {
		Expect(subject.RowSum(-1)).To(Equal(0.0))
		Expect(subject.RowSum(0)).To(Equal(10.0))
		Expect(subject.RowSum(1)).To(Equal(28.0))
		Expect(subject.RowSum(2)).To(Equal(0.0))
	})

	It("should count non-zero in rows", func() {
		Expect(subject.RowNNZ(-1)).To(Equal(0))
		Expect(subject.RowNNZ(0)).To(Equal(4))
		Expect(subject.RowNNZ(1)).To(Equal(5))
		Expect(subject.RowNNZ(2)).To(Equal(0))
	})

	It("should calculate col sums", func() {
		Expect(subject.ColSum(-1)).To(Equal(0.0))
		Expect(subject.ColSum(0)).To(Equal(3.0))
		Expect(subject.ColSum(1)).To(Equal(5.0))
		Expect(subject.ColSum(2)).To(Equal(2.0))
		Expect(subject.ColSum(3)).To(Equal(6.0))
		Expect(subject.ColSum(4)).To(Equal(10.0))
		Expect(subject.ColSum(5)).To(Equal(8.0))
		Expect(subject.ColSum(6)).To(Equal(4.0))
		Expect(subject.ColSum(7)).To(Equal(0.0))
	})

	It("should count non-zero in cols", func() {
		Expect(subject.ColNNZ(-1)).To(Equal(0))
		Expect(subject.ColNNZ(0)).To(Equal(2))
		Expect(subject.ColNNZ(1)).To(Equal(1))
		Expect(subject.ColNNZ(2)).To(Equal(1))
		Expect(subject.ColNNZ(3)).To(Equal(1))
		Expect(subject.ColNNZ(4)).To(Equal(2))
		Expect(subject.ColNNZ(5)).To(Equal(1))
		Expect(subject.ColNNZ(6)).To(Equal(1))
		Expect(subject.ColNNZ(7)).To(Equal(0))
	})

	It("should calculate weight sum", func() {
		Expect(subject.WeightSum()).To(Equal(38.0))
	})
})
