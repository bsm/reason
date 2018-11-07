package observer_test

import (
	"github.com/bsm/reason/common/observer"
	"github.com/bsm/reason/common/split"
	"github.com/bsm/reason/testdata"
	util "github.com/bsm/reason/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClassificationCategorical", func() {
	var subject *observer.ClassificationCategorical
	var pre *util.Vector

	model := testdata.SimpleModel
	target := model.Feature("play")
	feat := model.Feature("outlook")
	crit := split.DefaultCriterion(target)

	BeforeEach(func() {
		subject = observer.NewClassificationCategorical()
		pre = util.NewVector()
		for _, x := range testdata.SimpleDataSet {
			subject.Observe(feat.Category(x), target.Category(x)*2)
			pre.Incr(int(target.Category(x)*2), 1.0)
		}
	})

	It("should observe", func() {
		Expect(subject.Dist.Data).To(Equal([]float64{
			2, 0, 3,
			4, 0, 0,
			3, 0, 2,
		}))
	})

	It("should evaluate splits", func() {
		merit, post := subject.EvaluateSplit(crit, pre)
		Expect(merit).To(BeNumerically("~", 0.246, 0.001))
		Expect(post.Data).To(Equal([]float64{
			2, 0, 3,
			4, 0, 0,
			3, 0, 2,
		}))
	})

	It("should calculate probabilty of a category given a target", func() {
		yes, no := target.CategoryOf("yes")*2, target.CategoryOf("no")*2
		sunny, overcast, rainy := feat.CategoryOf("sunny"), feat.CategoryOf("overcast"), feat.CategoryOf("rainy")

		Expect(subject.Prob(sunny, -1)).To(Equal(0.0))
		Expect(subject.Prob(sunny, 3)).To(Equal(0.0))

		Expect(subject.Prob(sunny, yes)).To(BeNumerically("~", 0.333, 0.001))
		Expect(subject.Prob(sunny, no)).To(BeNumerically("~", 0.429, 0.001))
		Expect(subject.Prob(overcast, yes)).To(BeNumerically("~", 0.417, 0.001))
		Expect(subject.Prob(overcast, no)).To(BeNumerically("~", 0.143, 0.001))
		Expect(subject.Prob(rainy, yes)).To(BeNumerically("~", 0.25, 0.001))
		Expect(subject.Prob(rainy, no)).To(BeNumerically("~", 0.571, 0.001))
	})
})
