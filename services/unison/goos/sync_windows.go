// +build windows

package goos

// UnisonSvc ...
type UnisonSvc struct {
	SyncName    string
	SyncDir     string
	Filename    string
	Filepath    string
	Launchdpath string
	Username    string
}

// NewUnisonSvc - Return a new instance of `UnisonSvc`.
func NewUnisonSvc(syncName, syncDir string) UnisonSvc {
	return UnisonSvc{}
}

var bgSyncFailMsg = ""

func getServiceFilename(syncName string) string {
	return ""
}

func getServicePath(syncName string) string {
	return ""
}

// CreateSyncFile - Create a .plist file for the unison service
func (s UnisonSvc) CreateSyncFile() {
}

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func (s UnisonSvc) CreateSyncService() {
}

// RegisterSyncService Register the unison sync service for launchd
func (s UnisonSvc) RegisterSyncService() {
}

// UnloadSyncService Remove the unison sync service from launchd
func (s UnisonSvc) UnloadSyncService() {
}

// StartSyncService Start the launchd service after it is created
func (s UnisonSvc) StartSyncService() {
}

// RestartSyncService ...
func (s UnisonSvc) RestartSyncService() {
}

// StopSyncService ...
func (s UnisonSvc) StopSyncService() {
}

// SyncServiceStatus ...
func (s UnisonSvc) SyncServiceStatus() string {
	return ""
}

// CheckSyncService a verbose sync status check used for tok status
func (s UnisonSvc) CheckSyncService() error {
	return nil
}
