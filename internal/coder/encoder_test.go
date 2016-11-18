package coder

import (
	"bytes"

	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Encoder",
	func(v interface{}, x []byte) {
		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)

		Expect(enc.Encode(v)).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())
		Expect(buf.Bytes()).To(Equal(x), "expected %#v, was %#v", x, buf.Bytes())
	},

	Entry("nil", nil, []byte{mnil}),
	Entry("short strings", "s", []byte{0xa1, 's'}),
	Entry("longer strings", "this is a longer message", append([]byte{0xb8}, []byte("this is a longer message")...)),

	Entry("int8", int8(43), []byte{0x2b}),
	Entry("int16", int16(4133), []byte{mint16, 0x10, 0x25}),
	Entry("int32", int32(4111333), []byte{mint32, 0x0, 0x3e, 0xbb, 0xe5}),
	Entry("int64", int64(41111133333), []byte{mint64, 0x0, 0x0, 0x0, 0x9, 0x92, 0x6a, 0x1c, 0x95}),
	Entry("negative", int64(-4133), []byte{mint16, 0xef, 0xdb}),

	Entry("uint8", uint8(43), []byte{0x2b}),
	Entry("uint16", uint16(4133), []byte{muint16, 0x10, 0x25}),
	Entry("uint32", uint32(4111333), []byte{muint32, 0x0, 0x3e, 0xbb, 0xe5}),
	Entry("uint64", uint64(41111133333), []byte{muint64, 0x0, 0x0, 0x0, 0x9, 0x92, 0x6a, 0x1c, 0x95}),

	Entry("float64", 12.34, []byte{mfloat64, 0x40, 0x28, 0xae, 0x14, 0x7a, 0xe1, 0x47, 0xae}),

	Entry("binary", []byte("message"), append([]byte{mbin8, 7}, []byte("message")...)),

	Entry("slice (numeric)", []int{1, 2, 3}, []byte{mfixarray + 3, 1, 2, 3}),
	Entry("slice (strings)", []string{"a", "b"}, []byte{mfixarray + 2, 0xa1, 'a', 0xa1, 'b'}),

	Entry("custom", &mockType{data: []int{3, 1}},
		[]byte{mext8, 15, 8, '*', 'c', 'o', 'd', 'e', 'r', '.', 'm', 'o', 'c', 'k', 'T', 'y', 'p', 'e', mfixarray + 2, 3, 1},
	),
)
