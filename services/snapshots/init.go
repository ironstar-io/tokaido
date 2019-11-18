package snapshots

import (
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"
)

// Init preforms any db snapshot initalization functions that need to be run with general tok init functions
func Init() (err error) {
	mkSnapshotDir()

	return nil
}

func mkSnapshotDir() {
	sd := filepath.Join(conf.GetProjectPath(), ".tok", "local", "snapshots")

	fs.Mkdir(sd)
}
