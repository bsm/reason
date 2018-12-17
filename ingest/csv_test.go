package ingest_test

import (
	"os"

	"github.com/bsm/reason"
	"github.com/bsm/reason/ingest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ ingest.DataStream = (*ingest.CSVStream)(nil)

var _ = Describe("CSVStream", func() {

	It("should read CSVs", func() {
		file, err := os.Open("../testdata/bigdata.csv")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		stream, err := ingest.NewCSVStream(file, nil, nil)
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		Expect(stream.Next()).To(Equal(reason.MapExample{
			"c1": "v1", "c2": "v4", "c3": "v3", "c4": "v4", "c5": "v5",
			"n1": 0.036235, "n2": 0.658867, "n3": 0.71074, "n4": 0.152736, "n5": 0.159578,
			"tc": "c1", "tn": 13.9,
		}))
		Expect(stream.Next()).To(Equal(reason.MapExample{
			"c1": "v5", "c2": "v3", "c3": "v3", "c4": "v2", "c5": "v5",
			"n1": 0.397174, "n2": 0.347518, "n3": 0.294057, "n4": 0.506484, "n5": 0.115967,
			"tc": "c2", "tn": 11.9,
		}))

		Expect(stream.Close()).To(Succeed())
	})

	It("should auto-detect models", func() {
		file, err := os.Open("../testdata/bigdata.csv")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		stream, err := ingest.NewCSVStream(file, nil, nil)
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		model := stream.Model()
		Expect(model.Features).To(HaveLen(12))
		Expect(model.Feature("c1")).To(Equal(&reason.Feature{
			Name:       "c1",
			Kind:       reason.Feature_CATEGORICAL,
			Vocabulary: []string{"v1", "v2", "v3", "v4", "v5"},
		}))
		Expect(model.Feature("n1")).To(Equal(&reason.Feature{
			Name: "n1",
			Kind: reason.Feature_NUMERICAL,
		}))

		count := 0
		for {
			if _, err := stream.Next(); err != nil {
				break
			}
			count++
		}
		Expect(count).To(Equal(100000))
		Expect(stream.Close()).To(Succeed())
	})

	It("should allow to pass models", func() {
		file, err := os.Open("../testdata/bigdata.csv")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		model := reason.NewModel(
			reason.NewCategoricalFeature("c1", []string{"v1", "v2", "v3", "v4", "v5"}),
			reason.NewNumericalFeature("n1"),
			reason.NewNumericalFeature("n4"),
		)

		stream, err := ingest.NewCSVStream(file, model, nil)
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		Expect(stream.Next()).To(Equal(reason.MapExample{
			"c1": "v1", "n1": 0.036235, "n4": 0.152736,
		}))
		Expect(stream.Next()).To(Equal(reason.MapExample{
			"c1": "v5", "n1": 0.397174, "n4": 0.506484,
		}))

		Expect(stream.Close()).To(Succeed())
	})

	It("should allow to pass custom headers", func() {
		file, err := os.Open("../testdata/bigdata.csv")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		stream, err := ingest.NewCSVStream(file, nil, &ingest.CSVOptions{
			Headers:  []string{"c1", "", "", "", "", "n1", "", "", "n4"},
			SkipRows: 1,
		})
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		Expect(stream.Next()).To(Equal(reason.MapExample{
			"c1": "v1", "n1": 0.036235, "n4": 0.152736,
		}))
		Expect(stream.Next()).To(Equal(reason.MapExample{
			"c1": "v5", "n1": 0.397174, "n4": 0.506484,
		}))

		Expect(stream.Close()).To(Succeed())
	})
})
