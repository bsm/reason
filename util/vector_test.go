package util_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/bsm/reason/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vector", func() {
	var subject *util.Vector

	BeforeEach(func() {
		subject = util.NewVector()
	})

	It("should convert between sparse & dense", func() {
		Expect(subject.IsSparse()).Should(BeFalse())
		subject.MakeSparse()
		Expect(subject.IsSparse()).Should(BeTrue())
		subject.MakeSparse()
		Expect(subject.IsSparse()).Should(BeTrue())
		subject.MakeDense()
		Expect(subject.IsSparse()).Should(BeFalse())
		subject.MakeDense()
		Expect(subject.IsSparse()).Should(BeFalse())
	})

	ItBehavesLikeAVector := func() {
		BeforeEach(func() {
			subject.Set(0, 2.0)
			subject.Set(2, 0.0)
			subject.Set(3, 7.0)
			subject.Set(4, 9.0)
		})

		It("should at", func() {
			Expect(subject.At(0)).To(Equal(2.0))
			Expect(subject.At(1)).To(Equal(0.0))
			Expect(subject.At(2)).To(Equal(0.0))
			Expect(subject.At(3)).To(Equal(7.0))
			Expect(subject.At(4)).To(Equal(9.0))
			Expect(subject.At(-1)).To(Equal(0.0))
		})

		It("should set", func() {
			Expect(subject.At(6)).To(Equal(0.0))

			subject.Set(6, 8.0)
			subject.Set(-2, 4.0)
			Expect(subject.At(6)).To(Equal(8.0))
		})

		It("should add", func() {
			Expect(subject.At(0)).To(Equal(2.0))
			Expect(subject.At(1)).To(Equal(0.0))

			subject.Add(0, 1.0)
			subject.Add(1, 7.0)
			subject.Add(-1, 6.0)
			Expect(subject.At(0)).To(Equal(3.0))
			Expect(subject.At(1)).To(Equal(7.0))
		})

		It("should normalize", func() {
			subject.Normalize()
			Expect(subject.At(0)).To(Equal(2.0 / 18))
			Expect(subject.At(3)).To(Equal(7.0 / 18))
			Expect(subject.At(4)).To(Equal(9.0 / 18))
		})

		It("should count", func() {
			Expect(subject.Len()).To(Equal(3))
		})

		It("should reset", func() {
			subject.Reset()
			Expect(subject).To(Equal(&util.Vector{}))
		})

		It("should calculate total weight", func() {
			Expect(subject.Weight()).To(Equal(18.0))
		})

		It("should calculate min", func() {
			i, w := subject.Min()
			Expect(i).To(Equal(0))
			Expect(w).To(Equal(2.0))
		})

		It("should calculate max", func() {
			i, w := subject.Max()
			Expect(i).To(Equal(4))
			Expect(w).To(Equal(9.0))
		})

		It("should calculate mean", func() {
			Expect(subject.Mean()).To(Equal(6.0))
		})

		It("should iterate", func() {
			indicies := make([]int, 0)
			values := make([]float64, 0)
			subject.ForEach(func(i int, v float64) bool {
				indicies = append(indicies, i)
				values = append(values, v)
				return true
			})
			Expect(indicies).To(ConsistOf(0, 3, 4))
			Expect(values).To(ConsistOf(2.0, 7.0, 9.0))

			values = values[:0]
			subject.ForEachValue(func(v float64) bool {
				values = append(values, v)
				return true
			})
			Expect(values).To(ConsistOf(2.0, 7.0, 9.0))
		})

		It("should calculate variances", func() {
			Expect(subject.Variance()).To(BeNumerically("~", 13.0, 0.01))
		})

		It("should calculate standard deviations", func() {
			Expect(subject.StdDev()).To(BeNumerically("~", 3.61, 0.01))
		})

		It("should calculate entropy", func() {
			Expect(subject.Entropy()).To(BeNumerically("~", 1.38, 0.01))
		})
	}

	Describe("Dense", func() {
		ItBehavesLikeAVector()
	})

	Describe("Sparse", func() {
		BeforeEach(func() {
			subject.MakeSparse()
		})
		ItBehavesLikeAVector()
	})

	Describe("Blank", func() {
		It("should count", func() {
			Expect(subject.Len()).To(Equal(0))
		})

		It("should calculate total weight", func() {
			Expect(subject.Weight()).To(Equal(0.0))
		})

		It("should calculate min", func() {
			i, w := subject.Min()
			Expect(i).To(Equal(-1))
			Expect(w).To(Equal(0.0))
		})

		It("should calculate max", func() {
			i, w := subject.Max()
			Expect(i).To(Equal(-1))
			Expect(w).To(Equal(0.0))
		})

		It("should calculate mean", func() {
			Expect(subject.Mean()).To(Equal(0.0))
		})

		It("should calculate variances", func() {
			Expect(math.IsNaN(subject.Variance())).To(BeTrue())
		})

		It("should calculate standard deviations", func() {
			Expect(math.IsNaN(subject.StdDev())).To(BeTrue())
		})
	})
})

func BenchmarkVector_ForEach(b *testing.B) {
	for _, t := range []struct {
		N, MaxIndex int
		Sparse      bool
	}{
		{10, 100, false},
		{10, 100, true},
		{100, 100, false},
		{100, 100, true},

		{10, 1000, false},
		{10, 1000, true},
		{100, 1000, false},
		{100, 1000, true},
		{200, 1000, false},
		{200, 1000, true},
		{500, 1000, false},
		{500, 1000, true},
		{1000, 1000, false},
		{1000, 1000, true},

		{10, 10000, false},
		{10, 10000, true},
		{100, 10000, false},
		{100, 10000, true},
		{500, 10000, false},
		{500, 10000, true},
		{1000, 10000, false},
		{1000, 10000, true},
		{2000, 10000, false},
		{2000, 10000, true},
	} {
		vv := util.NewVector()
		label := "dense"
		if t.Sparse {
			vv.MakeSparse()
			label = "sparse"
		}

		for i := 0; i < t.N; i++ {
			offset := int(float64(i) * float64(t.MaxIndex) / float64(t.N))
			vv.Set(t.MaxIndex-offset, 1.0)
		}

		b.Run(fmt.Sprintf("%d/%d %s", t.N, t.MaxIndex, label), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				n, m := 0, 0
				vv.ForEach(func(i int, w float64) bool {
					n++
					if i > m {
						m = i
					}
					return true
				})
				if n != t.N {
					b.Fatalf("expected N to be %d, but was %d", t.N, n)
				}
				if m != t.MaxIndex {
					b.Fatalf("expected MaxIndex to be %d, but was %d", t.MaxIndex, m)
				}
			}
		})
	}
}

func BenchmarkVector_Add(b *testing.B) {
	for _, t := range []struct {
		N      int
		Sparse bool
	}{
		{1000, false},
		{1000, true},

		{10000, false},
		{10000, true},
	} {
		vv := util.NewVector()
		label := "dense"
		if t.Sparse {
			vv.MakeSparse()
			label = "sparse"
		}

		rnd := rand.New(rand.NewSource(10))
		b.ResetTimer()

		b.Run(fmt.Sprintf("%d %s", t.N, label), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				vv.Add(rnd.Intn(t.N), 1.0)
			}
		})
	}
}
