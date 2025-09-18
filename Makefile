.PHONY: build lint run test

build:
	go build -o /dev/null -a .

lint:
	golangci-lint run

run:
	go run .

test:
	go test -v ./...
