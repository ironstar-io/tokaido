package tok

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/git"
	"github.com/ironstar-io/tokaido/services/tok/goos"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
	"github.com/ironstar-io/tokaido/utils"

	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	yaml "gopkg.in/yaml.v2"
)

// Templates is a list of drupal templates available for download
type Templates struct {
	Template []Template `yaml:"templates"`
}

// Template ...
type Template struct {
	Description     string   `yaml:"description"`
	DrupalVersion   int      `yaml:"drupal_version"`
	Maintainer      string   `yaml:"maintainer"`
	Name            string   `yaml:"name"`
	PackageFilename string   `yaml:"package_filename"`
	PostUpCommands  []string `yaml:"post_up_commands,omitempty"`
}

func buildProjectFrame(name string) {
	if fs.CheckExists(name) == true {
		log.Fatalf("\n\nThe directory '%s' already exists. Exiting...", name)
	}

	fs.Mkdir(name)
	err := os.Chdir(name)
	if err != nil {
		panic(err)
	}

	pr, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	conf.RegisterProject(name, pr)

	git.Init(pr)
}

func resolveProjectName(args []string) (name string) {
	if len(args) == 1 {
		name = args[0]
	} else if len(args) < 1 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("")
		fmt.Printf("Enter a directory name for your new project: ")
		name, _ = reader.ReadString('\n')
	} else {
		fmt.Println("ERROR: Received too many arguments. New project name must be a one word string.")
		os.Exit(1)
	}

	// Rationalise the name to strip out any problematic characters
	reg, err := regexp.Compile("[^a-zA-Z0-9\\-]+")
	if err != nil {
		log.Fatal(err)
	}
	name = reg.ReplaceAllString(name, "")

	return name
}

func marshalTemplates() Templates {
	req, err := http.NewRequest(http.MethodGet, "https://downloads.tokaido.io/templates.yaml", nil)
	if err != nil {
		log.Fatalf("Error while trying to retrieve list of available templates: %v", err)
	}

	client := http.Client{
		Timeout: time.Second * 15,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error while trying to retrieve list of available templates: %v", err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatalf("Error while trying to read list of available templates: %v", err)
	}

	templates := Templates{}
	err = yaml.Unmarshal(body, &templates)
	if err != nil {
		log.Fatalf("Error while trying to unmarshal list of available templates: %v", err)
	}

	return templates
}

func chooseTemplate(tp *Templates) (template Template) {
	var templates []Template
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
URL: https://downloads.tokaido.io/packages/{{ .PackageFilename | faint }}
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

func resolveTemplateName(requestTemplate string) (template Template) {
	// Get a list of available templates
	templates := marshalTemplates()

	// If the user didn't specify a template name, give them a list to choose from
	if requestTemplate == "" {
		template = chooseTemplate(&templates)
		return template
	}

	// If the user did specify a template name, see if that template is available
	for _, t := range templates.Template {
		if requestTemplate == t.Name {
			utils.DebugString("[new] found matching template: [" + t.Name + "] at [" + t.PackageFilename + "]")
			return t
		}
	}

	log.Fatalf("Error: could not find a template to launch")
	return
}

func unpackTemplate(template Template, name string) {
	pr, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	utils.CheckCmdHard("curl")
	utils.StdoutStreamCmdDebugContext(pr, "curl", "-O", "https://downloads.tokaido.io/packages/"+template.PackageFilename)
	utils.StdoutStreamCmdDebugContext(pr, "tar", "xzf", template.PackageFilename)
	utils.StdoutStreamCmdDebugContext(pr, "rm", template.PackageFilename)
}

// New - The core run sheet of `tok new {project}`
func New(args []string, requestTemplate string) {

	// Determine the name of the new project
	name := resolveProjectName(args)

	// Determine the template to be used to launch this project
	template := resolveTemplateName(requestTemplate)

	fmt.Printf("\n🍚  Launching new project [%s] with Drupal template [%s]\n", name, template.Name)

	buildProjectFrame(name)

	// Download and untar the install package
	wo := console.SpinStart("Downloading and unpacking template...")
	unpackTemplate(template, name)
	console.SpinPersist(wo, "  ", "Unpacked the install template")

	// Start the environment
	initialize.TokConfig("up")
	Init(true, false)

	// Run the post-up custom installation commands
	l := len(template.PostUpCommands)
	if l > 0 {
		fmt.Println("🏃‍♀  Running post-up installation commands...")
		for _, v := range template.PostUpCommands {
			fmt.Printf("%s", aurora.Blue("    - "))
			fmt.Printf(v + "\n")
			s := strings.Fields(v)
			ssh.StreamConnectCommand(s)
		}
		fmt.Println()
		fmt.Printf("    Ran %d post-up installation commands\n", l)
	}

	wo = console.SpinStart("Setting up a new Git repository...")
	git.AddAll()
	git.Commit("Initial Tokaido Configuration")
	console.SpinPersist(wo, "  ", "Saved everything to a new Git repository")

	InitMessage()

	fmt.Println()
	fmt.Println(aurora.Yellow("You need to move into the project directory before you can these commands: 'cd " + conf.GetConfig().Tokaido.Project.Name + "'"))

}

// NewSuccessMessage ...
func NewSuccessMessage() {
	goos.InitMessage()
}
