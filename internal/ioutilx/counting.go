package ioutilx

import (
	"io"
)

// CountingReader wraps an io.Reader
type CountingReader struct {
	N int64
	R io.Reader
}

// Read implements io.Reader
func (r *CountingReader) Read(buf []byte) (int, error) {
	n, err := r.R.Read(buf)
	r.N += int64(n)
	return n, err
}

// CountingWriter wraps an io.Writer
type CountingWriter struct {
	N int64
	W io.Writer
}

// Writer implements io.Writer
func (w *CountingWriter) Write(buf []byte) (int, error) {
	n, err := w.W.Write(buf)
	w.N += int64(n)
	return n, err
}
