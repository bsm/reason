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
	var subject, blank *util.Vector

	BeforeEach(func() {
		subject = util.NewVector()
		subject.Set(0, 2.0)
		subject.Set(2, 0.0)
		subject.Set(3, 7.0)
		subject.Set(4, 9.0)

		blank = util.NewVector()
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

	It("should observe", func() {
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

	It("should count non-zero", func() {
		Expect(subject.NNZ()).To(Equal(3))
		Expect(blank.NNZ()).To(Equal(0))
	})

	It("should reset", func() {
		subject.Reset()
		Expect(subject).To(Equal(&util.Vector{}))
	})

	It("should calculate weight sum", func() {
		Expect(subject.WeightSum()).To(Equal(18.0))
		Expect(blank.WeightSum()).To(Equal(0.0))
	})

	It("should calculate min", func() {
		i, w := subject.Min()
		Expect(i).To(Equal(0))
		Expect(w).To(Equal(2.0))

		i, w = blank.Min()
		Expect(i).To(Equal(-1))
		Expect(w).To(Equal(0.0))
	})

	It("should calculate max", func() {
		i, w := subject.Max()
		Expect(i).To(Equal(4))
		Expect(w).To(Equal(9.0))

		i, w = blank.Max()
		Expect(i).To(Equal(-1))
		Expect(w).To(Equal(0.0))
	})

	It("should calculate mean", func() {
		Expect(subject.Mean()).To(Equal(6.0))
		Expect(blank.Mean()).To(Equal(0.0))
	})

	It("should calculate variance", func() {
		Expect(subject.Variance()).To(BeNumerically("~", 13.0, 0.01))
		Expect(math.IsNaN(blank.Variance())).To(BeTrue())
	})

	It("should calculate standard deviation", func() {
		Expect(subject.StdDev()).To(BeNumerically("~", 3.61, 0.01))
		Expect(math.IsNaN(blank.StdDev())).To(BeTrue())
	})

	It("should calculate entropy", func() {
		Expect(subject.Entropy()).To(BeNumerically("~", 1.38, 0.01))
	})
})

func BenchmarkVector_NNZ(b *testing.B) {
	for _, n := range []int{100, 1000, 10000} {
		vv := util.NewVector()
		for i := 0; i < n; i += 10 {
			vv.Set(i, 1.0)
		}

		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if nnz := vv.NNZ(); nnz != n/10 {
					b.Fatalf("expected NNZ to be %d, but was %d", n/10, nnz)
				}
			}
		})
	}
}

func BenchmarkVector_Add(b *testing.B) {
	for _, n := range []int{100, 1000, 10000} {
		vv := util.NewVector()
		rnd := rand.New(rand.NewSource(10))
		b.ResetTimer()

		b.Run(fmt.Sprintf("%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				vv.Add(rnd.Intn(n), 1.0)
			}
		})
	}
}
