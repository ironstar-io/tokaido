package unison

import (
	"fmt"
	"os"
	"runtime"

	"github.com/ironstar-io/tokaido/system/fs"
	. "github.com/logrusorgru/aurora"
)

// SystemCompatibilityChecks ...
func SystemCompatibilityChecks() {
	switch opsys := runtime.GOOS; opsys {
	case "darwin":
		if !fs.CheckExists("/usr/local/opt/unox") {
			fmt.Println(Red("Could not find Unison dependency 'unox' on your system. Please install with `brew install eugenmayer/dockersync/unox"))
			os.Exit(1)
		}
		if !fs.CheckExists("/usr/local/bin/unison-fsmonitor") {
			fmt.Println(Red("Could not find Unison dependency 'unison-fsmonitor' on your system. Please install with `brew install eugenmayer/dockersync/unox"))
			os.Exit(1)
		}
		// case "linux":
		// 	if c.Global.Syncservice != "unison" {
		// 		fmt.Println(Sprintf(Yellow("Warning: The syncservice '%s' is not compatible with Linux. Tokaido will automatically be set to use Unison\n\n"), Bold(c.Global.Syncservice)))
		// 		conf.SetGlobalConfigValueByArgs([]string{"syncservice", "unison"})
		// 	}
		// }
	}
}
