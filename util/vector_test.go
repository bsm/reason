package util

import (
	"bytes"
	"math"

	"github.com/bsm/reason/internal/coder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SparseVector", func() {
	var subject SparseVector

	BeforeEach(func() {
		subject = NewSparseVector()
		subject.Set(0, 2.0)
		subject.Set(3, 7.0)
		subject.Set(4, 9.0)
	})

	It("should set", func() {
		Expect(subject).To(Equal(SparseVector{0: 2, 3: 7, 4: 9}))
		subject.Set(6, 8.0)
		subject.Set(-2, 4.0)
		Expect(subject).To(Equal(SparseVector{0: 2, 3: 7, 4: 9, 6: 8}))
	})

	It("should get", func() {
		Expect(subject.Get(0)).To(Equal(2.0))
		Expect(subject.Get(1)).To(Equal(0.0))
		Expect(subject.Get(-1)).To(Equal(0.0))
	})

	It("should incr", func() {
		subject.Incr(0, 1.0)
		subject.Incr(1, 7.0)
		subject.Incr(-1, 6.0)
		Expect(subject).To(Equal(SparseVector{0: 3, 1: 7, 3: 7, 4: 9}))
	})

	It("should normalize", func() {
		subject.Normalize()
		Expect(subject).To(Equal(SparseVector{0: (1.0 / 9.0), 3: (7.0 / 18.0), 4: 0.5}))
	})

	It("should count", func() {
		Expect(subject.Count()).To(Equal(3))
		Expect(NewSparseVector().Count()).To(Equal(0))
	})

	It("should clear", func() {
		subject.Clear()
		Expect(subject.Count()).To(Equal(0))
	})

	It("should calculate sum", func() {
		Expect(subject.Sum()).To(Equal(18.0))
		Expect(NewSparseVector().Sum()).To(Equal(0.0))
	})

	It("should calculate min", func() {
		Expect(subject.Min()).To(Equal(2.0))
		Expect(NewSparseVector().Min()).To(Equal(math.MaxFloat64))
	})

	It("should calculate max", func() {
		Expect(subject.Max()).To(Equal(9.0))
		Expect(NewSparseVector().Max()).To(Equal(-math.MaxFloat64))
	})

	It("should calculate mean", func() {
		Expect(subject.Mean()).To(Equal(6.0))
		Expect(NewSparseVector().Mean()).To(Equal(0.0))
	})

	It("should iterate", func() {
		indicies := make([]int, 0)
		values := make([]float64, 0)
		subject.ForEach(func(i int, v float64) {
			indicies = append(indicies, i)
			values = append(values, v)
		})
		Expect(indicies).To(ConsistOf(0, 3, 4))
		Expect(values).To(ConsistOf(2.0, 7.0, 9.0))

		values = values[:0]
		subject.ForEachValue(func(v float64) {
			values = append(values, v)
		})
		Expect(values).To(ConsistOf(2.0, 7.0, 9.0))
	})

	It("should calculate variances", func() {
		Expect(subject.Variance()).To(BeNumerically("~", 8.67, 0.01))
		Expect(NewSparseVector().Variance()).To(Equal(0.0))
		Expect(subject.SampleVariance()).To(BeNumerically("~", 13.0, 0.01))
		Expect(NewSparseVector().SampleVariance()).To(Equal(0.0))
	})

	It("should calculate standard deviations", func() {
		Expect(subject.StdDev()).To(BeNumerically("~", 2.94, 0.01))
		Expect(subject.SampleStdDev()).To(BeNumerically("~", 3.61, 0.01))
	})

	It("should calculate entropy", func() {
		Expect(subject.Entropy()).To(BeNumerically("~", 1.38, 0.01))
	})

	It("should encode/decode", func() {
		buf := new(bytes.Buffer)
		enc := coder.NewEncoder(buf)
		err := enc.Encode(subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		var out SparseVector
		err = coder.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})

})

var _ = Describe("DenseVector", func() {
	var subject *DenseVector

	BeforeEach(func() {
		subject = NewDenseVector()
		subject.Set(0, 2.0)
		subject.Set(3, 7.0)
		subject.Set(4, 9.0)
	})

	It("should set", func() {
		Expect(subject).To(Equal(&DenseVector{vv: []float64{2, 0, 0, 7, 9}}))
		subject.Set(6, 8.0)
		subject.Set(-2, 4.0)
		Expect(subject).To(Equal(&DenseVector{vv: []float64{2, 0, 0, 7, 9, 0, 8}}))
	})

	It("should get", func() {
		Expect(subject.Get(0)).To(Equal(2.0))
		Expect(subject.Get(1)).To(Equal(0.0))
		Expect(subject.Get(-1)).To(Equal(0.0))
	})

	It("should incr", func() {
		subject.Incr(0, 1.0)
		subject.Incr(1, 7.0)
		subject.Incr(-1, 6.0)
		Expect(subject).To(Equal(&DenseVector{vv: []float64{3, 7, 0, 7, 9}}))
	})

	It("should normalize", func() {
		subject.Normalize()
		Expect(subject).To(Equal(&DenseVector{vv: []float64{(1.0 / 9.0), 0, 0, (7.0 / 18.0), 0.5}}))
	})

	It("should count", func() {
		Expect(subject.Count()).To(Equal(3))
		Expect(NewDenseVector().Count()).To(Equal(0))
	})

	It("should clear", func() {
		subject.Clear()
		Expect(subject.Count()).To(Equal(0))
	})

	It("should calculate sum", func() {
		Expect(subject.Sum()).To(Equal(18.0))
		Expect(NewDenseVector().Sum()).To(Equal(0.0))
	})

	It("should calculate min", func() {
		Expect(subject.Min()).To(Equal(2.0))
		Expect(NewDenseVector().Min()).To(Equal(math.MaxFloat64))
	})

	It("should calculate max", func() {
		Expect(subject.Max()).To(Equal(9.0))
		Expect(NewDenseVector().Max()).To(Equal(-math.MaxFloat64))
	})

	It("should calculate mean", func() {
		Expect(subject.Mean()).To(Equal(6.0))
		Expect(NewDenseVector().Mean()).To(Equal(0.0))
	})

	It("should iterate", func() {
		indicies := make([]int, 0)
		values := make([]float64, 0)
		subject.ForEach(func(i int, v float64) {
			indicies = append(indicies, i)
			values = append(values, v)
		})
		Expect(indicies).To(ConsistOf(0, 3, 4))
		Expect(values).To(ConsistOf(2.0, 7.0, 9.0))

		values = values[:0]
		subject.ForEachValue(func(v float64) {
			values = append(values, v)
		})
		Expect(values).To(ConsistOf(2.0, 7.0, 9.0))
	})

	It("should calculate variances", func() {
		Expect(subject.Variance()).To(BeNumerically("~", 8.67, 0.01))
		Expect(NewSparseVector().Variance()).To(Equal(0.0))
		Expect(subject.SampleVariance()).To(BeNumerically("~", 13.0, 0.01))
		Expect(NewSparseVector().SampleVariance()).To(Equal(0.0))
	})

	It("should calculate standard deviations", func() {
		Expect(subject.StdDev()).To(BeNumerically("~", 2.94, 0.01))
		Expect(subject.SampleStdDev()).To(BeNumerically("~", 3.61, 0.01))
	})

	It("should calculate entropy", func() {
		Expect(subject.Entropy()).To(BeNumerically("~", 1.38, 0.01))
	})

	It("should encode/decode", func() {
		buf := new(bytes.Buffer)
		enc := coder.NewEncoder(buf)
		err := enc.Encode(subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		out := new(DenseVector)
		err = coder.NewDecoder(buf).Decode(out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})

})

var _ = Describe("VectorDistribution", func() {
	var subject VectorDistribution

	BeforeEach(func() {
		subject = VectorDistribution{
			0: SparseVector{1: 2, 2: 3},
			2: SparseVector{2: 7, 9: 4},
			3: NewDenseVectorFromSlice(0, 0, 0, 4),
		}
	})

	It("should get", func() {
		Expect(subject.Get(0)).To(Equal(SparseVector{1: 2, 2: 3}))
		Expect(subject.Get(1)).To(BeNil())
		Expect(subject.Get(2)).To(Equal(SparseVector{2: 7, 9: 4}))
		Expect(subject.Get(3)).To(Equal(&DenseVector{vv: []float64{0, 0, 0, 4}}))
		Expect(subject.Get(4)).To(BeNil())
		Expect(subject.Get(-1)).To(BeNil())
	})

	It("should count", func() {
		Expect(subject.NumPredicates()).To(Equal(3))
		Expect(subject.NumTargets()).To(Equal(2))
	})

	It("should return weights", func() {
		Expect(subject.Weights()).To(Equal(map[int]float64{0: 5, 2: 11, 3: 4}))
		Expect(subject.TargetWeights()).To(Equal(map[int]float64{1: 2, 2: 10, 3: 4, 9: 4}))
	})

	It("should increment", func() {
		m := VectorDistribution{}
		m.Incr(2, 3, 4)
		Expect(m).To(Equal(VectorDistribution{
			2: NewDenseVectorFromSlice(0, 0, 0, 4),
		}))

		m.Incr(1, 2, 3)
		m.Incr(1, 2, 4)
		Expect(m).To(Equal(VectorDistribution{
			1: NewDenseVectorFromSlice(0, 0, 7),
			2: NewDenseVectorFromSlice(0, 0, 0, 4),
		}))
	})

	It("should encode/decode", func() {
		buf := new(bytes.Buffer)
		enc := coder.NewEncoder(buf)
		err := enc.Encode(subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(enc.Close()).NotTo(HaveOccurred())

		out := make(VectorDistribution)
		err = coder.NewDecoder(buf).Decode(&out)
		Expect(err).NotTo(HaveOccurred())
		Expect(out).To(Equal(subject))
	})

})
