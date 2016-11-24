package msgpack

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math"
	"reflect"
)

type Decoder struct {
	r   *bufio.Reader
	ctx context.Context
}

// NewDecoder opens a new encoder
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:   bufio.NewReader(r),
		ctx: context.Background(),
	}
}

func (d *Decoder) Context() context.Context { return d.ctx }
func (d *Decoder) SetContext(ctx context.Context) {
	if ctx != nil {
		d.ctx = ctx
	}
}
func (d *Decoder) WithContext(cb func(context.Context) context.Context) *Decoder {
	d.SetContext(cb(d.ctx))
	return d
}

// Decode decodes values
func (d *Decoder) Decode(vv ...interface{}) error {
	for _, v := range vv {
		if err := d.DecodeValue(reflect.ValueOf(v)); err != nil {
			return err
		}
	}
	return nil
}

// Decode decodes a value
func (d *Decoder) DecodeValue(rv reflect.Value) error {
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("msgpack: not a pointer type %q", rv.Type())
	}

	customType, err := d.readCustomType()
	if err != nil {
		return err
	}
	return d.decodeValue(rv, customType)
}

// --------------------------------------------------------------------

func (d *Decoder) decodeValue(rv reflect.Value, customType reflect.Type) error {
	elem := rv.Elem()

	switch elem.Kind() {
	case reflect.Bool:
		if n, err := d.readBool(); err != nil {
			return err
		} else {
			elem.SetBool(n)
		}
		return nil
	case reflect.Float32:
		if n, err := d.readFloat32(); err != nil {
			return err
		} else {
			elem.SetFloat(float64(n))
		}
		return nil
	case reflect.Float64:
		if n, err := d.readFloat64(); err != nil {
			return err
		} else {
			elem.SetFloat(n)
		}
		return nil
	case reflect.String:
		if s, err := d.readString(); err != nil {
			return err
		} else {
			elem.SetString(s)
		}
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n, err := d.readInt64(); err != nil {
			return err
		} else {
			elem.SetInt(n)
		}
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if n, err := d.readUint64(); err != nil {
			return err
		} else {
			elem.SetUint(n)
		}
		return nil
	case reflect.Slice, reflect.Array:
		if elem.Type() != sliceOfBytes {
			return d.decodeSlice(elem)
		}
		if b, err := d.readBytes(); err != nil {
			return err
		} else {
			elem.SetBytes(b)
		}
		return nil
	case reflect.Map:
		return d.decodeMap(elem)
	case reflect.Struct:
		if vv, ok := rv.Interface().(Decodable); ok {
			return vv.DecodeFrom(d)
		}
	case reflect.Ptr:
		if isNil, err := d.isNil(); isNil || err != nil {
			return err
		}

		if _, ok := elem.Interface().(Decodable); ok {
			x := reflect.New(elem.Type().Elem())
			if err := x.Interface().(Decodable).DecodeFrom(d); err != nil {
				return err
			}
			elem.Set(x)
			return nil
		}
	case reflect.Interface:
		if customType != nil {
			if !customType.AssignableTo(elem.Type()) {
				return fmt.Errorf("msgpack: %q is does not implement %q", customType, elem.Type())
			}

			cp := reflect.New(customType)
			if err := d.decodeValue(cp, customType); err != nil {
				return err
			}
			elem.Set(cp.Elem())
			return nil
		}
	}

	return errTypeNotSupported(elem.Type())
}

func (d *Decoder) readFloat64() (float64, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	} else if p != mfloat64 {
		return 0, errBadPrefix(p, "float64")
	}
	b, err := d.read(8)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(bigEndian.Uint64(b)), nil
}

func (d *Decoder) readFloat32() (float32, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	} else if p != mfloat32 {
		return 0, errBadPrefix(p, "float32")
	}
	b, err := d.read(4)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(bigEndian.Uint32(b)), nil
}

func (d *Decoder) readInt64() (int64, error) {
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
	return 0, errBadPrefix(p, "uint")
}

func (d *Decoder) readUint64() (uint64, error) {
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
		return uint64(bigEndian.Uint16(b)), nil
	case mint32, muint32:
		b, err := d.read(4)
		if err != nil {
			return 0, err
		}
		return uint64(bigEndian.Uint32(b)), nil
	case mint64, muint64:
		b, err := d.read(8)
		if err != nil {
			return 0, err
		}
		return bigEndian.Uint64(b), nil
	}
	return 0, errBadPrefix(p, "int")
}

func (d *Decoder) readString() (string, error) {
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
			return "", errBadPrefix(p, "string")
		}
	}

	b, err := d.read(n)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (d *Decoder) readBytes() ([]byte, error) {
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
		return nil, errBadPrefix(p, "[]byte")
	}

	b, err := d.read(n)
	if err != nil {
		return nil, err
	}

	bin := make([]byte, len(b))
	copy(bin, b)
	return bin, nil
}

func (d *Decoder) decodeSlice(rv reflect.Value) error {
	sz, err := d.readSliceSize()
	if err != nil || sz < 0 {
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

func (d *Decoder) decodeMap(rv reflect.Value) error {
	sz, err := d.readMapSize()
	if err != nil || sz < 0 {
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

func (d *Decoder) readSliceSize() (int, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}

	if p == mnil {
		return -1, nil
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
	return 0, errBadPrefix(p, "slice header")
}

func (d *Decoder) readMapSize() (int, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return 0, err
	}

	if p == mnil {
		return -1, nil
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
	return 0, errBadPrefix(p, "map header")
}

func (d *Decoder) readBool() (bool, error) {
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
	return false, errBadPrefix(p, "bool")
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

func (d *Decoder) isNil() (bool, error) {
	p, err := d.r.ReadByte()
	if err != nil {
		return false, err
	} else if p == mnil {
		return true, nil
	}
	return false, d.r.UnreadByte()
}

func (d *Decoder) readCustomType() (reflect.Type, error) {
	b, err := d.r.Peek(4)
	if err == io.EOF {
		return nil, nil // ignore errors
	} else if err != nil {
		return nil, err
	}

	if b[0] != mfixext2 || int8(b[1]) != customExtType {
		return nil, nil
	}

	if _, err := d.r.Discard(4); err != nil {
		return nil, err
	}

	code := bigEndian.Uint16(b[2:])
	rt, ok := concreteTypeByCode(code)
	if !ok {
		return nil, errCodeNotRegistered(code)
	}
	return rt, nil
}
