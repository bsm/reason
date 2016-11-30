package helpers

import (
	"bytes"

	"github.com/bsm/reason/classifiers"
	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/msgpack"
	"github.com/bsm/reason/testdata"
	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("nominalCObserver", func() {
	var subject CObserver

	model := testdata.ClassificationModel()
	predictor := model.Predictor("outlook")
	target := model.Target()
	instances := testdata.ClassificationData()

	BeforeEach(func() {
		subject = NewNominalCObserver()
		for _, inst := range instances {
			subject.Observe(target.Value(inst), predictor.Value(inst), inst.GetInstanceWeight())
		}
	})

	It("should observe", func() {
		o := subject.(*nominalCObserver)
		Expect(o.PostSplit).To(HaveLen(2))
		Expect(o.ByteSize()).To(BeNumerically("~", 190, 20))
	})

	It("should encode/decode", func() {
		buf := new(bytes.Buffer)
		enc := msgpack.NewEncoder(buf)
		err := enc.Encode(subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		var out CObserver
		err = msgpack.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})

	DescribeTable("should calculate probability",
		func(tv, pv string, p float64) {
			Expect(subject.Probability(
				target.ValueOf(tv),
				predictor.ValueOf(pv),
			)).To(BeNumerically("~", p, 0.001))
		},

		// 0.333 + 0.417 + 0.250 = 1.0
		Entry("play if sunny", "yes", "sunny", 0.333),
		Entry("play if overcast", "yes", "overcast", 0.417),
		Entry("play if rainy", "yes", "rainy", 0.250),

		// 0.375 + 0.125 + 0.500 = 1.0
		Entry("don't play if sunny", "no", "sunny", 0.375),
		Entry("don't play if overcast", "no", "overcast", 0.125),
		Entry("don't play if rainy", "no", "rainy", 0.500),
	)

	It("should calculate best split", func() {
		s := subject.BestSplit(
			classifiers.InfoGainSplitCriterion{MinBranchFrac: 0.1},
			predictor,
			util.SparseVector{0: 9.0, 1: 5.0},
		)
		Expect(s.Merit()).To(BeNumerically("~", 0.247, 0.001))
		Expect(s.Range()).To(Equal(1.0))
		Expect(s.Condition()).To(BeAssignableToTypeOf(&nominalMultiwaySplitCondition{}))
		Expect(s.Condition().Predictor()).To(Equal("outlook"))

		postStats := s.PostStats()
		Expect(postStats).To(HaveLen(3))
		Expect(postStats[0].State()).To(ConsistOf(core.Prediction{
			{AttributeValue: 0, Votes: 2},
			{AttributeValue: 1, Votes: 3},
		}))
		Expect(postStats[1].State()).To(ConsistOf(core.Prediction{
			{AttributeValue: 0, Votes: 4},
		}))
		Expect(postStats[2].State()).To(ConsistOf(core.Prediction{
			{AttributeValue: 0, Votes: 3},
			{AttributeValue: 1, Votes: 2},
		}))
	})

	It("should require at least two observed values for best split", func() {
		inst := instances[0]

		o := NewNominalCObserver()
		o.Observe(target.Value(inst), predictor.Value(inst), inst.GetInstanceWeight())
		o.Observe(target.Value(inst), predictor.Value(inst), inst.GetInstanceWeight())
		o.Observe(target.Value(inst), predictor.Value(inst), inst.GetInstanceWeight())

		Expect(o.BestSplit(
			classifiers.InfoGainSplitCriterion{MinBranchFrac: 0.1},
			predictor,
			util.SparseVector{0: 3.0},
		)).To(BeNil())
	})

})

var _ = Describe("gaussianCObserver", func() {
	var subject CObserver

	predictor := &core.Attribute{Name: "len", Kind: core.AttributeKindNumeric}
	target := &core.Attribute{Name: "class", Kind: core.AttributeKindNominal}
	instances := []core.Instance{
		core.MapInstance{"len": 1.4, "class": "a"},
		core.MapInstance{"len": 1.3, "class": "a"},
		core.MapInstance{"len": 1.5, "class": "a"},
		core.MapInstance{"len": 4.1, "class": "b"},
		core.MapInstance{"len": 3.7, "class": "b"},
		core.MapInstance{"len": 4.9, "class": "b"},
		core.MapInstance{"len": 4.0, "class": "b"},
		core.MapInstance{"len": 3.3, "class": "b"},
		core.MapInstance{"len": 6.3, "class": "c"},
		core.MapInstance{"len": 5.8, "class": "c"},
		core.MapInstance{"len": 5.1, "class": "c"},
		core.MapInstance{"len": 5.3, "class": "c"},
	}

	BeforeEach(func() {
		subject = NewNumericCObserver(4)
		for _, inst := range instances {
			tv := target.Value(inst)
			pv := predictor.Value(inst)
			subject.Observe(tv, pv, inst.GetInstanceWeight())
		}
	})

	It("should observe", func() {
		o := subject.(*gaussianCObserver)
		Expect(o.Range.SplitPoints(4)).To(Equal([]float64{2.3, 3.3, 4.3, 5.3}))
		Expect(o.PostSplit).To(HaveLen(3))
		Expect(o.ByteSize()).To(BeNumerically("~", 1140, 20))
	})

	It("should not calculate probability", func() {
		Expect(subject.Probability(
			target.ValueOf("b"),
			predictor.ValueOf(4.5),
		)).To(BeNumerically("~", 0.47, 0.01))

		Expect(subject.Probability(
			target.ValueOf("a"),
			predictor.ValueOf(4.5),
		)).To(BeNumerically("~", 0.00, 0.01))

		Expect(subject.Probability(
			target.ValueOf("a"),
			predictor.ValueOf(1.7),
		)).To(BeNumerically("~", 0.04, 0.01))
	})

	It("should calculate best split", func() {
		s := subject.BestSplit(
			classifiers.InfoGainSplitCriterion{MinBranchFrac: 0.1},
			predictor,
			util.SparseVector{0: 3.0, 1: 5.0, 2: 4.0},
		)
		Expect(s.Merit()).To(BeNumerically("~", 0.811, 0.001))
		Expect(s.Range()).To(BeNumerically("~", 1.585, 0.001))
		Expect(s.Condition()).To(BeAssignableToTypeOf(&numericBinarySplitCondition{}))
		Expect(s.Condition().Predictor()).To(Equal("len"))
		Expect(s.Condition().(*numericBinarySplitCondition).SplitValue).To(Equal(2.30))

		postStats := s.PostStats()
		Expect(postStats).To(HaveLen(2))
		Expect(postStats[0].State()).To(ConsistOf(core.Prediction{
			{AttributeValue: 0, Votes: 3},
		}))
		Expect(postStats[1].State()).To(ConsistOf(core.Prediction{
			{AttributeValue: 1, Votes: 5},
			{AttributeValue: 2, Votes: 4},
		}))
	})

	It("should encode/decode", func() {
		buf := new(bytes.Buffer)
		enc := msgpack.NewEncoder(buf)
		err := enc.Encode(subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		var out CObserver
		err = msgpack.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})
})

var _ = Describe("nominalRObserver", func() {
	var subject RObserver
	var preSplit *util.NumSeries

	model := testdata.RegressionModel()
	predictor := model.Predictor("outlook")
	target := model.Target()
	instances := testdata.RegressionData()

	BeforeEach(func() {
		subject = NewNominalRObserver()
		preSplit = new(util.NumSeries)

		for _, inst := range instances {
			tv := target.Value(inst)
			pv := predictor.Value(inst)
			subject.Observe(tv, pv, inst.GetInstanceWeight())
			preSplit.Append(tv.Value(), inst.GetInstanceWeight())
		}
	})

	It("should observe", func() {
		o := subject.(*nominalRObserver)
		Expect(o.PostSplit).To(HaveLen(3))
		a, b, c := o.PostSplit[0], o.PostSplit[1], o.PostSplit[2]
		Expect(a.StdDev()).To(BeNumerically("~", 7.78, 0.01))
		Expect(b.StdDev()).To(BeNumerically("~", 3.49, 0.01))
		Expect(c.StdDev()).To(BeNumerically("~", 10.87, 0.01))
		Expect(o.ByteSize()).To(BeNumerically("~", 1050, 20))
	})

	It("should calculate best split", func() {
		s := subject.BestSplit(
			classifiers.VarReductionSplitCriterion{},
			predictor,
			preSplit,
		)
		Expect(s.Merit()).To(BeNumerically("~", 19.572, 0.001))
		Expect(s.Range()).To(Equal(1.0))
		Expect(s.Condition()).To(BeAssignableToTypeOf(&nominalMultiwaySplitCondition{}))
		Expect(s.Condition().Predictor()).To(Equal("outlook"))
	})

	It("should require at least two observed values for best split", func() {
		inst := instances[0]

		o := NewNominalRObserver()
		o.Observe(target.Value(inst), predictor.Value(inst), inst.GetInstanceWeight())
		o.Observe(target.Value(inst), predictor.Value(inst), inst.GetInstanceWeight())
		o.Observe(target.Value(inst), predictor.Value(inst), inst.GetInstanceWeight())

		Expect(o.BestSplit(
			classifiers.VarReductionSplitCriterion{},
			predictor,
			preSplit,
		)).To(BeNil())
	})

	It("should encode/decode", func() {
		buf := new(bytes.Buffer)
		enc := msgpack.NewEncoder(buf)
		err := enc.Encode(subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		var out RObserver
		err = msgpack.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})
})

var _ = Describe("gaussianRObserver", func() {
	var subject RObserver
	var preSplit *util.NumSeries

	predictor := &core.Attribute{Name: "area", Kind: core.AttributeKindNumeric}
	target := &core.Attribute{Name: "price", Kind: core.AttributeKindNumeric}
	instances := []core.MapInstance{
		{"area": 1.1, "price": 4.5},
		{"area": 1.2, "price": 4.5},
		{"area": 1.5, "price": 5.0},
		{"area": 0.9, "price": 3.8},
		{"area": 1.3, "price": 5.8},
		{"area": 1.5, "price": 5.6},
		{"area": 0.8, "price": 3.2},
		{"area": 2.6, "price": 8.2},
		{"area": 1.0, "price": 3.9},
		{"area": 1.6, "price": 5.1},
		{"area": 1.8, "price": 8.7},
		{"area": 1.6, "price": 6.0},
	}

	BeforeEach(func() {
		subject = NewNumericRObserver(5)
		preSplit = new(util.NumSeries)

		for _, inst := range instances {
			tv := target.Value(inst)
			pv := predictor.Value(inst)
			subject.Observe(tv, pv, inst.GetInstanceWeight())
			preSplit.Append(tv.Value(), inst.GetInstanceWeight())
		}
	})

	It("should observe", func() {
		o := subject.(*gaussianRObserver)
		Expect(o.Range.SplitPoints(5)).To(Equal([]float64{1.1, 1.4, 1.7, 2, 2.3}))
		Expect(o.Observations).To(HaveLen(12))
		Expect(o.ByteSize()).To(Equal(368))
	})

	It("should calculate best split", func() {
		s := subject.BestSplit(
			classifiers.VarReductionSplitCriterion{},
			predictor,
			preSplit,
		)
		Expect(s.Merit()).To(BeNumerically("~", 1.911, 0.001))
		Expect(s.Range()).To(Equal(1.0))
		Expect(s.Condition()).To(BeAssignableToTypeOf(&numericBinarySplitCondition{}))
		Expect(s.Condition().Predictor()).To(Equal("area"))
		Expect(s.Condition().(*numericBinarySplitCondition).SplitValue).To(Equal(1.7))
	})

	It("should encode/decode", func() {
		buf := new(bytes.Buffer)
		enc := msgpack.NewEncoder(buf)
		err := enc.Encode(subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		var out RObserver
		err = msgpack.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})
})
