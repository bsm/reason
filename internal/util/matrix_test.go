package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NumMatrix", func() {
	var subject NumMatrix

	BeforeEach(func() {
		subject = NumMatrix{
			{1.0, 5.0, 3.0},
			{4.0, 6.0, 2.0},
		}
	})

	It("should count rows/cols", func() {
		Expect(subject.NumRows()).To(Equal(2))
		Expect(subject.NumCols()).To(Equal(3))
		Expect(NumMatrix{{1.0, 2.0}, {}}.NumCols()).To(Equal(2))
	})

	It("should sum rows/cols", func() {
		Expect(subject.SumCols()).To(Equal([]float64{5.0, 11.0, 5.0}))
		Expect(subject.SumRows()).To(Equal([]float64{9.0, 12.0}))
	})

	It("should sum rows and total", func() {
		rsum, tsum := subject.SumRowsPlusTotal()
		Expect(rsum).To(HaveLen(2))
		Expect(rsum[0]).To(BeNumerically("~", 9.0))
		Expect(rsum[1]).To(BeNumerically("~", 12.0))
		Expect(tsum).To(Equal(21.0))
	})

})
