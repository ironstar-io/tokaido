package tok

import (
	"fmt"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/services/docker"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/logrusorgru/aurora"
	"github.com/ryanuber/columnize"
)

// List will return a list of all known local Tokaido projects and their current status
func List() {
	missing := false // If set to true, we'll output helpful information

	o := []string{}
	o = append(o, fmt.Sprintf("%s|%s|%s", "Name", "Path", "Status"))
	o = append(o, "----|----|------")

	for _, v := range conf.GetConfig().Global.Projects {
		ok := docker.StatusCheck("", v.Name)
		if ok {
			o = append(o, fmt.Sprintf("%s|%s|%s", v.Name, v.Path, aurora.Green("running")))
		} else {
			x := fs.CheckExists(v.Path)
			if !x {
				o = append(o, fmt.Sprintf("%s|%s|%s", v.Name, v.Path, aurora.Red("missing")))
				missing = true
			} else {
				o = append(o, fmt.Sprintf("%s|%s|%s", v.Name, v.Path, aurora.Yellow("offline")))
			}
		}
	}
	fmt.Println()

	cc := columnize.DefaultConfig()
	cc.Delim = "|"
	cc.Glue = "  "
	cc.Prefix = ""
	cc.Empty = ""
	cc.NoTrim = false

	result := columnize.Format(o, cc)
	fmt.Println(result)

	if missing {
		fmt.Println()
		fmt.Println(fmt.Sprintf("'missing' projects indicate that the project was deleted without using '%s'", aurora.Bold("tok destroy")))
		fmt.Println("See https://docs.tokaido.io/tokaido/cleaning-up-tokaido-environments")
	}
	fmt.Println()

}
