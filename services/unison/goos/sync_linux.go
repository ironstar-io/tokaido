// +build linux

package goos

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/unison/templates"
	"bitbucket.org/ironstar/tokaido-cli/system/console"
	"bitbucket.org/ironstar/tokaido-cli/system/daemon"

	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
)

var bgSyncFailMsg = `
😓  The background sync service is not running

Tokaido will run, but your environment and local host will not be synchronised
Use 'tok up' to repair, or 'tok sync' to sync manually
		`

type service struct {
	ProjectName string
	ProjectPath string
}

func getServiceName() string {
	return "tokaido-sync-" + conf.GetConfig().Project + ".service"
}

func getServicePath() string {
	return conf.GetConfig().SystemdPath + getServiceName()
}

func createSyncFile() {
	c := conf.GetConfig()

	s := service{
		ProjectName: c.Project,
		ProjectPath: c.Path,
	}

	serviceFilename := getServiceName()

	tmpl := template.New(serviceFilename)
	tmpl, err := tmpl.Parse(unisontmpl.SyncTemplateStr)

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
	daemon.ReloadServices()
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

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func CreateSyncService() {
	fmt.Println()
	console.Println("🔄  Creating a background process to sync your local repo into the Tokaido environment", "")

	RegisterSyncService()
	StartSyncService()
}

// RegisterSystemdService Register the unison sync service for systemd
func RegisterSyncService() {
	createSyncFile()
}

// StartSystemdService Start the systemd service after it is created
func StartSyncService() {
	daemon.StartService(getServiceName())
}

// SyncServiceStatus ...
func SyncServiceStatus() string {
	return daemon.ServiceStatus(getServiceName())
}

// StopSyncService ...
func StopSyncService() {
	daemon.StopService(getServiceName())
	daemon.DeleteService(getServicePath())
}

// CheckSyncService a verbose sync status check used for tok status
func CheckSyncService() {
	c := conf.GetConfig()
	if c.CreateSyncService != true {
		return
	}

	s := SyncServiceStatus()
	if s == "running" {
		console.Println("✅  Background sync service is running", "√")
		return
	}

	fmt.Println(bgSyncFailMsg)
}
