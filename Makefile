.DEFAULT_GOAL := build

COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`

GO = go
BINARY_DIR=bin

BUILD_DEPS:= github.com/alecthomas/gometalinter
GODIRS_NOVENDOR = $(shell go list ./... | grep -v /vendor/)
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")


help:
	@echo "build      - go build"
	@echo "test       - go test"
	@echo "checkstyle - gofmt+golint+misspell"

vendor: ## Install govendor and sync Hugo's vendored dependencies
	go get github.com/kardianos/govendor
	govendor sync

get-build-deps:
	$(GO) get $(BUILD_DEPS)
	gometalinter --install

test: vendor
	govendor test +local

checkstyle test:
	gometalinter --vendor ./... --deadline 1m --disable=gas --disable=errcheck

fmt:
	gofmt -l -w ${GOFILES_NOVENDOR}

# Builds gorpRoot
build-app-root: test checkstyle
	$(GO) build -o ${BINARY_DIR}/gorpRoot ./gorpRoot

# Builds gorpUI
build-app-ui: test checkstyle
	$(GO) build -o ${BINARY_DIR}/gorpUI ./gorpUI

# Builds the project
build: build-app-root build-app-ui

clean:
	if [ -d ${BINARY_DIR} ] ; then rm -r ${BINARY_DIR} ; fi
