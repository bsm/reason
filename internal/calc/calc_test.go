package calc

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func() {

	It("should calc norm-probability", func() {
		Expect(NormProb(1.23)).To(BeNumerically("~", 0.891, 0.001))
		Expect(NormProb(0.12)).To(BeNumerically("~", 0.548, 0.001))
		Expect(NormProb(-0.76)).To(BeNumerically("~", 0.224, 0.001))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/calc")
}
