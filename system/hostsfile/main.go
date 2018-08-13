package hostsfile

import (
	"github.com/ironstar-io/tokaido/utils"

	"fmt"
)

const localhost = "127.0.0.1"

func confirmAmend() bool {
	c := utils.ConfirmationPrompt("Would you like Tokaido to automatically update your hostsfile? You may be prompted for elevated access.", "n")
	if c == false {
		fmt.Println(`Your hostsfile can be amended manually in order to enable this feature. See https://tokaido.io/docs/config/#updating-your-hostsfile for more information.`)
	}

	return c
}

// AddEntry - Add an entry to /etc/hosts or equivalent
func AddEntry(hostname string) error {
	hosts, err := NewHosts()
	if err != nil {
		return err
	}

	if !hosts.Has(localhost, hostname) {
		if confirmAmend() == false {
			return nil
		}

		hosts.Add(localhost, hostname)
		if hosts.IsWritable() == false {
			err := hosts.WriteElevated()
			if err != nil {
				return err
			}

			return nil
		}

		if err := hosts.Flush(); err != nil {
			return err
		}

		return nil
	}

	return nil
}

// RemoveEntry - Remove an entry from /etc/hosts or equivalent
func RemoveEntry(hostname string) error {
	hosts, err := NewHosts()
	if err != nil {
		return err
	}

	if hosts.Has(localhost, hostname) {
		if confirmAmend() == false {
			return nil
		}

		hosts.Remove(localhost, hostname)
		if hosts.IsWritable() == false {
			err := hosts.WriteElevated()
			if err != nil {
				return err
			}

			return nil
		}

		if err := hosts.Flush(); err != nil {
			return err
		}

		return nil
	}

	return nil
}
