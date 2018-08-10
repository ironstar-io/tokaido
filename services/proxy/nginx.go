package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/fs"

	"path/filepath"
	"strings"
)

const proxyNetwork = "proxy_proxy"

// RebuildNginxConfigFile ...
func RebuildNginxConfigFile() {
	pn := conf.GetConfig().Tokaido.Project.Name
	ng := docker.GetGateway(proxyNetwork)
	nc := generateNginxConf(pn, strings.Replace(ng, `"`, "", -1))
	np := filepath.Join(getProxyClientConfdDir(), pn+".conf")

	fs.Replace(np, nc)
}

// generateNginxConf ...
func generateNginxConf(projectName, networkGateway string) []byte {
	return []byte(`server {
  listen          5154 ssl;
  server_name     ` + projectName + `.tokaido.local;
  server_tokens   off;

  ssl_certificate           /tokaido/proxy/config/client/tls/tokaido.pem;
  ssl_certificate_key       /tokaido/proxy/config/client/tls/tokaido-key.pem;

  location / {
    proxy_pass https://` + networkGateway + `:8443;
  }
}
`)
}
