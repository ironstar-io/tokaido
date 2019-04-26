package tok

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/git"
	"github.com/ironstar-io/tokaido/services/tok/goos"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/utils"
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

	git.Init()
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
		Timeout: time.Second * 3,
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
	utils.CheckCmdHard("wget")
	utils.StdoutStreamCmdDebug("wget", "https://downloads.tokaido.io/packages/"+template.PackageFilename)
	utils.StdoutStreamCmdDebug("tar", "xzf", template.PackageFilename)
	utils.StdoutStreamCmdDebug("rm", template.PackageFilename)
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
	console.SpinPersist(wo, "🚛", "Unpacked the install template")

	// Start the environment
	initialize.TokConfig("up")
	Init(true, false)
	InitMessage()

	wo = console.SpinStart("    Adding everything to a new Git reop...")
	git.AddAll()
	git.Commit("Initial Tokaido Configuration")
	console.SpinPersist(wo, "🚛", "Added everything to a new Git reop")

	// console.Println("\n🍚  Creating a brand new Drupal 8 site with Tokaido!", "")

	// // Create Tokaido configuration
	// conf.SetDrupalConfig("DEFAULT")
	// docker.FindOrCreateTokCompose()
	// docker.CreateDatabaseVolume()
	// docker.CreateSiteVolume()
	// docker.CreateComposerCacheVolume()
	// ssh.GenerateKeys()

	// // Initial directory sync
	// siteVolName := "tok_" + conf.GetConfig().Tokaido.Project.Name + "_tokaido_site"
	// wo := console.SpinStart("Performing an initial sync...")
	// utils.StdoutStreamCmdDebug("docker", "run", "-e", "AUTO_SYNC=false", "-v", conf.GetConfig().Tokaido.Project.Path+":/tokaido/host-volume", "-v", siteVolName+":/tokaido/site", "tokaido/sync:stable")
	// console.SpinPersist(wo, "🚛", "Initial sync completed")

	// // Lift drush container and configure
	// drush.DockerUp()
	// drupal.ConfigureSSH()
	// xdebug.Configure()

	// // Lift background sync container
	// docker.UpMulti([]string{"sync"})

	// console.Println("🎼  Using composer to generate the new Drupal project files. This might take some time", "")
	// // `composer create-project` inside the drush container
	// composerCreateProject()
	// // Pull all required docker images
	// docker.PullImages()

	// wo = console.SpinStart("Waiting for sync to complete. This might take a few minutes...")
	// filepath := drupal.SettingsPath()
	// err := fs.WaitForSync(filepath, 180)
	// if err != nil {
	// 	console.Println("\n🙅‍  Your new Drupal site failed to sync from the Tokaido environment to your local disk", "")
	// 	panic(err)
	// }
	// console.SpinPersist(wo, "🚋", "Secondary sync completed successfully")

	// // Fire up the full Tokaido environment
	// fmt.Println()
	// wo = console.SpinStart("Tokaido is starting your containers")
	// docker.Up()
	// console.SpinPersist(wo, "🚅", "Tokaido started your containers")

	// // Create Tokaido configuration for drupal post composer install

	// // Setup HTTPS proxy service. Retain if statement to preserve Tokaido level enable/disable defaults
	// setupProxy()

	// wo = console.SpinStart("Waiting for settings.tok.php to sync. This might take a few minutes...")
	// filepath = drupal.SettingsTokPath()
	// err = fs.WaitForSync(filepath, 120)
	// if err != nil {
	// 	console.Println("\n🙅‍  Your new Drupal site failed to sync from the Tokaido environment to your local disk", "")
	// 	panic(err)
	// }

	// // Induce a final wait for the sync process to complete
	// time.Sleep(120 * time.Second)

	// console.SpinPersist(wo, "🚋", "Final sync completed successfully")

	// // Drush site install, add additional packages
	// console.Println(`💧  Running drush site-install for your new project`, "")
	// drushSiteInstall(profile)

	// // Generate a new .gitignore file
	// git.NewGitignore()
	// // Git stage all applicable files and commit
	// git.AddAll()
	// git.Commit("Initial Tokaido Configuration")

	// // TODO: `tok status`
	// // TODO: `tok test`

	// // Provide the user a success message
	// goos.InitMessage()

	// // Open the main site at `https://<project-name>.local.tokaido.io:5154`
	// system.OpenTokaidoProxy(false)

}

// NewSuccessMessage ...
func NewSuccessMessage() {
	goos.InitMessage()
}
