package internal

import (
	"bufio"
	"fmt"
	"io"

	"github.com/bsm/reason"
	"github.com/bsm/reason/internal/ioutilx"
	"github.com/bsm/reason/internal/protoio"
	"github.com/gogo/protobuf/proto"
)

// New inits a new optimiser from model with defaults.
func New(model *reason.Model, target string, size int) *Optimizer {
	return &Optimizer{
		Model:   model,
		Target:  target,
		Sums:    make([]float64, size),
		Weights: make([]float64, size),
	}
}

// WriteTo writes a tree to a Writer.
func (o *Optimizer) WriteTo(w io.Writer) (int64, error) {
	wc := &ioutilx.CountingWriter{W: w}
	wp := &protoio.Writer{Writer: bufio.NewWriter(wc)}

	if err := wp.WriteMessageField(1, o.Model); err != nil {
		return wc.N, err
	}
	if err := wp.WriteStringField(2, o.Target); err != nil {
		return wc.N, err
	}
	if err := wp.WriteDoubleSliceField(3, o.Sums); err != nil {
		return wc.N, err
	}
	if err := wp.WriteDoubleSliceField(4, o.Weights); err != nil {
		return wc.N, err
	}
	return wc.N, wp.Flush()
}

// ReadFrom reads an optimizer from a Reader.
func (o *Optimizer) ReadFrom(r io.Reader) (int64, error) {
	rc := &ioutilx.CountingReader{R: r}
	rp := &protoio.Reader{Reader: bufio.NewReader(rc)}

	for {
		tag, wire, err := rp.ReadField()
		if err == io.EOF {
			return rc.N, nil
		} else if err != nil {
			return rc.N, err
		}

		switch tag {
		case 1: // model
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			model := new(reason.Model)
			if err := rp.ReadMessage(model); err != nil {
				return rc.N, err
			}
			o.Model = model
		case 2: // target
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			str, err := rp.ReadString()
			if err != nil {
				return rc.N, err
			}
			o.Target = str
		case 3: // sums
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			slice, err := rp.ReadDoubleSlice(nil)
			if err != nil {
				return rc.N, err
			}
			o.Sums = slice
		case 4: // weights
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			slice, err := rp.ReadDoubleSlice(nil)
			if err != nil {
				return rc.N, err
			}
			o.Weights = slice
		default:
			return rc.N, fmt.Errorf("ftrl: unexpected field tag %d", tag)
		}
	}
}
