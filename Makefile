.PHONY: test test-verbose cover vet

## test: run all tests
test:
	go test ./...

## test-verbose: run all tests with verbose output
test-verbose:
	go test -v ./...

## cover: run all tests with a coverage summary
cover:
	go test -cover ./...

## vet: run go vet on all packages
vet:
	go vet ./...
