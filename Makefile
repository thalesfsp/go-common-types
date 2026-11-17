###
# Params.
###

PROJECT_FULL_NAME := go-common-types

HAS_GODOC := $(shell command -v godoc;)
HAS_DOCKER_COMPOSE := $(shell command -v docker-compose;)
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0

default: ci

###
# Entries.
###

ci: lint test coverage

coverage:
	@go tool cover -func=coverage.out && echo "Coverage OK"

deps:
	@go install golang.org/x/tools/cmd/godoc@latest

doc:
ifndef HAS_GODOC
	@echo "Could not find godoc, installing it and any other missing tool(s)"
	@make deps
endif
	@echo "Open localhost:6060/pkg/github.com/thalesfsp/$(PROJECT_FULL_NAME)/ in your browser\n"
	@godoc -http :6060

lint:
	@$(GOLANGCI_LINT) fmt --diff -c .golangci.yml
	@$(GOLANGCI_LINT) run --timeout 1m -v -c .golangci.yml && echo "Lint OK"

test:
	@go test -timeout 30s -short -v -race -cover -coverprofile=coverage.out ./... && echo "Test OK"

.PHONY: ci \
	coverage \
	deps \
	doc \
	lint \
	test
