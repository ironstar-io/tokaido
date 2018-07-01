package osx

// RegisterLaunchdService Register the unison sync service for systemd
func RegisterLaunchdService() bool {
	return true
}

// StartLaunchdService ...
func StartLaunchdService() bool {
	return true
}

// StopLaunchdService ...
func StopLaunchdService() bool {
	return true
}

// CheckSyncService ...
func CheckSyncService() bool {
	return true
}
