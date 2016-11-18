package coder

import (
	"fmt"
	"reflect"
	"sync"
)

const codableExtType int8 = 8

var (
	registerLock       sync.RWMutex
	nameToConcreteType = make(map[string]reflect.Type)
	concreteTypeToName = make(map[reflect.Type]string)
)

// Register registers a concrete type.
func Register(value Codable) {
	registerLock.Lock()
	defer registerLock.Unlock()

	rt := reflect.TypeOf(value)
	name := rt.String()

	if t, ok := nameToConcreteType[name]; ok && t != rt {
		panic(fmt.Sprintf("coder: registering duplicate types for %q: %s != %s", name, t, rt))
	}
	if n, ok := concreteTypeToName[rt]; ok && n != name {
		panic(fmt.Sprintf("coder: registering duplicate names for %s: %q != %q", rt, n, name))
	}
	nameToConcreteType[name] = rt
	concreteTypeToName[rt] = name
}

func concreteTypeByName(name string) (reflect.Type, error) {
	registerLock.RLock()
	rv, ok := nameToConcreteType[name]
	registerLock.RUnlock()

	if !ok {
		return nil, fmt.Errorf("coder: not a registered name %q", name)
	}
	return rv, nil
}

func concreteNameByType(rt reflect.Type) (string, error) {
	registerLock.RLock()
	name, ok := concreteTypeToName[rt]
	registerLock.RUnlock()

	if !ok {
		return "", fmt.Errorf("coder: not a registered type %q", rt)
	}
	return name, nil
}
