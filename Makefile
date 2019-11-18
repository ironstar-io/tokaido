VERSION_PATH ?= github.com/ironstar-io/tokaido/system/version
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
	" -o ./dist/tok-osx

build-installer:
	cd installer && make build-osx
	cd installer && make build-windows
	cd installer && make build-linux
	cd installer && make build-docker-images
	make build-osx && cp -R ./dist/tok-osx ./installer/dist/tokaido/tok-osx
	make build-linux && cp -R ./dist/tok-linux-amd64 ./installer/dist/tokaido/tok-linux-amd64
	make build-windows && cp -R ./dist/tok-windows-amd64.exe ./installer/dist/tokaido/tok-windows-amd64.exe
	cp -R ./installer/README.md ./installer/dist/README.md

test:
	ginkgo test ./...

clean:
	rm -rf ./dist/*

.PHONY: build build-windows build-linux build-osx test clean