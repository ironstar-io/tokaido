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
func CreateSyncService(syncName, syncDir string) {
	s := goos.NewUnisonSvc(syncName, syncDir)
	s.CreateSyncService()
}

// StopSyncService stop *and* remove the systemd sync service
func StopSyncService(syncName string) {
	if conf.GetConfig().System.Syncsvc.Enabled != true {
		return
	}

	s := goos.NewUnisonSvc(syncName, "")
	s.StopSyncService()
}

// SyncServiceStatus ...
func SyncServiceStatus(syncName string) string {
	s := goos.NewUnisonSvc(syncName, "")
	return s.SyncServiceStatus()
}

// CheckSyncService a verbose sync status check used for tok status
func CheckSyncService(syncName string) {
	s := goos.NewUnisonSvc(syncName, "")
	s.CheckSyncService()
}
