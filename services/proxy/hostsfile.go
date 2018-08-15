package proxy

import (
	"github.com/ironstar-io/tokaido/system/hostsfile"

	"fmt"
)

// ConfigureHostsfile ...
func ConfigureHostsfile(domain string) {
	err := hostsfile.AddEntry(domain)
	if err != nil {
		fmt.Println("There was an issue updating your hostsfile. Your hostsfile can be amended manually in order to enable this feature. See https://tokaido.io/docs/config/#updating-your-hostsfile for more information.")
		fmt.Println(err)
	}
}
