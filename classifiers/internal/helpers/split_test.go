package helpers

import (
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SplitCondition", func() {

	Describe("nominalMultiwaySplitCondition", func() {
		var subject SplitCondition
		model := testdata.ClassificationModel()

		BeforeEach(func() {
			subject = NewNominalMultiwaySplitCondition(model.Attribute("outlook"))
		})

		It("should calculate branch", func() {
			Expect(subject.Branch(core.MapInstance{"outlook": "overcast"})).To(Equal(1))
			Expect(subject.Branch(core.MapInstance{"outlook": nil})).To(Equal(-1))
		})

	})

	Describe("numericBinarySplitCondition", func() {
		var subject SplitCondition
		model := testdata.RegressionModel()

		BeforeEach(func() {
			subject = &numericBinarySplitCondition{
				predictor:  model.Attribute("hours"),
				splitValue: 25,
			}
		})

		It("should calculate branch", func() {
			Expect(subject.Branch(core.MapInstance{"hours": 24})).To(Equal(0))
			Expect(subject.Branch(core.MapInstance{"hours": 25})).To(Equal(0))
			Expect(subject.Branch(core.MapInstance{"hours": 26})).To(Equal(1))
			Expect(subject.Branch(core.MapInstance{"hours": nil})).To(Equal(-1))
		})

	})

})
