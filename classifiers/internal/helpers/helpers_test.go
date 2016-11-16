package helpers

import (
	"bytes"
	"encoding/gob"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MinMaxRange", func() {

	It("should split", func() {
		Expect(NewMinMaxRange().SplitPoints(2)).To(BeEmpty())
		Expect((&MinMaxRange{Min: 1, Max: 5}).SplitPoints(3)).To(Equal([]float64{2, 3, 4}))
		Expect((&MinMaxRange{Min: 1, Max: 5}).SplitPoints(30)).To(HaveLen(30))
	})

	It("should update", func() {
		r := NewMinMaxRange()
		r.Update(2)
		Expect(r).To(Equal(&MinMaxRange{Min: 2, Max: 2}))
		r.Update(6)
		Expect(r).To(Equal(&MinMaxRange{Min: 2, Max: 6}))
	})

	It("should gob marshal/unmarshal", func() {
		subject := &MinMaxRange{Min: 1, Max: 5}
		buf := new(bytes.Buffer)
		err := gob.NewEncoder(buf).Encode(subject)
		Expect(err).NotTo(HaveOccurred())

		var out *MinMaxRange
		err = gob.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})

})

var _ = Describe("MinMaxRanges", func() {
	var subject *MinMaxRanges

	BeforeEach(func() {
		subject = NewMinMaxRanges()
		subject.Update(0, 1)
		subject.Update(0, 2)
		subject.Update(0, 3)
		subject.Update(1, 6)
		subject.Update(1, 7)
		subject.Update(1, 8)
		subject.Update(2, 14)
		subject.Update(2, 16)
		subject.Update(2, 18)
	})

	It("should split", func() {
		Expect(NewMinMaxRanges().SplitPoints(2)).To(BeEmpty())
		Expect(subject.SplitPoints(3)).To(Equal([]float64{5.25, 9.5, 13.75}))
		Expect(subject.SplitPoints(30)).To(HaveLen(30))
	})

	It("should gob marshal/unmarshal", func() {
		buf := new(bytes.Buffer)
		err := gob.NewEncoder(buf).Encode(subject)
		Expect(err).NotTo(HaveOccurred())

		var out *MinMaxRanges
		err = gob.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/helpers")
}
