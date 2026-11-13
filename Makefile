.DEFAULT_GOAL := check

FUZZTIME ?= 30s

# Keep in sync with .github/workflows/ci.yml. A different version can report
# different issues, so lint and fmt refuse to run with any other version.
GOLANGCI_LINT_VERSION := 2.15.0

.PHONY: check build vet test race bench profile cover fuzz lint fmt fmt-check tidy-check golangci-lint-version

## check: everything CI runs (go.mod tidiness, formatting, vet, lint, race-enabled tests)
check: tidy-check fmt-check vet lint race

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
