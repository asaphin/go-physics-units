.DEFAULT_GOAL := check

FUZZTIME ?= 30s

.PHONY: check build vet test race bench profile cover fuzz lint fmt fmt-check

## check: everything CI runs (formatting, vet, lint, race-enabled tests)
check: fmt-check vet lint race

build:
	go build ./...

vet:
	go vet ./...

test:
	go test -count=1 ./...

race:
	go test -race -count=1 ./...

## bench: benchmarks live in the root package only
bench:
	go test -run='^$$' -bench=. -benchmem .

profile:
	go test -run='^$$' -bench=. -benchmem -cpuprofile=profile.out .

## cover: coverage across all packages, summary in the terminal, report in coverage/coverage.html
cover:
	mkdir -p coverage
	go test -count=1 -coverpkg=./... -coverprofile=coverage/coverage.out ./...
	go tool cover -func=coverage/coverage.out | tail -n 1
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

fuzz:
	go test -run='^$$' -fuzz=FuzzParseString -fuzztime=$(FUZZTIME) .

lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

fmt-check:
	@test -z "$$(gofmt -l .)" || { echo "gofmt needed on:"; gofmt -l .; exit 1; }
