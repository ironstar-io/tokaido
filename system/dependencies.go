package system

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/goos"
)

// CheckDependencies - Root executable
func CheckDependencies() {
	if conf.GetConfig().Tokaido.Dependencychecks == true {
		goos.CheckDependencies()
	}
}
