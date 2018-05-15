package internal

import (
	"bufio"
	"fmt"
	"io"

	"github.com/bsm/reason/core"
	"github.com/bsm/reason/internal/iocount"
	"github.com/bsm/reason/internal/protoio"
	"github.com/gogo/protobuf/proto"
)

// NewOptimizer inits a new model with defaults
func NewOptimizer(model *core.Model, target string, size int) *Optimizer {
	return &Optimizer{
		Model:   model,
		Target:  target,
		Sums:    make([]float64, size),
		Weights: make([]float64, size),
	}
}

// WriteTo writes a tree to a Writer.
func (o *Optimizer) WriteTo(w io.Writer) (int64, error) {
	wc := &iocount.Writer{W: w}
	wp := protoio.Writer{Writer: bufio.NewWriter(wc)}

	if err := wp.WriteMessageField(1, o.Model); err != nil {
		return wc.N, err
	}
	if err := wp.WriteStringField(2, o.Target); err != nil {
		return wc.N, err
	}
	for _, f := range o.Sums {
		if err := wp.WriteDoubleField(3, f); err != nil {
			return wc.N, err
		}
	}
	for _, f := range o.Weights {
		if err := wp.WriteDoubleField(4, f); err != nil {
			return wc.N, err
		}
	}
	return wc.N, wp.Flush()
}

// ReadFrom reads an optimizer from a Reader.
func (o *Optimizer) ReadFrom(r io.Reader) (int64, error) {
	rc := &iocount.Reader{R: r}
	rp := protoio.Reader{Reader: bufio.NewReader(rc)}

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

			model := new(core.Model)
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
			if wire != proto.WireFixed64 {
				return rc.N, proto.ErrInternalBadWireType
			}

			f, err := rp.ReadDouble()
			if err != nil {
				return rc.N, err
			}
			o.Sums = append(o.Sums, f)
		case 4: // weights
			if wire != proto.WireFixed64 {
				return rc.N, proto.ErrInternalBadWireType
			}

			f, err := rp.ReadDouble()
			if err != nil {
				return rc.N, err
			}
			o.Weights = append(o.Weights, f)
		default:
			return rc.N, fmt.Errorf("hoeffding: unexpected field tag %d", tag)
		}
	}
}
