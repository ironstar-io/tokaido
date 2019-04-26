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

	"github.com/ironstar-io/tokaido/conf"
	"github.com/ironstar-io/tokaido/constants"
	"github.com/ironstar-io/tokaido/initialize"
	"github.com/ironstar-io/tokaido/services/git"
	"github.com/ironstar-io/tokaido/services/proxy"
	"github.com/ironstar-io/tokaido/services/tok/goos"
	"github.com/ironstar-io/tokaido/system/console"
	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/ssh"
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
	Description    string   `yaml:"description"`
	DrupalVersion  int      `yaml:"drupal_version"`
	Maintainer     string   `yaml:"maintainer"`
	Name           string   `yaml:"name"`
	PackageURL     string   `yaml:"package_url"`
	PostUpCommands []string `yaml:"post_up_commands,omitempty"`
}

var profileFlag string

func buildProjectFrame(projectName string) {
	if fs.CheckExists(projectName) == true {
		log.Fatal("The project '" + projectName + "' already exists in this directory. Exiting...")
	}

	fs.Mkdir(projectName)
	err := os.Chdir(projectName)
	if err != nil {
		panic(err)
	}

	git.Init()
	initialize.TokConfig("new")
}

func composerCreateProject() {
	ssh.ConnectCommand([]string{"mkdir", "/tmp/composer"})
	ssh.StreamConnectCommand([]string{"composer", "create-project", "ironstar-io/d8-template:0.4", "/tmp/composer", "--stability", "dev", "--no-interaction"})
	ssh.ConnectCommand([]string{"cp", "-R", "/tmp/composer/*", "/tokaido/site"})
}

func setupProxy() {
	c := conf.GetConfig()
	if c.System.Syncsvc.Enabled && c.System.Proxy.Enabled {
		console.Println("🔐  Setting up HTTPS for your local development environment", "")
		proxy.Setup()
	}
}
func drushSiteInstall(profile string) {
	switch profile {
	case "":
		profile = "standard"
		fallthrough
	case
		"standard",
		"minimal",
		"demo_umami":
		ssh.StreamConnectCommand([]string{"drush", "site-install", profile, "-y"})
		ssh.StreamConnectCommand([]string{"drush", "en", "swiftmailer", "password_policy", "password_policy_character_types", "password_policy_characters", "password_policy_username", "memcache", "health_check"})
		return
	}
	log.Fatalf("Error: The install profile specified was not supported. Possible values are 'standard', 'minimal', and 'demo_umami'")
}

func deduceProjectName(args []string) string {
	if len(args) == 0 {
		return constants.DefaultDrupalProjectName
	}

	return args[0]
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
PackageURL: {{ .PackageURL | faint }}
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
			utils.DebugString("[new] found matching template: [" + t.Name + "] at [" + t.PackageURL + "]")
			return t
		}
	}

	fmt.Println("found no template")
	return
}

// New - The core run sheet of `tok new {project}`
func New(args []string, requestTemplate string) {

	// Determine the name of the new project
	name := resolveProjectName(args)

	// Determine the template to be used to launch this project
	template := resolveTemplateName(requestTemplate)

	fmt.Printf("Launching new project [%s] with Drupal template [%s]", name, template.Name)

	// pn := deduceProjectName(args)
	// buildProjectFrame(pn)

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
	// drupal.CheckSettings("FORCE")

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
