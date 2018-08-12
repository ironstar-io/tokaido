package proxy

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/hostsfile"

	"fmt"
	"os"
	"path/filepath"
)

// DestroyProject ...
func DestroyProject() {
	var _, err = os.Stat(getComposePath())

	if !os.IsNotExist(err) {
		DockerComposeDown()
		RemoveProjectFromDockerCompose()
	}

	RemoveNginxConf()
	RemoveFromHostsfile()
}

// RemoveFromHostsfile ...
func RemoveFromHostsfile() {
	pn := conf.GetConfig().Tokaido.Project.Name
	err := hostsfile.RemoveEntry(pn + "." + constants.ProxyDomain)
	if err != nil {
		fmt.Println("There was an issue updating your hostsfile. Your hostsfile can be amended manually in order to enable this feature. See XXXXXXX for more information.")
		fmt.Println(err)
	}
}

// RemoveNginxConf ...
func RemoveNginxConf() {
	pn := conf.GetConfig().Tokaido.Project.Name
	cf := filepath.Join(getProxyClientConfdDir(), pn+".conf")

	fs.Remove(cf)
}
