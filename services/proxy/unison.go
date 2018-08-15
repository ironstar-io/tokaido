package proxy

import (
	"github.com/ironstar-io/tokaido/services/unison"
)

// ConfigureUnison ...
func ConfigureUnison() {
	unison.CreateOrUpdatePrf(UnisonPort(), proxy, getProxyClientDir())
	unison.CreateSyncService(proxy, getProxyClientDir())
}
