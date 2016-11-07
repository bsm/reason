package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NumRange", func() {

	It("should split", func() {
		Expect(NewNumRange().SplitPoints(2)).To(BeEmpty())
		Expect((&NumRange{Min: 1, Max: 5}).SplitPoints(3)).To(Equal([]float64{2, 3, 4}))
		Expect((&NumRange{Min: 1, Max: 5}).SplitPoints(30)).To(HaveLen(30))
	})

	It("should update", func() {
		r := NewNumRange()
		r.Update(2)
		Expect(r).To(Equal(&NumRange{Min: 2, Max: 2}))
		r.Update(6)
		Expect(r).To(Equal(&NumRange{Min: 2, Max: 6}))
	})

})
