package stats

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NormDist", func() {

	It("should calc PDF", func() {
		Expect(StdNormal.PDF(2.34)).To(BeNumerically("~", 0.026, 0.001))
		Expect(StdNormal.PDF(1.23)).To(BeNumerically("~", 0.187, 0.001))
		Expect(StdNormal.PDF(0.12)).To(BeNumerically("~", 0.396, 0.001))
		Expect(StdNormal.PDF(0.0)).To(BeNumerically("~", 0.399, 0.001))
		Expect(StdNormal.PDF(-0.76)).To(BeNumerically("~", 0.299, 0.001))
		Expect(StdNormal.PDF(-1.23)).To(BeNumerically("~", 0.187, 0.001))
		Expect(StdNormal.PDF(-2.34)).To(BeNumerically("~", 0.026, 0.001))
	})

	It("should calc CDF", func() {
		Expect(StdNormal.CDF(2.34)).To(BeNumerically("~", 0.990, 0.001))
		Expect(StdNormal.CDF(1.23)).To(BeNumerically("~", 0.891, 0.001))
		Expect(StdNormal.CDF(0.12)).To(BeNumerically("~", 0.548, 0.001))
		Expect(StdNormal.CDF(0.0)).To(BeNumerically("~", 0.500, 0.001))
		Expect(StdNormal.CDF(-0.76)).To(BeNumerically("~", 0.224, 0.001))
		Expect(StdNormal.CDF(-1.23)).To(BeNumerically("~", 0.109, 0.001))
		Expect(StdNormal.CDF(-2.34)).To(BeNumerically("~", 0.010, 0.001))
	})

})
