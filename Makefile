default: test

deps:
	go get -t ./...

test:
	go test ./internal/... ./... -short

test-full:
	go test ./internal/... ./...

bench:
	go test ./internal/... ./... -test.run=NONE -test.bench=.

bench-race:
	go test ./internal/... ./... -test.run=NONE -test.bench=. -race
