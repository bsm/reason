package ftrl_test

import (
	"bytes"
	"testing"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/regression"
	"github.com/bsm/reason/regression/ftrl"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Optimizer", func() {

	var train = func(n int) (*ftrl.Optimizer, *core.Model, []core.Example) {
		stream, model, err := testdata.OpenRegression("../../testdata")
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
		Expect(t1.Predict(examples[4001])).To(BeNumerically("~", 0.252, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(BeNumerically("~", 18891000, 1000))
		Expect(b1.Len()).To(BeNumerically("~", 18891000, 1000))

		t2, err := ftrl.Load(b1)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Predict(examples[4001])).To(BeNumerically("~", 0.252, 0.001))
	})

	DescribeTable("should train & predict",
		func(n int, expR2, expRMSE float64) {
			if testing.Short() && n > 1000 {
				return
			}

			opt, model, examples := train(n)
			eval := regression.NewEvaluator()
			for _, x := range examples[n:] {
				prediction := opt.Predict(x)
				actual := model.Feature("target").Number(x)
				eval.Record(prediction, actual, 1.0)
			}
			Expect(eval.R2()).To(BeNumerically("~", expR2, 0.001))
			Expect(eval.RMSE()).To(BeNumerically("~", expRMSE, 0.001))
		},

		Entry("1,000", 1000, 0.044, 0.836),
		Entry("5,000", 5000, 0.032, 1.047),
		Entry("10,000", 10000, 0.053, 0.970),
		Entry("20,000", 20000, 0.177, 0.459),
	)
})
