// Package msgpack implements serialisation helpers for large objects with
// circular references and many interfaces.
//
// Heavily "borrowed" from https://github.com/tinylib/msgp
// Copyright (c) 2014 Philip Hofer
package msgpack

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

// Encodable can encode to an encoder
type Encodable interface {
	EncodeTo(*Encoder) error
}

// Decodable can decode to a decoder
type Decodable interface {
	DecodeFrom(*Decoder) error
}

var (
	bigEndian    = binary.BigEndian
	sliceOfBytes = reflect.TypeOf(([]byte)(nil))
)

const (
	last4  = 0x0f
	first4 = 0xf0
	last5  = 0x1f
	first3 = 0xe0
	last7  = 0x7f
)

// These are all the byte
// prefixes defined by the
// msgpack standard
const (

	// 111XXXXX
	mnfixint uint8 = 0xe0

	// 1000XXXX
	mfixmap uint8 = 0x80

	// 1001XXXX
	mfixarray uint8 = 0x90

	// 101XXXXX
	mfixstr uint8 = 0xa0

	mnil      uint8 = 0xc0
	mfalse    uint8 = 0xc2
	mtrue     uint8 = 0xc3
	mbin8     uint8 = 0xc4
	mbin16    uint8 = 0xc5
	mbin32    uint8 = 0xc6
	mext8     uint8 = 0xc7
	mext16    uint8 = 0xc8
	mext32    uint8 = 0xc9
	mfloat32  uint8 = 0xca
	mfloat64  uint8 = 0xcb
	muint8    uint8 = 0xcc
	muint16   uint8 = 0xcd
	muint32   uint8 = 0xce
	muint64   uint8 = 0xcf
	mint8     uint8 = 0xd0
	mint16    uint8 = 0xd1
	mint32    uint8 = 0xd2
	mint64    uint8 = 0xd3
	mfixext1  uint8 = 0xd4
	mfixext2  uint8 = 0xd5
	mfixext4  uint8 = 0xd6
	mfixext8  uint8 = 0xd7
	mfixext16 uint8 = 0xd8
	mstr8     uint8 = 0xd9
	mstr16    uint8 = 0xda
	mstr32    uint8 = 0xdb
	marray16  uint8 = 0xdc
	marray32  uint8 = 0xdd
	mmap16    uint8 = 0xde
	mmap32    uint8 = 0xdf
)

func isfixint(b byte) bool {
	return b>>7 == 0
}

func isnfixint(b byte) bool {
	return b&first3 == mnfixint
}

func isfixmap(b byte) bool {
	return b&first4 == mfixmap
}

func isfixarray(b byte) bool {
	return b&first4 == mfixarray
}

func isfixstr(b byte) bool {
	return b&first3 == mfixstr
}

func wfixint(u uint8) byte {
	return u & last7
}

func rfixint(b byte) uint8 {
	return b
}

func wnfixint(i int8) byte {
	return byte(i) | mnfixint
}

func rnfixint(b byte) int8 {
	return int8(b)
}

func rfixmap(b byte) uint8 {
	return b & last4
}

func wfixmap(u uint8) byte {
	return mfixmap | (u & last4)
}

func rfixstr(b byte) uint8 {
	return b & last5
}

func wfixstr(u uint8) byte {
	return (u & last5) | mfixstr
}

func rfixarray(b byte) uint8 {
	return (b & last4)
}

func wfixarray(u uint8) byte {
	return (u & last4) | mfixarray
}

func errBadPrefix(p byte, t string) error {
	return fmt.Errorf("msgpack: unexpected type prefix 0x%x when reading %s", p, t)
}

func errTypeNotRegistered(t reflect.Type) error {
	return fmt.Errorf("msgpack: type %q not registered", t)
}

func errCodeNotRegistered(code uint16) error {
	return fmt.Errorf("msgpack: code %d not registered", code)
}

func errTypeNotSupported(t reflect.Type) error {
	return fmt.Errorf("msgpack: type %q not supported", t)
}
