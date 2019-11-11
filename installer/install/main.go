package install

import (
	"fmt"
	"runtime"

	"github.com/ironstar-io/tokaido-installer/docker"
	"github.com/logrusorgru/aurora"
)

// Init - The core run sheet of the tok installer
func Init() {
	// System readiness checks
	docker.CheckClientVersion()

	fmt.Println(aurora.Cyan("\n🚀  Tokaido installer is starting up!"))

	docker.LoadTokImages()

	switch runtime.GOOS {
	case "darwin":

	case "linux":

	case "windows":

	}

	fmt.Println("Completed!")
}
