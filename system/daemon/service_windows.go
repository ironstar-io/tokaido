//go:build windows
// +build windows

package daemon

type service struct {
	ProjectName string
	ProjectPath string
	Username    string
}

// LoadService ...
func LoadService(servicePath string) {
}

// UnloadService ...
func UnloadService(servicePath string) {
}

// StartService ...
func StartService(serviceName string) {
}

// StopService ...
func StopService(serviceName string) {
}

// DeleteService ...
func DeleteService(servicePath string) {
}

// ServiceStatus checks if the unison background process is running
func ServiceStatus(serviceName string) string {
	return ""
}

// KillService ...
func KillService(serviceName string, servicePath string) {
}
