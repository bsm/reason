package ioutilx

import (
	"io/ioutil"
	"os"
)

// TempFile is a file that is automatically deleted on Close.
type TempFile struct{ *os.File }

// NewTempFile inits a new tempfile in dir.
func NewTempFile(dir string) (*TempFile, error) {
	f, err := ioutil.TempFile(dir, "reason")
	if err != nil {
		return nil, err
	}
	return &TempFile{File: f}, nil
}

// ReOpen closes and reopens the file for reading.
func (f *TempFile) ReOpen() error {
	if err := f.File.Close(); err != nil {
		return err
	}

	r, err := os.Open(f.Name())
	if err != nil {
		return err
	}

	f.File = r
	return nil
}

// Close closes the file and deletes it.
func (f *TempFile) Close() error {
	err := f.File.Close()
	_ = os.Remove(f.Name())
	return err
}
