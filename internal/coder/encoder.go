package coder

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"reflect"
)

type Encoder struct {
	w *bufio.Writer
}

// NewEncoder opens a new encoder
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: bufio.NewWriter(w)}
}

// Close closes the encoder and flushes the buffer
func (e *Encoder) Close() error {
	return e.w.Flush()
}

// MapHeader writes a map header
func (e *Encoder) MapHeader(sz uint32) error {
	switch {
	case sz <= 15:
		return e.push(wfixmap(uint8(sz)))
	case sz <= math.MaxUint16:
		return e.prefix16(mmap16, uint16(sz))
	default:
		return e.prefix32(mmap32, sz)
	}
}

// ArrayHeader writes an array header of the
// given size
func (e *Encoder) ArrayHeader(sz uint32) error {
	switch {
	case sz <= 15:
		return e.push(wfixarray(uint8(sz)))
	case sz <= math.MaxUint16:
		return e.prefix16(marray16, uint16(sz))
	default:
		return e.prefix32(marray32, sz)
	}
}

// Nil writes a nil byte to the buffer
func (e *Encoder) Nil() error {
	return e.push(mnil)
}

// Float64 writes a float64
func (e *Encoder) Float64(f float64) error {
	return e.prefix64(mfloat64, math.Float64bits(f))
}

// Float32 writes a float32
func (e *Encoder) Float32(f float32) error {
	return e.prefix32(mfloat32, math.Float32bits(f))
}

// Int64 writes an int64
func (e *Encoder) Int64(i int64) error {
	if i >= 0 {
		switch {
		case i <= math.MaxInt8:
			return e.push(wfixint(uint8(i)))
		case i <= math.MaxInt16:
			return e.prefix16(mint16, uint16(i))
		case i <= math.MaxInt32:
			return e.prefix32(mint32, uint32(i))
		default:
			return e.prefix64(mint64, uint64(i))
		}
	}
	switch {
	case i >= -32:
		return e.push(wnfixint(int8(i)))
	case i >= math.MinInt8:
		return e.prefix8(mint8, uint8(i))
	case i >= math.MinInt16:
		return e.prefix16(mint16, uint16(i))
	case i >= math.MinInt32:
		return e.prefix32(mint32, uint32(i))
	default:
		return e.prefix64(mint64, uint64(i))
	}
}

// Uint64 writes a uint64
func (e *Encoder) Uint64(u uint64) error {
	switch {
	case u <= (1<<7)-1:
		return e.push(wfixint(uint8(u)))
	case u <= math.MaxUint8:
		return e.prefix8(muint8, uint8(u))
	case u <= math.MaxUint16:
		return e.prefix16(muint16, uint16(u))
	case u <= math.MaxUint32:
		return e.prefix32(muint32, uint32(u))
	default:
		return e.prefix64(muint64, u)
	}
}

// Bytes writes binary
func (e *Encoder) Bytes(b []byte) error {
	sz := uint32(len(b))
	var err error
	switch {
	case sz <= math.MaxUint8:
		err = e.prefix8(mbin8, uint8(sz))
	case sz <= math.MaxUint16:
		err = e.prefix16(mbin16, uint16(sz))
	default:
		err = e.prefix32(mbin32, sz)
	}
	if err != nil {
		return err
	}
	_, err = e.w.Write(b)
	return err
}

// String writes a string
func (e *Encoder) String(s string) error {
	sz := uint32(len(s))
	var err error
	switch {
	case sz <= 31:
		err = e.push(wfixstr(uint8(sz)))
	case sz <= math.MaxUint8:
		err = e.prefix8(mstr8, uint8(sz))
	case sz <= math.MaxUint16:
		err = e.prefix16(mstr16, uint16(sz))
	default:
		err = e.prefix32(mstr32, sz)
	}
	if err != nil {
		return err
	}
	_, err = e.w.WriteString(s)
	return err
}

// Bool writes a bool
func (e *Encoder) Bool(b bool) error {
	if b {
		return e.push(mtrue)
	}
	return e.push(mfalse)
}

// Type stores the encodable type
func (e *Encoder) Type(t Encodable) error {
	name, err := concreteNameByType(reflect.TypeOf(t))
	if err != nil {
		return err
	}

	if err := e.extHeader(codableExtType, uint32(len(name))); err != nil {
		return err
	}
	_, err = e.w.WriteString(name)
	return err
}

// Encode write a value
func (e *Encoder) Encode(v interface{}) error {
	switch vv := v.(type) {
	case nil:
		return e.Nil()
	case bool:
		return e.Bool(vv)
	case int:
		return e.Int64(int64(vv))
	case int8:
		return e.Int64(int64(vv))
	case int16:
		return e.Int64(int64(vv))
	case int32:
		return e.Int64(int64(vv))
	case int64:
		return e.Int64(vv)
	case uint:
		return e.Uint64(uint64(vv))
	case uint8:
		return e.Uint64(uint64(vv))
	case uint16:
		return e.Uint64(uint64(vv))
	case uint32:
		return e.Uint64(uint64(vv))
	case uint64:
		return e.Uint64(vv)
	case float64:
		return e.Float64(vv)
	case float32:
		return e.Float32(vv)
	case []byte:
		return e.Bytes(vv)
	case string:
		return e.String(vv)
	case Encodable:
		if err := e.Type(vv); err != nil {
			return err
		}
		return vv.EncodeTo(e)
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		return e.writeMap(rv)
	case reflect.Slice, reflect.Array:
		return e.writeSlice(rv)
	}
	return fmt.Errorf("coder: type %q not supported", rv.Type())
}

// --------------------------------------------------------------------

func (e *Encoder) writeSlice(v reflect.Value) (err error) {
	if v.IsNil() {
		return e.Nil()
	}

	sz := v.Len()
	if err = e.ArrayHeader(uint32(sz)); err != nil {
		return
	}
	for i := 0; i < sz; i++ {
		if err = e.Encode(v.Index(i).Interface()); err != nil {
			return
		}
	}
	return
}

func (e *Encoder) writeMap(v reflect.Value) (err error) {
	if v.IsNil() {
		return e.Nil()
	}

	keys := v.MapKeys()
	if err = e.MapHeader(uint32(len(keys))); err != nil {
		return
	}

	for _, key := range keys {
		if err = e.Encode(key.Interface()); err != nil {
			return
		}

		val := v.MapIndex(key)
		if err = e.Encode(val.Interface()); err != nil {
			return
		}
	}
	return
}

func (e *Encoder) push(b byte) error {
	return e.w.WriteByte(b)
}

func (e *Encoder) uint8(u uint8) error {
	if err := e.push(byte(u)); err != nil {
		return err
	}
	return nil
}

func (e *Encoder) uint16(u uint16) error {
	if err := e.push(byte(u >> 8)); err != nil {
		return err
	}
	return e.uint8(uint8(u))
}

func (e *Encoder) uint32(u uint32) error {
	if err := e.push(byte(u >> 24)); err != nil {
		return err
	}
	if err := e.push(byte(u >> 16)); err != nil {
		return err
	}
	return e.uint16(uint16(u))
}

func (e *Encoder) uint64(u uint64) error {
	if err := e.push(byte(u >> 56)); err != nil {
		return err
	}
	if err := e.push(byte(u >> 48)); err != nil {
		return err
	}
	if err := e.push(byte(u >> 40)); err != nil {
		return err
	}
	if err := e.push(byte(u >> 32)); err != nil {
		return err
	}
	return e.uint32(uint32(u))
}

func (e *Encoder) prefix8(b byte, u uint8) error {
	if err := e.push(b); err != nil {
		return err
	}
	return e.uint8(u)
}

func (e *Encoder) prefix16(b byte, u uint16) error {
	if err := e.push(b); err != nil {
		return err
	}
	return e.uint16(u)
}

func (e *Encoder) prefix32(b byte, u uint32) error {
	if err := e.push(b); err != nil {
		return err
	}
	return e.uint32(u)
}

func (e *Encoder) prefix64(b byte, u uint64) error {
	if err := e.push(b); err != nil {
		return err
	}
	return e.uint64(u)
}

func (e *Encoder) extHeader(n int8, sz uint32) error {
	switch sz {
	case 1:
		if err := e.push(mfixext1); err != nil {
			return err
		}
	case 2:
		if err := e.push(mfixext2); err != nil {
			return err
		}
	case 4:
		if err := e.push(mfixext4); err != nil {
			return err
		}
	case 8:
		if err := e.push(mfixext8); err != nil {
			return err
		}
	case 16:
		if err := e.push(mfixext16); err != nil {
			return err
		}
	default:
		switch {
		case sz < math.MaxUint8:
			if err := e.prefix8(mext8, uint8(sz)); err != nil {
				return err
			}
		case sz < math.MaxUint16:
			if err := e.prefix16(mext16, uint16(sz)); err != nil {
				return err
			}
		default:
			if err := e.prefix32(mext32, uint32(sz)); err != nil {
				return err
			}
		}
	}
	return e.push(byte(n))
}
