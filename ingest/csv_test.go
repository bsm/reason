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
			"FC1": "v1", "FC2": "v4", "FC3": "v3", "FC4": "v4", "FC5": "v5",
			"FN1": 0.036235, "FN2": 0.658867, "FN3": 0.71074, "FN4": 0.152736, "FN5": 0.159578,
			"CC": "c1", "CN": 13.9,
		}))
		Expect(stream.Next()).To(Equal(reason.MapExample{
			"FC1": "v5", "FC2": "v3", "FC3": "v3", "FC4": "v2", "FC5": "v5",
			"FN1": 0.397174, "FN2": 0.347518, "FN3": 0.294057, "FN4": 0.506484, "FN5": 0.115967,
			"CC": "c2", "CN": 11.9,
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
		Expect(model.Feature("FC1")).To(Equal(&reason.Feature{
			Name:       "FC1",
			Kind:       reason.Feature_CATEGORICAL,
			Vocabulary: []string{"v1", "v2", "v3", "v4", "v5"},
		}))
		Expect(model.Feature("FN1")).To(Equal(&reason.Feature{
			Name: "FN1",
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
			reason.NewCategoricalFeature("FC1", []string{"v1", "v2", "v3", "v4", "v5"}),
			reason.NewNumericalFeature("FN1"),
			reason.NewNumericalFeature("FN4"),
		)

		stream, err := ingest.NewCSVStream(file, model, nil)
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		Expect(stream.Next()).To(Equal(reason.MapExample{
			"FC1": "v1", "FN1": 0.036235, "FN4": 0.152736,
		}))
		Expect(stream.Next()).To(Equal(reason.MapExample{
			"FC1": "v5", "FN1": 0.397174, "FN4": 0.506484,
		}))

		Expect(stream.Close()).To(Succeed())
	})

	It("should allow to pass custom headers", func() {
		file, err := os.Open("../testdata/bigdata.csv")
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		stream, err := ingest.NewCSVStream(file, nil, &ingest.CSVOptions{
			Headers:  []string{"FC1", "", "", "", "", "FN1", "", "", "FN4"},
			SkipRows: 1,
		})
		Expect(err).NotTo(HaveOccurred())
		defer stream.Close()

		Expect(stream.Next()).To(Equal(reason.MapExample{
			"FC1": "v1", "FN1": 0.036235, "FN4": 0.152736,
		}))
		Expect(stream.Next()).To(Equal(reason.MapExample{
			"FC1": "v5", "FN1": 0.397174, "FN4": 0.506484,
		}))

		Expect(stream.Close()).To(Succeed())
	})
})
