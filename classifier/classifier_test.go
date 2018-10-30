package classifier_test

import (
	"testing"

	"github.com/bsm/reason/classifier"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Problem", func() {

	It("should parse names", func() {
		Expect(classifier.ParseProblem("c")).To(Equal(classifier.Classification))
		Expect(classifier.ParseProblem("cls")).To(Equal(classifier.Classification))
		Expect(classifier.ParseProblem("r")).To(Equal(classifier.Regression))
		Expect(classifier.ParseProblem("reg")).To(Equal(classifier.Regression))
		_, err := classifier.ParseProblem("unk")
		Expect(err).To(MatchError(`reason: unable to parse problem "unk"`))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "classifier")
}
