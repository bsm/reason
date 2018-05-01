package hoeffding_test

import (
	"bytes"
	"testing"

	common "github.com/bsm/reason/common/hoeffding"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/regression"
	"github.com/bsm/reason/regression/hoeffding"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {

	var train = func(n int, config *hoeffding.Config) (*hoeffding.Tree, *core.Model, []core.Example) {
		stream, model, err := testdata.OpenRegression("../../testdata")
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		examples, err := stream.ReadN(n * 2)
		Expect(err).NotTo(HaveOccurred())

		tree, err := hoeffding.New(model, "target", nil)
		Expect(err).NotTo(HaveOccurred())

		for _, x := range examples[:n] {
			tree.Train(x, 1.0)
		}
		return tree, model, examples
	}

	It("should dump/load", func() {
		c := &hoeffding.Config{
			Config: common.Config{GracePeriod: 50},
		}

		t1, _, examples := train(3000, c)
		Expect(t1.Info()).To(Equal(&common.TreeInfo{NumNodes: 1474, NumLearning: 1473, MaxDepth: 2}))
		Expect(t1.Predict(nil, examples[4001]).Best().Mean()).To(BeNumerically("~", 0.472, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(BeNumerically("~", 188800, 1000))
		Expect(b1.Len()).To(BeNumerically("~", 188800, 1000))

		t2, err := hoeffding.Load(b1, c)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Info()).To(Equal(&common.TreeInfo{NumNodes: 1474, NumLearning: 1473, MaxDepth: 2}))
		Expect(t2.Predict(nil, examples[4001]).Best().Mean()).To(BeNumerically("~", 0.472, 0.001))
	})

	It("should prune", func() {
		t, _, _ := train(3000, nil)
		Expect(t.Info()).To(Equal(&common.TreeInfo{
			NumNodes:    1474,
			NumLearning: 1473,
			NumDisabled: 0,
			MaxDepth:    2,
		}))

		t.Prune(10)
		Expect(t.Info()).To(Equal(&common.TreeInfo{
			NumNodes:    1474,
			NumLearning: 10,
			NumDisabled: 1463,
			MaxDepth:    2,
		}))
	})

	It("should write TXT", func() {
		t, _, _ := train(3000, nil)

		b := new(bytes.Buffer)
		Expect(t.WriteText(b)).To(Equal(int64(63782)))
		Expect(b.Len()).To(Equal(63782))

		s := b.String()
		Expect(s).To(ContainSubstring(`ROOT [weight:2400, mean:0, variance: 1]`))
		Expect(s).To(ContainSubstring("\tc4 = v4 [weight:23, mean:0, variance: 0]"))
	})

	It("should write DOT", func() {
		t, _, _ := train(3000, nil)

		b := new(bytes.Buffer)
		Expect(t.WriteDOT(b)).To(Equal(int64(80763)))
		Expect(b.Len()).To(Equal(80763))

		s := b.String()
		Expect(s).To(ContainSubstring(`N [label="weight: 2400"];`))
		Expect(s).To(ContainSubstring(`N_4 [label="c4 = v4\nweight: 23"];`))
	})

	DescribeTable("should train & predict",
		func(n int, expInfo *common.TreeInfo, expR2, expRMSE float64) {
			if testing.Short() && n > 1000 {
				return
			}

			tree, model, examples := train(n, nil)
			Expect(tree.Info()).To(Equal(expInfo))

			eval := regression.NewEvaluator()
			for _, x := range examples[n:] {
				prediction := tree.Predict(nil, x).Best().Mean()
				actual := model.Feature("target").Number(x)
				eval.Record(prediction, actual, 1.0)
			}
			Expect(eval.R2()).To(BeNumerically("~", expR2, 0.001))
			Expect(eval.RMSE()).To(BeNumerically("~", expRMSE, 0.001))
		},

		Entry("1,000", 1000, &common.TreeInfo{
			NumNodes:    1,
			NumLearning: 1,
			MaxDepth:    1,
		}, 0.002, 0.854),

		Entry("5,000", 5000, &common.TreeInfo{
			NumNodes:    2224,
			NumLearning: 2223,
			MaxDepth:    2,
		}, 0.170, 0.970),

		Entry("10,000", 10000, &common.TreeInfo{
			NumNodes:    3688,
			NumLearning: 3687,
			NumDisabled: 0,
			MaxDepth:    2,
		}, 0.211, 0.885),
	)

})
