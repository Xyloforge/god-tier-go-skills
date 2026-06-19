# Starter Makefile for a cmd/-based Go module.
# Assumes binaries live in cmd/<name>/ and the module builds with the Go toolchain.

# Discover binaries from cmd/* (each subdir is one binary).
CMDS := $(notdir $(wildcard cmd/*))
GO   ?= go

.DEFAULT_GOAL := build

.PHONY: build
build: ## Build all binaries in cmd/* into ./bin
	@mkdir -p bin
	@for c in $(CMDS); do \
		echo "building $$c"; \
		$(GO) build -o bin/$$c ./cmd/$$c; \
	done

.PHONY: test
test: ## Run tests with the race detector
	$(GO) test -race ./...

.PHONY: cover
cover: ## Run tests and open a coverage report
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out | tail -1

.PHONY: bench
bench: ## Run benchmarks with allocation stats
	$(GO) test -run='^$$' -bench=. -benchmem ./...

.PHONY: vet
vet: ## Static analysis
	$(GO) vet ./...

.PHONY: lint
lint: ## golangci-lint (install: https://golangci-lint.run)
	golangci-lint run ./...

.PHONY: tidy
tidy: ## Sync go.mod/go.sum
	$(GO) mod tidy

.PHONY: check
check: vet lint test ## The gate to run before committing

.PHONY: clean
clean:
	rm -rf bin coverage.out

.PHONY: help
help: ## List targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-12s\033[0m %s\n",$$1,$$2}'
