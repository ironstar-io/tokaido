package auth

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
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
