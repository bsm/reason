package split_test

import (
	"github.com/bsm/reason/common/split"
	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InformationGain", func() {
	var subject = split.InformationGain{MinBranchFraction: 0.1}
	var _ split.Criterion = subject

	It("should evaluate split (classification)", func() {
		Expect(subject.ClassificationRange(nil)).To(Equal(1.0))
		Expect(subject.ClassificationRange(clspre)).To(Equal(1.0))
		Expect(subject.ClassificationRange(util.NewVectorFromSlice(1, 2, 3))).To(BeNumerically("~", 1.585, 0.001))

		Expect(subject.ClassificationMerit(nil, nil)).To(Equal(0.0))
		Expect(subject.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))
		Expect(subject.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.336, 0.001))
		Expect(subject.ClassificationMerit(clspre, clspost3)).To(Equal(0.0))

		x := split.InformationGain{MinBranchFraction: 0.3}
		Expect(x.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))

		x = split.InformationGain{MinBranchFraction: 0.35}
		Expect(x.ClassificationMerit(clspre, clspost1)).To(Equal(0.0))
	})

	It("should evaluate split (regression)", func() {
		Expect(subject.RegressionRange(nil)).To(Equal(0.0))
		Expect(subject.RegressionRange(regpre)).To(Equal(0.0))

		Expect(subject.RegressionMerit(nil, nil)).To(Equal(0.0))
		Expect(subject.RegressionMerit(regpre, regpost1)).To(Equal(0.0))
		Expect(subject.RegressionMerit(regpre, regpost2)).To(Equal(0.0))
	})
})
