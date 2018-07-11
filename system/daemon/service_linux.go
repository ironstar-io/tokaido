package daemon

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

type service struct {
	ProjectName string
	ProjectPath string
}

func ReloadServices() {
	_, err := utils.CommandSubSplitOutput("systemctl", "--user", "daemon-reload")
	if err != nil {
		log.Fatal("Unable to reload the systemd config: ", err)
	}
}

func StartService(service string) {
	_, err := utils.CommandSubSplitOutput("systemctl", "--user", "start", service)
	if err != nil {
		log.Fatal("Unable to start the sync service: ", err)
	}
}

func StopService(service string) {
	_, errStatus := utils.CommandSubSplitOutput("systemctl", "--user", "status", service)
	if errStatus == nil {
		_, errStop := utils.CommandSubSplitOutput("systemctl", "--user", "stop", service)
		if errStop != nil {
			log.Fatal("Unable to stop the sync service: ", errStop)
		}
	}
}

func DeleteService(service string) {
	c := conf.GetConfig()
	err := os.Remove(filepath.Join(c.SystemdPath, service))
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		log.Fatal("Unable to remove service configuration: ", err)
	}

	ReloadServices()

	fmt.Println(`
🔄  Removed the background sync process
	`)
}

// CheckSyncService checks if the unison background process is running
func CheckService(service string) string {
	_, err := utils.CommandSubSplitOutput("systemctl", "--user", "status", service)
	if err == nil {
		return "running"
	}

	if c.Debug == true {
		fmt.Printf("\033[33m%s\033[0m\n", err)
	}

	return "stopped"
}

// KillService ...
func KillService(service string) {
	StopService(service)
	DeleteService(service)
}
