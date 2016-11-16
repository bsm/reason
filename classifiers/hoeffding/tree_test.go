package hoeffding

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/eval"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tree", func() {

	It("should dump/load", func() {
		model := testdata.BigClassificationModel()
		stream, err := testdata.Open("../../testdata/bigcls.csv", model)
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		tree := New(model, nil)
		n := 0
		for stream.Next() {
			tree.Train(stream.Instance())
			if n++; n >= 500 {
				break
			}
		}
		Expect(stream.Err()).NotTo(HaveOccurred())

		file, err := ioutil.TempFile("", "bsm-reason-test")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		err = tree.WriteTo(file)
		Expect(err).NotTo(HaveOccurred())
		Expect(file.Close()).NotTo(HaveOccurred())

		file, err = os.Open(file.Name())
		Expect(err).NotTo(HaveOccurred())

		info, err := file.Stat()
		Expect(err).NotTo(HaveOccurred())
		Expect(info.Size()).To(BeNumerically("~", 3600, 20))

		tree2, err := Load(file, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(tree2.root).To(Equal(tree.root))
		Expect(tree2.model).To(Equal(tree.model))
		Expect(tree2.conf).To(Equal(tree.conf))
	})

	DescribeTable("should perform classification",
		func(n int, expInfo *TreeInfo, expCorrect, expKappa float64) {
			if testing.Short() && n > 1000 {
				return
			}

			model := testdata.BigClassificationModel()
			stats := eval.NewClassification(model)
			info, err := runBigDataTest(model, stats, n, "../../testdata/bigcls.csv", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(info).To(Equal(expInfo))
			Expect(stats.Correct() * 100).To(BeNumerically("~", expCorrect, 0.1))
			Expect(stats.Kappa() * 100).To(BeNumerically("~", expKappa, 0.1))
		},

		Entry("1,000", 1000, &TreeInfo{
			NumNodes:        6,
			NumActiveLeaves: 5,
			MaxDepth:        2,
		}, 71.1, 34.8),

		Entry("10,000", 10000, &TreeInfo{
			NumNodes:        38,
			NumActiveLeaves: 30,
			MaxDepth:        4,
		}, 80.3, 59.1),

		Entry("20,000", 20000, &TreeInfo{
			NumNodes:        63,
			NumActiveLeaves: 47,
			MaxDepth:        4,
		}, 84.6, 68.2),
	)

	DescribeTable("should perform regression",
		func(n int, c *Config, expInfo *TreeInfo, expR2, expRMSE float64) {
			if testing.Short() && n > 1000 {
				return
			}

			model := testdata.BigRegressionModel()
			stats := eval.NewRegression(model)
			info, err := runBigDataTest(model, stats, n, "../../testdata/bigreg.csv", c)
			Expect(err).NotTo(HaveOccurred())
			Expect(info.NumNodes).To(BeNumerically("~", expInfo.NumNodes, 100))
			Expect(info.NumActiveLeaves).To(BeNumerically("~", expInfo.NumActiveLeaves, 100))
			Expect(info.NumInactiveLeaves).To(BeNumerically("~", expInfo.NumInactiveLeaves, 100))
			Expect(info.MaxDepth).To(Equal(expInfo.MaxDepth))
			Expect(stats.R2()).To(BeNumerically("~", expR2, 0.01))
			Expect(stats.RMSE()).To(BeNumerically("~", expRMSE, 0.01))
		},

		Entry("1,000", 1000, nil, &TreeInfo{
			NumNodes:        1,
			NumActiveLeaves: 1,
			MaxDepth:        1,
		}, 0.00, 0.85),

		Entry("2,000", 2000, nil, &TreeInfo{
			NumNodes:        1071,
			NumActiveLeaves: 1070,
			MaxDepth:        2,
		}, 0.22, 0.70),

		Entry("10,000", 10000, nil, &TreeInfo{
			NumNodes:          3690,
			NumActiveLeaves:   3690,
			NumInactiveLeaves: 0,
			MaxDepth:          2,
		}, 0.21, 0.88),
	)

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "classifiers/hoeffding")
}

func runBigDataTest(model *core.Model, stats eval.Evaluator, n int, fname string, config *Config) (*TreeInfo, error) {
	stream, err := testdata.Open(fname, model)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	insts, err := stream.ReadN(n * 2)
	if err != nil {
		return nil, err
	}

	tree := New(model, config)
	for _, inst := range insts[:n] {
		tree.Train(inst)
	}
	if stats != nil {
		for _, inst := range insts[n:] {
			stats.Record(inst, tree.Predict(inst))
		}
	}

	return tree.Info(), nil
}
