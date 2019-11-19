package version

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ironstar-io/tokaido/services/github"
	"github.com/ironstar-io/tokaido/system/version/goos"
)

// GetInstallPath - Check if tok version is installed or not
func GetInstallPath(version string) string {
	return goos.GetInstallPath(version)
}

// DownloadAndInstall - Download and install a selected tok version
func DownloadAndInstall(version string) (string, error) {
	// Check version exists in GH
	// Get the URL for the version
	ghr := []github.ReleaseBody{}
	body, err := github.GetAllReleases()
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(body, &ghr)
	if err != nil {
		return "", err
	}

	var releaseTag string
	for _, r := range ghr {
		if r.TagName == version {
			if r.Draft == true {
				fmt.Println("\nWarning: The selected version is a draft and may not work as intended")
			}

			if r.Prerelease == true {
				fmt.Println("\nWarning: The selected version is a prerelease and may not work as intended")
			}

			releaseTag = r.TagName
		}
	}

	if releaseTag == "" {
		return "", errors.New("There is no release available for version " + version)
	}

	return goos.DownloadTokBinary(releaseTag)
}
