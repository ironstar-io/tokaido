package snapshots

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
)

// New preforms a DB Snapshot and saves it to .tok/local/snapshots
func New(args []string) {
	var name string

	// Check that this Tokaido environment is running
	err := docker.StatusCheck()
	if err != nil {
		console.Println("\n😦  Your Tokaido environment appears to be offline. Please run `tok up` to start it.", "")
		return
	}

	fmt.Println()
	wo := console.SpinStart("Backing up your Tokaido database")

	if len(args) == 0 {
		name = "tokaido"
	} else {
		reg, err := regexp.Compile("[^a-zA-Z0-9\\-]+")
		if err != nil {
			log.Fatal(err)
		}

		name = strings.Join(args, "-")
		name = "tokaido_" + reg.ReplaceAllString(name, "")

	}
	filename, _ := createSnapshot(name)
	filepath := filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/.tok/local/snapshots/", filename)

	err = fs.WaitForSync(filepath, 180)
	if err != nil {
		console.Println("\n🙅‍  Your backup failed to sync from the Tokaido environment to your local disk", "")
		panic(err)
	}

	console.SpinPersist(wo, "💾", "Your backup is available at .tok/local/snapshots/"+filename+"\n")

	return
}

func preSnapshotChecks() (err error) {
	p := filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/.tok/local/snapshots")
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		mkSnapshotDir()
	}

	return nil
}

func createSnapshot(name string) (filename string, err error) {
	now := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
	filename = name + "_" + now + ".sql.gz"

	args := []string{
		"mysqldump",
		"-u",
		"root",
		"-ptokaido",
		"-h",
		"mysql",
		"tokaido",
		"--ignore-table tokaido.cache",
		"--ignore-table tokaido.cache_block",
		"--ignore-table tokaido.cache_filter",
		"--ignore-table tokaido.cache_form",
		"--ignore-table tokaido.menu",
		"--ignore-table tokaido.cache_page",
		"--ignore-table tokaido.cache_update",
		"--ignore-table tokaido.history",
		"--ignore-table tokaido.watchdog",
		"--max_allowed_packet=1073741824",
		"|",
		"gzip",
		"-9",
		">",
		"/tokaido/site/.tok/local/snapshots/" + filename,
	}
	ssh.StreamConnectCommand(args)

	return filename, nil
}

// waitForSync polls the snapshots dir locally until the sync job finishes or times out
func waitForSync(filename string) (err error) {
	retries := 180

	for i := 1; i <= retries; i++ {
		p := filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/.tok/local/snapshots/", filename)
		_, err = os.Stat(p)
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return err

}
