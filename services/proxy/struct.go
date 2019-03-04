package proxy

import (
	"github.com/ironstar-io/tokaido/constants"
)

// DockerCompose ...
type DockerCompose struct {
	Version  string `yaml:"version"`
	Services struct {
		Proxy struct {
			Hostname    string   `yaml:"hostname,omitempty"`
			Entrypoint  []string `yaml:"entrypoint,omitempty"`
			User        string   `yaml:"user,omitempty"`
			Cmd         string   `yaml:"cmd,omitempty"`
			Dependson   []string `yaml:"depends_on,omitempty"`
			Environment []string `yaml:"environment,omitempty"`
			Volumes     []string `yaml:"volumes,omitempty"`
			Image       string   `yaml:"image"`
			Ports       []string `yaml:"ports"`
			Networks    []string `yaml:"networks"`
		} `yaml:"proxy"`
	} `yaml:"services"`
}

// ComposeDefaults - Template byte array for proxy `docker-compose.yml`
func ComposeDefaults() []byte {
	return []byte(`version: "2"
services:
  proxy:
    image: tokaido/proxy:latest
    ports:
      - "` + constants.ProxyPort + `:` + constants.ProxyPort + `"
    volumes:
      - ./client:/tokaido/proxy/config/client
    networks:
      - proxy
`)
}
