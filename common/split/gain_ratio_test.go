package split_test

import (
	"github.com/bsm/reason/common/split"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GainRatio", func() {
	var subject split.GainRatio
	var _ split.Criterion = subject

	It("should reduce merit of 'super-attributes' (classification)", func() {
		parent := split.InformationGain{MinBranchFraction: 0.1}
		subject.Criterion = parent

		Expect(parent.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.280, 0.001))
		Expect(subject.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.179, 0.001))

		Expect(parent.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.337, 0.001))
		Expect(subject.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.117, 0.001))

		Expect(parent.ClassificationMerit(clspre, clspost3)).To(BeNumerically("~", 0.0, 0.001))
		Expect(subject.ClassificationMerit(clspre, clspost3)).To(BeNumerically("~", 0.0, 0.001))
	})

	It("should reduce merit of 'super-attributes' (regression)", func() {
		parent := split.VarianceReduction{MinWeight: 1.0}
		subject.Criterion = parent

		Expect(parent.RegressionMerit(regpre, regpost1)).To(BeNumerically("~", 7.808, 0.001))
		Expect(subject.RegressionMerit(regpre, regpost1)).To(BeNumerically("~", 8.181, 0.001))

		Expect(parent.RegressionMerit(regpre, regpost2)).To(BeNumerically("~", 4.576, 0.001))
		Expect(subject.RegressionMerit(regpre, regpost2)).To(BeNumerically("~", 1.664, 0.001))
	})
})
