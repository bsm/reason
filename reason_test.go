package reason_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "reason")
}

type (
	customNumeric int
	customString  string
)

func intPtr(n int) *int { return &n }
