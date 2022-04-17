//go:build darwin
// +build darwin

package goos

import (
	"bytes"
	"errors"
	"log"
	"os/user"
	"path/filepath"
	"text/template"

	unisontmpl "github.com/ironstar-io/tokaido/services/unison/templates"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/daemon"
	"github.com/ironstar-io/tokaido/system/fs"
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
	u, uErr := user.Current()
	if uErr != nil {
		log.Fatal(uErr)
	}

	s := UnisonSvc{
		SyncName:    syncName,
		SyncDir:     syncDir,
		Filename:    getServiceFilename(syncName),
		Filepath:    getServicePath(syncName),
		Launchdpath: filepath.Join(fs.HomeDir(), "/Library/LaunchAgents/"),
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
	return filepath.Join(filepath.Join(fs.HomeDir(), "/Library/LaunchAgents/"), getServiceFilename(syncName))
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

	fs.TouchOrReplace(filepath.Join(s.Launchdpath, s.Filename), tpl.Bytes())
}

// CreateSyncService Register a launchd or systemctl service for Unison active sync
func (s UnisonSvc) CreateSyncService() {
	s.RegisterSyncService()
	s.StartSyncService()
}

// RegisterSyncService Register the unison sync service for launchd
func (s UnisonSvc) RegisterSyncService() {
	s.CreateSyncFile()

	daemon.LoadService(s.Filepath)
}

// UnloadSyncService Remove the unison sync service from launchd
func (s UnisonSvc) UnloadSyncService() {
	daemon.UnloadService(s.Filepath)
}

// StartSyncService Start the launchd service after it is created
func (s UnisonSvc) StartSyncService() {
	daemon.StartService(s.Filename)
}

// RestartSyncService ...
func (s UnisonSvc) RestartSyncService() {
	daemon.StopService(s.Filename)
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
func (s UnisonSvc) CheckSyncService() error {
	c := s.SyncServiceStatus()
	if c == "running" {
		console.Println("✅  Background sync service is running", "√")
		return nil
	}

	console.Println(bgSyncFailMsg, "")
	return errors.New(bgSyncFailMsg)
}
