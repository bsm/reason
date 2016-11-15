package calc

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KappaStat", func() {
	var subject *KappaStat

	BeforeEach(func() {
		subject = NewKappaStat()
		subject.Record(0, 0, 22)
		subject.Record(0, 1, 7)
		subject.Record(1, 0, 9)
		subject.Record(1, 1, 13)
	})

	It("should record", func() {
		Expect(subject).To(Equal(&KappaStat{
			m: [][]float64{
				{22, 9},
				{7, 13},
			},
			ncols: 2,
		}))

		k := NewKappaStat()
		k.Record(1, 3, 1.1)
		Expect(k).To(Equal(&KappaStat{
			m: [][]float64{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 1.1, 0, 0},
			},
			ncols: 4,
		}))
	})

	It("should calculate kappa", func() {
		Expect(subject.Kappa()).To(BeNumerically("~", 0.353, 0.001))
	})

})
