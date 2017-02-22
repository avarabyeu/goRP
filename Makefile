.DEFAULT_GOAL := build

COMMIT_HASH = `git rev-parse --short HEAD 2>/dev/null`
BUILD_DATE = `date +%FT%T%z`

GO = go
BINARY_DIR=bin
BINARY=${BINARY_DIR}/goRP
BUILD_DEPS:= github.com/alecthomas/gometalinter
GODIRS_NOVENDOR = $(shell go list ./... | grep -v /vendor/)
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")


help:
	@echo "build      - go build"
	@echo "test       - go test"
	@echo "checkstyle - gofmt+golint+misspell"

vendor: ## Install govendor and sync Hugo's vendored dependencies
	go get github.com/kardianos/govendor
	govendor sync ${PACKAGE}

get-build-deps:
	$(GO) get $(BUILD_DEPS)
	gometalinter --install


get-deps: get-build-deps
	$(GO) get ./...

test: get-deps
	$(GO) test -v ${GODIRS_NOVENDOR}

#checkstyle: test
#	./checkstyle.sh

checkstyle:
	gometalinter ${GODIRS_NOVENDOR} --deadline 1m

fmt:
	gofmt -l -w ${GOFILES_NOVENDOR}

# Builds the project
build: test checkstyle
	$(GO) build -o ${BINARY}

clean:
	if [ -d ${BINARY_DIR} ] ; then rm -r ${BINARY_DIR} ; fi
