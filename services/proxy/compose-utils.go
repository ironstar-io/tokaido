package proxy

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

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

func buildNetworks(networks []string) []byte {
	n := `networks:
  proxy:
`

	for _, v := range networks {
		if v == "proxy" {
			continue
		}

		n = n + `  ` + strings.Replace(v, ".", "", -1) + `:
    external: true
`
	}

	return []byte(n)
}

// SetUnisonVersion ...
func SetUnisonVersion(version string) []byte {
	return []byte(`services:
  unison:
    image: tokaido/unison:` + version)
}
