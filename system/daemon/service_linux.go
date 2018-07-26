package daemon

import (
	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/system/console"
	"bitbucket.org/ironstar/tokaido-cli/utils"

	"fmt"
	"log"
	"os"
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

func StartService(serviceName string) {
	_, err := utils.CommandSubSplitOutput("systemctl", "--user", "start", serviceName)
	if err != nil {
		log.Fatal("Unable to start the sync service: ", err)
	}
}

func StopService(serviceName string) {
	_, errStatus := utils.CommandSubSplitOutput("systemctl", "--user", "status", serviceName)
	if errStatus == nil {
		_, errStop := utils.CommandSubSplitOutput("systemctl", "--user", "stop", serviceName)
		if errStop != nil {
			log.Fatal("Unable to stop the sync service: ", errStop)
		}
	}
}

func DeleteService(servicePath string) {
	err := os.Remove(servicePath)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		log.Fatal("Unable to remove service configuration: ", err)
	}

	ReloadServices()

	console.Println(`
🔄  Removed the background sync process
	`, "")
}

// ServiceStatus checks if the unison background process is running
func ServiceStatus(serviceName string) string {
	_, err := utils.CommandSubSplitOutput("systemctl", "--user", "status", serviceName)
	if err == nil {
		return "running"
	}

	if conf.GetConfig().Debug == true {
		fmt.Printf("\033[33m%s\033[0m\n", err)
	}

	return "stopped"
}

// KillService ...
func KillService(serviceName string, servicePath string) {
	StopService(serviceName)
	DeleteService(servicePath)
}
