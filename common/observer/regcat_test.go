package observer_test

import (
	"github.com/bsm/reason/common/observer"
	"github.com/bsm/reason/common/split"
	"github.com/bsm/reason/testdata"
	util "github.com/bsm/reason/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegressionCategorical", func() {
	var subject *observer.RegressionCategorical
	var pre *util.NumStream

	model := testdata.SimpleModel
	target := model.Feature("hours")
	feat := model.Feature("outlook")
	crit := split.VarianceReduction{MinWeight: 1}

	BeforeEach(func() {
		subject = observer.NewRegressionCategorical()
		pre = util.NewNumStream()
		for _, x := range testdata.SimpleDataSet {
			subject.Observe(feat.Category(x), target.Number(x))
			pre.Observe(target.Number(x))
		}
	})

	It("should observe", func() {
		Expect(subject.Dist.Data).To(Equal([]util.NumStream{
			{Weight: 5, Sum: 176, SumSquares: 6498, Min: 25, Max: 48},
			{Weight: 4, Sum: 185, SumSquares: 8605, Min: 43, Max: 52},
			{Weight: 5, Sum: 196, SumSquares: 8274, Min: 23, Max: 52},
		}))
	})

	It("should evaluate splits", func() {
		merit, post := subject.EvaluateSplit(crit, pre)
		Expect(merit).To(BeNumerically("~", 9.137, 0.001))
		Expect(post.Data).To(Equal([]util.NumStream{
			{Weight: 5, Sum: 176, SumSquares: 6498, Min: 25, Max: 48},
			{Weight: 4, Sum: 185, SumSquares: 8605, Min: 43, Max: 52},
			{Weight: 5, Sum: 196, SumSquares: 8274, Min: 23, Max: 52},
		}))
	})
})
