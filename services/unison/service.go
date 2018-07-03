package unison

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/linux"
	"bitbucket.org/ironstar/tokaido-cli/system/osx"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
)

var bgSyncFailMsg = `
😓  The background sync service is not running

Tokaido will run, but your environment and local host will not be synchronised
Use 'tok up' to repair, or 'tok sync' to sync manually
		`

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func CreateSyncService() {
	GOOS := utils.CheckOS()

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

// CheckSyncService a verbose sync status check used for tok status
func CheckSyncService() error {
	GOOS := utils.CheckOS()
	c := conf.GetConfig()
	if c.CreateSyncService != true {
		return nil
	}

	if GOOS == "linux" {
		linuxErr := linux.CheckSyncService()
		if linuxErr == nil {
			fmt.Println("✅  Background sync service is running")
			return linuxErr
		}

		fmt.Println(bgSyncFailMsg)

		if conf.GetConfig().Debug == true {
			fmt.Printf("\033[33m%s\033[0m\n", linuxErr)
		}
	}

	if GOOS == "osx" {
		osxErr := osx.CheckSyncService()
		if osxErr == nil {
			fmt.Println("✅  Background sync service is running")
			return osxErr
		}
		if conf.GetConfig().Debug == true {
			fmt.Printf("\033[33m%s\033[0m\n", osxErr)
		}

		fmt.Println(bgSyncFailMsg)

		return osxErr
	}

	return nil

}

// CheckSyncServiceSilent status check for internal use, produces no output
func CheckSyncServiceSilent() error {
	GOOS := utils.CheckOS()
	c := conf.GetConfig()
	if c.CreateSyncService != true {
		return nil
	}

	if GOOS == "linux" {
		linuxErr := linux.CheckSyncService()
		if conf.GetConfig().Debug == true {
			fmt.Printf("\033[33m%s\033[0m\n", linuxErr)
		}
		return linuxErr
	}

	if GOOS == "osx" {
		osxErr := osx.CheckSyncService()
		if conf.GetConfig().Debug == true {
			fmt.Printf("\033[33m%s\033[0m\n", osxErr)
		}
		return osxErr
	}

	return nil

}
