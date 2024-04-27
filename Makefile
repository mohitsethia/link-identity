include ./scripts/make/golangci-lint.mk

#include .env if exists
-include .env

.DEFAULT_GOAL:=help

# runtime options
COMMIT_HASH = $(shell git rev-parse --short HEAD)
TAG         = $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)
GOPACKAGES  = $(shell go list ./...)
GOFILES		= $(shell find . -type f -name '*.go' -not -path "./vendor/*")
HAS_GOLINT  = $(shell command -v golint)

# app specific
APP_NAME    = link-identity-api

# go options
GO       ?= go
LDFLAGS  =  -X "main.tag=$(TAG)" \
			-X "main.commitHash=$(COMMIT_HASH)"

.PHONY: help
help: ## this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

tools: ## Install general tools globally (not in the project)
#	go install github.com/pressly/goose/cmd/goose@latest

build: ## Build binaries
	$(GO) build -ldflags '$(LDFLAGS)' -o ./link-identity-api cmd/link-identity-api/main.go

build-static: ## Build binaries statically
	CGO_ENABLED=0 $(GO) build -ldflags '$(LDFLAGS)' -v -installsuffix cgo -o ./link-identity-api cmd/link-identity-api/main.go

run: ## Run application
	$(GO) run cmd/link-identity-api/main.go

# running application
run-link-identity-api: ## Run application
	docker-compose up -d mysql
	$(GO) run cmd/link-identity-api/main.go

unit-tests: ## Run unit tests
	$(DC) $(GO) test --short -race -v ./...

tests: ## Run all tests
	$(DC) $(GO) test -race -v ./... -coverprofile=coverage.out

tests-deps: ## Install test dependencies
	GO111MODULE=off $(DC) $(GO) get github.com/stretchr/testify/mock
	GO111MODULE=off $(DC) $(GO) get -u github.com/schrej/godacov

coverage: ## View test coverage
	$(GO) tool cover -html=coverage.out

specific-unit-test:
	$(DC) $(GO) test -race -v -run $(TEST_NAME) ./... -count 1

generate-coverage:
	$(DC) $(GO) test --short -race -v ./... -count 1 -coverprofile cover.out . && go tool cover -html=cover.out -o cover.html && echo 'open in the browser the cover.html to see you branch coverage details'

lint-full: .golangci-lint-full ## Run golangci-lint

run-docker:
	chmod +x run-docker.sh && bash run-docker.sh
