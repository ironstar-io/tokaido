package proxy

import (
	"github.com/ironstar-io/tokaido/services/unison"
)

// ConfigureUnison ...
func ConfigureUnison() {
	unison.CreateOrUpdatePrf(UnisonPort(), proxy, getProxyClientDir())

	s := unison.SyncServiceStatus(proxy)
	if s == "stopped" {
		unison.Sync(proxy)
	}

	unison.CreateSyncService(proxy, getProxyClientDir())
}
