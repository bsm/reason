package observer_test

import (
	"github.com/bsm/reason/common/observer"
	"github.com/bsm/reason/common/split"
	"github.com/bsm/reason/testdata"
	util "github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegressionNumerical", func() {
	var subject *observer.RegressionNumerical
	var pre *util.NumStream

	model := testdata.SimpleModel
	target := model.Feature("hours")
	feat := model.Feature("humidex")
	crit := split.VarianceReduction{MinWeight: 1}

	BeforeEach(func() {
		subject = observer.NewRegressionNumerical(4)
		pre = util.NewNumStream()
		for _, x := range testdata.SimpleDataSet {
			subject.Observe(feat.Number(x), target.Number(x))
			pre.Observe(target.Number(x))
		}
	})

	It("should observe", func() {
		Expect(subject.Dist).To(Equal([]observer.RegressionNumerical_Bucket{
			{Threshold: 36.0, NumStream: util.NumStream{Weight: 2, Sum: 100, SumSquares: 5008, Min: 48, Max: 52}},
			{Threshold: 40.8, NumStream: util.NumStream{Weight: 5, Sum: 194, SumSquares: 7874, Min: 23, Max: 46}},
			{Threshold: 60.5, NumStream: util.NumStream{Weight: 6, Sum: 233, SumSquares: 9595, Min: 25, Max: 52}},
			{Threshold: 64.0, NumStream: util.NumStream{Weight: 1, Sum: 30, SumSquares: 900, Min: 30, Max: 30}},
		}))
	})

	It("should evaluate splits", func() {
		merit, pivot, post := subject.EvaluateSplit(crit, pre)
		Expect(merit).To(BeNumerically("~", 17.235, 0.001))
		Expect(pivot).To(BeNumerically("~", 40.8, 0.001))
		Expect(post.Data).To(Equal([]util.NumStream{
			{Weight: 2, Sum: 100, SumSquares: 5008, Min: 0, Max: 52},
			{Weight: 12, Sum: 457, SumSquares: 18369, Min: 0, Max: 52},
		}))
	})
})
