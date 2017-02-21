package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Prediction", func() {

	It("should recycle", func() {
		p1 := NewPrediction(10)
		p1.Release()

		p2 := NewPrediction(5)
		Expect(p2).To(HaveCap(10))
	})

})
