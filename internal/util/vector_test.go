package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NumVector", func() {

	It("should normalize", func() {
		nv := NumVector{2, 5, 9}.Normalize()
		Expect(nv).To(Equal(NumVector{0.125, 0.3125, 0.5625}))
	})

	It("should get", func() {
		vv := NumVector{1.1, 2.2}
		Expect(vv.Get(0)).To(Equal(1.1))
		Expect(vv.Get(1)).To(Equal(2.2))
		Expect(vv.Get(2)).To(Equal(0.0))
		Expect(vv.Get(-1)).To(Equal(0.0))
	})

	It("should set", func() {
		vv := make(NumVector, 0)
		vv = vv.Set(2, 3.3)
		vv = vv.Set(0, 1.1)
		Expect(vv).To(Equal(NumVector{1.1, 0.0, 3.3}))
		Expect(vv.Set(3, 4.4)).To(Equal(NumVector{1.1, 0.0, 3.3, 4.4}))
		Expect(vv).To(Equal(NumVector{1.1, 0.0, 3.3}))
	})

	It("should incr", func() {
		vv := make(NumVector, 0)
		vv = vv.Incr(2, 3.3)
		vv = vv.Incr(0, 1.1)
		vv = vv.Incr(0, 0.2)
		Expect(vv).To(Equal(NumVector{1.3, 0.0, 3.3}))
	})

})
