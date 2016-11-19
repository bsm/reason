package msgpack

import (
	"sort"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	Register((*mockSliceType)(nil))
	Register((mockSetType)(nil))
	Register(mockNone{})
}

type mockInterface interface {
	All() []int
}

type mockSliceType struct{ data []int }

func (x *mockSliceType) All() []int                    { return x.data }
func (x *mockSliceType) EncodeTo(enc *Encoder) error   { return enc.Encode(x.data) }
func (x *mockSliceType) DecodeFrom(dec *Decoder) error { return dec.Decode(&x.data) }

type mockSetType map[int]mockNone

func (x mockSetType) All() []int {
	nn := make([]int, 0, len(x))
	for n := range x {
		nn = append(nn, n)
	}
	sort.Ints(nn)
	return nn
}

type mockNone struct{}

func (mockNone) EncodeTo(enc *Encoder) error   { return nil }
func (mockNone) DecodeFrom(dec *Decoder) error { return nil }

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/msgpack")
}
