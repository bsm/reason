package coder

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/coder")
}

func init() {
	Register((*mockType)(nil))
}

type mockType struct {
	data []int
}

func (x *mockType) EncodeTo(enc *Encoder) error   { return enc.Encode(x.data) }
func (x *mockType) DecodeFrom(dec *Decoder) error { return dec.Decode(&x.data) }
