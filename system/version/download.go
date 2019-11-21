package version

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ironstar-io/tokaido/services/github"
	"github.com/ironstar-io/tokaido/system/version/goos"
)

// GetInstallPath - Check if tok version is installed or not
func GetInstallPath(version string) string {
	return goos.GetInstallPath(version)
}

// DownloadAndInstall - Download and install a selected tok version
func DownloadAndInstall(version string) (string, error) {
	releaseTag := GetReleaseTagFromVersion(version)
	return goos.DownloadTokBinary(releaseTag)
}

// GetReleaseTagFromVersion returns a githab-ready release tag from a Tokaido version string
func GetReleaseTagFromVersion(version string) (releaseTag string) {
	// Check version exists in GH
	// Get the URL for the version
	ghr := []github.ReleaseBody{}
	body, err := github.GetAllReleases()
	if err != nil {
		fmt.Println("Unexpected error retrieving list of available Tokaido releases: " + err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(body, &ghr)
	if err != nil {
		fmt.Println("Unexpected error assembling list of available Tokaido releases: " + err.Error())
		os.Exit(1)
	}

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
		fmt.Println("There is no release information available for version [" + version + "]")
		os.Exit(1)
	}

	return releaseTag
}
