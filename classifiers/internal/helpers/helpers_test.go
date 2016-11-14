package helpers

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("minMaxRange", func() {

	It("should split", func() {
		Expect(newMinMaxRange().SplitPoints(2)).To(BeEmpty())
		Expect((&minMaxRange{Min: 1, Max: 5}).SplitPoints(3)).To(Equal([]float64{2, 3, 4}))
		Expect((&minMaxRange{Min: 1, Max: 5}).SplitPoints(30)).To(HaveLen(30))
	})

	It("should update", func() {
		r := newMinMaxRange()
		r.Update(2)
		Expect(r).To(Equal(&minMaxRange{Min: 2, Max: 2}))
		r.Update(6)
		Expect(r).To(Equal(&minMaxRange{Min: 2, Max: 6}))
	})

})

var _ = Describe("minMaxRanges", func() {

	var subject *minMaxRanges

	BeforeEach(func() {
		subject = newMinMaxRanges()
		subject.Update(0, 1)
		subject.Update(0, 2)
		subject.Update(0, 3)
		subject.Update(1, 6)
		subject.Update(1, 7)
		subject.Update(1, 8)
		subject.Update(2, 14)
		subject.Update(2, 16)
		subject.Update(2, 18)
	})

	It("should split", func() {
		Expect(newMinMaxRanges().SplitPoints(2)).To(BeEmpty())
		Expect(subject.SplitPoints(3)).To(Equal([]float64{5.25, 9.5, 13.75}))
		Expect(subject.SplitPoints(30)).To(HaveLen(30))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/helpers")
}
