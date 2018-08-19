package proxy

import (
	"github.com/ironstar-io/tokaido/constants"
)

const proxyNetwork = "proxy_proxy"

// GenerateNginxConf ...
func GenerateNginxConf(domain, proxyPassDomain string) []byte {
	return []byte(`server {
  listen          ` + constants.ProxyPort + ` ssl;
  server_name     ` + domain + `;
  server_tokens   off;

  ssl_certificate           /tokaido/proxy/config/client/tls/tokaido.pem;
  ssl_certificate_key       /tokaido/proxy/config/client/tls/tokaido-key.pem;

  location / {
    proxy_pass ` + proxyPassDomain + `;
  }
}
`)
}
