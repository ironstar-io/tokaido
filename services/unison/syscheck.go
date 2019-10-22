package unison

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ironstar-io/tokaido/system/fs"
	// "github.com/ironstar-io/tokaido/system/wsl"
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
		// w := wsl.IsWSL()

		// // Can't use docker volumes on Linux, except for in WSL
		// if c.Global.Syncservice != "unison" && !w {
		// 	fmt.Println(aurora.Sprintf(aurora.Yellow("Warning: The syncservice '%s' is not compatible with Linux. Tokaido will automatically be set to use Unison\n\n"), aurora.Bold(c.Global.Syncservice)))
		// 	conf.SetGlobalConfigValueByArgs([]string{"syncservice", "unison"})
		// }
		// // Must use docker volumes on WSL
		// if c.Global.Syncservice != "docker" && w {
		// 	fmt.Println(aurora.Sprintf(aurora.Yellow("Warning: The syncservice '%s' is not compatible with WSL. Tokaido will automatically be set to use Docker Volumes\n\n"), aurora.Bold(c.Global.Syncservice)))
		// 	conf.SetGlobalConfigValueByArgs([]string{"syncservice", "docker"})
		// }
	}
}
