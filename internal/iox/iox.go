package iox

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Open opens an io.ReadCloser that is either a plain file,
// a gzipped file or stdout.
func Open(name string) (io.ReadCloser, error) {
	if name == "-" {
		return ioutil.NopCloser(os.Stdin), nil
	}

	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(file)
	sig, err := buf.Peek(2)
	if err != nil {
		return nil, err
	}

	if len(sig) == 2 && sig[0] == 0x1f && sig[1] == 0x8b {
		z, err := gzip.NewReader(buf)
		if err != nil {
			return nil, err
		}
		return &zReader{Reader: z, file: file}, nil
	}

	return &bReader{Reader: buf, file: file}, nil
}

type bReader struct {
	*bufio.Reader
	file *os.File
}

// Close implements io.ReadCloser
func (b *bReader) Close() error {
	return b.file.Close()
}

type zReader struct {
	*gzip.Reader
	file *os.File
}

// Close implements io.ReadCloser
func (z *zReader) Close() error {
	err := z.Reader.Close()
	if e := z.file.Close(); e != nil {
		err = e
	}
	return err
}

// --------------------------------------------------------------------

// Create opens an io.WriteCloser that is either a plain file,
// a gzipped file or stdout.
func Create(name string) (io.WriteCloser, error) {
	if name == "-" {
		return &nopCloserWriter{Writer: os.Stdout}, nil
	}

	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(filepath.Ext(name), "z") {
		return &zWriter{
			Writer: gzip.NewWriter(file),
			file:   file,
		}, nil
	}

	return file, nil
}

type nopCloserWriter struct{ io.Writer }

// Close implements io.WriteCloser
func (*nopCloserWriter) Close() error { return nil }

type zWriter struct {
	*gzip.Writer
	file *os.File
}

// Close implements io.WriteCloser
func (z *zWriter) Close() error {
	err := z.Writer.Close()
	if e := z.file.Close(); e != nil {
		err = e
	}
	return err
}
