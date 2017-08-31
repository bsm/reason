PKG=$(shell go list ./... | grep -v 'vendor')

default: vet test

vet:
	go vet $(PKG)

test:
	go test $(PKG) -short

test-full:
	go test $(PKG)

bench:
	go test $(PKG) -test.run=NONE -test.bench=. -benchmem -benchtime=1s

bench-race:
	go test $(PKG) -test.run=NONE -test.bench=. -race
