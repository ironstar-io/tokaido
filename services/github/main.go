package github

import (
	"io/ioutil"
	"net/http"
	"time"
)

// ReleaseBody - Used properties from the GH API GET call
type ReleaseBody struct {
	TagName    string `json:"tag_name"`
	HTMLURL    string `json:"html_url"`
	Prerelease bool   `json:"prerelease"`
	Draft      bool   `json:"draft"`
	Assets     []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	}
}

// GetRelease - Get a Tokaido release
func GetRelease(version string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/repos/ironstar-io/tokaido/releases"+version, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return body, nil
}

// GetLatestRelease - Get latest Tokaido release
func GetLatestRelease() ([]byte, error) {
	return GetRelease("/latest")
}

// GetAllReleases - Get all Tokaido releases
func GetAllReleases() ([]byte, error) {
	return GetRelease("")
}
