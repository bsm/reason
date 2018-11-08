package bayes_test

import (
	"bytes"
	"testing"

	"github.com/bsm/mlmetrics"
	"github.com/bsm/reason/classifier/bayes"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("NaiveBayes", func() {
	It("should dump/load", func() {
		t1, _, examples := runClassification(3000)
		p1 := t1.PredictMC(examples[4001])
		Expect(p1.Prob(p1.Category())).To(BeNumerically("~", 0.824, 1e-3))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		t2, err := bayes.LoadFrom(b1, nil)
		Expect(err).NotTo(HaveOccurred())
		p2 := t2.PredictMC(examples[4001])
		Expect(p2.Prob(p2.Category())).To(BeNumerically("~", 0.824, 1e-3))
	})

	DescribeTable("should train & predict",
		func(n int, exp *testdata.ClassificationScore) {
			cls, model, examples := runClassification(n)
			met := mlmetrics.NewConfusionMatrix()
			target := model.Feature("target")
			for _, x := range examples[n:] {
				actual := target.Category(x)
				predicted := cls.PredictMC(x).Category()
				met.Observe(int(actual), int(predicted))

			}
			Expect(met.Accuracy()).To(BeNumerically("~", exp.Accuracy, 1e-3))
			Expect(met.Kappa()).To(BeNumerically("~", exp.Kappa, 1e-3))
		},

		Entry("1,000", 1000, &testdata.ClassificationScore{
			Accuracy: 0.733,
			Kappa:    0.419,
		}),
		Entry("5,000", 5000, &testdata.ClassificationScore{
			Accuracy: 0.728,
			Kappa:    0.423,
		}),
		Entry("10,000", 10000, &testdata.ClassificationScore{
			Accuracy: 0.732,
			Kappa:    0.434,
		}),
		Entry("20,000", 20000, &testdata.ClassificationScore{
			Accuracy: 0.736,
			Kappa:    0.441,
		}),
	)
})

// --------------------------------------------------------------------
func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "classifier/bayes")
}

func runClassification(n int) (*bayes.NaiveBayes, *core.Model, []core.Example) {
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
