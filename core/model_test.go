package core

import (
	"bytes"

	"github.com/bsm/reason/internal/msgpack"

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

	It("should encode/decode", func() {
		buf := new(bytes.Buffer)
		enc := msgpack.NewEncoder(buf)
		err := enc.Encode(subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		var out *Model
		err = msgpack.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})
})
