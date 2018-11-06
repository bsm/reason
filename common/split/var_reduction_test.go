package split_test

import (
	"github.com/bsm/reason/common/split"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("VarianceReduction", func() {
	var subject = split.VarianceReduction{MinWeight: 1.0}
	var _ split.Criterion = subject

	It("should evaluate split (classification)", func() {
		Expect(subject.ClassificationRange(nil)).To(Equal(0.0))
		Expect(subject.ClassificationRange(clspre)).To(Equal(0.0))

		Expect(subject.ClassificationMerit(nil, nil)).To(Equal(0.0))
		Expect(subject.ClassificationMerit(clspre, clspost1)).To(Equal(0.0))
		Expect(subject.ClassificationMerit(clspre, clspost2)).To(Equal(0.0))
		Expect(subject.ClassificationMerit(clspre, clspost3)).To(Equal(0.0))
	})

	It("should evaluate split (regression)", func() {
		Expect(subject.RegressionRange(nil)).To(Equal(1.0))
		Expect(subject.RegressionRange(regpre)).To(Equal(1.0))

		Expect(subject.RegressionMerit(nil, nil)).To(Equal(0.0))
		Expect(subject.RegressionMerit(regpre, regpost1)).To(BeNumerically("~", 7.808, 0.001))

		c := split.VarianceReduction{MinWeight: 4.0}
		Expect(c.RegressionMerit(regpre, regpost1)).To(Equal(0.0))
	})
})
