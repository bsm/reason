package ingest_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bsm/reason/ingest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IntroReader", func() {
	var subject *ingest.IntroReader
	var dir string
	var src *os.File

	BeforeEach(func() {
		line := append(bytes.Repeat([]byte{'x'}, 159), '\n')

		var err error
		dir, err = ioutil.TempDir("", "reason-test")
		Expect(err).NotTo(HaveOccurred())

		srw, err := os.Create(filepath.Join(dir, "source.txt"))
		Expect(err).NotTo(HaveOccurred())
		defer srw.Close()

		for i := 0; i < 200; i++ {
			_, err := srw.Write(line)
			Expect(err).NotTo(HaveOccurred())
		}
		Expect(srw.Close()).To(Succeed())

		src, err = os.Open(filepath.Join(dir, "source.txt"))
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(subject.Close()).To(Succeed())
		Expect(src.Close()).To(Succeed())
		Expect(os.RemoveAll(dir)).To(Succeed())
	})

	It("should proxy", func() {
		subject = ingest.NewIntroReader(src, dir, 4096)
		Expect(ioutil.ReadAll(subject)).To(HaveLen(32000))

		_, err := subject.Intro()
		Expect(err).To(MatchError(`reason: intro unavailable, stream already read`))
	})

	It("should generate intro", func() {
		subject = ingest.NewIntroReader(src, dir, 4096)

		intro, err := subject.Intro()
		Expect(err).NotTo(HaveOccurred())
		defer intro.Close()

		Expect(ioutil.ReadAll(intro)).To(HaveLen(4096))
		Expect(ioutil.ReadAll(subject)).To(HaveLen(32000))
	})

	It("should allow intro to overflow content length", func() {
		subject = ingest.NewIntroReader(src, dir, 64000)

		intro, err := subject.Intro()
		Expect(err).NotTo(HaveOccurred())
		defer intro.Close()

		Expect(ioutil.ReadAll(intro)).To(HaveLen(32000))
		Expect(ioutil.ReadAll(subject)).To(HaveLen(32000))
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ingest")
}
