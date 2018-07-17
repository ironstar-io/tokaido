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
	reloadSystemd()
}

func reloadSystemd() {
	_, err := utils.CommandSubSplitOutput("systemctl", "--user", "daemon-reload")
	if err != nil {
		log.Fatal("Unable to reload the systemd config: ", err)
	}
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
	_, errStatus := utils.CommandSubSplitOutput("systemctl", "--user", "status", "tokaido.sync."+c.Project+".service")
	if errStatus == nil {
		_, errStop := utils.CommandSubSplitOutput("systemctl", "--user", "stop", "tokaido-sync-"+c.Project+".service")
		if errStop != nil {
			log.Fatal("Unable to stop the sync service: ", errStop)
		}
	}
}

func deleteSyncService() {
	c := conf.GetConfig()
	rmErr := os.Remove(c.SystemdPath + "/tokaido-sync-" + c.Project + ".service")
	if os.IsNotExist(rmErr) {
		return
	} else if rmErr != nil {
		log.Fatal("Unable to remove service configuration: ", rmErr)
	}

	reloadSystemd()

	fmt.Println(`🔄  Removed the background sync process`)
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
	stopSyncService()
	deleteSyncService()
}
