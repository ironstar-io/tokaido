package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/hostsfile"

	"fmt"
	"path/filepath"
)

// DestroyProject ...
func DestroyProject() {
	if fs.CheckExists(getComposePath()) == true {
		RemoveProjectFromDockerCompose()
	}

	RemoveNginxConf()
	RemoveFromHostsfile()

	DetachFromGatsbyEnvFile()

	RestartContainer("yamanote")
	RestartContainer("proxy")
}

// RemoveFromHostsfile ...
func RemoveFromHostsfile() {
	pn := conf.GetConfig().Tokaido.Project.Name
	err := hostsfile.RemoveEntry(pn + "." + constants.ProxyDomain)
	if err != nil {
		fmt.Println("There was an issue updating your hostsfile. Your hosts file can be amended manually in order to enable this feature. See https://tokaido.io/docs/config/#updating-your-hostsfile for more information.")
		fmt.Println(err)
	}
}

// RemoveNginxConf ...
func RemoveNginxConf() {
	pn := conf.GetConfig().Tokaido.Project.Name
	cf := filepath.Join(getProxyClientConfdDir(), pn+".conf")

	fs.Remove(cf)
}
