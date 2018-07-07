package linux

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/linux/templates"
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

type service struct {
	ProjectName string
	ProjectPath string
}

func createSyncFile() {
	c := conf.GetConfig()

	s := service{
		ProjectName: c.Project,
		ProjectPath: c.Path,
	}

	serviceFilename := "tokaido-sync-" + s.ProjectName + ".service"

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

	writeSyncFile(tpl.String(), c.SystemdPath, serviceFilename)
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
	_, err := utils.CommandSubSplitOutput("systemctl", "--user", "start", "tokaido-sync-"+c.Project+".service")
	if err != nil {
		log.Fatal("Unable to start the sync service: ", err)
	}
}

func stopSyncService() {
	c := conf.GetConfig()
	_, err := utils.CommandSubSplitOutput("systemctl", "--user", "stop", "tokaido-sync-"+c.Project+".service")
	if err != nil {
		log.Fatal("Unable to stop the sync service: ", err)
	}
}

func deleteSyncService() {
	c := conf.GetConfig()
	rmErr := os.Remove(c.SystemdPath + "/tokaido-sync-" + c.Project + ".service")
	if rmErr != nil {
		log.Fatal("Unable to remove the sync service: ", rmErr)
	}

	_, reloadErr := utils.CommandSubSplitOutput("systemctl", "--user", "daemon-reload")
	if reloadErr != nil {
		log.Fatal("Unable to reload the sync service: ", reloadErr)
	}
}

// RegisterSystemdService Register the unison sync service for systemd
func RegisterSystemdService() {
	createSyncFile()
}

// StartSystemdService Start the systemd service after it is created
func StartSystemdService() {
	startSyncService()
	CheckSyncService()
}

// CheckSyncService checks if the unison background process is running
func CheckSyncService() string {
	c := conf.GetConfig()
	_, err := utils.CommandSubSplitOutput("systemctl", "--user", "status", "tokaido-sync-"+c.Project+".service")
	if err == nil {
		return "running"
	}

	if c.Debug == true {
		fmt.Printf("\033[33m%s\033[0m\n", err)
	}

	return "stopped"
}

// StopSystemdService ...
func StopSystemdService() {
	fmt.Println(`
🔄  Removing the background sync process
	`)
	stopSyncService()
	deleteSyncService()
}
