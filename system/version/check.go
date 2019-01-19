package version

import (
	"github.com/ironstar-io/tokaido/system/console"
)

// Check - checks if the current tokaido version is the latest available
func Check() {
	// c := conf.GetConfig()
	v := Get()
	console.Println("current version is "+v.Version, "")

}
