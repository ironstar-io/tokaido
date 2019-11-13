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
	fmt.Println("")

	// docker.LoadTokImages()
	docker.RestoreComposerCache()

	switch runtime.GOOS {
	case "darwin":

	case "linux":

	case "windows":

	}

	fmt.Println("")
	fmt.Println("Tokaido has been successfully installed on your machine!")
	fmt.Println("")
}
