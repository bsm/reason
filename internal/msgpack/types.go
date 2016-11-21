package msgpack

import (
	"fmt"
	"reflect"
	"sync"
)

const customExtType int8 = 8

var (
	registerLock       sync.RWMutex
	codeToConcreteType = make(map[uint16]reflect.Type)
	concreteTypeToCode = make(map[reflect.Type]uint16)
)

// Register registers a concrete type.
func Register(code uint16, v interface{}) {
	registerLock.Lock()
	defer registerLock.Unlock()

	rt := reflect.TypeOf(v)
	if t, ok := codeToConcreteType[code]; ok && t != rt {
		panic(fmt.Sprintf("msgpack: registering duplicate types for %d: %s != %s", code, t, rt))
	}
	if n, ok := concreteTypeToCode[rt]; ok && n != code {
		panic(fmt.Sprintf("msgpack: registering duplicate codes for %s: %d != %d", rt, n, code))
	}
	codeToConcreteType[code] = rt
	concreteTypeToCode[rt] = code
}

func concreteTypeByCode(code uint16) (reflect.Type, bool) {
	registerLock.RLock()
	rv, ok := codeToConcreteType[code]
	registerLock.RUnlock()
	return rv, ok
}

func concreteCodeByType(rt reflect.Type) (uint16, bool) {
	registerLock.RLock()
	code, ok := concreteTypeToCode[rt]
	registerLock.RUnlock()
	return code, ok
}
