package auth

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/system/fs"

	"github.com/olekukonko/tablewriter"
	yaml "gopkg.in/yaml.v2"
)

var ShowCredentialsErrorMsg = "Unable to display credentials"

func IronstarShowCredentials(args []string) error {
	credSet, err := ReadInCredentials()
	if err != nil {
		return err
	}

	globals, err := ReadInGlobals()
	if err != nil {
		return err
	}

	if globals.DefaultLogin == "" {
		globals.DefaultLogin = "UNSET"
	}
	fmt.Println()
	fmt.Println("Default Login: " + globals.DefaultLogin)
	fmt.Println()
	fmt.Println("Project Credentials:")

	pcreds := make([][]string, len(globals.Projects))

	for _, v := range globals.Projects {
		if v.Login == "" {
			v.Login = "UNSET"
		}
		pcreds = append(pcreds, []string{v.Name, v.Path, v.Login})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Path", "Login"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(pcreds)
	table.Render()

	fmt.Println()

	fmt.Println("Available Credentials:")

	acreds := make([][]string, len(credSet))

	for _, v := range credSet {
		acreds = append(acreds, []string{v.Login, v.Expiry.String()})
	}

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Login", "Expiry"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(acreds)
	table.Render()

	return nil
}

func ReadInGlobals() (conf.Global, error) {
	globals := conf.Global{}
	gp := filepath.Join(fs.HomeDir(), ".tok", "global.yml")

	err := SafeTouchGlobalConfigYAML("global")
	if err != nil {
		return globals, err
	}

	gBytes, err := ioutil.ReadFile(gp)
	if err != nil {
		return globals, err
	}

	err = yaml.Unmarshal(gBytes, &globals)
	if err != nil {
		return globals, err
	}

	return globals, nil
}
