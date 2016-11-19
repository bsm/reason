package msgpack

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	_ "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decoder", func() {

	It("should decode", func() {
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

		dec := NewDecoder(buf)
		var s string
		Expect(dec.Decode(&s)).NotTo(HaveOccurred())
		Expect(s).To(Equal("string"))

		var b []byte
		Expect(dec.Decode(&b)).NotTo(HaveOccurred())
		Expect(string(b)).To(Equal("bytes"))

		var i int
		Expect(dec.Decode(&i)).NotTo(HaveOccurred())
		Expect(i).To(Equal(1))

		var bl bool
		Expect(dec.Decode(&bl)).NotTo(HaveOccurred())
		Expect(bl).To(BeTrue())

		var f float64
		Expect(dec.Decode(&f)).NotTo(HaveOccurred())
		Expect(f).To(BeNumerically("~", 8.2, 0.001))

		var si []int
		Expect(dec.Decode(&si)).NotTo(HaveOccurred())
		Expect(si).To(Equal([]int{7, 8}))

		var mp map[string]int
		Expect(dec.Decode(&mp)).NotTo(HaveOccurred())
		Expect(mp).To(Equal(map[string]int{"xy": 1}))
	})

	It("should decode nils", func() {
		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)
		Expect(enc.Encode(([]int)(nil))).NotTo(HaveOccurred())
		Expect(enc.Encode((mockSetType)(nil))).NotTo(HaveOccurred())
		Expect(enc.Encode((*mockSliceType)(nil))).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		dec := NewDecoder(buf)
		var sl []int
		Expect(dec.Decode(&sl)).NotTo(HaveOccurred())
		Expect(sl).To(BeNil())

		var mset mockSetType
		Expect(dec.Decode(&mset)).NotTo(HaveOccurred())
		Expect(mset).To(BeNil())

		var msl *mockSliceType
		Expect(dec.Decode(&msl)).NotTo(HaveOccurred())
		Expect(msl).To(BeNil())
	})

	FIt("should decode custom types", func() {
		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)
		Expect(enc.Encode(mockNone{})).NotTo(HaveOccurred())
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

		dec := NewDecoder(buf)

		var mn mockNone
		Expect(dec.Decode(&mn)).NotTo(HaveOccurred())
		Expect(mn).To(Equal(mockNone{}))

		var slt *mockSliceType
		Expect(dec.Decode(&slt)).NotTo(HaveOccurred())
		Expect(slt).To(Equal(&mockSliceType{data: []int{4, 5}}))

		var set *mockSliceType
		Expect(dec.Decode(&set)).NotTo(HaveOccurred())
		Expect(set).To(Equal(mockSetType{6: mockNone{}}))

		var mps map[int]*mockSliceType
		Expect(dec.Decode(&mps)).NotTo(HaveOccurred())
		Expect(mps).To(Equal(map[int]*mockSliceType{
			5: {data: []int{9}},
		}))

		var sli []mockInterface
		Expect(dec.Decode(&mps)).NotTo(HaveOccurred())
		Expect(sli).To(Equal([]mockInterface{
			mockSetType{8: mockNone{}},
			&mockSliceType{data: []int{9}},
		}))
	})

})
