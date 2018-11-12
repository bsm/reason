package ftrl_test

import (
	"bytes"
	"testing"

	"github.com/bsm/mlmetrics"
	"github.com/bsm/reason"
	"github.com/bsm/reason/classifier/ftrl"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("FTRL", func() {
	It("should dump/load", func() {
		t1, _, examples := runTraining("classification", 1000)
		Expect(t1.Predict(examples[1101])).To(BeNumerically("~", 0.640, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		t2, err := ftrl.LoadFrom(b1, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Predict(examples[1101])).To(BeNumerically("~", 0.640, 0.001))
	})

	DescribeTable("classification",
		func(n int, exp *testdata.ClassificationScore) {
			opt, target, examples := runTraining("classification", n)
			m1 := mlmetrics.NewConfusionMatrix()
			m2 := mlmetrics.NewLogLoss()
			for _, x := range examples[n:] {
				prediction := opt.Predict(x)
				actual := target.Category(x)
				m1.Observe(int(actual), int(prediction.Category()))
				m2.Observe(prediction.Prob(actual))
			}
			Expect(m1.Accuracy()).To(BeNumerically("~", exp.Accuracy, 0.001))
			Expect(m1.Kappa()).To(BeNumerically("~", exp.Kappa, 0.001))
			Expect(m2.Score()).To(BeNumerically("~", exp.LogLoss, 0.001))
		},

		Entry("1,000", 1000, &testdata.ClassificationScore{
			Accuracy: 0.730,
			Kappa:    0.399,
			LogLoss:  0.568,
		}),
		Entry("5,000", 5000, &testdata.ClassificationScore{
			Accuracy: 0.731,
			Kappa:    0.432,
			LogLoss:  0.547,
		}),
		Entry("10,000", 10000, &testdata.ClassificationScore{
			Accuracy: 0.735,
			Kappa:    0.441,
			LogLoss:  0.546,
		}),
		Entry("20,000", 20000, &testdata.ClassificationScore{
			Accuracy: 0.738,
			Kappa:    0.447,
			LogLoss:  0.540,
		}),
	)

	DescribeTable("regression",
		func(n int, exp *testdata.RegressionScore) {
			opt, target, examples := runTraining("regression", n)
			metric := mlmetrics.NewRegression()
			for _, x := range examples[n:] {
				prediction := opt.PredictNum(x).Number()
				actual := target.Number(x)
				metric.Observe(actual, prediction)
			}
			Expect(metric.R2()).To(BeNumerically("~", exp.R2, 0.001))
			Expect(metric.RMSE()).To(BeNumerically("~", exp.RMSE, 0.001))
		},

		Entry("1,000", 1000, &testdata.RegressionScore{
			R2:   0.372,
			RMSE: 3.247,
		}),
		Entry("5,000", 5000, &testdata.RegressionScore{
			R2:   0.798,
			RMSE: 1.815,
		}),
		Entry("10,000", 10000, &testdata.RegressionScore{
			R2:   0.881,
			RMSE: 1.396,
		}),
		Entry("20,000", 20000, &testdata.RegressionScore{
			R2:   0.923,
			RMSE: 1.119,
		}),
	)
})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "classifier/ftrl")
}

func runTraining(kind string, n int) (*ftrl.FTRL, *reason.Feature, []reason.Example) {
	stream, model, err := testdata.OpenBigData(kind, "../../testdata")
	Expect(err).NotTo(HaveOccurred())
	defer stream.Close()

	examples, err := stream.ReadN(n * 2)
	Expect(err).NotTo(HaveOccurred())

	opt, err := ftrl.New(model, "target", nil)
	Expect(err).NotTo(HaveOccurred())

	for _, x := range examples[:n] {
		opt.Train(x)
	}
	return opt, model.Feature("target"), examples
}
