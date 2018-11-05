package ftrl_test

import (
	"bytes"

	"github.com/bsm/mlmetrics"

	"github.com/bsm/reason/classifier/ftrl"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Optimizer", func() {
	train := func(n int) (*ftrl.Optimizer, *core.Model, []core.Example) {
		stream, model, err := testdata.OpenBigData("classification", "../../testdata")
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		examples, err := stream.ReadN(n * 2)
		Expect(err).NotTo(HaveOccurred())

		opt, err := ftrl.New(model, "target", nil)
		Expect(err).NotTo(HaveOccurred())

		for _, x := range examples[:n] {
			opt.Train(x)
		}
		return opt, model, examples
	}

	It("should dump/load", func() {
		t1, _, examples := train(3000)
		Expect(t1.Predict(examples[4001])).To(BeNumerically("~", 0.785, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		t2, err := ftrl.Load(b1, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Predict(examples[4001])).To(BeNumerically("~", 0.785, 0.001))
	})

	DescribeTable("should train & predict",
		func(n int, exp *testdata.ClassificationScore) {
			opt, model, examples := train(n)
			m1 := mlmetrics.NewConfusionMatrix()
			m2 := mlmetrics.NewLogLoss()
			for _, x := range examples[n:] {
				prob := opt.Predict(x)
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

		Entry("1,000", 1000, &testdata.ClassificationScore{
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
