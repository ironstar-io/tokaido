// +build darwin

package goos

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/user"
	"text/template"

	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/unison/templates"
	"bitbucket.org/ironstar/tokaido-cli/system/console"
	"bitbucket.org/ironstar/tokaido-cli/system/daemon"
)

type service struct {
	ProjectName string
	ProjectPath string
	Username    string
}

var bgSyncFailMsg = `
😓  The background sync service is not running

Tokaido will run, but your environment and local host will not be synchronised
Use 'tok up' to repair, or 'tok sync' to sync manually
		`

func getServiceName() string {
	return "tokaido.sync." + conf.GetConfig().Tokaido.Project.Name + ".plist"
}

func getServicePath() string {
	return conf.GetConfig().System.SyncSvc.LaunchdPath + getServiceName()
}

func createSyncFile() {
	c := conf.GetConfig()
	u, uErr := user.Current()
	if uErr != nil {
		log.Fatal(uErr)
	}

	s := service{
		ProjectName: c.Tokaido.Project.Name,
		ProjectPath: c.Tokaido.Project.Path,
		Username:    u.Username,
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

	writeSyncFile(tpl.String(), c.System.SyncSvc.LaunchdPath, serviceFilename)
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

func stopSyncService() {
	daemon.StopService(getServiceName())
}

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func CreateSyncService() {
	fmt.Println()
	console.Println("🔄  Creating a background process to sync your local repo into the Tokaido environment", "")

	RegisterSyncService()
	StartSyncService()
}

// RegisterSyncService Register the unison sync service for launchd
func RegisterSyncService() {
	createSyncFile()

	daemon.LoadService(getServicePath())
}

// StartSyncService Start the launchd service after it is created
func StartSyncService() {
	daemon.StartService(getServiceName())
}

// StopSyncService ...
func StopSyncService() {
	daemon.KillService(getServiceName(), getServicePath())
}

// SyncServiceStatus ...
func SyncServiceStatus() string {
	return daemon.ServiceStatus(getServiceName())
}

// CheckSyncService a verbose sync status check used for tok status
func CheckSyncService() {
	if conf.GetConfig().System.SyncSvc.Enabled != true {
		return
	}

	s := SyncServiceStatus()
	if s == "running" {
		console.Println("✅  Background sync service is running", "√")
		return
	}

	console.Println(bgSyncFailMsg, "")
}
