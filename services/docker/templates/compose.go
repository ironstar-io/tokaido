package dockertmpl

import "log"

func calcPhpVersionString(version string) string {
	var v string
	switch version {
	case "7.1":
		v = "71"
	case "7.2":
		v = "72"
	default:
		log.Fatalf("PHP version %s is not supported. Must use '7.1' or '7.2'", version)
	}

	return v
}

// DrupalSettings ...
func DrupalSettings(drupalRoot string, projectName string) []byte {
	return []byte(`services:
  fpm:
    environment:
      DRUPAL_ROOT: ` + drupalRoot + `
  nginx:
    environment:
      DRUPAL_ROOT: ` + drupalRoot + `
  drush:
    environment:
      DRUPAL_ROOT: ` + drupalRoot + `
      PROJECT_NAME: ` + projectName + `
  kishu:
      environment:
        DRUPAL_ROOT: ` + drupalRoot)
}

// StabilityLevel ...
func StabilityLevel(phpVersion, stability string) []byte {
	v := calcPhpVersionString(phpVersion)

	return []byte(`services:
  syslog:
    image: tokaido/syslog:` + stability + `
  haproxy:
    image: tokaido/haproxy:` + stability + `
  varnish:
    image: tokaido/varnish:` + stability + `
  nginx:
    image: tokaido/nginx:` + stability + `
  fpm:
    image: tokaido/php` + v + `-fpm:` + stability + `
  drush:
    image: tokaido/admin` + v + `-heavy:` + stability + ``)
}

// EnableSolr ...
func EnableSolr(version string) []byte {
	return []byte(`services:
  solr:
    image: tokaido/solr:` + version + `
    ports:
      - "8983"
    entrypoint:
      - solr-precreate
      - drupal
      - /opt/solr/server/solr/configsets/search-api-solr/`)
}

// EnableRedis ...
func EnableRedis(version string) []byte {
	return []byte(`services:
  redis:
    image: redis:` + version + `
    ports:
      - "6379"`)
}

// EnableMailhog ...
func EnableMailhog(version string) []byte {
	return []byte(`services:
  mailhog:
    image: mailhog/mailhog:` + version + `
    ports:
      - "1025"
      - "8025"`)
}

// EnableAdminer ...
func EnableAdminer(version string) []byte {
	return []byte(`services:
  adminer:
    image: adminer:` + version + `
    ports:
      - "8080"`)
}

// SetUnisonVersion ...
func SetUnisonVersion(version string) []byte {
	return []byte(`services:
  unison:
    image: tokaido/unison:` + version)
}

// EnableMemcache ...
func EnableMemcache(version string) []byte {
	return []byte(`services:
  memcache:
    image: memcached:` + version)
}

// EnableXdebug ...
func EnableXdebug(phpVersion, xdebugImageVersion string) []byte {
	v := calcPhpVersionString(phpVersion)
	return []byte(`services:
  fpm:
    image: tokaido/php` + v + `-fpm-xdebug:` + xdebugImageVersion)
}

// WindowsAjustments ...
func WindowsAjustments() []byte {
	return []byte(`services:
  unison:
    image: onnimonni/unison:2.48.4`)
}

// ExternalVolumeDeclare ...
func ExternalVolumeDeclare(name string) []byte {
	return []byte(`volumes:
  ` + name + `:
    external: true
`)
}

// MysqlVolumeAttach ...
func MysqlVolumeAttach(name string) []byte {
	return []byte(`services:
  mysql:
    volumes:
      - ` + name + `:/var/lib/mysql
`)
}

// ComposerCacheVolumeAttach ...
func ComposerCacheVolumeAttach() []byte {
	return []byte(`services:
  drush:
    volumes:
      - tok_composer_cache:/home/tok/.composer/cache
`)
}

// ModWarning - Displayed at the top of `docker-compose.tok.yml`
var ModWarning = []byte(`
# WARNING: THIS FILE IS MANAGED DIRECTLY BY TOKAIDO.
# DO NOT MAKE MODIFICATIONS HERE, THEY WILL BE OVERWRITTEN

`)

// ComposeTokDefaults - Template byte array for `docker-compose.tok.yml`
var ComposeTokDefaults = []byte(`
version: "2"
services:
  unison:
    image: tokaido/unison:2.51.2
    environment:
      - UNISON_DIR=/tokaido/site
      - UNISON_UID=1001
      - UNISON_GID=1001
    ports:
      - "5000"
    volumes:
      - /tokaido/site
  syslog:
    image: tokaido/syslog:stable
    volumes:
      - /tokaido/logs
  haproxy:
    user: "1005"
    image: tokaido/haproxy:stable
    ports:
      - "8080"
      - "8443"
    depends_on:
      - varnish
      - nginx
  varnish:
    user: "1004"
    image: tokaido/varnish:stable
    ports:
      - "8081"
    depends_on:
      - nginx
    volumes_from:
      - syslog
  nginx:
    user: "1002"
    image: tokaido/nginx:stable
    volumes_from:
      - unison
      - syslog
    depends_on:
      - fpm
    ports:
      - "8082"
    environment:
      DRUPAL_ROOT: docroot
  fpm:
    user: "1001"
    image: tokaido/php71-fpm:stable
    working_dir: /tokaido/site/
    volumes_from:
      - unison
      - syslog
    depends_on:
      - syslog
    ports:
      - "9000"
    environment:
      PHP_DISPLAY_ERRORS: "yes"
  mysql:
    image: mysql:5.7
    volumes_from:
      - syslog
    volumes:
      - waiting
    ports:
      - "3306"
    command: --max_allowed_packet=67108864 --ignore-db-dir=lost+found
    environment:
      MYSQL_DATABASE: tokaido
      MYSQL_USER: tokaido
      MYSQL_PASSWORD: tokaido
      MYSQL_ROOT_PASSWORD: tokaido
  drush:
    image: tokaido/admin71-heavy:stable
    hostname: 'tokaido'
    ports:
      - "22"
    working_dir: /tokaido/site
    volumes_from:
      - unison
      - syslog
    environment:
      SSH_AUTH_SOCK: /ssh/auth/sock
      APP_ENV: local
      PROJECT_NAME: tokaido
  kishu:
    image: tokaido/kishu:stable
    volumes_from:
      - unison
    environment:
      DRUPAL_ROOT: docroot
`)
