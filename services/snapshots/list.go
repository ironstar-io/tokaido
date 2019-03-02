package snapshots

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ryanuber/columnize"
)

// List returns a list of all snapshots with id numbers.
// IDs are based on the list of snapshots in descending order
// IDs will change as more snapshots are added
func List() {
	fmt.Println()
	p := filepath.Join(conf.GetConfig().Tokaido.Project.Path, "/.tok/local/snapshots")

	list := getSortedSnapshotList(p)
	if len(list) == 0 {
		fmt.Println("💢  No snapshots were found")
		fmt.Println("")
		return
	}

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

	fmt.Println("")

	return
}

func getSortedSnapshotList(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		console.Println("\n🙅‍  There was an error getting a list of database snapshots", "")
		log.Fatal(err)
	}
	if len(files) == 0 {
		return []string{}
	}

	// Sort the list by last modified time
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Unix() < files[j].ModTime().Unix()
	})

	// Reverse the order of the slice, so that the newest snapshots appear first
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}

	// Populate our slice with the list of files that are now in ascending order
	list := make([]string, len(files))
	for k, f := range files {
		list[k] = f.Name()
	}

	return list

}
