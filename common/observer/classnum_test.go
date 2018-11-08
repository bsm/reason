package observer_test

import (
	"github.com/bsm/reason/common/observer"
	"github.com/bsm/reason/common/split"
	"github.com/bsm/reason/testdata"
	util "github.com/bsm/reason/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClassificationNumerical", func() {
	var subject *observer.ClassificationNumerical
	var pre *util.Vector

	model := testdata.SimpleModel
	target := model.Feature("play")
	feat := model.Feature("humidex")
	crit := split.DefaultCriterion(target)

	BeforeEach(func() {
		subject = observer.NewClassificationNumerical(0)
		pre = util.NewVector()
		for _, x := range testdata.SimpleDataSet {
			subject.Observe(feat.Number(x), target.Category(x)*2)
			pre.Incr(int(target.Category(x)*2), 1.0)
		}
	})

	It("should observe", func() {
		Expect(subject.Dist.Data).To(Equal([]util.NumStream{
			{Weight: 9, Sum: 416, SumSquares: 20242, Min: 35, Max: 62},
			{Weight: 0},
			{Weight: 5, Sum: 287, SumSquares: 16751, Min: 43, Max: 64},
		}))
	})

	It("should calculate probabilty of a value given a target", func() {
		yes, no := target.CategoryOf("yes")*2, target.CategoryOf("no")*2

		Expect(subject.Prob(40, -1)).To(Equal(0.0))
		Expect(subject.Prob(40, 3)).To(Equal(0.0))

		Expect(subject.Prob(40, yes)).To(BeNumerically("~", 0.030, 1e-3))
		Expect(subject.Prob(40, no)).To(BeNumerically("~", 0.005, 1e-3))
		Expect(subject.Prob(50, yes)).To(BeNumerically("~", 0.034, 1e-3))
		Expect(subject.Prob(50, no)).To(BeNumerically("~", 0.032, 1e-3))
		Expect(subject.Prob(60, yes)).To(BeNumerically("~", 0.017, 1e-3))
		Expect(subject.Prob(60, no)).To(BeNumerically("~", 0.046, 1e-3))
	})

	It("should evaluate splits", func() {
		merit, pivot, post := subject.EvaluateSplit(crit, pre)
		Expect(merit).To(BeNumerically("~", 0.176, 1e-3))
		Expect(pivot).To(BeNumerically("~", 42.25, 1e-2))
		Expect(post.Data).To(Equal([]float64{
			3.2587268705731476, 0, 0,
			5.741273129426852, 0, 5,
		}))
	})
})
