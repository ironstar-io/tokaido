package hostsfile

import (
	"github.com/ironstar-io/tokaido/utils"

	"log"

	"github.com/lextoumbourou/goodhosts"
)

// AddEntry - Add an entry to /etc/hosts or equivalent
func AddEntry(hostname string) {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		log.Fatal(err)
	}

	if !hosts.Has("127.0.0.1", hostname) {
		utils.GainSudo()

		hosts.Add("127.0.0.1", hostname)

		if err := hosts.Flush(); err != nil {
			log.Fatal(err)
		}
		return
	}
}

// RemoveEntry - Remove an entry from /etc/hosts or equivalent
func RemoveEntry(hostname string) {
	hosts, err := goodhosts.NewHosts()
	if err != nil {
		log.Fatal(err)
	}

	if !hosts.Has("127.0.0.1", hostname) {
		utils.GainSudo()

		hosts.Remove("127.0.0.1", hostname)

		if err := hosts.Flush(); err != nil {
			log.Fatal(err)
		}
		return
	}
}
