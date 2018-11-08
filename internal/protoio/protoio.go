package protoio

import (
	"bufio"
	"encoding/binary"
	"io"
	"math"

	"github.com/gogo/protobuf/proto"
)

// Reader reads raw protobuf protocol from streams.
type Reader struct {
	*bufio.Reader

	tmp []byte
}

// ReadField reads the next field.
func (r *Reader) ReadField() (uint32, uint8, error) {
	u, err := r.ReadVarint()
	if err != nil {
		return 0, 0, err
	}
	return uint32(u) >> 3, uint8(u) & 0x07, nil
}

// ReadVarint reads a number.
func (r *Reader) ReadVarint() (uint64, error) {
	return binary.ReadUvarint(r)
}

// ReadString reads a string.
func (r *Reader) ReadString() (string, error) {
	b, err := r.ReadBytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ReadBytes reads raw bytes.
func (r *Reader) ReadBytes() ([]byte, error) {
	u, err := r.ReadVarint()
	if err != nil {
		return nil, err
	}

	b := r.getBuffer(int(u))
	_, err = io.ReadFull(r, b)
	return b, err
}

// ReadDouble reads a float64.
func (r *Reader) ReadDouble() (float64, error) {
	b := r.getBuffer(8)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return 0, err
	}

	return math.Float64frombits(binary.LittleEndian.Uint64(b)), nil
}

// ReadDoubleSlice reads a slice of float64s.
func (r *Reader) ReadDoubleSlice(dst []float64) ([]float64, error) {
	u, err := r.ReadVarint()
	if err != nil {
		return dst, err
	}

	n := int(u / 8)
	if dst == nil {
		dst = make([]float64, 0, n)
	}

	for i := 0; i < n; i++ {
		f, err := r.ReadDouble()
		if err != nil {
			return dst, err
		}
		dst = append(dst, f)
	}
	return dst, nil
}

// ReadMessage reads a message.
func (r *Reader) ReadMessage(m proto.Message) error {
	b, err := r.ReadBytes()
	if err != nil {
		return err
	}
	return proto.Unmarshal(b, m)
}

func (r *Reader) getBuffer(n int) []byte {
	if cap(r.tmp) < n {
		r.tmp = make([]byte, 0, n)
	}
	return r.tmp[:n]
}

// ---------------------------------------------------------------------

// Writer writes protobuf to streams.
type Writer struct {
	*bufio.Writer

	tmp []byte
}

// WriteField writes field info.
func (w *Writer) WriteField(tag uint32, wire uint8) error {
	n := (uint64(tag) << 3) | uint64(wire&0x07)
	return w.WriteVarint(n)
}

// WriteVarint writes a number.
func (w *Writer) WriteVarint(u uint64) error {
	b := w.getBuffer(binary.MaxVarintLen64)
	n := binary.PutUvarint(b, u)
	_, err := w.Write(b[:n])
	return err
}

// WriteDouble writes a float64.
func (w *Writer) WriteDouble(v float64) error {
	b := w.getBuffer(8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(v))
	_, err := w.Write(b)
	return err
}

// WriteVarintField writes a number field.
func (w *Writer) WriteVarintField(tag uint32, u uint64) error {
	if u == 0 {
		return nil
	}

	if err := w.WriteField(tag, proto.WireVarint); err != nil {
		return err
	}
	return w.WriteVarint(u)
}

// WriteStringField writes a string.
func (w *Writer) WriteStringField(tag uint32, s string) error {
	if s == "" {
		return nil
	}

	if err := w.WriteField(tag, proto.WireBytes); err != nil {
		return err
	}

	if err := w.WriteVarint(uint64(len(s))); err != nil {
		return err
	}

	_, err := w.Writer.WriteString(s)
	return err
}

// WriteBinaryField writes bytes.
func (w *Writer) WriteBinaryField(tag uint32, p []byte) error {
	if len(p) == 0 {
		return nil
	}

	if err := w.WriteField(tag, proto.WireBytes); err != nil {
		return err
	}

	if err := w.WriteVarint(uint64(len(p))); err != nil {
		return err
	}

	_, err := w.Writer.Write(p)
	return err
}

// WriteMessageField writes a message.
func (w *Writer) WriteMessageField(tag uint32, m proto.Message) error {
	if m == nil {
		return nil
	}

	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	if err := w.WriteField(tag, proto.WireBytes); err != nil {
		return err
	}

	if err := w.WriteVarint(uint64(len(data))); err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// WriteDoubleSliceField writes a slice of float64s.
func (w *Writer) WriteDoubleSliceField(tag uint32, slice []float64) error {
	if len(slice) == 0 {
		return nil
	}

	if err := w.WriteField(tag, proto.WireBytes); err != nil {
		return err
	}
	if err := w.WriteVarint(uint64(len(slice) * 8)); err != nil {
		return err
	}
	for _, f := range slice {
		if err := w.WriteDouble(f); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) getBuffer(n int) []byte {
	if cap(w.tmp) < n {
		w.tmp = make([]byte, n)
	}
	return w.tmp[:n]
}
