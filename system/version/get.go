package version

import (
	"fmt"
	"runtime"
)

var (
	version   string
	buildDate string
	goVersion string
	compiler  string
	platform  string
)

// Info - Data model for version information
type Info struct {
	Version   string
	BuildDate string
	GoVersion string
	Compiler  string
	Platform  string
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in pkg/version/base.go
	return Info{
		Version:   version,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
