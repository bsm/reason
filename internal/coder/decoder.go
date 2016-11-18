package coder

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"reflect"
)

type Decoder struct {
	r *bufio.Reader
}

// NewDecoder opens a new encoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func (d *Decoder) Decode(v interface{}) error {
	return d.DecodeValue(reflect.ValueOf(v))
}

func (d *Decoder) DecodeValue(rv reflect.Value) error {
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("coder: non-pointer type %q", rv.Type())
	}

	el := rv.Elem()
	switch el.Kind() {
	case reflect.Float32:
		if n, err := d.Float32(); err != nil {
			return err
		} else {
			el.SetFloat(float64(n))
		}
		return nil
	case reflect.Float64:
		if n, err := d.Float64(); err != nil {
			return err
		} else {
			el.SetFloat(n)
		}
		return nil
	case reflect.String:
		if s, err := d.String(); err != nil {
			return err
		} else {
			el.SetString(s)
		}
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n, err := d.Int64(); err != nil {
			return err
		} else {
			el.SetInt(n)
		}
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if n, err := d.Uint64(); err != nil {
			return err
		} else {
			el.SetUint(n)
		}
		return nil
	case reflect.Map:
		return d.readMap(el)
	case reflect.Slice, reflect.Array:
		if el.Type() != sliceOfBytes {
			return d.readSlice(el)
		}
		if b, err := d.Bytes(); err != nil {
			return err
		} else {
			el.SetBytes(b)
		}
		return nil
	}

	if val, ok := rv.Interface().(Decodable); ok {
		regType, err := d.Type()
		if err != nil {
			return err
		}
		if rv.Type() != regType {
			return fmt.Errorf("coder: type %q does not match encoded type %q", rv.Type(), regType)
		}
		return val.DecodeFrom(d)
	}

	return fmt.Errorf("coder: type %q is not decodable", rv.Type())
}

// Type reads a type
func (d *Decoder) Type() (reflect.Type, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return nil, err
	}

	n := 0
	switch p {
	case mfixext1:
		n = 1
	case mfixext2:
		n = 2
	case mfixext4:
		n = 4
	case mfixext8:
		n = 8
	case mfixext16:
		n = 16
	case mext8:
		b, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		n = int(b)
	case mext16:
		b, err := d.read(2)
		if err != nil {
			return nil, err
		}
		n = int(bigEndian.Uint16(b))
	case mext32:
		b, err := d.read(4)
		if err != nil {
			return nil, err
		}
		n = int(bigEndian.Uint32(b))
	default:
		return nil, badPrefix(p)
	}

	x, err := d.r.ReadByte()
	if err != nil {
		return nil, err
	}
	if int8(x) != codableExtType {
		return nil, fmt.Errorf("coder: not a type extension")
	}

	name, err := d.read(n)
	if err != nil {
		return nil, err
	}
	return concreteTypeByName(string(name))
}

// PeekType peeks a type
func (d *Decoder) PeekType() (reflect.Type, error) {
	b, err := d.r.Peek(2)
	if err != nil {
		return nil, err
	}

	n := 0
	o := 2
	switch b[0] {
	case mfixext1:
		n = 1
	case mfixext2:
		n = 2
	case mfixext4:
		n = 4
	case mfixext8:
		n = 8
	case mfixext16:
		n = 16
	case mext8:
		n = int(b[1])
		o = 2
	case mext16:
		b, err := d.r.Peek(3)
		if err != nil {
			return nil, err
		}
		n = int(bigEndian.Uint16(b[1:]))
		o = 3
	case mext32:
		b, err := d.r.Peek(5)
		if err != nil {
			return nil, err
		}
		n = int(bigEndian.Uint32(b[1:]))
		o = 5
	default:
		return nil, badPrefix(b[0])
	}

	if b, err = d.r.Peek(o + n + 1); err != nil {
		return nil, err
	}
	if int8(b[o]) != codableExtType {
		return nil, fmt.Errorf("coder: not a type extension")
	}

	return concreteTypeByName(string(b[o+1:]))
}

// Float64 reads a float64
func (d *Decoder) Float64() (float64, error) {
	if err := d.validatePrefix(mfloat64); err != nil {
		return 0, err
	}
	b, err := d.read(8)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(bigEndian.Uint64(b)), nil
}

// Float32 reads a float32
func (d *Decoder) Float32() (float32, error) {
	if err := d.validatePrefix(mfloat32); err != nil {
		return 0, err
	}
	b, err := d.read(4)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(bigEndian.Uint32(b)), nil
}

// Int64 reads an int64
func (d *Decoder) Int64() (int64, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}

	if isfixint(p) {
		return int64(rfixint(p)), nil
	} else if isnfixint(p) {
		return int64(rnfixint(p)), nil
	}

	switch p {
	case mint8, muint8:
		b, err := d.r.ReadByte()
		if err != nil {
			return 0, err
		}
		return int64(b), nil
	case mint16, muint16:
		b, err := d.read(2)
		if err != nil {
			return 0, err
		}
		return int64(int16(b[0])<<8 | int16(b[1])), nil
	case mint32, muint32:
		b, err := d.read(4)
		if err != nil {
			return 0, err
		}
		return int64(int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])), nil
	case mint64, muint64:
		b, err := d.read(8)
		if err != nil {
			return 0, err
		}
		return int64(b[0])<<56 | int64(b[1])<<48 |
			int64(b[2])<<40 | int64(b[3])<<32 |
			int64(b[4])<<24 | int64(b[5])<<16 |
			int64(b[6])<<8 | int64(b[7]), nil
	}
	return 0, badPrefix(p)
}

// Uint64 reads a uint64
func (d *Decoder) Uint64() (uint64, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}

	if isfixint(p) {
		return uint64(rfixint(p)), nil
	}

	switch p {
	case mint8, muint8:
		b, err := d.r.ReadByte()
		if err != nil {
			return 0, err
		}
		return uint64(b), nil
	case mint16, muint16:
		b, err := d.read(2)
		if err != nil {
			return 0, err
		}
		return uint64(uint16(b[0])<<8 | uint16(b[1])), nil
	case mint32, muint32:
		b, err := d.read(4)
		if err != nil {
			return 0, err
		}
		return uint64(uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])), nil
	case mint64, muint64:
		b, err := d.read(8)
		if err != nil {
			return 0, err
		}
		return uint64(b[0])<<56 | uint64(b[1])<<48 |
			uint64(b[2])<<40 | uint64(b[3])<<32 |
			uint64(b[4])<<24 | uint64(b[5])<<16 |
			uint64(b[6])<<8 | uint64(b[7]), nil
	}
	return 0, badPrefix(p)
}

func (d *Decoder) ArrayHeader() (int, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}

	if isfixarray(p) {
		return int(rfixarray(p)), nil
	}
	switch p {
	case marray16:
		b, err := d.read(2)
		if err != nil {
			return 0, err
		}
		return int(bigEndian.Uint16(b)), nil
	case marray32:
		b, err := d.read(4)
		if err != nil {
			return 0, err
		}
		return int(bigEndian.Uint32(b)), nil
	}
	return 0, badPrefix(p)
}

func (d *Decoder) MapHeader() (int, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}

	if isfixmap(p) {
		return int(rfixmap(p)), nil
	}
	switch p {
	case mmap16:
		b, err := d.read(2)
		if err != nil {
			return 0, err
		}
		return int(bigEndian.Uint16(b)), nil
	case mmap32:
		b, err := d.read(4)
		if err != nil {
			return 0, err
		}
		return int(bigEndian.Uint32(b)), nil
	}
	return 0, badPrefix(p)
}

// Bool reads a bool
func (d *Decoder) Bool() (bool, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return false, err
	}

	switch p {
	case mtrue:
		return true, nil
	case mfalse:
		return false, nil
	}
	return false, badPrefix(p)
}

// String reads a string
func (d *Decoder) String() (string, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return "", err
	}

	var n int
	if isfixstr(p) {
		n = int(rfixstr(p))
	} else {
		switch p {
		case mstr8:
			p, err := d.r.ReadByte()
			if err != nil {
				return "", err
			}
			n = int(p)
		case mstr16:
			b, err := d.read(2)
			if err != nil {
				return "", err
			}
			n = int(bigEndian.Uint16(b))
		case mstr32:
			b, err := d.read(4)
			if err != nil {
				return "", err
			}
			n = int(bigEndian.Uint32(b))
		default:
			return "", badPrefix(p)
		}
	}

	b, err := d.read(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Bytes reads bytes
func (d *Decoder) Bytes() ([]byte, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return nil, err
	}

	var n int
	switch p {
	case mnil:
		return nil, nil
	case mbin8:
		p, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		n = int(p)
	case mbin16:
		b, err := d.read(2)
		if err != nil {
			return nil, err
		}
		n = int(bigEndian.Uint16(b))
	case mbin32:
		b, err := d.read(4)
		if err != nil {
			return nil, err
		}
		n = int(bigEndian.Uint32(b))
	default:
		return nil, badPrefix(p)
	}

	b, err := d.read(n)
	if err != nil {
		return nil, err
	}

	bin := make([]byte, len(b))
	copy(bin, b)
	return bin, nil
}

// --------------------------------------------------------------------

func (d *Decoder) validatePrefix(prefix byte) error {
	p, err := d.r.ReadByte()
	if err != nil {
		return err
	} else if p != prefix {
		return badPrefix(p)
	}
	return nil
}

func (d *Decoder) read(n int) ([]byte, error) {
	b, err := d.r.Peek(n)
	if err != nil {
		return nil, err
	}

	if _, err := d.r.Discard(n); err != nil {
		return nil, err
	}
	return b, nil
}

func (d *Decoder) readMap(rv reflect.Value) error {
	sz, err := d.MapHeader()
	if err != nil {
		return err
	}

	mp := reflect.MakeMap(rv.Type())
	kt := rv.Type().Key()
	vt := rv.Type().Elem()

	for i := 0; i < sz; i++ {
		key := reflect.New(kt)
		if err := d.DecodeValue(key); err != nil {
			return err
		}
		val := reflect.New(vt)
		if err := d.DecodeValue(val); err != nil {
			return err
		}
		mp.SetMapIndex(key.Elem(), val.Elem())
	}

	rv.Set(mp)
	return nil
}

func (d *Decoder) readSlice(rv reflect.Value) error {
	sz, err := d.ArrayHeader()
	if err != nil {
		return err
	}

	sl := reflect.MakeSlice(rv.Type(), sz, sz)
	for i := 0; i < sz; i++ {
		if err := d.DecodeValue(sl.Index(i).Addr()); err != nil {
			return err
		}
	}
	rv.Set(sl)
	return nil
}
