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

build-windows:
	env GOOS=windows GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(INSTALLER_PATH).buildDate=$(BUILD_DATE) \
	-X $(INSTALLER_PATH).version=$(VERSION) \
	-X $(INSTALLER_PATH).tokVersion=$(TOK_VERSION) \
	" -o ./dist/tok-windows-amd64.exe

build-linux:
	env GOOS=linux GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(INSTALLER_PATH).buildDate=$(BUILD_DATE) \
	-X $(INSTALLER_PATH).version=$(VERSION) \
	-X $(INSTALLER_PATH).tokVersion=$(TOK_VERSION) \
	" -o ./dist/tok-linux-amd64

build-macos:
	env GOOS=darwin GOARCH=amd64 \
	go build \
	-ldflags "\
	-X $(INSTALLER_PATH).buildDate=$(BUILD_DATE) \
	-X $(INSTALLER_PATH).version=$(VERSION) \
	-X $(INSTALLER_PATH).tokVersion=$(TOK_VERSION) \
	" -o ./dist/tok-macos

build-docker-images:
	docker save loomchild/volume-backup > dist/images/volume-backup.tar
	docker save mysql:5.7 > dist/images/mysql.tar
	docker save testcafe/testcafe > dist/images/testcafe.tar
	docker save mailhog/mailhog > dist/images/mailhog.tar
	docker save tokaido/syslog:1.11.0 > dist/images/tokaido-syslog.tar
	docker save tokaido/haproxy:1.11.0 > dist/images/tokaido-haproxy.tar
	docker save tokaido/varnish:1.11.0 > dist/images/tokaido-varnish.tar
	docker save tokaido/nginx:1.11.0 > dist/images/tokaido-nginx.tar
	docker save tokaido/php72-fpm:1.11.0 > dist/images/tokaido-php72-fpm.tar
	docker save tokaido/admin72-heavy:1.11.0 > dist/images/tokaido-admin72-heavy.tar
	docker save tokaido/unison:1.11.0 > dist/images/tokaido-unison.tar
	docker save tokaido/kishu:1.11.0 > dist/images/tokaido-kishu.tar

.PHONY: build build-windows build-linux build-macos test clean