.DEFAULT_GOAL := check

FUZZTIME ?= 30s

# Keep in sync with .github/workflows/ci.yml. A different version can report
# different issues, so lint and fmt refuse to run with any other version.
GOLANGCI_LINT_VERSION := 2.15.0

# bench-compare: a pinned benchstat, run with go run (a tool, not a module dependency).
BENCHSTAT := golang.org/x/perf/cmd/benchstat@v0.0.0-20261009192801-be2c69fb417e
BENCH_BASE ?= $(shell git describe --tags --abbrev=0)
BENCH_COUNT ?= 6
BENCH_TIME ?= 1s
# Both sides run on the toolchain this module selects, so Go releases don't skew the comparison.
GO_VERSION := $(shell go env GOVERSION)

# api: a pinned gorelease, run with go run. Keep in sync with .github/workflows/ci.yml.
GORELEASE := golang.org/x/exp/cmd/gorelease@v0.0.0-20261014143641-f0251eb4aa3e

.PHONY: check api build vet test race allocs bench bench-compare profile cover fuzz lint fmt fmt-check tidy-check golangci-lint-version

## check: everything CI runs (go.mod tidiness, formatting, vet, lint, race-enabled tests, allocation guards)
check: tidy-check fmt-check vet lint race allocs

## api: exported API changes since the latest release tag, and the suggested next version (needs a clean working tree: commit first)
api:
	GOTOOLCHAIN=$(GO_VERSION) go run $(GORELEASE)

build:
	go build ./...

vet:
	go vet ./...

test:
	go test -count=1 ./...

race:
	go test -race -count=1 ./...

## allocs: TestAllocations, which -race builds exclude because the race detector changes allocation counts
allocs:
	go test -count=1 -run='^TestAllocations$$' .

## bench: benchmarks live in the root package only
bench:
	go test -run='^$$' -bench=. -benchmem .

## bench-compare: benchmarks BENCH_BASE (default: the latest tag) and the working tree, then compares them
bench-compare:
	rm -rf .bench && mkdir -p .bench/base
	git archive $(BENCH_BASE) | tar -x -C .bench/base
	cd .bench/base && GOTOOLCHAIN=$(GO_VERSION) go test -run='^$$' -bench=. -benchmem \
		-count=$(BENCH_COUNT) -benchtime=$(BENCH_TIME) . > ../base.txt
	GOTOOLCHAIN=$(GO_VERSION) go test -run='^$$' -bench=. -benchmem \
		-count=$(BENCH_COUNT) -benchtime=$(BENCH_TIME) . > .bench/head.txt
	go run $(BENCHSTAT) .bench/base.txt .bench/head.txt

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

lint: golangci-lint-version
	golangci-lint run ./...

## fmt: format with the formatters in .golangci.yaml (gofumpt, goimports, golines)
fmt: golangci-lint-version
	golangci-lint fmt ./...

fmt-check: golangci-lint-version
	@golangci-lint fmt --diff ./... || { echo "formatting needed: run make fmt"; exit 1; }

tidy-check:
	@go mod tidy -diff || { echo "go.mod is not tidy: run go mod tidy"; exit 1; }

golangci-lint-version:
	@golangci-lint version 2>/dev/null | grep -q 'version $(GOLANGCI_LINT_VERSION) ' || { \
		echo "golangci-lint v$(GOLANGCI_LINT_VERSION) is required; found: $$(golangci-lint version 2>&1 | head -n 1)"; \
		echo "install: curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v$(GOLANGCI_LINT_VERSION)"; \
		exit 1; }
