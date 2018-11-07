package bayes_test

import (
	"testing"

	"github.com/bsm/mlmetrics"
	"github.com/bsm/reason/classifier/bayes"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Naive", func() {
	PIt("should dump/load", func() {
		// t1, _, examples := train(3000)
		// Expect(t1.Predict(examples[4001])).To(BeNumerically("~", 0.785, 0.001))

		// b1 := new(bytes.Buffer)
		// Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		// t2, err := ftrl.LoadFrom(b1, nil)
		// Expect(err).NotTo(HaveOccurred())
		// Expect(t2.Predict(examples[4001])).To(BeNumerically("~", 0.785, 0.001))
	})

	PDescribeTable("should train & predict",
		func(n int, exp *testdata.ClassificationScore) {
			cls, model, examples := runClassification(n)
			cls.Info()

			m1 := mlmetrics.NewConfusionMatrix()
			m2 := mlmetrics.NewLogLoss()
			for _, x := range examples[n:] {
				prob := cls.Predict(x)
				actual := int(model.Feature("target").Category(x))
				if prob <= 0.5 {
					m1.Observe(actual, 0)
				} else {
					m1.Observe(actual, 1)
				}
				if actual == 0 {
					m2.Observe(prob)
				} else if actual == 1 {
					m2.Observe(1 - prob)
				}
			}
			Expect(m1.Accuracy()).To(BeNumerically("~", exp.Accuracy, 0.001))
			Expect(m1.Kappa()).To(BeNumerically("~", exp.Kappa, 0.001))
			Expect(m2.Score()).To(BeNumerically("~", exp.LogLoss, 0.001))
		},

		FEntry("1,000", 1000, &testdata.ClassificationScore{
			Accuracy: 0.730,
			Kappa:    0.399,
			LogLoss:  0.971,
		}),
		Entry("5,000", 5000, &testdata.ClassificationScore{
			Accuracy: 0.731,
			Kappa:    0.432,
			LogLoss:  1.105,
		}),
		Entry("10,000", 10000, &testdata.ClassificationScore{
			Accuracy: 0.735,
			Kappa:    0.441,
			LogLoss:  1.196,
		}),
		Entry("20,000", 20000, &testdata.ClassificationScore{
			Accuracy: 0.738,
			Kappa:    0.447,
			LogLoss:  1.224,
		}),
	)
})

// --------------------------------------------------------------------
func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "classifier/bayes")
}

func runClassification(n int) (*bayes.Naive, *core.Model, []core.Example) {
	stream, model, err := testdata.OpenBigData("classification", "../../testdata")
	Expect(err).NotTo(HaveOccurred())
	defer stream.Close()

	examples, err := stream.ReadN(n * 2)
	Expect(err).NotTo(HaveOccurred())

	cls, err := bayes.New(model, "target", nil)
	Expect(err).NotTo(HaveOccurred())

	for _, x := range examples[:n] {
		cls.Train(x)
	}
	return cls, model, examples
}
