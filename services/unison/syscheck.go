package unison

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/wsl"
	"github.com/logrusorgru/aurora"
)

// SystemCompatibilityChecks ...
func SystemCompatibilityChecks() {
	switch opsys := runtime.GOOS; opsys {
	case "darwin":
		if !fs.CheckExists("/usr/local/opt/unox") {
			fmt.Println(aurora.Red("Could not find Unison dependency 'unox' on your system. Please install with `brew install eugenmayer/dockersync/unox"))
			os.Exit(1)
		}
		if !fs.CheckExists("/usr/local/bin/unison-fsmonitor") {
			fmt.Println(aurora.Red("Could not find Unison dependency 'unison-fsmonitor' on your system. Please install with `brew install eugenmayer/dockersync/unox"))
			os.Exit(1)
		}
	case "linux":
		w := wsl.IsWSL()
		c := conf.GetConfig()

		// Can't use docker volumes on Linux, except for in WSL
		if c.Global.Syncservice != "unison" && !w {
			conf.SetGlobalConfigValueByArgs([]string{"global", "syncservice", "unison"})
		}
		// Must use docker volumes on WSL
		if c.Global.Syncservice != "docker" && w {
			conf.SetGlobalConfigValueByArgs([]string{"global", "syncservice", "docker"})
		}
	case "windows":
		c := conf.GetConfig()

		// Must use docker volumes on Windows
		if c.Global.Syncservice != "docker" {
			conf.SetGlobalConfigValueByArgs([]string{"global", "syncservice", "docker"})
		}
	}
}
