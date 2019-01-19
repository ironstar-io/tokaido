package version

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
	"time"

	"github.com/blang/semver"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/utils"
)

// Latest is the latest version response from release.tokaido.io/latest
type Latest struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

// Check - checks if the current tokaido version is the latest available
func Check() {
	// c := conf.GetConfig()
	info := Get()
	re := regexp.MustCompile("[^-]*")
	match := re.FindStringSubmatch(info.Version)
	current, _ := semver.Make(match[0])

	latestVersion, latestURL, err := getLatestVersion()
	if err != nil {
		utils.DebugErrOutput(err)
		return
	}
	latestSemver, _ := semver.Make(latestVersion)
	if err != nil {
		utils.DebugErrOutput(err)
		return
	}

	if latestSemver.GT(current) {
		console.Println("\n👵🏻  You're running an old version of Tokaido. Please consider upgrading to Tokaido "+latestVersion, "")

		if runtime.GOOS == "darwin" {
			console.Println("    You can upgrade easily by running 'brew upgrade tokaido'", "")
		} else {
			console.Println("    You can get the latest version at "+latestURL, "")
		}
	}
}

func getLatestVersion() (ver, url string, err error) {
	req, err := http.NewRequest(http.MethodGet, "https://releases.tokaido.io/latest", nil)
	if err != nil {
		return "", "", err
	}

	client := http.Client{
		Timeout: time.Second * 3,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", "", err
	}

	latest := Latest{}
	err = json.Unmarshal(body, &latest)
	if err != nil {
		return "", "", err
	}

	return latest.Version, latest.URL, nil

}
