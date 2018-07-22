package unison

import (
	"log"
	"strings"

	"bitbucket.org/ironstar/tokaido-cli/utils"
)

// CheckVersion and fail if we can't
func CheckVersion() {
	v := utils.CommandSubstitution("unison", "-version")

	if strings.Contains(v, "2.48.4") || strings.Contains(v, "2.51.2") {
		return
	}

	log.Fatalf("Unsupported Unison version. You need Unison 2.48.4 or 2.51.2 on your local system.")
}
