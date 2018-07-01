package dockertmpl

// ComposeDotTok ...
type ComposeDotTok struct {
	Version  string
	Services struct {
		Unison struct {
			Image       string
			Environment []string
			Ports       []string
			Volumes     []string
		}
		Syslog struct {
			Image   string
			Volumes []string
		}
		Haproxy struct {
			User      string
			Image     string
			Ports     []string
			DependsOn []string `yaml:"depends_on"`
		}
		Varnish struct {
			User        string
			Image       string
			DependsOn   []string `yaml:"depends_on"`
			VolumesFrom []string `yaml:"volumes_from"`
		}
		Nginx struct {
			User        string
			Image       string
			VolumesFrom []string `yaml:"volumes_from"`
			DependsOn   []string `yaml:"depends_on"`
			Ports       []string
		}
		Fpm struct {
			User        string
			Image       string
			WorkingDir  string   `yaml:"working_dir"`
			VolumesFrom []string `yaml:"volumes_from"`
			DependsOn   []string `yaml:"depends_on"`
			Ports       []string
			Environment map[string]string
		}
		Memcache struct {
			Image string
		}
		Mysql struct {
			Image       string
			VolumesFrom []string `yaml:"volumes_from"`
			Ports       []string
			Environment map[string]string
		}
		Drush struct {
			Image       string
			Hostname    string
			Ports       []string
			WorkingDir  string   `yaml:"working_dir"`
			VolumesFrom []string `yaml:"volumes_from"`
			Environment map[string]string
		}
		Solr struct {
			Image      string
			Ports      []string
			Entrypoint []string
		}
	}
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
    image: onnimonni/unison
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
  memcache:
    image: memcached:1.5-alpine
  mysql:
    image: mysql:5.7
    volumes_from:
      - syslog
    ports:
      - "3306"
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
  solr:
    image: tokaido/solr:6.6
    ports:
      - "8983"
    entrypoint:
      - solr-precreate
      - drupal
      - /opt/solr/server/solr/configsets/search-api-solr/
`)
