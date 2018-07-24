package version

import (
	"log"
	"strings"

	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// GetUnisonVersion retrieves the current version of Unison and returns it.
// If the current version isn't supported, we'll error right here.
func GetUnisonVersion() string {
	v := utils.CommandSubstitution("unison", "-version")

	if strings.Contains(v, "2.48.4") {
		v = "2.48.4"
	} else if strings.Contains(v, "2.51.2") {
		v = "2.51.2"
	} else {
		log.Fatalf("Error matching Unison version. You need Unison 2.48.4 or 2.51.2 on your local system.")
	}

	return v
}
