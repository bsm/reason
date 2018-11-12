package ingest

import (
	"io"

	"github.com/bsm/reason"
)

// CSVOptions contain optional settings for CSV sources.
type CSVOptions struct {
	// Row separator character. Default: ','
	Separator rune

	// Column names. To ignore a column, set it's name to "".
	// Default: autodetected from the first row.
	Headers []string

	// Number of rows to skip. Allows to e.g. skip the header row
	// while providing custom header names.
	SkipRows int
}

func (o *CSVOptions) norm() {
	if o.Separator == 0 {
		o.Separator = ','
	}
}

// DetectModelCSV detects a model from a Reader.
func DetectModelCSV(r io.Reader, opt *CSVOptions) (*reason.Model, error) {
	return nil, nil
}
