package hoeffding_test

import (
	"bytes"

	"github.com/bsm/mlmetrics"
	"github.com/bsm/reason/classifier/hoeffding"
	"github.com/bsm/reason/common/split"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {

	It("should validate target", func() {
		model := testdata.SimpleModel

		_, err := hoeffding.New(model, "unknown", nil)
		Expect(err).To(MatchError(`hoeffding: unknown feature "unknown"`))

		_, err = hoeffding.New(model, "play", &hoeffding.Config{
			SplitCriterion: split.VarianceReduction{},
		})
		Expect(err).To(MatchError(`hoeffding: split criterion is incompatible with target "play"`))
	})

	It("should dump/load", func() {
		t1, _, examples := runTraining("classification", 3000)
		Expect(t1.Info()).To(Equal(&hoeffding.TreeInfo{NumNodes: 11, NumLearning: 9, MaxDepth: 3}))
		Expect(t1.PredictMC(examples[4001]).Prob(0)).To(BeNumerically("~", 0.273, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		t2, err := hoeffding.LoadFrom(b1, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Info()).To(Equal(&hoeffding.TreeInfo{NumNodes: 11, NumLearning: 9, MaxDepth: 3}))
		Expect(t2.PredictMC(examples[4001]).Prob(0)).To(BeNumerically("~", 0.273, 0.001))
	})

	It("should prune", func() {
		t, _, _ := runTraining("classification", 3000)
		Expect(t.Info()).To(Equal(&hoeffding.TreeInfo{
			NumNodes:    11,
			NumLearning: 9,
			NumDisabled: 0,
			MaxDepth:    3,
		}))

		t.Prune(5)
		Expect(t.Info()).To(Equal(&hoeffding.TreeInfo{
			NumNodes:    11,
			NumLearning: 5,
			NumDisabled: 4,
			MaxDepth:    3,
		}))
	})

	It("should write TXT", func() {
		t, _, _ := runTraining("classification", 3000)

		b := new(bytes.Buffer)
		Expect(t.WriteText(b)).To(Equal(int64(b.Len())))

		s := b.String()
		Expect(s).To(ContainSubstring(`ROOT [weight:600]`))
		Expect(s).To(ContainSubstring("\tc5 = v1 [weight:644]"))
	})

	It("should write DOT", func() {
		t, _, _ := runTraining("classification", 3000)

		b := new(bytes.Buffer)
		Expect(t.WriteDOT(b)).To(Equal(int64(b.Len())))

		s := b.String()
		Expect(s).To(ContainSubstring(`N [label="weight: 600"];`))
		Expect(s).To(ContainSubstring(`N_0 [label="c5 = v1\nweight: 644"];`))
	})

	DescribeTable("classification",
		func(n int, expInfo *hoeffding.TreeInfo, exp *testdata.ClassificationScore) {
			tree, model, examples := runTraining("classification", n)
			Expect(tree.Info()).To(Equal(expInfo))

			m1 := mlmetrics.NewConfusionMatrix()
			m2 := mlmetrics.NewLogLoss()
			for _, x := range examples[n:] {
				prediction := tree.PredictMC(x)
				actual := model.Feature("target").Category(x)

				m1.Observe(int(actual), int(prediction.Category()))
				m2.Observe(prediction.Prob(actual))
			}

			Expect(m1.Accuracy()).To(BeNumerically("~", exp.Accuracy, 0.001))
			Expect(m1.Kappa()).To(BeNumerically("~", exp.Kappa, 0.001))
			Expect(m2.Score()).To(BeNumerically("~", exp.LogLoss, 0.001))
		},

		Entry("1,000", 1000, &hoeffding.TreeInfo{
			NumNodes:    6,
			NumLearning: 5,
			MaxDepth:    2,
		}, &testdata.ClassificationScore{
			Accuracy: 0.711,
			Kappa:    0.348,
			LogLoss:  0.561,
		}),
		Entry("10,000", 10000, &hoeffding.TreeInfo{
			NumNodes:    38,
			NumLearning: 30,
			MaxDepth:    4,
		}, &testdata.ClassificationScore{
			Accuracy: 0.803,
			Kappa:    0.594,
			LogLoss:  0.449,
		}),
		Entry("20,000", 20000, &hoeffding.TreeInfo{
			NumNodes:    65,
			NumLearning: 48,
			MaxDepth:    4,
		}, &testdata.ClassificationScore{
			Accuracy: 0.850,
			Kappa:    0.690,
			LogLoss:  0.372,
		}),
	)

	DescribeTable("regression",
		func(n int, expInfo *hoeffding.TreeInfo, exp *testdata.RegressionScore) {
			tree, model, examples := runTraining("regression", n)
			Expect(tree.Info()).To(Equal(expInfo))

			metric := mlmetrics.NewRegression()
			for _, x := range examples[n:] {
				prediction := tree.PredictNum(x).Number()
				actual := model.Feature("target").Number(x)
				metric.Observe(actual, prediction)
			}
			Expect(metric.R2()).To(BeNumerically("~", exp.R2, 0.001))
			Expect(metric.RMSE()).To(BeNumerically("~", exp.RMSE, 0.001))
		},

		Entry("1,000", 1000, &hoeffding.TreeInfo{
			NumNodes:    6,
			NumLearning: 5,
			MaxDepth:    2,
		}, &testdata.RegressionScore{
			R2:   0.837,
			RMSE: 1.652,
		}),
		Entry("5,000", 5000, &hoeffding.TreeInfo{
			NumNodes:    31,
			NumLearning: 25,
			MaxDepth:    3,
		}, &testdata.RegressionScore{
			R2:   0.899,
			RMSE: 1.280,
		}),
		FEntry("10,000", 10000, &hoeffding.TreeInfo{
			NumNodes:    81,
			NumLearning: 50,
			NumDisabled: 0,
			MaxDepth:    4,
		}, &testdata.RegressionScore{
			R2:   0.932,
			RMSE: 1.056,
		}),
	)
})

func runTraining(kind string, n int) (*hoeffding.Tree, *core.Model, []core.Example) {
	stream, model, err := testdata.OpenBigData(kind, "../../testdata")
	Expect(err).NotTo(HaveOccurred())
	defer stream.Close()

	examples, err := stream.ReadN(n * 2)
	Expect(err).NotTo(HaveOccurred())

	tree, err := hoeffding.New(model, "target", nil)
	Expect(err).NotTo(HaveOccurred())

	for _, x := range examples[:n] {
		tree.Train(x)
	}
	return tree, model, examples
}
