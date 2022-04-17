//go:build !windows
// +build !windows

package goos

import (
	"fmt"
	"log"

	"github.com/ironstar-io/tokaido/services/tok/types"

	"github.com/manifoldco/promptui"
)

// ChooseTemplate - Use promptui to select templates for unix systems
func ChooseTemplate(tp *types.Templates) (template types.Template) {
	var templates []types.Template
	size := 0

	// Convert our templates struct into a more useable format
	for _, t := range tp.Template {
		size = size + 1
		templates = append(templates, t)
	}

	// Define the menu's visual template
	menuTemplate := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }}`,
		Inactive: `   {{ .Name | cyan }}`,
		Selected: "{{ .Name | blue | cyan }}",
		Details: `---------
{{ .Description | faint  }}

Maintainer: {{ .Maintainer | faint }}
`,
	}

	fmt.Println("Please choose the Drupal template you'd like to launch")

	prompt := promptui.Select{
		Label:     "Templates >>",
		Items:     templates,
		Templates: menuTemplate,
		Size:      size,
	}

	i, _, err := prompt.Run()

	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)

	}

	return templates[i]
}
