BINARY := pokedex

.DEFAULT_GOAL := help

.PHONY: help run build test test-race cover fmt fmt-check vet vet-shadow tidy tidy-check ci clean

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the REPL
	go run .

build: ## Build the binary
	go build -o $(BINARY)

test: ## Run unit tests
	go test ./...

test-race: ## Run tests with the race detector (matches CI)
	go test -race ./...

cover: ## Run tests with race + coverage and print the function report
	go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

fmt: ## Format all Go files in place
	gofmt -w .

fmt-check: ## Fail if any file needs formatting (matches CI)
	@files="$$(gofmt -l .)"; \
	if [ -n "$$files" ]; then \
		echo "The following files need formatting (run 'make fmt'):"; \
		echo "$$files"; \
		exit 1; \
	fi

vet: ## Run go vet
	go vet ./...

vet-shadow: ## Install shadow analyzer (if missing) and run it
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
	go vet -vettool=$$(go env GOPATH)/bin/shadow ./...

tidy: ## Tidy the module
	go mod tidy

tidy-check: ## Fail if go.mod/go.sum would change (matches CI)
	@go mod tidy; \
	if [ -n "$$(git status --porcelain go.mod go.sum)" ]; then \
		git diff go.mod go.sum 2>/dev/null || true; \
		echo "go.mod or go.sum changed. Run 'make tidy' and commit the result."; \
		exit 1; \
	fi

ci: fmt-check vet vet-shadow cover build tidy-check ## Run the same checks CI runs

clean: ## Remove build and coverage artifacts
	rm -f $(BINARY) coverage.out coverage2.out
