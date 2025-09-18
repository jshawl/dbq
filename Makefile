.PHONY: build lint run test

build:
	go build -o /dev/null -a .

lint:
	golangci-lint run

run:
	go run .

test:
	go test -v . -coverprofile coverage.out

test-cover: test
	go tool cover -html=coverage.out
