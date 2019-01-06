package proxy

import (
	"fmt"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/hostsfile"
)

const proxyNetworkName = "proxy_proxy_1"

// DestroyProject ...
func DestroyProject() {
	if fs.CheckExists(getComposePath()) == true {
		RemoveProjectFromDockerCompose()
	}

	RemoveNetwork()

	RemoveNginxConf()
	RemoveFromHostsfile()

	RestartContainer("proxy")
}

// RemoveNetwork ...
func RemoveNetwork() {
	n := conf.GetConfig().Tokaido.Project.Name + "_default"

	err := DisconnectNetworkEndpoint(n, proxyNetworkName)
	if err != nil {
		fmt.Println("There was an issue disconnecting the docker network from the proxy containers. These can be removed manually with the command `docker network disconnect " + n + " " + proxyNetworkName + "`")
		fmt.Println(err)
	}
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
