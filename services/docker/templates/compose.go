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
	case "7.4":
		v = "74"
	case "8.0":
		v = "80"
	case "8.1":
		v = "81"
	default:
		log.Fatalf("PHP version '%s' is not supported. Please use 7.4, 8.0, 8.1", version)
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
    ssh:
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
    ssh:
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

	return []byte(`services:
  nginx:
    image: tokaido/nginx:` + imageVersion + `
  fpm:
    image: tokaido/php` + v + `:` + imageVersion + `
  ssh:
    image: tokaido/ssh` + v + `:` + imageVersion + ``)
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

// LogsVolumeAttach ...
func LogsVolumeAttach(name string) []byte {
	return []byte(`services:
  ssh:
    volumes:
      - ` + name + `:/app/logs
  fpm:
    volumes:
      - ` + name + `:/app/logs
  nginx:
    volumes:
      - ` + name + `:/app/logs
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

	logsVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_logs"

	vols := `services:
  nginx:
    volumes:
      - ` + path + `:/app/site` + diskMode + `
      - ` + tlsPath + `wildcard.crt:/app/config/nginx/runtime/tls/default.crt
      - ` + tlsPath + `wildcard.key:/app/config/nginx/runtime/tls/default.key
      - ` + logsVolName + `:/app/logs
  fpm:
    volumes:
      - ` + path + `:/app/site` + diskMode + `
      - ` + logsVolName + `:/app/logs
  ssh:
    volumes:
      - ` + path + `:/app/site` + diskMode + `
      - ` + logsVolName + `:/app/logs
      - tok_composer_cache:/home/app/.composer/cache`

	// We'll mount the .gitconfig and .drush paths if they exist
	gp := h + "/.gitconfig"
	dp := h + "/.drush"

	if fs.CheckExists(gp) {
		vols = vols + `
      - ` + gp + `:/home/app/.gitconfig`
	}

	if fs.CheckExists(dp) {
		vols = vols + `
      - ` + dp + `:/home/app/.drush`
	}

	return []byte(vols)
}

// ComposerCacheVolumeAttach ...
func ComposerCacheVolumeAttach() []byte {
	return []byte(`services:
  ssh:
    volumes:
      - tok_composer_cache:/home/app/.composer/cache
`)
}

// ModWarning - Displayed at the top of `docker-compose.tok.yml`
var ModWarning = []byte(`
# WARNING: THIS FILE IS MANAGED DIRECTLY BY TOKAIDO.
# DO NOT MAKE MODIFICATIONS HERE, THEY WILL BE OVERWRITTEN

`)

// ComposeTokDefaultsDockerVolume - Template byte array for `docker-compose.tok.yml`
var ComposeTokDefaultsDockerVolume = []byte(`
version: "3"
services:
  nginx:
    user: "1002"
    image: tokaido/nginx:stable
    volumes:
      - waiting
    depends_on:
      - fpm
    ports:
      - "8082"
      - "8443"
    environment:
      DRUPAL_ROOT: docroot
    depends_on:
      - ssh
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  fpm:
    user: "1001"
    image: tokaido/php81:stable
    working_dir: /app/site/
    volumes:
      - waiting
    ports:
      - "9000"
    environment:
      PHP_DISPLAY_ERRORS: "yes"
    depends_on:
      - ssh
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
  mysql:
    image: mariadb:10.8
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
  ssh:
    image: tokaido/ssh74:stable
    hostname: 'tokaido'
    ports:
      - "22"
    working_dir: /app/site
    volumes:
      - waiting
    environment:
      SSH_AUTH_SOCK: /ssh/auth/sock
      APP_ENV: local
      PROJECT_NAME: tokaido
    labels:
      io.tokaido.managed: local
      io.tokaido.project: ` + conf.GetConfig().Tokaido.Project.Name + `
`)
