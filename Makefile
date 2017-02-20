.DEFAULT_GOAL := build

GO = go
BINARY_DIR=bin
BINARY=${BINARY_DIR}/goRP
BUILD_DEPS:= github.com/golang/lint/golint \
             github.com/client9/misspell/cmd/misspell
GODIRS_NOVENDOR = $(shell go list ./... | grep -v /vendor/)
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")


help:
	@echo "build      - go build"
	@echo "test       - go test"
	@echo "checkstyle - gofmt+golint+misspell"

get-build-deps:
	$(GO) get $(BUILD_DEPS)

get-deps: get-build-deps
	$(GO) get ./...

test: get-deps
	$(GO) test -v ${GODIRS_NOVENDOR}

#checkstyle: test
#	./checkstyle.sh
checkstyle:
	@gofmt -l ${GOFILES_NOVENDOR} | read && echo "Code differs from gofmt's style" 1>&2 && exit 1 || true
	go vet ${GOPACKAGES}

fmt:
	gofmt -l -w ${GOFILES_NOVENDOR}

# Builds the project
build: test checkstyle
	$(GO) build -o ${BINARY}

clean:
	if [ -d ${BINARY_DIR} ] ; then rm -r ${BINARY_DIR} ; fi
