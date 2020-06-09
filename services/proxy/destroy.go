package proxy

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"

	"github.com/logrusorgru/aurora"
)

// DecommissionProxy entirely removes the proxy network from this host
// This is part of the 1.13 upgrade process
func DecommissionProxy() {
	// check if the proxy exists
	home := fs.HomeDir()
	if !fs.CheckExists(filepath.Join(home, ".tok/proxy")) {
		return
	}

	fmt.Printf("%s", aurora.Blue(`
Welcome to Tokaido 1.13! This version removes the Tokaido proxy service
(which lives at *.local.tokaido.io:5154) and replaces it with direct connections to each project.

The Tokaido Proxy service offered a fixed port number (5154) but with the side-effect of masking problems
that might be happening with your site. To make it easier for your get real insight into your Drupal
application faster, we've decided to remove the Proxy service starting from this version.

Tokaido will now attempt to remove the proxy service from your system. This requires stopping and cleaning up
all of your running Tokaido environments. Can Tokaido do this now? [Y/n] `))
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(string([]rune(response)[0]))
	if response == "n" {
		fmt.Println(aurora.Red("Upgrade process interrupted. Tokaido can not start."))
		os.Exit(1)
	}

	// Search for a list of Tokaido projects - see which ones are running
	fmt.Printf("\nSearching for Tokaido projects... \n")
	c := conf.GetConfig()

	for _, v := range c.Global.Projects {
		// Check if the project is running
		fmt.Printf("Stopping project %s...\n", v.Name)

		cmd := exec.Command("docker-compose", "-f", "docker-compose.tok.yml", "rm", "-sf")
		cmd.Dir = v.Path
		cmd.CombinedOutput()
	}

	fmt.Println("Removing the Tokaido Proxy")
	cmd := exec.Command("docker-compose", "-f", "docker-compose.yml", "rm", "-sf")
	cmd.Dir = filepath.Join(home, ".tok/proxy")
	cmd.CombinedOutput()

	cmd = exec.Command("rm", "-rf", filepath.Join(home, ".tok/proxy"))
	cmd.CombinedOutput()

	cmd = exec.Command("docker", "network", "rm", "tokaido_proxy")
	cmd.CombinedOutput()

	cmd = exec.Command("docker", "network", "rm", "proxy_default")
	cmd.CombinedOutput()

	fmt.Println(aurora.Green("Upgrade complete. Thank you for using Tokaido!"))

	return

}
