VERSION_PATH ?= github.com/ironstar-io/tokaido/system/version
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= $(shell git describe --tags)
TOK_VERSION := $(or $(TOK_VERSION),experimental)
TOK_ARCH := $(or $(TOK_ARCH),amd64)

build:
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/tok

build-all: build-macos build-windows build-linux

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

build-macos:
	env GOOS=darwin GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(VERSION_PATH).buildDate=$(BUILD_DATE) \
	-X $(VERSION_PATH).version=$(VERSION) \
	" -o ./dist/tok-macos

usb-installer:
	cd installer && make build-macos
	cd installer && make build-windows
	cd installer && make build-linux
	cd installer && make build-docker-images
	make build-macos && cp -R ./dist/tok-macos ./installer/dist/tokaido/tok-macos
	make build-linux && cp -R ./dist/tok-linux-amd64 ./installer/dist/tokaido/tok-linux-amd64
	make build-windows && cp -R ./dist/tok-windows-amd64.exe ./installer/dist/tokaido/tok-windows-amd64.exe
	cp -R ./installer/README.md ./installer/dist/README.md

test:
	ginkgo test ./...

clean:
	rm -rf ./dist/*

######
# Docker Images
# ####

.PHONY: images
	@$(MAKE) build-images-all
	@$(MAKE) push-images-all

##################
# BASE
# The 'base' image is the root image for most Tokaido containers. It sets up common users, permissions, and shared utilities
##################

.PHONY: base
base: ## Build and release the base image
	@$(MAKE) build-base

.PHONY: build-base
build-base:
	cd images/base && docker build . ${DOCKER_CLI_FLAGS} -t tokaido/base:${TOK_VERSION}

.PHONY: push-base
push-base:
	docker push tokaido/base:${TOK_VERSION}

##################
# PHP
# The 'php' images contain all the code needed for Drupal to run and are sometimes referred to as the "app" servers or instances
##################

.PHONY: php74
php74: ## Build and release the php74 image
	@$(MAKE) build-php74

.PHONY: build-php74
build-php74:
	cd images/php74 && docker buildx build --platform linux/amd64 . ${DOCKER_CLI_FLAGS} --build-arg TOK_VERSION=${TOK_VERSION} --build-arg LIBRAY_PATH=/usr/lib/x86_64-linux-gnu -t tokaido/php74:${TOK_VERSION}
	cd images/php74 && docker buildx build --platform linux/amd64 . ${DOCKER_CLI_FLAGS} --build-arg TOK_VERSION=${TOK_VERSION} --build-arg LIBRAY_PATH=/usr/lib/arm64-linux-gnu -t tokaido/php74:${TOK_VERSION}

.PHONY: push-php74
push-base:
	docker push tokaido/php74:${TOK_VERSION}






.PHONY: build build-windows build-linux build-macos test clean
