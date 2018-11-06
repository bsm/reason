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
})
