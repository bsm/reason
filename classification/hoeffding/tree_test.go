package hoeffding_test

import (
	"bytes"
	"testing"

	"github.com/bsm/reason/classification/eval"
	"github.com/bsm/reason/classification/hoeffding"
	common "github.com/bsm/reason/common/hoeffding"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {

	var train = func(n int, config *hoeffding.Config) (*hoeffding.Tree, *core.Model, []core.Example) {
		stream, model, err := testdata.OpenClassification("../../testdata")
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

		t1, _, examples := train(3000, nil)
		Expect(t1.Info()).To(Equal(&common.TreeInfo{NumNodes: 11, NumLearning: 9, MaxDepth: 3}))
		Expect(t1.Predict(nil, examples[4001]).Best().P(0)).To(BeNumerically("~", 0.273, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(BeNumerically("~", 13371, 100))
		Expect(b1.Len()).To(BeNumerically("~", 13371, 1000))

		t2, err := hoeffding.Load(b1, c)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Info()).To(Equal(&common.TreeInfo{NumNodes: 11, NumLearning: 9, MaxDepth: 3}))
		Expect(t2.Predict(nil, examples[4001]).Best().P(0)).To(BeNumerically("~", 0.273, 0.001))
	})

	It("should prune", func() {
		t, _, _ := train(3000, nil)
		Expect(t.Info()).To(Equal(&common.TreeInfo{
			NumNodes:    11,
			NumLearning: 9,
			NumDisabled: 0,
			MaxDepth:    3,
		}))

		t.Prune(5)
		Expect(t.Info()).To(Equal(&common.TreeInfo{
			NumNodes:    11,
			NumLearning: 5,
			NumDisabled: 4,
			MaxDepth:    3,
		}))
	})

	It("should write TXT", func() {
		t, _, _ := train(3000, nil)

		b := new(bytes.Buffer)
		Expect(t.WriteText(b)).To(Equal(int64(241)))
		Expect(b.Len()).To(Equal(241))

		s := b.String()
		Expect(s).To(ContainSubstring(`ROOT [weight:600]`))
		Expect(s).To(ContainSubstring("\tc5 = v1 [weight:644]"))
	})

	It("should write DOT", func() {
		t, _, _ := train(3000, nil)

		b := new(bytes.Buffer)
		Expect(t.WriteDOT(b)).To(Equal(int64(625)))
		Expect(b.Len()).To(Equal(625))

		s := b.String()
		Expect(s).To(ContainSubstring(`N [label="weight: 600"];`))
		Expect(s).To(ContainSubstring(`N_0 [label="c5 = v1\nweight: 644"];`))
	})

	DescribeTable("should train & predict",
		func(n int, expInfo *common.TreeInfo, expAccuracy, expKappa float64) {
			if testing.Short() && n > 1000 {
				return
			}

			tree, model, examples := train(n, nil)
			Expect(tree.Info()).To(Equal(expInfo))

			accuracy := eval.NewAccuracy()
			kappa := eval.NewKappa()

			for _, x := range examples[n:] {
				predicted, _ := tree.Predict(nil, x).Best().Top()
				actual := model.Feature("target").Category(x)

				accuracy.Record(predicted, actual, 1.0)
				kappa.Record(predicted, actual, 1.0)
			}

			Expect(accuracy.Accuracy() * 100).To(BeNumerically("~", expAccuracy, 0.1))
			Expect(kappa.Score() * 100).To(BeNumerically("~", expKappa, 0.1))
		},

		Entry("1,000", 1000, &common.TreeInfo{
			NumNodes:    6,
			NumLearning: 5,
			MaxDepth:    2,
		}, 71.1, 34.8),

		Entry("10,000", 10000, &common.TreeInfo{
			NumNodes:    38,
			NumLearning: 30,
			MaxDepth:    4,
		}, 80.3, 59.4),

		Entry("20,000", 20000, &common.TreeInfo{
			NumNodes:    65,
			NumLearning: 48,
			MaxDepth:    4,
		}, 85.0, 69.0),
	)

})
