package split_test

import (
	"github.com/bsm/reason/common/split"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GiniImpurity", func() {
	var subject = split.GiniImpurity{}
	var _ split.Criterion = subject

	It("should evaluate split (classification)", func() {
		Expect(subject.ClassificationRange(nil)).To(Equal(1.0))
		Expect(subject.ClassificationRange(clspre)).To(Equal(1.0))

		Expect(subject.ClassificationMerit(nil, nil)).To(Equal(0.0))
		Expect(subject.ClassificationMerit(clspre, clspost1)).To(BeNumerically("~", 0.338, 0.001))
		Expect(subject.ClassificationMerit(clspre, clspost2)).To(BeNumerically("~", 0.311, 0.001))
		Expect(subject.ClassificationMerit(clspre, clspost3)).To(BeNumerically("~", 0.480, 0.001))
	})

	It("should evaluate split (regression)", func() {
		Expect(subject.RegressionRange(nil)).To(Equal(0.0))
		Expect(subject.RegressionRange(regpre)).To(Equal(0.0))

		Expect(subject.RegressionMerit(nil, nil)).To(Equal(0.0))
		Expect(subject.RegressionMerit(regpre, regpost1)).To(Equal(0.0))
		Expect(subject.RegressionMerit(regpre, regpost2)).To(Equal(0.0))
	})
})
