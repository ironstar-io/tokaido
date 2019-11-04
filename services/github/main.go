package github

import (
	"io/ioutil"
	"net/http"
	"time"
)

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
