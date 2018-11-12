package ingest

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/bsm/reason"
)

// DataStream defines a common interface for data streams.
type DataStream interface {
	// Next returns the next example. It will return io.EOF when there are no more examples.
	Next() (reason.Example, error)
}

var errNoIntro = errors.New("reason: intro unavailable, stream already read")

// IntroReader wraps a reader and allows to separately read a pre-buffere intro chunk.
type IntroReader struct {
	r io.Reader
	t *os.File

	dir string
	lim int64

	err error
	nn  int64
}

// NewIntroReader wraps a reader.
func NewIntroReader(r io.Reader, dir string, lim int64) *IntroReader {
	return &IntroReader{r: r, dir: dir, lim: lim}
}

// Intro returns an intro reader. Can only be called if nothing has been read yet.
func (r *IntroReader) Intro() (io.ReadCloser, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.t == nil && r.nn != 0 {
		return nil, errNoIntro
	}
	if r.t == nil {
		if err := r.initBuffer(); err != nil {
			r.err = err
			return nil, err
		}
	}
	return os.Open(r.t.Name())
}

// Read implements io.Reader interface.
func (r *IntroReader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}

	n, err := r.r.Read(p)
	r.nn += int64(n)
	return n, err
}

// Close implements io.Closer interface.
func (r *IntroReader) Close() error {
	var err error
	if r.t != nil {
		err = r.t.Close()
		_ = os.Remove(r.t.Name())
	}
	return err
}

func (r *IntroReader) initBuffer() error {
	tw, err := ioutil.TempFile(r.dir, "reason")
	if err != nil {
		return err
	}

	if _, err := io.CopyN(tw, r.r, r.lim); err != nil && err != io.EOF {
		_ = tw.Close()
		_ = os.Remove(tw.Name())
		return err
	}

	if err := tw.Close(); err != nil {
		_ = tw.Close()
		_ = os.Remove(tw.Name())
		return err
	}

	tr, err := os.Open(tw.Name())
	if err != nil {
		return err
	}

	r.t = tr
	r.r = io.MultiReader(r.t, r.r)
	return nil
}
