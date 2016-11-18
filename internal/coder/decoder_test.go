package coder

import (
	"bytes"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decoder", func() {
	var buf *bytes.Buffer
	var subject *Decoder

	BeforeEach(func() {
		buf = new(bytes.Buffer)
		subject = NewDecoder(buf)
	})

	It("should decode strings", func() {
		buf.Write([]byte{0xa1, 's'})
		buf.Write(append([]byte{0xb8}, []byte("this is a longer message")...))

		s, err := subject.String()
		Expect(err).NotTo(HaveOccurred())
		Expect(s).To(Equal("s"))

		err = subject.Decode(&s)
		Expect(err).NotTo(HaveOccurred())
		Expect(s).To(Equal("this is a longer message"))
	})

	It("should decode bytes", func() {
		buf.Write(append([]byte{mbin8, 7}, []byte("message")...))
		buf.Write([]byte{mnil})
		buf.Write(append([]byte{mbin8, 7}, []byte("message")...))

		b, err := subject.Bytes()
		Expect(err).NotTo(HaveOccurred())
		Expect(b).To(Equal([]byte("message")))

		b, err = subject.Bytes()
		Expect(err).NotTo(HaveOccurred())
		Expect(b).To(BeNil())

		err = subject.Decode(&b)
		Expect(err).NotTo(HaveOccurred())
		Expect(b).To(Equal([]byte("message")))
	})

	It("should decode ints", func() {
		buf.Write([]byte{0x2b})
		buf.Write([]byte{mint16, 0x10, 0x25})
		buf.Write([]byte{mint32, 0x0, 0x3e, 0xbb, 0xe5})
		buf.Write([]byte{mint64, 0x0, 0x0, 0x0, 0x9, 0x92, 0x6a, 0x1c, 0x95})
		buf.Write([]byte{mint16, 0xef, 0xdb})

		n, err := subject.Int64()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(43)))

		n, err = subject.Int64()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(4133)))

		n, err = subject.Int64()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(4111333)))

		n, err = subject.Int64()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(41111133333)))

		n, err = subject.Int64()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(-4133)))
	})

	It("should decode uints", func() {
		buf.Write([]byte{0x2b})
		buf.Write([]byte{mint16, 0x10, 0x25})
		buf.Write([]byte{mint32, 0x0, 0x3e, 0xbb, 0xe5})
		buf.Write([]byte{mint64, 0x0, 0x0, 0x0, 0x9, 0x92, 0x6a, 0x1c, 0x95})

		n, err := subject.Uint64()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(uint64(43)))

		n, err = subject.Uint64()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(uint64(4133)))

		n, err = subject.Uint64()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(uint64(4111333)))

		n, err = subject.Uint64()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(uint64(41111133333)))
	})

	It("should decode floats", func() {
		buf.Write([]byte{mfloat64, 0x40, 0x28, 0xae, 0x14, 0x7a, 0xe1, 0x47, 0xae})

		var n float64
		err := subject.Decode(&n)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(BeNumerically("~", 12.34, 0.001))

	})

	It("should decode slices", func() {
		in := []int32{1, 2, 3, 8, 7, 5}
		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)
		err := enc.Encode(in)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		var out []int32
		err = NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(in))
	})

	It("should decode maps", func() {
		in := map[int]float64{
			1: 1.1,
			3: 3.3,
			7: 7.7,
		}
		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)
		err := enc.Encode(in)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		var out map[int]float64
		err = NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(in))
	})

	It("should peek custom types", func() {
		in := &mockType{data: []int{3, 1, 9, 4, 2}}
		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)
		err := enc.Encode(in)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		dec := NewDecoder(buf)
		typ, err := dec.PeekType()
		Expect(err).NotTo(HaveOccurred())
		Expect(typ).To(Equal(reflect.TypeOf(&mockType{})))

		typ, err = dec.PeekType()
		Expect(err).NotTo(HaveOccurred())
		Expect(typ).To(Equal(reflect.TypeOf(&mockType{})))
	})

	It("should decode custom", func() {
		in := &mockType{data: []int{3, 1, 9, 4, 2}}
		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)
		err := enc.Encode(in)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		out := new(mockType)
		err = NewDecoder(buf).Decode(out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(in))
	})

})
