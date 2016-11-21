package msgpack

import (
	"bufio"
	"io"
	"math"
	"reflect"
)

// Encoder encodes structs
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

// Encode writes values
func (e *Encoder) Encode(vv ...interface{}) error {
	for _, v := range vv {
		if err := e.EncodeValue(reflect.ValueOf(v)); err != nil {
			return err
		}
	}
	return nil
}

// EncodeValue writes a value
func (e *Encoder) EncodeValue(v reflect.Value) error {
	if v.Kind() == reflect.Interface {
		return e.EncodeValue(v.Elem())
	}

	code, isCustomType := concreteCodeByType(v.Type())
	if isCustomType {
		if err := e.writeTypeCode(code); err != nil {
			return err
		}
	}

	switch v.Kind() {
	case reflect.Bool:
		return e.writeBool(v.Bool())
	case reflect.Float32:
		return e.writeFloat32(float32(v.Float()))
	case reflect.Float64:
		return e.writeFloat64(v.Float())
	case reflect.String:
		return e.writeString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.writeInt64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.writeUint64(v.Uint())
	case reflect.Map:
		return e.writeMap(v)
	case reflect.Slice, reflect.Array:
		return e.writeSlice(v)
	}

	switch vv := v.Interface().(type) {
	case Encodable:
		if !isCustomType {
			return errTypeNotRegistered(v.Type())
		}
		if v.Kind() == reflect.Ptr && v.IsNil() {
			return e.writeNil()
		}
		return vv.EncodeTo(e)
	}

	return errTypeNotSupported(v.Type())
}

func (e *Encoder) writeMapHeader(sz uint32) error {
	switch {
	case sz <= 15:
		return e.push(wfixmap(uint8(sz)))
	case sz <= math.MaxUint16:
		return e.prefix16(mmap16, uint16(sz))
	default:
		return e.prefix32(mmap32, sz)
	}
}

func (e *Encoder) writeSliceHeader(sz uint32) error {
	switch {
	case sz <= 15:
		return e.push(wfixarray(uint8(sz)))
	case sz <= math.MaxUint16:
		return e.prefix16(marray16, uint16(sz))
	default:
		return e.prefix32(marray32, sz)
	}
}

func (e *Encoder) writeNil() error { return e.push(mnil) }

func (e *Encoder) writeFloat64(f float64) error {
	return e.prefix64(mfloat64, math.Float64bits(f))
}

func (e *Encoder) writeFloat32(f float32) error {
	return e.prefix32(mfloat32, math.Float32bits(f))
}

func (e *Encoder) writeInt64(i int64) error {
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

func (e *Encoder) writeUint64(u uint64) error {
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

func (e *Encoder) writeBytes(b []byte) error {
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

func (e *Encoder) writeString(s string) error {
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

func (e *Encoder) writeBool(b bool) error {
	if b {
		return e.push(mtrue)
	}
	return e.push(mfalse)
}

func (e *Encoder) writeTypeCode(code uint16) error {
	if err := e.w.WriteByte(byte(mfixext2)); err != nil {
		return err
	}
	if err := e.w.WriteByte(byte(customExtType)); err != nil {
		return err
	}
	return e.uint16(code)
}

func (e *Encoder) writeSlice(v reflect.Value) (err error) {
	if v.IsNil() {
		return e.writeNil()
	}

	if v.Type() == sliceOfBytes {
		return e.writeBytes(v.Bytes())
	}

	sz := v.Len()
	if err = e.writeSliceHeader(uint32(sz)); err != nil {
		return
	}
	for i := 0; i < sz; i++ {
		if err = e.EncodeValue(v.Index(i)); err != nil {
			return
		}
	}
	return
}

func (e *Encoder) writeMap(v reflect.Value) (err error) {
	if v.IsNil() {
		return e.writeNil()
	}

	keys := v.MapKeys()
	if err = e.writeMapHeader(uint32(len(keys))); err != nil {
		return
	}

	for _, key := range keys {
		if err = e.EncodeValue(key); err != nil {
			return
		}

		val := v.MapIndex(key)
		if err = e.EncodeValue(val); err != nil {
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
