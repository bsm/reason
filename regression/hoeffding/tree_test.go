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

		t1, _, examples := train(3000)
		Expect(t1.Info()).To(Equal(&common.TreeInfo{NumNodes: 31, NumLearning: 25, NumDisabled: 0, MaxDepth: 3}))
		Expect(t1.Predict(nil, examples[4001]).Best().Mean()).To(BeNumerically("~", 17.736, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		t2, err := hoeffding.Load(b1, c)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Info()).To(Equal(&common.TreeInfo{NumNodes: 31, NumLearning: 25, NumDisabled: 0, MaxDepth: 3}))
		Expect(t2.Predict(nil, examples[4001]).Best().Mean()).To(BeNumerically("~", 17.736, 0.001))
	})

	It("should prune", func() {
		t, _, _ := train(3000)
		Expect(t.Info()).To(Equal(&common.TreeInfo{NumNodes: 31, NumLearning: 25, NumDisabled: 0, MaxDepth: 3}))

		t.Prune(5)
		Expect(t.Info()).To(Equal(&common.TreeInfo{NumNodes: 31, NumLearning: 5, NumDisabled: 20, MaxDepth: 3}))
	})

	It("should write TXT", func() {
		t, _, _ := train(3000)

		b := new(bytes.Buffer)
		Expect(t.WriteText(b)).To(Equal(int64(b.Len())))

		s := b.String()
		Expect(s).To(ContainSubstring(`ROOT [weight:200 mean:12.1 variance:17.6]`))
		Expect(s).To(ContainSubstring("\t\tc4 = v1 [weight:109 mean:5.6 variance:0.9]"))
	})

	It("should write DOT", func() {
		t, _, _ := train(3000)

		b := new(bytes.Buffer)
		Expect(t.WriteDOT(b)).To(Equal(int64(b.Len())))

		s := b.String()
		Expect(s).To(ContainSubstring(`N [label="weight: 200"];`))
		Expect(s).To(ContainSubstring(`N_0_3 [label="c4 = v4\nweight: 120"];`))
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
			NumNodes:    6,
			NumLearning: 5,
			MaxDepth:    2,
		}, &testdata.RegressionScore{
			R2:   0.837,
			RMSE: 1.652,
		}),
		Entry("5,000", 5000, &common.TreeInfo{
			NumNodes:    31,
			NumLearning: 25,
			MaxDepth:    3,
		}, &testdata.RegressionScore{
			R2:   0.899,
			RMSE: 1.280,
		}),
		Entry("10,000", 10000, &common.TreeInfo{
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
