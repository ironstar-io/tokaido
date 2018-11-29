// +build linux

package goos

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/unison/templates"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/daemon"
	"github.com/ironstar-io/tokaido/system/fs"
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

// CreateSyncFile creates the systemd path (if necessary) and file
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

	// Create the systemd path if it doesn't not exist
	if _, pErr := os.Stat(s.Systemdpath); os.IsNotExist(pErr) {
		mErr := os.MkdirAll(s.Systemdpath, os.ModePerm)
		if mErr != nil {
			log.Fatal("There was an error creating the systemd path:", mErr)
		}
	}

	fs.TouchOrReplace(filepath.Join(s.Systemdpath, s.Filename), tpl.Bytes())

	daemon.ReloadServices()
}

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func (s UnisonSvc) CreateSyncService() {
	s.RegisterSyncService()
	s.StartSyncService()
}

// RegisterSyncService Register the unison sync service for systemd
func (s UnisonSvc) RegisterSyncService() {
	s.CreateSyncFile()
}

// StartSyncService Start the systemd service after it is created
func (s UnisonSvc) StartSyncService() {
	daemon.StartService(s.Filename)
}

// UnloadSyncService Remove the unison sync service
func (s UnisonSvc) UnloadSyncService() {
	return
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
func (s UnisonSvc) CheckSyncService() error {
	if conf.GetConfig().System.Syncsvc.Enabled != true {
		return nil
	}

	c := s.SyncServiceStatus()
	if c == "running" {
		console.Println("✅  Background sync service is running", "√")
		return nil
	}

	fmt.Println(bgSyncFailMsg)
	return errors.New(bgSyncFailMsg)
}
