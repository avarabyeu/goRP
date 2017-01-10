.DEFAULT_GOAL := build

GO = go
BINARY_DIR=bin
BINARY=${BINARY_DIR}/goRP
BUILD_DEPS:= github.com/golang/lint/golint \
             github.com/client9/misspell/cmd/misspell


help:
	@echo "build      - go build"
	@echo "test       - go test"
	@echo "checkstyle - gofmt+golint+misspell"

get-build-deps:
	$(GO) get -u $(BUILD_DEPS)

get-deps: get-build-deps
	$(GO) get ./...

test: get-deps
	$(GO) test -v ./...

checkstyle: test
	./checkstyle.sh

# Builds the project
build: test checkstyle
	$(GO) build -o ${BINARY}

clean:
	if [ -d ${BINARY_DIR} ] ; then rm -r ${BINARY_DIR} ; fi
