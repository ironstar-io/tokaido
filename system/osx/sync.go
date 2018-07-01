package osx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"

	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/osx/templates"
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

type service struct {
	ProjectName string
	ProjectPath string
}

type status struct {
	Pid int64 `json:"PID"`
}

func createSyncFile() {
	c := conf.GetConfig()

	s := service{
		ProjectName: c.Project,
		ProjectPath: c.Path,
	}

	serviceFilename := "tokaido.sync." + s.ProjectName + ".plist"

	tmpl := template.New(serviceFilename)
	tmpl, err := tmpl.Parse(synctmpl.SyncTemplateStr)

	if err != nil {
		log.Fatal("Parse: ", err)
		return
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, s); err != nil {
		log.Fatal("Parse: ", err)
		return
	}

	writeSyncFile(tpl.String(), c.LaunchdPath, serviceFilename)
}

func writeSyncFile(body string, path string, filename string) {
	mkdErr := os.MkdirAll(path, os.ModePerm)
	if mkdErr != nil {
		log.Fatal("Mkdir: ", mkdErr)
	}

	var file, err = os.Create(path + filename)
	if err != nil {
		log.Fatal("Create: ", err)
	}

	_, _ = file.WriteString(body)

	defer file.Close()
}

func startSyncService() {
	c := conf.GetConfig()
	_, err := utils.CommandSubSplitOutput("launchctl", "start", "tokaido.sync."+c.Project+".plist")
	if err != nil {
		log.Fatal("Unable to start sync service: ", err)
	}
}

func stopSyncService() {
	c := conf.GetConfig()
	_, err := utils.CommandSubSplitOutput("launchctl", "stop", "tokaido.sync."+c.Project+".plist")
	if err != nil {
		log.Fatal("Unable to start sync service: ", err)
	}
}

func deleteSyncService() {
	c := conf.GetConfig()
	rmErr := os.Remove(c.LaunchdPath + "/tokaido.sync." + c.Project + ".plist")
	if rmErr != nil {
		log.Fatal("Unable to start sync service: ", rmErr)
	}
}

// RegisterLaunchdService Register the unison sync service for launchd
func RegisterLaunchdService() {
	fmt.Println(`
🔄  Creating a background process to sync your local repo into the Tokaido environment
	`)
	createSyncFile()
}

// StartLaunchdService Start the launchd service after it is created
func StartLaunchdService() {
	startSyncService()
	CheckSyncService()
}

// CheckSyncService checks if the unison background process is running
func CheckSyncService() error {
	c := conf.GetConfig()
	s, err := utils.CommandSubSplitOutput("launchctl", "list", "tokaido.sync."+c.Project+".plist")
	if err != nil {
		return err
	}

	// Marhsal the result into JSON so we can analyse the pid number.
	// If launchd returns a pid, that's how we know a service is running
	var j status
	jsonErr := json.Unmarshal([]byte(s), &j)
	if jsonErr != nil {
		fmt.Println("in JSON err")
		return jsonErr
	}
	fmt.Printf("PID = %d", j.Pid)

	return nil

}

// StopLaunchdService ...
func StopLaunchdService() {
	fmt.Println(`
🔄  Removing the background sync process
	`)
	stopSyncService()
	deleteSyncService()
}
