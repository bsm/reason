package ingest

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/bsm/reason"
	"github.com/bsm/reason/internal/ioutilx"
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

	// AutoDetectBufferSize is the maximum number of bytes that
	// will be used for the on-disk buffer as part of model
	// auto-detection.
	// Default: 256MiB
	AutoDetectBufferSize int64

	// TempDir directory to use for temporary files.
	// Default: system setting (usually /tmp).
	TempDir string
}

func (o *CSVOptions) norm() {
	if o.Separator == 0 {
		o.Separator = ','
	}
	if o.AutoDetectBufferSize == 0 {
		o.AutoDetectBufferSize = 256 * 1024 * 1024
	}
}

// CSVStream can ingest examples from a CSV source.
type CSVStream struct {
	r *csv.Reader
	m *reason.Model
	t *ioutilx.TempFile
	o CSVOptions
}

// NewCSVStream returns a CSV stream for a given model. The model will be auto-detected
// if not provided.
func NewCSVStream(r io.Reader, m *reason.Model, opt *CSVOptions) (*CSVStream, error) {
	// normalize options
	var o CSVOptions
	if opt != nil {
		o = *opt
	}
	o.norm()

	br := bufio.NewReader(r)

	// init a CSV reader
	cr := newCSVReader(br, o.Separator)

	// skip N rows
	for i := 0; i < o.SkipRows; i++ {
		if _, err := cr.Read(); err != nil {
			return nil, err
		}
	}

	// read row to auto-detect header names
	if len(o.Headers) == 0 {
		row, err := cr.Read()
		if err != nil {
			return nil, fmt.Errorf("reason: unable to detect CSV headers: %v", err)
		}
		o.Headers = append(o.Headers, row...)
	}

	// ensure header names are unique & present
	present := false
	for i, name := range o.Headers {
		if name == "" {
			continue
		}
		for _, next := range o.Headers[i+1:] {
			if name == next {
				return nil, fmt.Errorf("reason: duplicate header %q", name)
			}
		}
		present = true
	}
	if !present {
		return nil, fmt.Errorf("reason: CSV streams require at least one header")
	}

	// detect mode if none given
	if m == nil {
		m, t, err := csvParseModel(br, &o)
		if err != nil {
			return nil, err
		}

		cr := newCSVReader(io.MultiReader(t, br), o.Separator)
		return &CSVStream{r: cr, m: m, t: t, o: o}, nil
	}
	return &CSVStream{r: cr, m: m, o: o}, nil
}

func csvParseModel(br *bufio.Reader, o *CSVOptions) (*reason.Model, *ioutilx.TempFile, error) {
	// create a tempfile
	t, err := ioutilx.NewTempFile(o.TempDir)
	if err != nil {
		return nil, nil, err
	}

	// buffer to tempfile
	if _, err := io.CopyN(t, br, o.AutoDetectBufferSize); err != nil && err != io.EOF {
		_ = t.Close()
		return nil, nil, err
	}

	// reopen the tempfile for reading
	if err := t.ReOpen(); err != nil {
		_ = t.Close()
		return nil, nil, err
	}

	// detect model
	m, err := csvDetectModel(t, o)
	if err != nil {
		_ = t.Close()
		return nil, nil, err
	}

	// reopen tempfile again
	if err := t.ReOpen(); err != nil {
		_ = t.Close()
		return nil, nil, err
	}

	// combine readers
	return m, t, nil
}

// Model returns the related model.
func (s *CSVStream) Model() *reason.Model {
	return s.m
}

// Next returns the next example.
func (s *CSVStream) Next() (reason.Example, error) {
	row, err := s.r.Read()
	if err != nil {
		return nil, err
	}

	mapx := make(reason.MapExample, len(s.o.Headers))
	for i, name := range s.o.Headers {
		if i >= len(row) {
			break
		} else if name == "" {
			continue
		}

		feat := s.m.Feature(name)
		if feat == nil {
			continue
		}

		switch feat.Kind {
		case reason.Feature_NUMERICAL:
			mapx[name] = feat.NumberOf(row[i])
		case reason.Feature_CATEGORICAL:
			mapx[name] = row[i]
		}
	}
	return mapx, nil
}

// Close closes the stream and releases resources.
func (s *CSVStream) Close() error {
	var err error
	if s.t != nil {
		err = s.t.Close()
	}
	return err
}

func csvDetectModel(r io.Reader, o *CSVOptions) (*reason.Model, error) {
	cr := newCSVReader(r, o.Separator)

	fads := make([]*featureAutoDetect, len(o.Headers))
	for i, name := range o.Headers {
		if name != "" {
			fads[i] = newFeatureAutoDetect(name)
		}
	}

	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		for i, val := range row {
			if i >= len(fads) {
				break
			}

			if fads[i] != nil {
				fads[i].ObserveString(val)
			}
		}
	}

	feats := make([]*reason.Feature, 0, len(fads))
	for _, fad := range fads {
		if fad == nil {
			continue
		}

		feat, err := fad.Feature()
		if err != nil {
			return nil, err
		}
		feats = append(feats, feat)
	}
	return reason.NewModel(feats...), nil
}

func newCSVReader(r io.Reader, sep rune) *csv.Reader {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true
	cr.Comma = sep
	cr.ReuseRecord = true
	return cr
}
