package msgpack

import (
	"fmt"
	"math"
	"reflect"
	"sync"
)

const customExtType int8 = 8

var (
	registerLock       sync.RWMutex
	nameToConcreteType = make(map[string]reflect.Type)
	concreteTypeToName = make(map[reflect.Type]string)
)

// Register registers a concrete type.
func Register(v interface{}) {
	registerLock.Lock()
	defer registerLock.Unlock()

	rt := reflect.TypeOf(v)
	name := rt.String()

	if len(name) >= math.MaxUint8 {
		panic(fmt.Sprintf("msgpack: type exceeds maximum name length %q", name))
	}
	if t, ok := nameToConcreteType[name]; ok && t != rt {
		panic(fmt.Sprintf("msgpack: registering duplicate types for %q: %s != %s", name, t, rt))
	}
	if n, ok := concreteTypeToName[rt]; ok && n != name {
		panic(fmt.Sprintf("msgpack: registering duplicate names for %s: %q != %q", rt, n, name))
	}
	nameToConcreteType[name] = rt
	concreteTypeToName[rt] = name
}

func concreteTypeByName(name string) (reflect.Type, bool) {
	registerLock.RLock()
	rv, ok := nameToConcreteType[name]
	registerLock.RUnlock()
	return rv, ok
}

func concreteNameByType(rt reflect.Type) (string, bool) {
	registerLock.RLock()
	name, ok := concreteTypeToName[rt]
	registerLock.RUnlock()
	return name, ok
}
