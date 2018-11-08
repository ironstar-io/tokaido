// +build windows

package goos

type UnisonSvc struct {
	SyncName string
	SyncDir  string
	Filename string
	Filepath string
}

// NewUnisonSvc - Return a new instance of `UnisonSvc`.
func NewUnisonSvc(syncName, syncDir string) UnisonSvc {
	// c := conf.GetConfig()

	s := UnisonSvc{}

	return s
}

// CreateSyncService - Unused placeholder
func (s UnisonSvc) CreateSyncService() {
	return
}

// RegisterSystemdService - Unused placeholder
func (s UnisonSvc) RegisterSyncService() {
	return
}

// StartSyncService - Unused placeholder
func (s UnisonSvc) StartSyncService() {
	return
}

// SyncServiceStatus - While there is no service functionality for Windows, this will always return "stopped"
func (s UnisonSvc) SyncServiceStatus() string {
	return "stopped"
}

// StopSyncService - Unused placeholder
func (s UnisonSvc) StopSyncService() {
	return
}

// CheckSyncService - Unused placeholder
func (s UnisonSvc) CheckSyncService() error {
	return nil
}
