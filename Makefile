PKG=./...
PROTOS=$(shell find .  -name '*.proto' -not -path './vendor/*')

default: vet test

vet:
	go vet $(PKG)

test:
	go test $(PKG)

test-short:
	go test $(PKG) -short

bench:
	go test $(PKG) -test.run=NONE -test.bench=. -benchmem -benchtime=1s

bench-race:
	go test $(PKG) -test.run=NONE -test.bench=. -race

.PHONY: vet test test-short bench bench-race

proto: proto.go
proto.go: $(patsubst %.proto,%.pb.go,$(PROTOS))

.PHONY: proto proto.go

# ---------------------------------------------------------------------

PROTO_PATH=.:$$GOPATH/src:$$GOPATH/src/github.com/gogo/protobuf/protobuf

### proto.go

%.pb.go: %.proto
	protoc --gogo_out=. --proto_path=${PROTO_PATH} $<

