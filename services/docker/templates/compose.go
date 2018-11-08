package dockertmpl

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

// EdgeContainers ...
func EdgeContainers() []byte {
	return []byte(`services:
  syslog:
    image: tokaido/syslog:edge
  haproxy:
    image: tokaido/haproxy:edge
  varnish:
    image: tokaido/varnish:edge
  nginx:
    image: tokaido/nginx:edge
  fpm:
    image: tokaido/fpm:edge
  drush:
    image: tokaido/drush-heavy:edge`)
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
func EnableXdebug(version string) []byte {
	return []byte(`services:
  fpm:
    image: tokaido/fpm-xdebug:` + version)
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
      - tok_composer_cache:~/cache
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
    image: tokaido/syslog:latest
    volumes:
      - /tokaido/logs
  haproxy:
    user: "1005"
    image: tokaido/haproxy:latest
    ports:
      - "8080"
      - "8443"
    depends_on:
      - varnish
      - nginx
  varnish:
    user: "1004"
    image: tokaido/varnish:latest
    ports:
      - "8081"
    depends_on:
      - nginx
    volumes_from:
      - syslog
  nginx:
    user: "1002"
    image: tokaido/nginx:latest
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
    image: tokaido/fpm:latest
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
    image: tokaido/drush-heavy:latest
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
    image: tokaido/kishu:latest
    volumes_from:
      - unison
    environment:
      DRUPAL_ROOT: docroot
`)
