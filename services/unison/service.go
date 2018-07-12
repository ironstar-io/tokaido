package unison

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/unison/goos"

	"fmt"
)

var bgSyncFailMsg = `
😓  The background sync service is not running

Tokaido will run, but your environment and local host will not be synchronised
Use 'tok up' to repair, or 'tok sync' to sync manually
		`

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func CreateSyncService() {
	fmt.Println("🔄  Creating a background process to sync your local repo into the Tokaido environment")

	goos.RegisterSyncService()
	goos.StartSyncService()
}

// StopSyncService stop *and* remove the systemd sync service
func StopSyncService() {
	if conf.GetConfig().CreateSyncService != true {
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
	c := conf.GetConfig()
	if c.CreateSyncService != true {
		return
	}

	s := goos.SyncServiceStatus()
	if s == "running" {
		fmt.Println("✅  Background sync service is running")
		return
	}

	fmt.Println(bgSyncFailMsg)
}
