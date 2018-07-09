package unison

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system"
	"bitbucket.org/ironstar/tokaido-cli/system/linux"
	"bitbucket.org/ironstar/tokaido-cli/system/osx"

	"fmt"
)

var bgSyncFailMsg = `
😓  The background sync service is not running

Tokaido will run, but your environment and local host will not be synchronised
Use 'tok up' to repair, or 'tok sync' to sync manually
		`

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func CreateSyncService() {
	GOOS := system.CheckOS()

	fmt.Println("🔄  Creating a background process to sync your local repo into the Tokaido environment")

	if GOOS == "linux" {
		linux.RegisterSystemdService()
		linux.StartSystemdService()
	}

	if GOOS == "osx" {
		osx.RegisterLaunchdService()
		osx.StartLaunchdService()
	}
}

// func StartSyncService() -- not needed yet

// StopSyncService stop *and* remove the systemd sync service
func StopSyncService() {
	GOOS := system.CheckOS()
	c := conf.GetConfig()
	if c.CreateSyncService != true {
		return
	}

	if GOOS == "linux" {
		linux.StopSystemdService()
	}

	if GOOS == "osx" {
		osx.StopLaunchdService()
	}
}

// SyncServiceStatus ...
func SyncServiceStatus() string {
	GOOS := system.CheckOS()

	if GOOS == "linux" {
		return linux.CheckSyncService()
	}

	if GOOS == "osx" {
		return osx.CheckSyncService()
	}

	return ""
}

// CheckSyncService a verbose sync status check used for tok status
func CheckSyncService() {
	c := conf.GetConfig()
	if c.CreateSyncService != true {
		return
	}

	s := SyncServiceStatus()
	if s == "running" {
		fmt.Println("✅  Background sync service is running")
		return
	}

	fmt.Println(bgSyncFailMsg)
}
