package ftrl_test

import (
	"bytes"
	"testing"

	"github.com/bsm/reason/classification/eval"
	"github.com/bsm/reason/classification/ftrl"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Optimizer", func() {

	var train = func(n int) (*ftrl.Optimizer, *core.Model, []core.Example) {
		stream, model, err := testdata.OpenClassification("../../testdata")
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		examples, err := stream.ReadN(n * 2)
		Expect(err).NotTo(HaveOccurred())

		opt, err := ftrl.New(model, "target", nil)
		Expect(err).NotTo(HaveOccurred())

		for _, x := range examples[:n] {
			opt.Train(x, 1.0)
		}
		return opt, model, examples
	}

	It("should dump/load", func() {
		t1, _, examples := train(3000)
		Expect(t1.Predict(examples[4001]).P(0)).To(BeNumerically("~", 0.215, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		t2, err := ftrl.Load(b1, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Predict(examples[4001]).P(0)).To(BeNumerically("~", 0.215, 0.001))
	})

	DescribeTable("should train & predict",
		func(n int, exp *testdata.ClassificationScore) {
			opt, model, examples := train(n)
			accuracy := eval.NewAccuracy()
			kappa := eval.NewKappa()
			logLoss := eval.NewLogLoss()

			for _, x := range examples[n:] {
				predicted, prob := opt.Predict(x).Top()
				actual := model.Feature("target").Category(x)

				accuracy.Record(predicted, actual, 1.0)
				kappa.Record(predicted, actual, 1.0)
				logLoss.Record(prob, 1.0)
			}

			Expect(accuracy.Accuracy() * 100).To(BeNumerically("~", exp.Accuracy, 0.1))
			Expect(kappa.Score()).To(BeNumerically("~", exp.Kappa, 0.001))
			Expect(logLoss.Value()).To(BeNumerically("~", exp.LogLoss, 0.001))
		},

		Entry("1,000", 1000, &testdata.ClassificationScore{
			Accuracy: 73.0,
			Kappa:    0.399,
			LogLoss:  0.434,
		}),
		Entry("10,000", 10000, &testdata.ClassificationScore{
			Accuracy: 73.5,
			Kappa:    0.441,
			LogLoss:  0.349,
		}),
		Entry("20,000", 20000, &testdata.ClassificationScore{
			Accuracy: 73.8,
			Kappa:    0.447,
			LogLoss:  0.343,
		}),
	)

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "classification/ftrl")
}
