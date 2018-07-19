// +build windows

package goos

import ()

type service struct {
	ProjectName string
	ProjectPath string
}

// CreateSyncService - Unused placeholder
func CreateSyncService() {
	return
}

// RegisterSystemdService - Unused placeholder
func RegisterSyncService() {
	return
}

// StartSyncService - Unused placeholder
func StartSyncService() {
	return
}

// SyncServiceStatus - While there is no service functionality for Windows, this will always return "stopped"
func SyncServiceStatus() string {
	return "stopped"
}

// StopSyncService - Unused placeholder
func StopSyncService() {
	return
}

// CheckSyncService - Unused placeholder
func CheckSyncService() {
	return
}
