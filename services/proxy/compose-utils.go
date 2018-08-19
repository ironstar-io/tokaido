package proxy

import (
	"io/ioutil"
	"log"
	"os"
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
	for _, v := range networks {
		if v == projectName+"_default" {
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

		n = n + `  ` + v + `:
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
