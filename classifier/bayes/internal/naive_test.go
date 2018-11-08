package internal_test

import (
	"bytes"

	"github.com/bsm/reason/util"

	"github.com/bsm/reason/classifier/bayes/internal"
	"github.com/bsm/reason/testdata"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NaiveBayes", func() {
	var subject *internal.NaiveBayes

	model := testdata.SimpleModel

	BeforeEach(func() {
		target := model.Feature("play")
		subject = internal.New(model, target.Name)
		for _, x := range testdata.SimpleDataSet {
			subject.ObserveWeight(x, target.Category(x), 1.0)
		}
	})

	It("should marshal to writer", func() {
		buf := new(bytes.Buffer)
		Expect(subject.WriteTo(buf)).To(Equal(int64(716)))

		t := new(internal.NaiveBayes)
		Expect(proto.Unmarshal(buf.Bytes(), t)).To(Succeed())
		Expect(t).To(Equal(subject))
	})

	It("should unmarshal from reader", func() {
		data, err := proto.Marshal(subject)
		Expect(err).NotTo(HaveOccurred())

		t := new(internal.NaiveBayes)
		Expect(t.ReadFrom(bytes.NewReader(data))).To(Equal(int64(len(data))))
		Expect(t).To(Equal(subject))
	})

	It("should observe", func() {
		Expect(subject.TargetStats.Data).To(Equal([]float64{9, 5}))
		Expect(subject.FeatureStats).To(HaveLen(6))
		Expect(subject.FeatureStats["outlook"].GetCat().Dist.Data).To(Equal([]float64{
			2, 3,
			4, 0,
			3, 2,
		}))
		Expect(subject.FeatureStats["humidex"].GetNum().Dist.Data).To(Equal([]util.NumStream{
			{Weight: 9, Sum: 416, SumSquares: 20242, Min: 35, Max: 62},
			{Weight: 5, Sum: 287, SumSquares: 16751, Min: 43, Max: 64},
		}))
	})
})
