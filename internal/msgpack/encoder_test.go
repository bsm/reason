package msgpack

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	_ "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Encode", func() {

	It("should encode", func() {
		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)
		Expect(enc.Encode("string")).NotTo(HaveOccurred())
		Expect(enc.Encode([]byte("bytes"))).NotTo(HaveOccurred())
		Expect(enc.Encode(1)).NotTo(HaveOccurred())
		Expect(enc.Encode(true)).NotTo(HaveOccurred())
		Expect(enc.Encode(8.2)).NotTo(HaveOccurred())
		Expect(enc.Encode([]int{7, 8})).NotTo(HaveOccurred())
		Expect(enc.Encode(map[string]int{"xy": 1})).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		exp := []byte{
			mfixstr + 6, 's', 't', 'r', 'i', 'n', 'g',
			mbin8, 5, 'b', 'y', 't', 'e', 's',
			1,
			mtrue,
			mfloat64, 0x40, 0x20, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66,
			mfixarray + 2, 7, 8,
			mfixmap + 1, mfixstr + 2, 'x', 'y', 1,
		}
		Expect(buf.Bytes()).To(Equal(exp), "expected: %#v\ngot:      %#v", exp, buf.Bytes())
	})

	It("should encode nils", func() {
		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)
		Expect(enc.Encode(([]int)(nil))).NotTo(HaveOccurred())
		Expect(enc.Encode((mockSetType)(nil))).NotTo(HaveOccurred())
		Expect(enc.Encode((*mockSliceType)(nil))).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		exp := []byte{
			mnil,
			mfixext2, 8, 0x0, 112, mnil,
			mfixext2, 8, 0x0, 111, mnil,
		}
		Expect(buf.Bytes()).To(Equal(exp), "expected: %#v\ngot:      %#v", exp, buf.Bytes())
	})

	It("should encode custom types", func() {
		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)
		Expect(enc.Encode(&mockSliceType{data: []int{4, 5}})).NotTo(HaveOccurred())
		Expect(enc.Encode(mockSetType{6: mockNone{}})).NotTo(HaveOccurred())
		Expect(enc.Encode(map[int]*mockSliceType{
			5: {data: []int{9}},
		})).NotTo(HaveOccurred())
		Expect(enc.Encode([]mockInterface{
			mockSetType{8: mockNone{}},
			&mockSliceType{data: []int{9}},
		})).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		exp := []byte{
			mfixext2, 8, 0x0, 111, mfixarray + 2, 4, 5,
			mfixext2, 8, 0x0, 112, mfixmap + 1, 6, mfixext2, 8, 0x0, 113,
			mfixmap + 1, 5, mfixext2, 8, 0x0, 111, mfixarray + 1, 9,

			mfixarray + 2,
			mfixext2, 8, 0x0, 112, mfixmap + 1, 8, mfixext2, 8, 0x0, 113,
			mfixext2, 8, 0x0, 0x6f, mfixarray + 1, 9,
		}
		Expect(buf.Bytes()).To(Equal(exp), "expected: %#v\ngot:      %#v", exp, buf.Bytes())
	})

})
