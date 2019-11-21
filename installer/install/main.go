package install

import (
	"fmt"

	"github.com/ironstar-io/tokaido-installer/docker"
	"github.com/ironstar-io/tokaido/system/version/goos"
	"github.com/logrusorgru/aurora"
)

// These values and populated at build time using ldflags
var (
	version    string
	buildDate  string
	goVersion  string
	tokVersion string
	compiler   string
	platform   string
)

// Init - The core run sheet of the tok installer
func Init() {
	// System readiness checks
	docker.CheckClientVersion()

	fmt.Println(aurora.Cyan("\n🚀  Tokaido installer is starting up!"))
	fmt.Println("")

	docker.LoadTokImages()
	docker.RestoreComposerCache()

	goos.SaveTokBinary(tokVersion)

	fmt.Println("")
	fmt.Println("Tokaido has been successfully installed on your machine!")
	fmt.Println("")
}
