package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KappaStat", func() {
	var subject KappaStat

	BeforeEach(func() {
		subject = subject.Record(0, 0, 22)
		subject = subject.Record(0, 1, 7)
		subject = subject.Record(1, 0, 9)
		subject = subject.Record(1, 1, 13)
	})

	It("should record", func() {
		Expect(subject).To(Equal(KappaStat{
			{22, 9},
			{7, 13},
		}))
	})

	It("should calculate kappa", func() {
		Expect(subject.Kappa()).To(BeNumerically("~", 0.353, 0.001))
	})

})
