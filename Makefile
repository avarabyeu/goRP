.DEFAULT_GOAL := build
BUILD_DATE = `date +%FT%T%z`
GO = go
BINARY_DIR=bin

GODIRS_NOVENDOR = $(shell go list ./... | grep -v /vendor/)
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
BUILD_INFO_LDFLAGS=-ldflags "-extldflags '"-static"' -X main.buildDate=${BUILD_DATE} -X main.version=${v}"

.PHONY: test build

help:
	@echo "build      - go build"
	@echo "test       - go test"
	@echo "checkstyle - gofmt+golint+misspell"

init-deps:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.57.2

#vendor:
#	dep ensure --vendor-only

test:
	$(GO) test -cover ${GODIRS_NOVENDOR}

lint:
	bin/golangci-lint run ./...

fmt:
	gofumpt -extra -l -w ${GOFILES_NOVENDOR}
	gci write --section Standard --section Default --section "Prefix(github.com/reportportal/goRP/v5)" ${GOFILES_NOVENDOR}

#build: checkstyle test
build:
	$(GO) build ${BUILD_INFO_LDFLAGS} -o ${BINARY_DIR}/gorp ./

cross-build:
	gox ${BUILD_INFO_LDFLAGS} -arch="amd64 arm64" -os="linux windows darwin" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}" ./cmd/gorp

clean:
	if [ -d ${BINARY_DIR} ] ; then rm -r ${BINARY_DIR} ; fi
	if [ -d 'build' ] ; then rm -r 'build' ; fi

tag:
	git tag -a v${v} -m "creating tag ${v}"
	git push origin "refs/tags/v${v}"

release:
	rm -rf dist
	goreleaser release

grpc-gen:
	#Learn here: https://jbrandhorst.com/post/go-protobuf-tips/
	protoc -I=. -I=vendor -I=${GOPATH}/src model/*.proto --go_out=plugins=grpc:.
