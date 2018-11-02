package hoeffding_test

import (
	"bytes"

	"github.com/bsm/mlmetrics"

	common "github.com/bsm/reason/common/hoeffding"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/regression/hoeffding"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {

	var train = func(n int) (*hoeffding.Tree, *core.Model, []core.Example) {
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

		t1, _, examples := train(5000)
		Expect(t1.Info()).To(Equal(&common.TreeInfo{NumNodes: 3, NumLearning: 2, NumDisabled: 0, MaxDepth: 2}))
		Expect(t1.Predict(nil, examples[6001]).Best().Mean()).To(BeNumerically("~", 1.976, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		t2, err := hoeffding.Load(b1, c)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Info()).To(Equal(&common.TreeInfo{NumNodes: 3, NumLearning: 2, NumDisabled: 0, MaxDepth: 2}))
		Expect(t2.Predict(nil, examples[6001]).Best().Mean()).To(BeNumerically("~", 1.976, 0.001))
	})

	It("should prune", func() {
		t, _, _ := train(10000)
		Expect(t.Info()).To(Equal(&common.TreeInfo{NumNodes: 9, NumLearning: 5, NumDisabled: 0, MaxDepth: 5}))

		t.Prune(3)
		Expect(t.Info()).To(Equal(&common.TreeInfo{NumNodes: 9, NumLearning: 3, NumDisabled: 2, MaxDepth: 5}))
	})

	It("should write TXT", func() {
		t, _, _ := train(10000)

		b := new(bytes.Buffer)
		Expect(t.WriteText(b)).To(Equal(int64(b.Len())))

		s := b.String()
		Expect(s).To(ContainSubstring(`ROOT [weight:3400 mean:3.3 variance:1.2]`))
		Expect(s).To(ContainSubstring("\tn5 <= 1.68 [weight:3362 mean:2.9 variance:1.8]"))
	})

	It("should write DOT", func() {
		t, _, _ := train(10000)

		b := new(bytes.Buffer)
		Expect(t.WriteDOT(b)).To(Equal(int64(b.Len())))

		s := b.String()
		Expect(s).To(ContainSubstring(`N [label="weight: 3400"];`))
		Expect(s).To(ContainSubstring(`N_0_0 [label="n5 <= 1.68\nweight: 3362"];`))
	})

	DescribeTable("should train & predict",
		func(n int, expInfo *common.TreeInfo, exp *testdata.RegressionScore) {
			tree, model, examples := train(n)
			Expect(tree.Info()).To(Equal(expInfo))

			metric := mlmetrics.NewRegression()
			for _, x := range examples[n:] {
				prediction := tree.Predict(nil, x).Best().Mean()
				actual := model.Feature("target").Number(x)
				metric.Observe(actual, prediction)
			}
			Expect(metric.R2()).To(BeNumerically("~", exp.R2, 0.001))
			Expect(metric.RMSE()).To(BeNumerically("~", exp.RMSE, 0.001))
		},

		Entry("1,000", 1000, &common.TreeInfo{
			NumNodes:    1,
			NumLearning: 1,
			MaxDepth:    1,
		}, &testdata.RegressionScore{
			R2:   0.005,
			RMSE: 1.095,
		}),
		Entry("5,000", 5000, &common.TreeInfo{
			NumNodes:    3,
			NumLearning: 2,
			MaxDepth:    2,
		}, &testdata.RegressionScore{
			R2:   0.186,
			RMSE: 0.960,
		}),
		Entry("10,000", 10000, &common.TreeInfo{
			NumNodes:    9,
			NumLearning: 5,
			NumDisabled: 0,
			MaxDepth:    5,
		}, &testdata.RegressionScore{
			R2:   0.618,
			RMSE: 0.616,
		}),
	)
})
