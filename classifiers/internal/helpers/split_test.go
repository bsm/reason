package helpers

import (
	"bytes"
	"context"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/msgpack"
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

		It("should encode/decode", func() {
			buf := new(bytes.Buffer)
			enc := msgpack.NewEncoder(buf)
			err := enc.Encode(subject)
			Expect(err).NotTo(HaveOccurred())
			Expect(enc.Close()).NotTo(HaveOccurred())

			var out SplitCondition
			err = msgpack.NewDecoder(buf).Decode(&out)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(Equal(subject))
		})

	})

	Describe("numericBinarySplitCondition", func() {
		var subject SplitCondition
		model := core.NewModel(
			&core.Attribute{Name: "target"},
			&core.Attribute{Name: "hours"},
		)

		BeforeEach(func() {
			subject = NewNumericBinarySplitCondition(model.Predictor("hours"), 25)
		})

		It("should calculate branch", func() {
			Expect(subject.Branch(core.MapInstance{"hours": 24})).To(Equal(0))
			Expect(subject.Branch(core.MapInstance{"hours": 25})).To(Equal(0))
			Expect(subject.Branch(core.MapInstance{"hours": 26})).To(Equal(1))
			Expect(subject.Branch(core.MapInstance{"hours": nil})).To(Equal(-1))
		})

		It("should encode/decode", func() {
			buf := new(bytes.Buffer)
			enc := msgpack.NewEncoder(buf)
			err := enc.Encode(subject)
			Expect(err).NotTo(HaveOccurred())
			Expect(enc.Close()).NotTo(HaveOccurred())

			var out SplitCondition
			dec := msgpack.NewDecoder(buf)
			dec.Ctx = context.WithValue(dec.Ctx, core.ModelContextKey, model)
			err = dec.Decode(&out)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(Equal(subject))
		})
	})

})
