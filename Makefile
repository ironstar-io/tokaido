VERSION_PATH ?= bitbucket.org/ironstar/tokaido-cli/system/version
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= $(shell git describe --tags)

build:
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/tok

build-windows:
	env GOOS=windows GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/tok-windows-amd64.exe

build-linux:
	env GOOS=linux GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/tok-linux-amd64

build-osx:
	env GOOS=darwin GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/tok-macos

test:
	ginkgo test ./...

clean:
	rm -rf ./dist/*

.PHONY: build build-windows build-linux build-osx test clean