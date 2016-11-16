package helpers

import (
	"bytes"
	"encoding/gob"

	"github.com/bsm/reason/classifiers"
	"github.com/bsm/reason/testdata"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ObservationStats", func() {
	var subject ObservationStats

	Describe("classification", func() {
		model := testdata.ClassificationModel()
		target := model.Target()
		instances := testdata.ClassificationData()

		BeforeEach(func() {
			subject = NewObservationStats(false)
			for _, inst := range instances {
				subject.UpdatePreSplit(target.Value(inst), inst.GetInstanceWeight())
			}
		})

		It("should estimate byte size", func() {
			Expect(subject.ByteSize()).To(BeNumerically("~", 80, 20))
		})

		It("should check sufficiency", func() {
			Expect(subject.IsSufficient()).To(BeTrue())

			blank := NewObservationStats(false)
			Expect(blank.IsSufficient()).To(BeFalse())
		})

		It("should return total weight", func() {
			Expect(subject.TotalWeight()).To(Equal(14.0))
		})

		It("should create new observers", func() {
			Expect(subject.NewObserver(true)).To(BeAssignableToTypeOf(&nominalCObserver{}))
			Expect(subject.NewObserver(false)).To(BeAssignableToTypeOf(&gaussianCObserver{}))
		})

		It("should return state", func() {
			state := subject.State()
			Expect(state).To(HaveLen(2))
			Expect(state.Top().Value.Index()).To(Equal(0))
			Expect(state.Top().Votes).To(Equal(9.0))
		})

		It("should calculate best-splits", func() {
			predictor := model.Predictor("outlook")
			obs := subject.NewObserver(true)
			for _, inst := range instances {
				obs.Observe(target.Value(inst), predictor.Value(inst), inst.GetInstanceWeight())
			}

			split := subject.BestSplit(classifiers.InfoGainSplitCriterion{}, obs, predictor)
			Expect(split.Merit()).To(BeNumerically("~", 0.247, 0.001))
		})

		It("should marshal/unmarshal", func() {
			buf := new(bytes.Buffer)
			err := gob.NewEncoder(buf).Encode(subject)
			Expect(err).NotTo(HaveOccurred())

			var out ObservationStats
			err = gob.NewDecoder(buf).Decode(out)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(Equal(subject))
		})

	})

	Describe("regression", func() {
		model := testdata.RegressionModel()
		target := model.Target()
		instances := testdata.RegressionData()

		BeforeEach(func() {
			subject = NewObservationStats(true)
			for _, inst := range instances {
				subject.UpdatePreSplit(target.Value(inst), inst.GetInstanceWeight())
			}
		})

		It("should check sufficiency", func() {
			Expect(subject.IsSufficient()).To(BeTrue())

			blank := NewObservationStats(true)
			Expect(blank.IsSufficient()).To(BeFalse())
		})

		It("should estimate byte size", func() {
			Expect(subject.ByteSize()).To(Equal(40))
		})

		It("should return total weight", func() {
			Expect(subject.TotalWeight()).To(Equal(14.0))
		})

		It("should create new observers", func() {
			Expect(subject.NewObserver(true)).To(BeAssignableToTypeOf(&nominalRObserver{}))
			Expect(subject.NewObserver(false)).To(BeAssignableToTypeOf(&gaussianRObserver{}))
		})

		It("should return state", func() {
			state := subject.State()
			Expect(state).To(HaveLen(1))
			Expect(state.Top().Value.Value()).To(BeNumerically("~", 39.7, 0.1))
			Expect(state.Top().Votes).To(Equal(14.0))
		})

		It("should calculate best-splits", func() {
			predictor := model.Predictor("outlook")
			obs := subject.NewObserver(true)
			for _, inst := range instances {
				obs.Observe(target.Value(inst), predictor.Value(inst), inst.GetInstanceWeight())
			}

			split := subject.BestSplit(classifiers.VarReductionSplitCriterion{}, obs, predictor)
			Expect(split.Merit()).To(BeNumerically("~", 19.572, 0.001))
		})

	})
})
