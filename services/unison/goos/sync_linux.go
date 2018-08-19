// +build linux

package goos

import (
	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/unison/templates"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/daemon"
	"github.com/ironstar-io/tokaido/system/fs"

	"bytes"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
)

var bgSyncFailMsg = `
😓  The background sync service is not running

Tokaido will run, but your environment and local host will not be synchronised
Use 'tok up' to repair, or 'tok sync' to sync manually
		`

// UnisonSvc ...
type UnisonSvc struct {
	SyncName    string
	SyncDir     string
	Filename    string
	Filepath    string
	Systemdpath string
}

// NewUnisonSvc - Return a new instance of `UnisonSvc`.
func NewUnisonSvc(syncName, syncDir string) UnisonSvc {
	c := conf.GetConfig()

	s := UnisonSvc{
		SyncName:    syncName,
		SyncDir:     syncDir,
		Filename:    getServiceFilename(syncName),
		Filepath:    getServicePath(syncName),
		Systemdpath: c.System.Syncsvc.Systemdpath,
	}

	return s
}

func getServiceFilename(syncName string) string {
	return "tokaido-sync-" + syncName + ".service"
}

func getServicePath(syncName string) string {
	return conf.GetConfig().System.Syncsvc.Systemdpath + getServiceFilename(syncName)
}

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

	fs.TouchOrReplace(filepath.Join(s.Systemdpath, s.Filename), tpl.Bytes())

	daemon.ReloadServices()
}

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func (s UnisonSvc) CreateSyncService() {
	s.RegisterSyncService()
	s.StartSyncService()
}

// RegisterSystemdService Register the unison sync service for systemd
func (s UnisonSvc) RegisterSyncService() {
	s.CreateSyncFile()
}

// StartSystemdService Start the systemd service after it is created
func (s UnisonSvc) StartSyncService() {
	daemon.StartService(s.Filename)
}

// SyncServiceStatus ...
func (s UnisonSvc) SyncServiceStatus() string {
	return daemon.ServiceStatus(s.Filename)
}

// StopSyncService ...
func (s UnisonSvc) StopSyncService() {
	daemon.StopService(s.Filename)
	daemon.DeleteService(s.Filepath)
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

	fmt.Println(bgSyncFailMsg)
}
