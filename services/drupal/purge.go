package drupal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/ssh"
)

// Purge will empty the Varnish cache and run a drush cache-rebuild
func Purge() {
	purgeVarnish()
	purgeDrush()
}

func purgeVarnish() {
	port := docker.LocalPort("varnish", "8081")
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://localhost:"+port, nil)
	if err != nil {
		log.Fatalf("    Error when opening HTTP client for Varnish purge: %v", err)
	}

	req.Header.Add("X-PURGEALL", "")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("    Error while trying to send Varnish purge request: %v", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("    Successfully purged Varnish cache")
	} else {
		log.Fatalf("    Sorry, Varnish purge failed. Response: %s", resp.Status)
	}
}

func purgeDrush() {
	var cmd []string
	if conf.GetConfig().Drupal.Majorversion == "7" {
		cmd = []string{"drush", "cache-clear", "all", "-y"}
		fmt.Println("    Sending drush cache-clear all request...")
	} else {
		cmd = []string{"drush", "cache-rebuild", "-y"}
		fmt.Println("    Sending drush cache-rebuild request...")
	}

	resp := ssh.ConnectCommand(cmd)
	fmt.Println(resp)
}
