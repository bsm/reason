package stats

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kappa", func() {
	var subject *Kappa

	BeforeEach(func() {
		subject = NewKappa()
		subject.Record(0, 0, 22)
		subject.Record(0, 1, 7)
		subject.Record(1, 0, 9)
		subject.Record(1, 1, 13)
	})

	It("should record", func() {
		Expect(subject).To(Equal(&Kappa{
			m: [][]float64{
				{22, 9},
				{7, 13},
			},
			ncols: 2,
		}))

		k := NewKappa()
		k.Record(1, 3, 1.1)
		Expect(k).To(Equal(&Kappa{
			m: [][]float64{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 1.1, 0, 0},
			},
			ncols: 4,
		}))
	})

	It("should calculate value", func() {
		Expect(subject.Value()).To(BeNumerically("~", 0.353, 0.001))
	})

})
