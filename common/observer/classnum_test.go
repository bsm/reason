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
		subject = observer.NewClassificationNumerical(11)
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

	It("should evaluate splits", func() {
		merit, pivot, post := subject.EvaluateSplit(crit, pre)
		Expect(merit).To(BeNumerically("~", 0.176, 0.001))
		Expect(pivot).To(BeNumerically("~", 42.25, 0.01))
		Expect(post.Data).To(Equal([]float64{
			3.2587268705731476, 0, 0,
			5.741273129426852, 0, 5,
		}))
	})
})
