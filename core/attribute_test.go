package core

import (
	"bytes"
	"encoding/gob"
	"math"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Attribute", func() {
	var nominal, numeric *Attribute

	BeforeEach(func() {
		nominal = &Attribute{Name: "cat", Kind: AttributeKindNominal}
		numeric = &Attribute{Name: "num", Kind: AttributeKindNumeric}
	})

	It("should extract values from instances", func() {
		Expect(nominal.Value(MapInstance{"cat": "x"})).To(Equal(AttributeValue(0.0)))
		Expect(nominal.Value(MapInstance{"cat": "y"})).To(Equal(AttributeValue(1.0)))
		Expect(nominal.Value(MapInstance{"cat": nil}).IsMissing()).To(BeTrue())
		Expect(nominal.Value(MapInstance{"cat": 1}).IsMissing()).To(BeTrue())

		Expect(numeric.Value(MapInstance{"num": 2.3})).To(Equal(AttributeValue(2.3)))
		Expect(numeric.Value(MapInstance{"num": -1})).To(Equal(AttributeValue(-1.0)))
		Expect(numeric.Value(MapInstance{"num": uint(8)})).To(Equal(AttributeValue(8.0)))
		Expect(numeric.Value(MapInstance{"num": nil}).IsMissing()).To(BeTrue())
	})

	It("should return values by value", func() {
		Expect(nominal.ValueOf("b")).To(Equal(AttributeValue(0)))
		Expect(nominal.ValueOf("a")).To(Equal(AttributeValue(1)))
		Expect(nominal.ValueOf("b")).To(Equal(AttributeValue(0)))

		Expect(nominal.ValueOf(nil).IsMissing()).To(BeTrue())
		Expect(nominal.ValueOf([]byte{'a'})).To(Equal(AttributeValue(1)))
		Expect(nominal.ValueOf(3).IsMissing()).To(BeTrue())

		Expect(numeric.ValueOf("x").IsMissing()).To(BeTrue())
		Expect(numeric.ValueOf(3.2)).To(Equal(AttributeValue(3.2)))
	})

	It("should gob marshal/unmarshal", func() {
		nominal.ValueOf("b")

		buf := new(bytes.Buffer)
		err := gob.NewEncoder(buf).Encode(nominal)
		Expect(err).NotTo(HaveOccurred())

		var out *Attribute
		err = gob.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(nominal))
	})
})

var _ = Describe("AttributeValue", func() {

	It("should check if missing", func() {
		Expect(AttributeValue(1).IsMissing()).To(BeFalse())
		Expect(MissingValue().IsMissing()).To(BeTrue())
		Expect(AttributeValue(math.NaN()).IsMissing()).To(BeTrue())
	})

	It("should return values", func() {
		Expect(AttributeValue(1).Value()).To(Equal(1.0))
		Expect(math.IsNaN(MissingValue().Value())).To(BeTrue())
	})

})

var _ = Describe("AttributeValues", func() {
	var subject *AttributeValues

	BeforeEach(func() {
		subject = NewAttributeValues("c", "a", "b")
	})

	It("should have a len", func() {
		Expect(subject.Len()).To(Equal(3))
	})

	It("should return values", func() {
		Expect(subject.Values()).To(Equal([]string{
			"c", "a", "b",
		}))
		Expect(subject.IndexOf("d")).To(Equal(3))
		Expect(subject.Values()).To(Equal([]string{
			"c", "a", "b", "d",
		}))
	})

	It("should fetch indices", func() {
		Expect(subject.IndexOf("a")).To(Equal(1))
		Expect(subject.IndexOf("b")).To(Equal(2))
		Expect(subject.IndexOf("c")).To(Equal(0))
		Expect(subject.IndexOf("d")).To(Equal(3))
		Expect(subject.IndexOf("e")).To(Equal(4))

		Expect(subject.Len()).To(Equal(5))
		Expect(subject.Values()).To(Equal([]string{
			"c", "a", "b", "d", "e",
		}))
	})

})
