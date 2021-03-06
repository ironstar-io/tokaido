INSTALLER_PATH ?= github.com/ironstar-io/tokaido-installer/install
TOK_VERSION ?= 1.12.0
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= $(shell git describe --tags)

build:
	go build \
	-ldflags "\
	-X $(INSTALLER_PATH).buildDate=$(BUILD_DATE) \
	-X $(INSTALLER_PATH).version=$(VERSION) \
	-X $(INSTALLER_PATH).tokVersion=$(TOK_VERSION) \
	" -o ./dist/tok

build-all: build-linux build-macos build-windows

build-windows:
	env GOOS=windows GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(INSTALLER_PATH).buildDate=$(BUILD_DATE) \
	-X $(INSTALLER_PATH).version=$(VERSION) \
	-X $(INSTALLER_PATH).tokVersion=$(TOK_VERSION) \
	" -o ./dist/tok-installer-windows-amd64.exe

build-linux:
	env GOOS=linux GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(INSTALLER_PATH).buildDate=$(BUILD_DATE) \
	-X $(INSTALLER_PATH).version=$(VERSION) \
	-X $(INSTALLER_PATH).tokVersion=$(TOK_VERSION) \
	" -o ./dist/tok-installer-linux-amd64

build-macos:
	env GOOS=darwin GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(INSTALLER_PATH).buildDate=$(BUILD_DATE) \
	-X $(INSTALLER_PATH).version=$(VERSION) \
	-X $(INSTALLER_PATH).tokVersion=$(TOK_VERSION) \
	" -o ./dist/tok-installer-macos

build-docker-images:
	docker save loomchild/volume-backup > dist/images/volume-backup.tar
	docker save mysql:5.7 > dist/images/mysql.tar
	docker save testcafe/testcafe > dist/images/testcafe.tar
	docker save mailhog/mailhog > dist/images/mailhog.tar
	docker save tokaido/unison:2.51.2 > dist/images/tokaido-unison.tar
	docker save tokaido/syslog:$(TOK_VERSION) > dist/images/tokaido-syslog.tar
	docker save tokaido/nginx:$(TOK_VERSION) > dist/images/tokaido-nginx.tar
	docker save tokaido/php74-fpm:$(TOK_VERSION) > dist/images/tokaido-php74-fpm.tar
	docker save tokaido/admin74-heavy:$(TOK_VERSION) > dist/images/tokaido-admin74-heavy.tar
	docker save tokaido/kishu:$(TOK_VERSION) > dist/images/tokaido-kishu.tar
	docker save tokaido/proxy:1.11.0 > dist/images/tokaido-proxy.tar

.PHONY: build build-windows build-linux build-macos test clean
