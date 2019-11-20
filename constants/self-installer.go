package constants

const (
	// BaseBinaryURL is where Tokaido releases are available
	BaseBinaryURL = "https://github.com/ironstar-io/tokaido/releases/download/"

	// BaseInstallPathLinux is where Tokaido binaries are installed on Linux
	BaseInstallPathLinux = ".tok/bin"

	// BaseInstallPathDarwin is where Tokaido binaries are installed on macOS/Darwin
	BaseInstallPathDarwin = ".tok/bin"

	// BaseInstallPathWindows is where Tokaido binaries are installed on Windows
	BaseInstallPathWindows = "AppData/Local/Ironstar/Tokaido"

	// BinaryNameLinux is the name of the Tokaido Linux binary on Github
	BinaryNameLinux = "tok-linux-amd64"

	// BinaryNameDarwin is the name of the Tokaido macOS binary on Github
	BinaryNameDarwin = "tok-macos"

	// BinaryNameWindows is the name of the Tokaido Windows binary on Github
	BinaryNameWindows = "tok-windows.exe"

	// GlobalBinaryPathDarwin is the location of of the 'tok' command which is a symlink to the active Tokaido version
	GlobalBinaryPathDarwin = "/usr/local/bin/tok"

	// GlobalBinaryPathLinux is the location of of the 'tok' command which is a symlink to the active Tokaido version
	GlobalBinaryPathLinux = "/usr/local/bin/tok"
)
