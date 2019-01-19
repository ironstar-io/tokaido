package unison

import (
	"fmt"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/unison/goos"
	"github.com/ironstar-io/tokaido/system/console"
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

// StopSyncService stop systemd sync service
func StopSyncService(syncName string) {
	if conf.GetConfig().System.Syncsvc.Enabled != true {
		return
	}

	s := goos.NewUnisonSvc(syncName, "")
	s.StopSyncService()
}

// UnloadSyncService will uninstall the sync service
func UnloadSyncService(syncName string) {
	if conf.GetConfig().System.Syncsvc.Enabled != true {
		return
	}

	s := goos.NewUnisonSvc(syncName, "")
	s.UnloadSyncService()
}

// SyncServiceStatus ...
func SyncServiceStatus(syncName string) string {
	s := goos.NewUnisonSvc(syncName, "")
	return s.SyncServiceStatus()
}

// CheckSyncService a verbose sync status check used for tok status
func CheckSyncService(syncName string) error {
	s := goos.NewUnisonSvc(syncName, "")
	return s.CheckSyncService()
}

// BackgroundServiceWarning - Check if the background sync service is running and warn if not
func BackgroundServiceWarning(syncName string) {
	if conf.GetConfig().System.Syncsvc.Enabled != true {
		return
	}

	s := goos.NewUnisonSvc(syncName, "")
	if s.SyncServiceStatus() != "running" {
		fmt.Println()
		console.Println("⚠️  The background sync service is not running. Manually sync with `tok watch` or try running `tok up` again", "")
		fmt.Println()
	}
}
