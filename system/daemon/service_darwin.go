package daemon

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"bitbucket.org/ironstar/tokaido-cli/conf"
	"bitbucket.org/ironstar/tokaido-cli/utils"
)

type service struct {
	ProjectName string
	ProjectPath string
	Username    string
}

// LoadService ...
func LoadService(service string) {
	_, err := utils.CommandSubSplitOutput("launchctl", "load", service)
	if err != nil {
		log.Fatal("Unable to load sync service: ", err)
	}
}

// UnloadService ...
func UnloadService(service string) {
	_, err := utils.CommandSubSplitOutput("launchctl", "unload", service)
	if err != nil {
		log.Fatal("Unable to unload sync service: ", err)
	}
}

// StartService ...
func StartService(service string) {
	_, err := utils.CommandSubSplitOutput("launchctl", "start", service)
	if err != nil {
		log.Fatal("Unable to start the sync service: ", err)
	}
}

// StopService ...
func StopService(service string) {
	_, err := utils.CommandSubSplitOutput("launchctl", "stop", service)
	if err != nil {
		log.Fatal("Unable to stop the sync service: ", err)
	}
}

// DeleteService ...
func DeleteService(service string) {
	c := conf.GetConfig()
	os.Remove(filepath.Join(c.LaunchdPath, service))
}

// CheckService checks if the unison background process is running
func CheckService(service string) string {
	c := conf.GetConfig()

	u, uErr := user.Current()
	if uErr != nil {
		log.Fatal(uErr)
	}

	o, _ := utils.CommandSubSplitOutput("launchctl", "print", "gui/"+u.Uid+service)

	if c.Debug == true {
		fmt.Printf("\033[33m%s\033[0m\n", o)
	}

	if strings.Contains(o, "state = running") == true {
		return "running"
	}

	return "stopped"
}

// KillService ...
func KillService(service string) {
	ps, _ := utils.CommandSubSplitOutput("launchctl", "list", service)
	if ps != "" {
		fmt.Println(`
🔄  Removing the background sync process
	`)
		StopService(service)
		UnloadService(service)
		DeleteService(service)
	}
}
