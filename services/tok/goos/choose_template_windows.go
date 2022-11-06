//go:build windows
// +build windows

package goos

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ironstar-io/tokaido/services/tok/types"
	"github.com/ironstar-io/tokaido/system/console"

	"github.com/ryanuber/columnize"
)

// ChooseTemplate - Use a simple numeric selector to select templates for Windows systems
func ChooseTemplate(tp *types.Templates) (template types.Template) {
	var templates []types.Template
	size := 0

	// Convert our templates struct into a more useable format
	for _, t := range tp.Template {
		size = size + 1
		templates = append(templates, t)
	}

	// Display the list of snapshots with ID numbers
	output := []string{
		"ID | Template | Description",
	}
	for k, t := range templates {
		output = append(output, fmt.Sprintf("%d|%s|%s", k, t.Name, t.Description))
	}

	fmt.Println("Please choose the Drupal template you'd like to launch")
	fmt.Println()

	result := columnize.Format(output, &columnize.Config{
		Delim: "|",
	})
	fmt.Println(result)

	// Ask the user to input an ID number

	fmt.Println("")
	fmt.Printf("Enter the ID of the template to use: ")
	var input string
	fmt.Scanln(&input)

	id, err := strconv.Atoi(input)
	if err != nil {
		console.Println("\n🙅‍  That wasn't a valid selection\n", "")
		log.Fatal(err)
	}

	if id > len(templates) {
		console.Println("\n🙅‍  The ID you specified was not found\n", "")
		return
	}

	return templates[id]
}
