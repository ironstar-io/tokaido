package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	tsurumi "github.com/ironstar-io/telemetry"
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/version"
	"github.com/ironstar-io/tokaido/utils"
	uuid "github.com/satori/go.uuid"
)

// GenerateTelemetryID will generate a unique telemetry ID and save it to global config
func GenerateTelemetryID() {
	c := conf.GetConfig()
	if c.Global.Telemetry.OptOut == true {
		return
	}

	c.Global.Telemetry.Identifier = uuid.NewV4().String()

	conf.WriteGlobalConfig(c.Global)
}

// GenerateProjectID will generate a unique project ID and save it to the project config
func GenerateProjectID() {
	c := conf.GetConfig()

	if len(c.Tokaido.Project.Identifier) > 0 {
		utils.DebugString("received request to generate project id, but skipping because one already exists")
		return
	}

	conf.SetConfigValueByArgs([]string{"tokaido", "project", "identifier", uuid.NewV4().String()}, "project")
}

// RequestOptIn will prompt the user to opt-in to telemetry if they haven't previously opted out
func RequestOptIn() {
	c := conf.GetConfig()

	if c.Global.Telemetry.OptOut == true {
		return
	}

	if len(c.Global.Telemetry.Identifier) > 1 {
		return
	}

	confirmation := utils.ConfirmationPrompt(`
To help us improve Tokaido, you can opt-in to sending anonymous usage data.

The data we collect can not be used to identify you and will never be sold, licensed, or
transferred to a third party. We commit to sharing summary reports on this data (but not
the raw data itself) with the community to help us all understand how Drupal is used.

You can learn more and read our privacy policy at https://docs.tokaido.io/tokaido/usage-data

Would you like to opt-in to sending anonymous usage data?`, "n")

	if !confirmation {
		c.Global.Telemetry.OptOut = true
		conf.WriteGlobalConfig(c.Global)
		return
	}

	fmt.Println()

	GenerateTelemetryID()
}

// SendGlobal posts the current global telemetry info
func SendGlobal() {
	c := conf.GetConfig()

	if c.Global.Telemetry.OptOut {
		return
	}

	pjs := 0
	for range c.Global.Projects {
		pjs++
	}

	checkin := tsurumi.UserCheckin{
		OperatingSystemVersion: runtime.GOOS,
		SyncStrategy:           c.Global.Syncservice,
		ProjectCount:           pjs,
		TokaidoVersion:         version.Get().Version,
	}

	url := "https://api.tokaido.io/v1/user/checkin/" + c.Global.Telemetry.Identifier
	utils.DebugString("sending global checkin to: " + url)

	data, err := json.Marshal(checkin)
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not prepare global telemetry data: %v", err.Error()))
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not prepare global telemetry data: %v", err.Error()))
		return
	}

	client := &http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not POST global telemetry data: %v", err.Error()))
		return
	}
	defer resp.Body.Close()

	utils.DebugString("received global checkin response to: " + strconv.Itoa(resp.StatusCode) + " " + resp.Status)

}

// SendProject posts the  current project telemetry info
// startTime is used to record the startup time and tags the checkin
// duration is how long Tokaido took to complete a successful startup
func SendProject(startTime time.Time, duration int) {
	c := conf.GetConfig()

	if c.Global.Telemetry.OptOut {
		return
	}

	if len(c.Global.Telemetry.Identifier) < 1 {
		utils.DebugString("unable to send command telemetry data because the global identifier is missing")
		return
	}

	if len(c.Tokaido.Project.Identifier) < 1 {
		utils.DebugString("unable to send command telemetry data because the project identifier is missing")
		return
	}

	checkin := tsurumi.DrupalProjectCheckin{
		TelemetryID:   c.Global.Telemetry.Identifier,
		Timestamp:     startTime,
		PhpVersion:    c.Tokaido.Phpversion,
		Mailhog:       c.Services.Mailhog.Enabled,
		Adminer:       c.Services.Adminer.Enabled,
		Solr:          c.Services.Solr.Enabled,
		Redis:         c.Services.Redis.Enabled,
		Memcache:      c.Services.Memcache.Enabled,
		Stability:     c.Tokaido.Stability,
		DrupalVersion: c.Drupal.Majorversion,
		PHPMemory:     c.Fpm.Phpmemorylimit,
		Duration:      duration,
	}

	url := "https://api.tokaido.io/v1/project/drupal/checkin/" + c.Tokaido.Project.Identifier
	utils.DebugString("sending drupal project checkin to: " + url)

	data, err := json.Marshal(checkin)
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not prepare project telemetry data: %v", err.Error()))
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not prepare command telemetry data: %v", err.Error()))
		return
	}

	client := &http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not POST project telemetry data: %v", err.Error()))
		return
	}
	defer resp.Body.Close()
	utils.DebugString("received project checkin response to: " + strconv.Itoa(resp.StatusCode) + " " + resp.Status)

}

// SendCommand posts that a command was executed
func SendCommand(command string) {
	c := conf.GetConfig()

	if c.Global.Telemetry.OptOut {
		return
	}

	checkin := tsurumi.CommandCheckin{
		TelemetryID: c.Global.Telemetry.Identifier,
		ProjectID:   c.Tokaido.Project.Identifier,
		Command:     command,
	}

	url := "https://api.tokaido.io/v1/command/checkin/" + command
	utils.DebugString("sending command checkin to: " + url)

	data, err := json.Marshal(checkin)
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not prepare command telemetry data: %v", err.Error()))
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not prepare command telemetry data: %v", err.Error()))
		return
	}

	client := &http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not POST command telemetry data: %v", err.Error()))
		return
	}
	defer resp.Body.Close()
	utils.DebugString("received command checkin response to: " + strconv.Itoa(resp.StatusCode) + " " + resp.Status)

}
