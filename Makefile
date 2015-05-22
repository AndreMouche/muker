all: build test

build:
	go build -o ./bin/muker ./cmd/muker.go

test:
	go test ./...

run: build start

start:
	./bin/muker


