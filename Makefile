VERSION_PATH ?= bitbucket.org/ironstar/tokaido-cli/system/version
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= $(shell git describe --tags)

build:
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/tok

test:
	ginkgo test ./...

clean:
	rm -rf ./dist/tok

.PHONY: build test clean