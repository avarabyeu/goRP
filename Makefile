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
	# installs gometalinter
	curl -L https://git.io/vp6lP | sh
#	gometalinter --install

#vendor:
#	dep ensure --vendor-only

test:
	$(GO) test -cover ${GODIRS_NOVENDOR}

checkstyle:
	gometalinter --vendor ./... --fast --deadline 10m

fmt:
	gofmt -l -w -s ${GOFILES_NOVENDOR}

#build: checkstyle test
build:
	$(GO) build ${BUILD_INFO_LDFLAGS} -o ${BINARY_DIR}/gorp ./

cross-build:
	gox ${BUILD_INFO_LDFLAGS} -arch="amd64 386" -os="linux windows darwin" -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

clean:
	if [ -d ${BINARY_DIR} ] ; then rm -r ${BINARY_DIR} ; fi
	if [ -d 'build' ] ; then rm -r 'build' ; fi

tag:
	git tag -a v${v} -m "creating tag ${v}"
	git push origin "refs/tags/v${v}"

release:
	rm -rf dist
	goreleaser release
