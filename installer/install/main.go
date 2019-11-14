package install

import (
	"fmt"

	"github.com/ironstar-io/tokaido-installer/docker"
	"github.com/ironstar-io/tokaido-installer/install/goos"
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

	goos.InstallTokBinary("1.12.0")

	fmt.Println("")
	fmt.Println("Tokaido has been successfully installed on your machine!")
	fmt.Println("")
}
