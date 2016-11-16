package core

import (
	"bytes"
	"encoding/gob"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Model", func() {
	var subject *Model

	BeforeEach(func() {
		subject = NewModel(
			&Attribute{
				Name:   "season",
				Kind:   AttributeKindNominal,
				Values: NewAttributeValues("winter", "spring", "summer", "autumn"),
			},
			&Attribute{
				Name: "temperature",
				Kind: AttributeKindNumeric,
			},
			&Attribute{
				Name: "humidity",
				Kind: AttributeKindNumeric,
			},
		)
	})

	It("should return target", func() {
		Expect(subject.Target().Name).To(Equal("season"))
	})

	It("should return (immutable) predictor", func() {
		Expect(subject.Predictor("humidity").Name).To(Equal("humidity"))
	})

	It("should gob marshal/unmarshal", func() {
		buf := new(bytes.Buffer)
		err := gob.NewEncoder(buf).Encode(subject)
		Expect(err).NotTo(HaveOccurred())

		var out *Model
		err = gob.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})
})
