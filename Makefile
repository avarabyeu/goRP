.DEFAULT_GOAL := build

GO = go
BINARY_DIR=bin
BINARY=${BINARY_DIR}/goRP

help:
	@echo "build      - go build"
	@echo "test       - go test"
	@echo "checkstyle - gofmt+golint+misspell"


test:
	go test -v ./...

checkstyle: test
	./checkstyle.sh

# Builds the project
build: test checkstyle
	$(GO) build -o ${BINARY}

clean:
	if [ -d ${BINARY_DIR} ] ; then rm -r ${BINARY_DIR} ; fi
