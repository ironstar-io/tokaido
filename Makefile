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

######################################################
## DOCKER IMAGES
## Official Tokaido docker images should be built by GitHub Actions
## These commands are for convenience of building and testing locally
######################################################

##################
# BASE
##################

.PHONY: base
base: ## Build and release the base image
	@$(MAKE) build-base
	@$(MAKE) push-base

.PHONY: build-base
build-base: ## Build and release the base image
	@$(MAKE) build-base-arm64

.PHONY: build-base-arm64
build-base-arm64:
	cd images/arm64/base && docker buildx build --platform linux/arm64 . ${DOCKER_CLI_FLAGS} --build-arg TOK_VERSION=${TOK_VERSION} -t tokaido/base:${TOK_VERSION}-arm64

.PHONY: push-base
push-base: ## push and release the base image
	@$(MAKE) push-base-arm64

.PHONY: push-base-arm64
push-base-arm64:
	docker push tokaido/base:${TOK_VERSION}-arm64

##################
# PHP
##################

.PHONY: php
php: ## Build and release the php images
	@$(MAKE) build-php
	@$(MAKE) push-php

.PHONY: build-php
build-php: ## Build and release the base image
	@$(MAKE) build-php74-arm64
	@$(MAKE) build-php80-arm64
	@$(MAKE) build-php81-arm64

.PHONY: build-php74-arm64
build-php74-arm64:
	cd images/arm64/php74 && docker buildx build --push --platform linux/arm64 . ${DOCKER_CLI_FLAGS} --build-arg TOK_VERSION=${TOK_VERSION} -t tokaido/php74:${TOK_VERSION}-arm64

.PHONY: build-php80-arm64
build-php80-arm64:
	cd images/arm64/php80 && docker buildx build --push --platform linux/arm64 . ${DOCKER_CLI_FLAGS} --build-arg TOK_VERSION=${TOK_VERSION} -t tokaido/php80:${TOK_VERSION}-arm64

.PHONY: build-php81-arm64
build-php81-arm64:
	cd images/arm64/php81 && docker buildx build --push --platform linux/arm64 . ${DOCKER_CLI_FLAGS} --build-arg TOK_VERSION=${TOK_VERSION} -t tokaido/php81:${TOK_VERSION}-arm64

.PHONY: push-php
push-php: ## push and release the php image
	@$(MAKE) push-php74-arm64
	@$(MAKE) push-php80-arm64
	@$(MAKE) push-php81-arm64

.PHONY: push-php74-arm64
push-php74-arm64:
	docker push tokaido/php74:${TOK_VERSION}-arm64

.PHONY: push-php80-arm64
push-php80-arm64:
	docker push tokaido/php80:${TOK_VERSION}-arm64

.PHONY: push-php81-arm64
push-php81-arm64:
	docker push tokaido/php81:${TOK_VERSION}-arm64

##################
# NGINX
##################

.PHONY: nginx
nginx: ## Build and release the nginx image
	@$(MAKE) build-nginx
	@$(MAKE) push-nginx

.PHONY: build-nginx
build-nginx: ## Build and release the nginx image
	@$(MAKE) build-nginx-arm64

.PHONY: build-nginx-arm64
build-nginx-arm64:
	echo "Building tokaido/nginx:${TOK_VERSION}-arm64"
	cd images/arm64/nginx && docker buildx build --push --platform linux/arm64 . ${DOCKER_CLI_FLAGS} --build-arg TOK_VERSION=${TOK_VERSION} -t tokaido/nginx:${TOK_VERSION}-arm64

.PHONY: push-nginx
push-nginx: ## push and release the nginx image
	@$(MAKE) push-nginx-arm64

.PHONY: push-nginx-arm64
push-nginx-arm64:
	docker push tokaido/nginx:${TOK_VERSION}-arm64

##################
# SSH
##################

.PHONY: ssh
ssh: ## Build and release the ssh images
	@$(MAKE) build-ssh
	@$(MAKE) push-ssh

.PHONY: build-ssh
build-ssh: ## Build and release the base image
	@$(MAKE) build-ssh74-arm64
	@$(MAKE) build-ssh80-arm64
	@$(MAKE) build-ssh81-arm64

.PHONY: build-ssh74-arm64
build-ssh74-arm64:
	cd images/arm64/ssh && docker buildx build --platform linux/arm64 --push . ${DOCKER_CLI_FLAGS} --build-arg TOK_VERSION=${TOK_VERSION} --build-arg PHP_VERSION=74 -t tokaido/ssh74:${TOK_VERSION}-arm64

.PHONY: build-ssh80-arm64
build-ssh80-arm64:
	cd images/arm64/ssh && docker buildx build --platform linux/arm64 --push . ${DOCKER_CLI_FLAGS} --build-arg TOK_VERSION=${TOK_VERSION} --build-arg PHP_VERSION=80 -t tokaido/ssh80:${TOK_VERSION}-arm64

.PHONY: build-ssh81-arm64
build-ssh81-arm64:
	cd images/arm64/ssh && docker buildx build --platform linux/arm64 --push . ${DOCKER_CLI_FLAGS} --build-arg TOK_VERSION=${TOK_VERSION} --build-arg PHP_VERSION=81 -t tokaido/ssh81:${TOK_VERSION}-arm64
