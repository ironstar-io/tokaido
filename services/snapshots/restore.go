package snapshots

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/drupal"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"
	"github.com/ryanuber/columnize"
)

// Restore restores the specified snapshot (or lists snapshots to be chosen)
func Restore(args []string) {
	fmt.Println()
	p := filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/.tok/local/snapshots")

	// Get a list of all available snapshots
	list := getSortedSnapshotList(p)
	if len(list) == 0 {
		fmt.Println("💢  No snapshots were found")
		fmt.Println("")
		return
	}

	// If the user specified a snapshot, just restore that one
	if len(args) > 0 {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			console.Println("\n🙅‍  There was an error looking up that database snapshot", "")
			log.Fatal(err)
		}
		restoreSnapshot(list[id])

		return
	}

	// Display the list of snapshots with ID numbers
	output := []string{
		"ID | Snapshot Name",
	}
	for k, v := range list {
		output = append(output, fmt.Sprintf("%d|%s", k, v))
	}

	result := columnize.Format(output, &columnize.Config{
		Delim: "|",
	})
	fmt.Println(result)

	// Ask the user to input an ID number

	fmt.Println("")
	fmt.Printf("Enter the ID of the snapshot to restore: ")
	var input string
	fmt.Scanln(&input)

	id, err := strconv.Atoi(input)
	if err != nil {
		console.Println("\n🙅‍  There was an error looking up that database snapshot\n", "")
		log.Fatal(err)
	}

	if id > len(list) {
		console.Println("\n🙅‍  The ID you specified was not found\n", "")
		return
	}

	restoreSnapshot(list[id])

	return

}

func restoreSnapshot(f string) {
	// Confirm Prompt
	confirmRestore := utils.ConfirmationPrompt("💀  Tokaido will delete your current database and replace it \n    with the snapshot ["+f+"].\n    Are you sure?", "n")
	if confirmRestore == false {
		fmt.Println("🕊  No databases were harmed")
		return
	}

	fmt.Printf("💾  Restoring snapshot: %s\n", f)

	p := "/tokaido/site/.tok/local/snapshots"
	var args []string
	// Determine if we're restoring a compressed backup or not
	if strings.Contains(f, "sql.gz") {
		args = []string{
			"gunzip",
			"<",
			filepath.Join(p, f),
			"|",
			"mysql",
			"-u",
			"root",
			"-ptokaido",
			"-h",
			"mysql",
			"tokaido",
		}
	} else if strings.Contains(f, ".sql") {
		args = []string{
			"mysql",
			"-u",
			"root",
			"-ptokaido",
			"-h",
			"mysql",
			"tokaido",
			"<",
			filepath.Join(p, f),
		}
	} else {
		console.Println("\n🙅‍  The specified snapshot is not in a supported file format. Please supply a .gz or .sql.gz file\n", "")
	}

	ssh.StreamConnectCommand(args)

	fmt.Printf("    Running `tok purge`\n")
	drupal.Purge()

	fmt.Println("")

	return
}
