package iox_test

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bsm/reason/internal/iox"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader/Writer", func() {
	var dir string
	var data = bytes.Repeat([]byte("abcd"), 1024)

	BeforeEach(func() {
		var err error
		dir, err = ioutil.TempDir("", "reason-iox-test")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(dir)).To(Succeed())
	})

	Describe("Open", func() {

		It("should open stdin", func() {
			rc, err := iox.Open("-")
			Expect(err).NotTo(HaveOccurred())
			Expect(rc.Close()).To(Succeed())
		})

		It("should open plain files", func() {
			fn := filepath.Join(dir, "plain.txt")
			Expect(ioutil.WriteFile(fn, data, 0777)).To(Succeed())

			rc, err := iox.Open(fn)
			Expect(err).NotTo(HaveOccurred())
			defer rc.Close()

			read := make([]byte, 14)
			Expect(rc.Read(read)).To(Equal(14))
			Expect(string(read)).To(Equal("abcdabcdabcdab"))
			Expect(rc.Close()).To(Succeed())
		})

		It("should open compressed files", func() {
			fn := filepath.Join(dir, "compressed.gz")
			f, err := os.Create(fn)
			Expect(err).NotTo(HaveOccurred())
			defer f.Close()
			z := gzip.NewWriter(f)
			defer z.Close()
			Expect(z.Write(data)).To(Equal(4096))
			Expect(z.Close()).To(Succeed())
			Expect(f.Close()).To(Succeed())

			rc, err := iox.Open(fn)
			Expect(err).NotTo(HaveOccurred())
			defer rc.Close()

			read := make([]byte, 14)
			Expect(rc.Read(read)).To(Equal(14))
			Expect(string(read)).To(Equal("abcdabcdabcdab"))
			Expect(rc.Close()).To(Succeed())
		})

	})

	Describe("Create", func() {

		It("should access stdout", func() {
			wc, err := iox.Create("-")
			Expect(err).NotTo(HaveOccurred())
			Expect(wc.Close()).To(Succeed())
		})

		It("should create plain files", func() {
			fn := filepath.Join(dir, "plain.txt")
			wc, err := iox.Create(fn)
			Expect(err).NotTo(HaveOccurred())
			defer wc.Close()

			Expect(wc.Write(data)).To(Equal(4096))
			Expect(wc.Close()).To(Succeed())

			info, err := os.Stat(fn)
			Expect(err).NotTo(HaveOccurred())
			Expect(info.Size()).To(Equal(int64(4096)))
		})

		It("should create compressed files", func() {
			fn := filepath.Join(dir, "compressed.gz")
			wc, err := iox.Create(fn)
			Expect(err).NotTo(HaveOccurred())
			defer wc.Close()

			Expect(wc.Write(data)).To(Equal(4096))
			Expect(wc.Close()).To(Succeed())

			info, err := os.Stat(fn)
			Expect(err).NotTo(HaveOccurred())
			Expect(info.Size()).To(BeNumerically("~", 47, 10))
		})

	})

})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "internal/iox")
}
