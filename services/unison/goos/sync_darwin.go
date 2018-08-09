// +build darwin

package goos

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/user"
	"text/template"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/unison/templates"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/daemon"
)

// UnisonSvc ...
type UnisonSvc struct {
	SyncName    string
	SyncDir     string
	Filename    string
	Filepath    string
	Launchdpath string
	Username    string
}

// NewUnisonSvc - Return a new instance of `UnisonSvc`.
func NewUnisonSvc(syncName, syncDir string) UnisonSvc {
	c := conf.GetConfig()
	u, uErr := user.Current()
	if uErr != nil {
		log.Fatal(uErr)
	}

	s := UnisonSvc{
		SyncName:    syncName,
		SyncDir:     syncDir,
		Filename:    getServiceFilename(syncName),
		Filepath:    getServicePath(syncName),
		Launchdpath: c.System.Syncsvc.Launchdpath,
		Username:    u.Username,
	}

	return s
}

var bgSyncFailMsg = `
😓  The background sync service is not running

Tokaido will run, but your environment and local host will not be synchronised
Use 'tok up' to repair, or 'tok sync' to sync manually
		`

func getServiceFilename(syncName string) string {
	return "tokaido.sync." + syncName + ".plist"
}

func getServicePath(syncName string) string {
	return conf.GetConfig().System.Syncsvc.Launchdpath + getServiceFilename(syncName)
}

// CreateSyncFile - Create a .plist file for the unison service
func (s UnisonSvc) CreateSyncFile() {
	tmpl := template.New(s.SyncName)
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

	writeSyncFile(tpl.String(), s.Launchdpath, s.Filename)
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
func (s UnisonSvc) CreateSyncService() {
	fmt.Println()
	console.Println("🔄  Creating a background process to sync your local repo into the Tokaido environment", "")

	s.RegisterSyncService()
	s.StartSyncService()
}

// RegisterSyncService Register the unison sync service for launchd
func (s UnisonSvc) RegisterSyncService() {
	s.CreateSyncFile()

	daemon.LoadService(s.Filepath)
}

// StartSyncService Start the launchd service after it is created
func (s UnisonSvc) StartSyncService() {
	daemon.StartService(s.Filename)
}

// StopSyncService ...
func (s UnisonSvc) StopSyncService() {
	daemon.KillService(s.Filename, s.Filepath)
}

// SyncServiceStatus ...
func (s UnisonSvc) SyncServiceStatus() string {
	return daemon.ServiceStatus(s.Filename)
}

// CheckSyncService a verbose sync status check used for tok status
func (s UnisonSvc) CheckSyncService() {
	if conf.GetConfig().System.Syncsvc.Enabled != true {
		return
	}

	c := s.SyncServiceStatus()
	if c == "running" {
		console.Println("✅  Background sync service is running", "√")
		return
	}

	console.Println(bgSyncFailMsg, "")
}
