package system

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/goos"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	if conf.GetConfig().Tokaido.Dependencychecks == true {
		goos.CheckDependencies()
	}
}
