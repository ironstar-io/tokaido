package proxy

import (
	"fmt"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/hostsfile"
)

const proxyNetworkName = "tokaido_proxy"

// DestroyProject ...
func DestroyProject() {
	RemoveNetwork()

	RemoveNginxConf()

	if conf.GetConfig().Global.Syncservice == "unison" {
		docker.DeleteVolume("tok_" + conf.GetConfig().Tokaido.Project.Name + "_tokaido_site")
		unison.UnloadSyncService(conf.GetConfig().Tokaido.Project.Name)
	}

	restartContainer("proxy")
}

// RemoveNetwork ...
func RemoveNetwork() {
	hc := conf.GetConfig().Tokaido.Project.Name + "_haproxy_1"

	err := disconnectNetworkEndpoint(proxyNetworkName, hc)
	if err != nil {
		fmt.Println("There was an issue disconnecting the docker network from the proxy containers. These can be removed manually with the command `docker network disconnect " + proxyNetworkName + " " + hc + "`")
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
