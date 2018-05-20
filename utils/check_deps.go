package utils

// CheckDeps - Root executable
func CheckDeps(p string) string {
	var GOOS = CheckOS()

	if GOOS == "linux" {
		// TODO
		return "linux.checkDeps()"
	}

	if GOOS == "osx" {
		// TODO
		return "osx.checkDeps()"
	}

	return ""
}
