package proxy

import (
	"github.com/ironstar-io/tokaido/services/unison"
)

// configureUnison ...
func configureUnison() {
	unison.CreateOrUpdatePrf(unisonPort(), proxy, getProxyClientDir())

	s := unison.SyncServiceStatus(proxy)
	if s == "stopped" {
		unison.Sync(proxy)
	}

	unison.CreateSyncService(proxy, getProxyClientDir())

	unison.RestartSyncService(proxy)

}
