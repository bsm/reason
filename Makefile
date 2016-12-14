default: vet test

deps:
	go get -t ./...

vet:
	go vet ./...

test:
	go test ./internal/... ./... -short

test-full:
	go test ./internal/... ./...

bench:
	go test ./internal/... ./... -test.run=NONE -test.bench=. -benchmem -benchtime=1s

bench-race:
	go test ./internal/... ./... -test.run=NONE -test.bench=. -race
