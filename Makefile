.DEFAULT_GOAL := build

COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`

GO = go
BINARY_DIR=bin

BUILD_DEPS:= github.com/alecthomas/gometalinter
GODIRS_NOVENDOR = $(shell go list ./... | grep -v /vendor/)
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: vendor test build

help:
	@echo "build      - go build"
	@echo "test       - go test"
	@echo "checkstyle - gofmt+golint+misspell"

vendor: ## Install govendor and sync Hugo's vendored dependencies
	go get github.com/kardianos/govendor
	govendor sync

get-build-deps: vendor
	$(GO) get $(BUILD_DEPS)
	gometalinter --install

test: vendor
	govendor test +local

checkstyle: get-build-deps
	gometalinter --vendor ./... --fast --disable=gas --disable=errcheck --disable=gotype #--deadline 5m

fmt:
	gofmt -l -w -s ${GOFILES_NOVENDOR}

# Builds gorpRoot
build-app-root: checkstyle test
	CGO_ENABLED=0 GOOS=linux $(GO) build -o ${BINARY_DIR}/gorpRoot ./gorpRoot

# Builds gorpUI
build-app-ui: checkstyle test
	CGO_ENABLED=0 GOOS=linux $(GO) build -o ${BINARY_DIR}/gorpUI ./gorpUI

# Builds the project
build: build-app-root build-app-ui

# Builds containers
docker: build
	docker build -t gorproot -f gorpRoot/Dockerfile .
	docker build -t gorpui -f gorpUI/Dockerfile .

clean:
	if [ -d ${BINARY_DIR} ] ; then rm -r ${BINARY_DIR} ; fi
