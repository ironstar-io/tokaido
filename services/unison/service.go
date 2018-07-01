package unison

import (
	"fmt"

	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/linux"
	"bitbucket.org/ironstar/tokaido-cli/system/osx"
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func CreateSyncService() {
	GOOS := utils.CheckOS()

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
	GOOS := utils.CheckOS()
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

// CheckSyncService ...
func CheckSyncService() {
	GOOS := utils.CheckOS()
	c := conf.GetConfig()
	if c.CreateSyncService != true {
		return
	}

	if GOOS == "linux" {
		err := linux.CheckSyncService()
		if err == nil {
			fmt.Println("✅  Background sync service is running")
			return
		}

		fmt.Println(`
😓  The background sync service is not running
    Tokaido will run, but your environment and local host will not be synchronised
    Use 'tok up' to repair, or 'tok sync' to sync manually
	`)

		if conf.GetConfig().Debug == true {
			fmt.Printf("\033[33m%s\033[0m\n", err)
		}
	}

	if GOOS == "osx" {
		osx.CheckSyncService()
	}

}
