package telemetry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	tsurumi "github.com/ironstar-io/telemetry"
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/version"
	"github.com/ironstar-io/tokaido/utils"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Survey ...
func Survey() {
	fmt.Println("")
	c := conf.GetConfig()
	nextFeatures := []string{
		"The ability to share my local environment on a public URL",
		"Improved Apache Solr support",
		"Support for PostgreSQL databases",
		"Windows support",
		"Bug fixes and stability improvements",
	}
	qs := []*survey.Question{
		{
			Name: "satisfaction",
			Prompt: &survey.Select{
				Message: "How happy are you with Tokaido? ",
				Options: []string{"Very Happy", "Somewhat Happy", "Neutral", "Somewhat Unhappy", "Very Unhappy"},
			},
		},
		{
			Name: "nextFeature",
			Prompt: &survey.Select{
				Message: "Which of the following upcoming features should we build first? ",
				Options: nextFeatures,
			},
		},

		{
			Name:   "message",
			Prompt: &survey.Input{Message: "Is there any other feedback you'd like tell us? (or press enter to submit) "},
		},
	}
	answers := tsurumi.SurveyResponse{}

	if !c.Global.Telemetry.OptOut {
		answers.TelemetryID = c.Global.Telemetry.Identifier
		answers.TokaidoVersion = version.Get().Version
	}

	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Printf("Unexpected error with the Survey: %s", err.Error())
	}

	url := "https://api.tokaido.io/v1/survey"
	utils.DebugString("sending survey to: " + url)

	data, err := json.Marshal(answers)
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not prepare survey data: %v", err.Error()))
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not prepare survey data: %v", err.Error()))
		return
	}

	client := &http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
	resp, err := client.Do(req)
	if err != nil {
		utils.DebugString(fmt.Sprintf("[warning] Could not POST survey data: %v", err.Error()))
		return
	}
	defer resp.Body.Close()

	utils.DebugString("received submission response: " + strconv.Itoa(resp.StatusCode) + " " + resp.Status)
	fmt.Println("")
}
