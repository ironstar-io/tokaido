package dockertmpl

import (
	"log"
	"runtime"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	homedir "github.com/mitchellh/go-homedir"
)

func calcPhpVersionString(version string) string {
	var v string
	switch version {
	case "7.2":
		v = "72"
	case "7.3":
		v = "73"
	case "7.4":
		v = "74"
	default:
		log.Fatalf("PHP version %s is not supported. Must use 7.2, 7.3, or 7.4", version)
	}

	return v
}

// DrupalSettings ...
func DrupalSettings(drupalRoot string, projectName string) []byte {
	if conf.GetConfig().Global.Syncservice == "fusion" {
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
	// Return without kishu for Docker Volume mounts
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
        PROJECT_NAME: ` + projectName)
}

// ImageVersion ...
func ImageVersion(phpVersion, stability string) []byte {
	v := calcPhpVersionString(phpVersion)

	imageVersion := constants.EdgeVersion
	if stability == "stable" {
		imageVersion = constants.StableVersion
	}

	if conf.GetConfig().Global.Syncservice == "fusion" {
		return []byte(`services:
  sync:
    image: tokaido/sync:` + imageVersion + `
  syslog:
    image: tokaido/syslog:` + imageVersion + `
  nginx:
    image: tokaido/nginx:` + imageVersion + `
  fpm:
    image: tokaido/php` + v + `:` + imageVersion + `
  drush:
    image: tokaido/admin` + v + `-heavy:` + imageVersion + ``)
	}

	return []byte(`services:
  syslog:
    image: tokaido/syslog:` + imageVersion + `
  nginx:
    image: tokaido/nginx:` + imageVersion + `
  fpm:
    image: tokaido/php` + v + `:` + imageVersion + `
  drush:
    image: tokaido/admin` + v + `-heavy:` + imageVersion + ``)
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
      - /opt/solr/server/solr/configsets/search-api-solr/
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name)
}

// EnableRedis ...
func EnableRedis(version string) []byte {
	return []byte(`services:
  redis:
    image: redis:` + version + `
    ports:
      - "6379"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name)
}

// EnableMailhog ...
func EnableMailhog(version string) []byte {
	return []byte(`services:
  mailhog:
    image: mailhog/mailhog:` + version + `
    ports:
      - "1025"
      - "8025"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name)
}

// EnableAdminer ...
func EnableAdminer(version string) []byte {
	return []byte(`services:
  adminer:
    image: adminer:` + version + `
    ports:
      - "8080"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name)
}

// EnableMemcache ...
func EnableMemcache(version string) []byte {
	return []byte(`services:
  memcache:
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
    image: memcached:` + version)
}

// EnableChromedriver ...
func EnableChromedriver() []byte {
	return []byte(`services:
  chromedriver:
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
    image: drupalci/chromedriver:production
    ports:
      - "9515"`)
}

// ExternalVolumeDeclare ...
func ExternalVolumeDeclare(name string) []byte {
	return []byte(`volumes:
  ` + name + `:
    external: true
`)
}

// InternalVolumeDeclare ...
func InternalVolumeDeclare(name string) []byte {
	return []byte(`volumes:
  ` + name + `:
    external: false
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

// SetDatabase sets the database engine and configuration
func SetDatabase(image, version string) []byte {
	return []byte(`services:
  mysql:
    image: ` + image + `:` + version + `
`)
}

// SetDatabasePort assigns a static local port for the database
func SetDatabasePort(port string) []byte {
	return []byte(`services:
  mysql:
    ports:
      - ` + port + `:3306`)
}

// SetUnisonVersion ...
func SetUnisonVersion(version string) []byte {
	return []byte(`services:
  unison:
    image: tokaido/unison:` + version)
}

// TokaidoFusionSiteVolumeAttach ...
func TokaidoFusionSiteVolumeAttach(path, name string) []byte {
	return []byte(`services:
  sync:
    volumes:
      - ` + path + `:/app/host-volume
      - ` + name + `:/app/site
  drush:
    volumes:
      - ` + name + `:/app/site
      - tok_composer_cache:/home/tok/.composer/cache
  nginx:
    volumes:
      - ` + name + `:/app/site
  testcafe:
    volumes:
      - ` + path + `/.tok/testcafe:/testcafe
  fpm:
    volumes:
      - ` + name + `:/app/site
  kishu:
    volumes:
      - ` + name + `:/app/site
`)
}

// TokaidoDockerSiteVolumeAttach ...
func TokaidoDockerSiteVolumeAttach(path string) []byte {
	h, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Could not resolve your home directory: %v", err)
	}

	// diskMode is set to ":cached" and appended to /app/site mounts
	// this improves osxfs performance by about 50%
	diskMode := ""
	if runtime.GOOS == "darwin" {
		diskMode = ":cached"
	}

	// use the tokaido proxy tls wildcard certificate
	tlsPath := h + "/.tok/tls/proxy/"

	vols := `services:
  nginx:
    volumes:
      - ` + path + `:/app/site` + diskMode + `
      - ` + tlsPath + `wildcard.crt:/app/config/nginx/runtime/tls/default.crt
      - ` + tlsPath + `wildcard.key:/app/config/nginx/runtime/tls/default.key
  fpm:
    volumes:
      - ` + path + `:/app/site` + diskMode + `
  testcafe:
    volumes:
      - ` + path + `/.tok/testcafe:/testcafe` + diskMode + `
  drush:
    volumes:
      - ` + path + `:/app/site` + diskMode + `
      - tok_composer_cache:/home/tok/.composer/cache`

	// We'll mount the .gitconfig and .drush paths if they exist
	gp := h + "/.gitconfig"
	dp := h + "/.drush"

	if fs.CheckExists(gp) {
		vols = vols + `
      - ` + gp + `:/home/tok/.gitconfig`
	}

	if fs.CheckExists(dp) {
		vols = vols + `
      - ` + dp + `:/home/tok/.drush`
	}

	return []byte(vols)
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

// ComposeTokDefaultsFusionSync - Template byte array for `docker-compose.tok.yml`
var ComposeTokDefaultsFusionSync = []byte(`
version: "2"
services:
  sync:
    image: tokaido/sync:stable
    volumes:
      - waiting
    environment:
      AUTO_SYNC: "true"
    restart: unless-stopped
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  syslog:
    image: tokaido/syslog:stable
    volumes:
      - /app/logs
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  nginx:
    user: "1002"
    image: tokaido/nginx:stable
    volumes:
      - waiting
    volumes_from:
      - syslog
    depends_on:
      - fpm
    ports:
      - "8082"
      - "8443"
    environment:
      DRUPAL_ROOT: docroot
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  testcafe:
    image: testcafe/testcafe
    working_dir: /testcafe
    command: tail -f /dev/null
    entrypoint:
      - tail
      - -f
      - /dev/null
    volumes:
      - waiting
    depends_on:
      - nginx
    ports:
      - "1337"
  fpm:
    user: "1001"
    image: tokaido/php71-fpm:stable
    working_dir: /app/site/
    volumes:
      - waiting
    volumes_from:
      - syslog
    depends_on:
      - syslog
    ports:
      - "9000"
    environment:
      PHP_DISPLAY_ERRORS: "yes"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  mysql:
    image: mysql:5.7
    volumes_from:
      - syslog
    volumes:
      - waiting
    ports:
      - "3306"
    command: --max_allowed_packet=1073741824 --ignore-db-dir=lost+found --bind-address=0.0.0.0
    environment:
      MYSQL_DATABASE: tokaido
      MYSQL_USER: tokaido
      MYSQL_PASSWORD: tokaido
      MYSQL_ROOT_PASSWORD: tokaido
      MYSQL_ROOT_HOST: "%"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  drush:
    image: tokaido/admin71-heavy:stable
    hostname: 'tokaido'
    ports:
      - "22"
    working_dir: /app/site
    volumes:
      - waiting
    volumes_from:
      - syslog
    environment:
      SSH_AUTH_SOCK: /ssh/auth/sock
      APP_ENV: local
      PROJECT_NAME: tokaido
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  kishu:
    image: tokaido/kishu:stable
    volumes:
      - waiting
    environment:
      DRUPAL_ROOT: docroot
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
`)

// ComposeTokDefaultsDockerVolume - Template byte array for `docker-compose.tok.yml`
var ComposeTokDefaultsDockerVolume = []byte(`
version: "2"
services:
  syslog:
    image: tokaido/syslog:stable
    volumes:
      - /app/logs
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  nginx:
    user: "1002"
    image: tokaido/nginx:stable
    volumes:
      - waiting
    volumes_from:
      - syslog
    depends_on:
      - fpm
    ports:
      - "8082"
      - "8443"
    environment:
      DRUPAL_ROOT: docroot
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  testcafe:
    image: testcafe/testcafe
    working_dir: /testcafe
    command: tail -f /dev/null
    entrypoint:
      - tail
      - -f
      - /dev/null
    volumes:
      - waiting
    depends_on:
      - nginx
    ports:
      - "1337"
  fpm:
    user: "1001"
    image: tokaido/php71-fpm:stable
    working_dir: /app/site/
    volumes:
      - waiting
    volumes_from:
      - syslog
    depends_on:
      - syslog
    ports:
      - "9000"
    environment:
      PHP_DISPLAY_ERRORS: "yes"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  mysql:
    image: mysql:5.7
    volumes_from:
      - syslog
    volumes:
      - waiting
    ports:
      - "3306"
    command: --max_allowed_packet=1073741824 --ignore-db-dir=lost+found --bind-address=0.0.0.0
    environment:
      MYSQL_DATABASE: tokaido
      MYSQL_USER: tokaido
      MYSQL_PASSWORD: tokaido
      MYSQL_ROOT_PASSWORD: tokaido
      MYSQL_ROOT_HOST: "%"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  drush:
    image: tokaido/admin71-heavy:stable
    hostname: 'tokaido'
    ports:
      - "22"
    working_dir: /app/site
    volumes:
      - waiting
    volumes_from:
      - syslog
    environment:
      SSH_AUTH_SOCK: /ssh/auth/sock
      APP_ENV: local
      PROJECT_NAME: tokaido
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
`)

// ComposeTokDefaultsUnison - Template byte array for `docker-compose.tok.yml`
var ComposeTokDefaultsUnison = []byte(`
version: "2"
services:
  unison:
    image: tokaido/unison:2.51.2
    environment:
      - UNISON_DIR=/app/site
      - UNISON_UID=1001
      - UNISON_GID=1001
    ports:
      - "5000"
    volumes:
      - /app/site
  syslog:
    image: tokaido/syslog:stable
    volumes:
      - /app/logs
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  nginx:
    user: "1002"
    image: tokaido/nginx:stable
    volumes_from:
      - syslog
      - unison
    depends_on:
      - fpm
    ports:
      - "8082"
      - "8443"
    environment:
      DRUPAL_ROOT: docroot
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  testcafe:
    image: testcafe/testcafe
    working_dir: /testcafe
    command: tail -f /dev/null
    entrypoint:
      - tail
      - -f
      - /dev/null
    volumes_from:
      - unison
    depends_on:
      - nginx
    ports:
      - "1337"
  fpm:
    user: "1001"
    image: tokaido/php71-fpm:stable
    working_dir: /app/site/
    volumes_from:
      - syslog
      - unison
    depends_on:
      - syslog
    ports:
      - "9000"
    environment:
      PHP_DISPLAY_ERRORS: "yes"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  mysql:
    image: mysql:5.7
    volumes_from:
      - syslog
    volumes:
      - waiting
    ports:
      - "3306"
    command: --max_allowed_packet=1073741824 --ignore-db-dir=lost+found --bind-address=0.0.0.0
    environment:
      MYSQL_DATABASE: tokaido
      MYSQL_USER: tokaido
      MYSQL_PASSWORD: tokaido
      MYSQL_ROOT_PASSWORD: tokaido
      MYSQL_ROOT_HOST: "%"
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  drush:
    image: tokaido/admin71-heavy:stable
    hostname: 'tokaido'
    ports:
      - "22"
    working_dir: /app/site
    volumes:
      - waiting
    volumes_from:
      - syslog
      - unison
    environment:
      SSH_AUTH_SOCK: /ssh/auth/sock
      APP_ENV: local
      PROJECT_NAME: tokaido
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
`)
