package goos

import (
	"bytes"
	"log"
	"os"
	"os/user"
	"text/template"

	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/services/unison/templates"
	"bitbucket.org/ironstar/tokaido-cli/system/daemon"
)

type service struct {
	ProjectName string
	ProjectPath string
	Username    string
}

func getServiceName() string {
	return "tokaido.sync." + conf.GetConfig().Project + ".plist"
}

func getServicePath() string {
	return conf.GetConfig().LaunchdPath + getServiceName()
}

func createSyncFile() {
	c := conf.GetConfig()
	u, uErr := user.Current()
	if uErr != nil {
		log.Fatal(uErr)
	}

	s := service{
		ProjectName: c.Project,
		ProjectPath: c.Path,
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

func stopSyncService() {
	daemon.StopService(getServiceName())
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
