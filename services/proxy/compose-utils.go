package proxy

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/utils"
)

func getExistingCompose() []byte {
	cp := getComposePath()
	var _, err = os.Stat(cp)

	if os.IsNotExist(err) {
		return nil
	}

	dcf, err := ioutil.ReadFile(cp)
	if err != nil {
		log.Fatalf("There was an issue reading in your proxy docker-compose.yml file\n%v", err)
	}

	return dcf
}

// projectInCompose identifies if the specified project exists in the
// provided list of networks
func projectInCompose(networks []string, projectName string) bool {
	// Periods being replaced in recent versions of Docker for network names
	n := strings.Replace(projectName, ".", "", -1)
	utils.DebugString("[proxy] searching to determine if the network '" + n + "' is already configured")

	for _, v := range networks {
		if v == n {
			return true
		}
	}

	return false
}

// buildNetworks compiles the provided list of networks into a docker-compose format
// while also checking those networks still exist in Docker
func buildNetworks(networks []string) []byte {
	n := `networks:
  proxy:
`

	for _, v := range networks {
		if v == "proxy" {
			continue
		}

		// Verify that this network still exists in Docker
		ok, err := validateNetwork(v)
		if err != nil {
			log.Fatal(err)
		}
		if !ok {
			// That network no longer exists, don't add it
			utils.DebugString("network [" + v + "] no longer exists, not adding it to proxy service")
			continue
		}

		n = n + `  ` + strings.Replace(v, ".", "", -1) + `:
    external: true
`
	}

	return []byte(n)
}

// validateNetwork confirms that a docker network exists
// and returns false if no match is found
func validateNetwork(name string) (ok bool, err error) {
	dcli := docker.GetAPIClient()

	filter := filters.NewArgs()
	filter.Add("name", name)

	networks, err := dcli.NetworkList(context.Background(), types.NetworkListOptions{
		Filters: filter,
	})
	if err != nil {
		return false, err
	}

	if len(networks) == 1 {
		return true, nil
	}

	if len(networks) == 0 {
		return false, nil
	}

	if len(networks) > 1 {
		return false, fmt.Errorf("error looking up network [%s]. Received [%d] matching networks when only expecting 1 or none", name, len(networks))
	}

	return false, fmt.Errorf("unexpected error looking up network [%s]", name)
}
