package hoeffding_test

import (
	"bytes"

	"github.com/bsm/reason/classifier/hoeffding"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	_ "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {

	It("should dump/load", func() {
		c := &hoeffding.Config{
			GracePeriod: 50,
		}

		t1, _, examples := trainClassification(3000)
		Expect(t1.Info()).To(Equal(&hoeffding.TreeInfo{NumNodes: 11, NumLearning: 9, MaxDepth: 3}))
		Expect(t1.Predict(nil, examples[4001]).Best().P(0)).To(BeNumerically("~", 0.273, 0.001))

		b1 := new(bytes.Buffer)
		Expect(t1.WriteTo(b1)).To(Equal(int64(b1.Len())))

		t2, err := hoeffding.LoadFrom(b1, c)
		Expect(err).NotTo(HaveOccurred())
		Expect(t2.Info()).To(Equal(&hoeffding.TreeInfo{NumNodes: 11, NumLearning: 9, MaxDepth: 3}))
		Expect(t2.Predict(nil, examples[4001]).Best().P(0)).To(BeNumerically("~", 0.273, 0.001))
	})

	// It("should prune", func() {
	// 	t, _, _ := trainClassification(3000)
	// 	Expect(t.Info()).To(Equal(&hoeffding.TreeInfo{
	// 		NumNodes:    11,
	// 		NumLearning: 9,
	// 		NumDisabled: 0,
	// 		MaxDepth:    3,
	// 	}))

	// 	t.Prune(5)
	// 	Expect(t.Info()).To(Equal(&hoeffding.TreeInfo{
	// 		NumNodes:    11,
	// 		NumLearning: 5,
	// 		NumDisabled: 4,
	// 		MaxDepth:    3,
	// 	}))
	// })

	// It("should write TXT", func() {
	// 	t, _, _ := trainClassification(3000)

	// 	b := new(bytes.Buffer)
	// 	Expect(t.WriteText(b)).To(Equal(int64(b.Len())))

	// 	s := b.String()
	// 	Expect(s).To(ContainSubstring(`ROOT [weight:600]`))
	// 	Expect(s).To(ContainSubstring("\tc5 = v1 [weight:644]"))
	// })

	// It("should write DOT", func() {
	// 	t, _, _ := trainClassification(3000)

	// 	b := new(bytes.Buffer)
	// 	Expect(t.WriteDOT(b)).To(Equal(int64(b.Len())))

	// 	s := b.String()
	// 	Expect(s).To(ContainSubstring(`N [label="weight: 600"];`))
	// 	Expect(s).To(ContainSubstring(`N_0 [label="c5 = v1\nweight: 644"];`))
	// })

	// DescribeTable("should train & predict",
	// 	func(n int, expInfo *hoeffding.TreeInfo, exp *testdata.ClassificationScore) {
	// 		tree, model, examples := trainClassification(n)
	// 		Expect(tree.Info()).To(Equal(expInfo))

	// 		accuracy := mlmetrics.NewAccuracy()
	// 		confusion := mlmetrics.NewConfusionMatrix()
	// 		logLoss := mlmetrics.NewLogLoss()

	// 		for _, x := range examples[n:] {
	// 			predicted, probability := tree.Predict(nil, x).Best().Top()
	// 			actual := model.Feature("target").Category(x)

	// 			accuracy.Observe(int(actual), int(predicted))
	// 			confusion.Observe(int(actual), int(predicted))
	// 			logLoss.Observe(probability)
	// 		}

	// 		Expect(accuracy.Rate() * 100).To(BeNumerically("~", exp.Accuracy, 0.1))
	// 		Expect(confusion.Kappa()).To(BeNumerically("~", exp.Kappa, 0.001))
	// 		Expect(logLoss.Score()).To(BeNumerically("~", exp.LogLoss, 0.001))
	// 	},

	// 	Entry("1,000", 1000, &hoeffding.TreeInfo{
	// 		NumNodes:    6,
	// 		NumLearning: 5,
	// 		MaxDepth:    2,
	// 	}, &testdata.ClassificationScore{
	// 		Accuracy: 71.1,
	// 		Kappa:    0.348,
	// 		LogLoss:  0.349,
	// 	}),
	// 	Entry("10,000", 10000, &hoeffding.TreeInfo{
	// 		NumNodes:    38,
	// 		NumLearning: 30,
	// 		MaxDepth:    4,
	// 	}, &testdata.ClassificationScore{
	// 		Accuracy: 80.3,
	// 		Kappa:    0.594,
	// 		LogLoss:  0.230,
	// 	}),
	// 	Entry("20,000", 20000, &hoeffding.TreeInfo{
	// 		NumNodes:    65,
	// 		NumLearning: 48,
	// 		MaxDepth:    4,
	// 	}, &testdata.ClassificationScore{
	// 		Accuracy: 85.0,
	// 		Kappa:    0.690,
	// 		LogLoss:  0.183,
	// 	}),
	// )

})

func trainClassification(n int) (*hoeffding.Tree, *core.Model, []core.Example) {
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
