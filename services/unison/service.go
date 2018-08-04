package unison

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/unison/goos"
)

var bgSyncFailMsg = `
😓  The background sync service is not running

Tokaido will run, but your environment and local host will not be synchronised
Use 'tok up' to repair, or 'tok sync' to sync manually
		`

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func CreateSyncService() {
	goos.CreateSyncService()
}

// StopSyncService stop *and* remove the systemd sync service
func StopSyncService() {
	if conf.GetConfig().System.Syncsvc.Enabled != true {
		return
	}

	goos.StopSyncService()
}

// SyncServiceStatus ...
func SyncServiceStatus() string {
	return goos.SyncServiceStatus()
}

// CheckSyncService a verbose sync status check used for tok status
func CheckSyncService() {
	goos.CheckSyncService()
}
