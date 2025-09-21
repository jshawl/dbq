.PHONY: build lint run test

build:
	go build -o /dev/null .

lint:
	golangci-lint run

run:
	go run .

test:
	go test ./... -coverprofile coverage.out

test-cover: test
	go tool cover -html=coverage.out
