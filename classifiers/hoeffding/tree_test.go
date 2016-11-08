package hoeffding

import (
	"testing"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/eval"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {

	DescribeTable("should perform classification",
		func(n int, expInfo *TreeInfo, expCorrect, expKappa float64) {
			if testing.Short() && n > 1000 {
				return
			}

			model := testdata.BigClassificationModel()
			stats := eval.NewClassification(model)
			info, err := runBigDataTest(model, stats, n, "../../testdata/bigcls.csv", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(info).To(Equal(expInfo))
			Expect(stats.Correct() * 100).To(BeNumerically("~", expCorrect, 0.1))
			Expect(stats.Kappa() * 100).To(BeNumerically("~", expKappa, 0.1))
		},

		Entry("1,000", 1000, &TreeInfo{
			NumNodes:        6,
			NumActiveLeaves: 5,
			MaxDepth:        2,
		}, 71.1, 34.8),

		Entry("10,000", 10000, &TreeInfo{
			NumNodes:        38,
			NumActiveLeaves: 30,
			MaxDepth:        4,
		}, 80.3, 59.1),

		Entry("20,000", 20000, &TreeInfo{
			NumNodes:        63,
			NumActiveLeaves: 47,
			MaxDepth:        4,
		}, 84.6, 68.2),
	)

	DescribeTable("should perform regression",
		func(n int, expInfo *TreeInfo, expR2, expRMSE float64) {
			if testing.Short() && n > 1000 {
				return
			}

			model := testdata.BigRegressionModel()
			stats := eval.NewRegression(model)
			info, err := runBigDataTest(model, stats, n, "../../testdata/bigreg.csv", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(info).To(Equal(expInfo))
			Expect(stats.R2()).To(BeNumerically("~", expR2, 0.01))
			Expect(stats.RMSE()).To(BeNumerically("~", expRMSE, 0.01))
		},

		Entry("1,000", 1000, &TreeInfo{
			NumNodes:        603,
			NumActiveLeaves: 602,
			MaxDepth:        2,
		}, 0.13, 0.80),
	)

	It("should prune", func() {
		model := testdata.BigRegressionModel()
		stats := eval.NewRegression(model)
		info, err := runBigDataTest(model, stats, 5000, "../../testdata/bigreg.csv", &Config{
			HeapTarget: 2 * 1024 * 1024,
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(info.NumInactiveLeaves).To(BeNumerically("~", 2100, 100))
		Expect(info.NumNodes).To(BeNumerically("~", 2500, 100))
		Expect(info.NumActiveLeaves).To(BeNumerically("~", 400, 100))
		Expect(stats.R2()).To(BeNumerically("~", 0.17, 0.01))
		Expect(stats.RMSE()).To(BeNumerically("~", 0.97, 0.01))
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "classifiers/hoeffding")
}

func runBigDataTest(model *core.Model, stats eval.Evaluator, n int, fname string, config *Config) (*TreeInfo, error) {
	stream, err := testdata.Open(fname, model)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	insts, err := stream.ReadN(n * 2)
	if err != nil {
		return nil, err
	}

	tree := New(model, config)
	for _, inst := range insts[:n] {
		tree.Train(inst)
	}
	for _, inst := range insts[n:] {
		stats.Record(inst, tree.Predict(inst))
	}

	return tree.Info(), nil
}
