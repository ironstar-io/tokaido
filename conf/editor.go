package conf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

// ConfigRoot ...
type ConfigRoot struct {
	Name        string
	Description string
	Detail      string
}

// ConfigGenericString ...
type ConfigGenericString struct {
	Name    string
	Default string
	Current string
	Type    string
	Detail  string
}

// reloadConfig will reload the viper project config from disk
func reloadConfig() {
	reloadProjectConfig()
	reloadGlobalConfig()
}

func reloadProjectConfig() {
	pr := fs.ProjectRoot()
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Join(pr, ".tok"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error reading project config file: %v", err)
	}
}

func reloadGlobalConfig() {
	h := fs.HomeDir()
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(filepath.Join(h, ".tok"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error reading global config file: %v", err)
	}
}

func newStringValue(label string) string {
	templates := &promptui.PromptTemplates{
		Prompt: "{{ . }} ",
	}

	prompt := promptui.Prompt{
		Label:     label,
		Templates: templates,
	}

	dp, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return dp

}

// MainMenu is the root menu object for the `tok config` command. It is used
// to view and edit the Tokaido config
func MainMenu() {
	menu := []ConfigRoot{
		{
			Name:        "Tokaido Project Settings »",
			Description: "Tokaido settings that define how this project built",
			Detail:      "These settings tell Tokaido how this project should be built and managed. \nIt includes essential items like if Tokaido should use beta Docker images and if\nyou want to self-manage your Docker Compose file.\n\n\n\n",
		},
		{
			Name:        "Database Engine Settings »",
			Description: "Configure your Database Engine and Version",
			Detail:      "Tokaido includes support for MySQL and MariaDB. The default is currently MySQL 5.7\n\n\n\n",
		},
		{
			Name:        "Drupal Settings »",
			Description: "Simple Drupal settings that Tokaido needs, like your document root",
			Detail:      "Tokaido needs to know a little bit about your project.\nYou almost never need to edit these settings, and doing \nso could break your installation",
		},
		{
			Name:        "Nginx Settings »",
			Description: "Nginx config settings that can be controlled from your codebase",
			Detail:      "Some Nginx-level settings can be read directly from your code \n(rather than being managed from the hosting environment).\nThis menu will show let you modify Nginx settings \nthat you define right here in the repo, and ship them into \nyour production Tokaido environment.",
		},
		{
			Name:        "PHP FPM Settings »",
			Description: "FPM config settings that can be controlled from your codebase",
			Detail:      "Some PHP FPM-level settings can be read directly from your code \n(rather than being managed from the hosting environment).\nThis menu will show let you modify PHP settings \nthat you define right here in the repo, and ship them into \nyour production Tokaido environment.",
		},
		{
			Name:        "Services »",
			Description: "Enable or disable additional Tokaido services like Solr and Mailhog",
			Detail:      "Tokaido is packed with added services like Solr, Mailhog, PimpMyLog, and \nmuch more. In this menu, you can enable or disable these services",
		},
		{
			Name:        "Exit",
			Description: "",
			Detail:      "",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }} {{ if ne .Name "Exit" }} - {{ .Description }} {{ end }}`,
		Inactive: `   {{ .Name | cyan }} {{ if ne .Name "Exit" }} - {{ .Description }} {{ end }}`,
		Selected: "{{ .Name | blue | cyan }}",
		Details: `
{{ if ne .Name "Exit" }}---------
{{ .Detail | faint  }}
{{ end }}`,
	}

	fmt.Println("Please choose the Tokaido config area you'd like to edit")

	prompt := promptui.Select{
		Label:     "Main Menu >>",
		Items:     menu,
		Templates: templates,
		Size:      8,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		TokaidoMenu()
	case 1:
		DatabaseMenu()
	case 2:
		DrupalMenu()
	case 3:
		NginxMenu()
	case 4:
		FpmMenu()
	case 5:
		ServicesMenu()
	case 6:
		fmt.Println("Please note that if you have made config changes, you need to run `tok rebuild`")
		os.Exit(0)
	}
}

// TokaidoMenu is exposes Tokaido-level config settings
func TokaidoMenu() {
	gprj, err := GetGlobalProjectSettings()
	if err != nil {
		panic(err)
	}
	menu := []ConfigGenericString{
		{
			Name:    "Use Custom Compose File",
			Type:    "value",
			Default: "false",
			Current: strconv.FormatBool(GetConfig().Tokaido.Customcompose),
			Detail:  "If true, Tokaido will no longer update the docker-compose.tok.yml file.\nUse this if you want complete control over your Docker environment,\nbut please note that this will stop Tokaido from being able to add\nnew features and most Tokaido config settings will stop having any effect.",
		},
		{
			Name:    "Stability Release Set",
			Default: "edge",
			Type:    "value",
			Current: GetConfig().Tokaido.Stability,
			Detail:  "Choose between 'edge', 'stable', or 'experimental'. Edge, the default, runs images that are considered to be stable but are being tested for up to one month before being promoted to production (stable).\nLearn more at docs.tokaido.io.",
		},
		{
			Name:    "PHP Version",
			Default: "7.4",
			Type:    "value",
			Current: GetConfig().Tokaido.Phpversion,
			Detail:  "Use the latest version of PHP 7.2, 7.3, or 7.4",
		},
		{
			Name:    "PHP XDebug Support",
			Default: "false",
			Type:    "value",
			Current: strconv.FormatBool(gprj.Xdebug.Enabled),
			Detail:  "Set to 'true' to enable Xdebug support in the FPM and Admin containers",
		},
		{
			Name:    "PHP FPM XDebug Port",
			Default: "9000",
			Type:    "value",
			Current: strconv.Itoa(gprj.Xdebug.FpmPort),
			Detail:  "Set to port number your IDE is listening for incoming Xdebug connections for this project",
		},
		{
			Name:    "Use Emojis",
			Type:    "value",
			Default: "true",
			Current: strconv.FormatBool(GetConfig().Tokaido.Enableemoji),
			Detail:  "You might have noticed we like to use emoji icons. Cool, huh? 😎\nBut some terminal emulators can't display emojis, so set this to false to stop Tokaido from being so cool.",
		},
		{
			Name:    "Depencendy Checks",
			Type:    "value",
			Default: "true",
			Current: strconv.FormatBool(GetConfig().Tokaido.Dependencychecks),
			Detail:  "Turn this feature off to stop Tokaido from performing system dependency checks\nwhenever it is run. This may be helpful if Tokaido is falsely reporting a \ndependency check failure, when your system is still capabpel of running Tokaido.",
		},
		{
			Name:    "« Main Menu",
			Type:    "menu",
			Default: "",
			Current: "",
			Detail:  "Go back to the Main Menu",
		},
		{
			Name:    "Exit",
			Type:    "menu",
			Default: "",
			Current: "",
			Detail:  "Stop editing your configuration",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }} {{ if ne .Type "menu" }} {{ if eq .Current .Default }} Using default value [{{ .Default | cyan }}] {{ else }} Using custom value [{{ .Current | green }}] {{ end }} {{ end }}`,
		Inactive: `   {{ .Name | cyan }} {{ if ne .Type "menu" }} {{ if eq .Current .Default }} Using default value [{{ .Default | cyan }}] {{ else }} Using custom value [{{ .Current | green }}] {{ end }} {{ end }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "menu" }}---------
{{ .Detail }}

Default Setting: [{{ .Default | cyan }}]
Current Setting: [{{ .Current | green }}]
{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » Tokaido Configuration",
		Items:     menu,
		Templates: templates,
		Size:      10,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		if GetConfig().Tokaido.Customcompose == true {
			SetConfigValueByArgs([]string{"tokaido", "customcompose", "false"}, "project")
		} else {
			SetConfigValueByArgs([]string{"tokaido", "customcompose", "true"}, "project")
		}
		reloadConfig()
		TokaidoMenu()
	case 1:
		TokaidoStabilityMenu()
	case 2:
		TokaidoPhpversionMenu()
	case 3:
		if gprj.Xdebug.Enabled == true {
			SetGlobalConfigValueByArgs([]string{"global", "project", "xdebug", "enabled", "false"})
		} else {
			SetGlobalConfigValueByArgs([]string{"global", "project", "xdebug", "enabled", "true"})
		}
		reloadConfig()
		TokaidoMenu()
	case 4:
		res := newStringValue("Specify the port that your IDE is listening on for xdebug from PHP FPM:")
		SetGlobalConfigValueByArgs([]string{"global", "project", "xdebug", "fpmport", res})
		reloadConfig()
		TokaidoMenu()
	case 5:
		if GetConfig().Tokaido.Enableemoji == true {
			SetConfigValueByArgs([]string{"tokaido", "enableemoji", "false"}, "project")
		} else {
			SetConfigValueByArgs([]string{"tokaido", "enableemoji", "true"}, "project")
		}
		reloadConfig()
		TokaidoMenu()
	case 6:
		if GetConfig().Tokaido.Dependencychecks == true {
			SetConfigValueByArgs([]string{"tokaido", "dependencychecks", "false"}, "project")
		} else {
			SetConfigValueByArgs([]string{"tokaido", "dependencychecks", "true"}, "project")
		}
		reloadConfig()
		TokaidoMenu()
	case 7:
		MainMenu()
	case 8:
		fmt.Println("Please note that if you have made config changes, you need to run `tok rebuild`")
		os.Exit(0)
	}
}

// TokaidoStabilityMenu is exposes Tokaido-level config settings
func TokaidoStabilityMenu() {
	menu := []ConfigGenericString{
		{
			Name:   "Use Stable Releases",
			Type:   "value",
			Detail: "Use only the stable Tokaido images. These are production-ready. It's recommended to use the 'edge' version so you can catch errors before going into production on Tokaido-based hosting platforms like Ironstar.",
		},
		{
			Name:   "Use Edge Releases",
			Type:   "value",
			Detail: "Use the Edge Tokaido releases, which is recommended. These are considred to be stable, and will be used for one month until being made 'stable' and deployed to production.",
		},
		{
			Name:   "Use Experimental Releases",
			Type:   "value",
			Detail: "Use the Experimental images. These are under active development, and are useful for testing new features before we release them to the wider public.",
		},
		{
			Name:   "« Tokaido Config",
			Type:   "menu",
			Detail: "Go back to the Main Menu",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }}`,
		Inactive: `   {{ .Name | cyan }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "menu" }}---------
{{ .Detail }}

{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » Tokaido Configuration » Stability Release Set",
		Items:     menu,
		Templates: templates,
		Size:      4,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		SetConfigValueByArgs([]string{"tokaido", "stability", "stable"}, "project")
		reloadConfig()
		TokaidoMenu()
	case 1:
		SetConfigValueByArgs([]string{"tokaido", "stability", "edge"}, "project")
		reloadConfig()
		TokaidoMenu()
	case 2:
		SetConfigValueByArgs([]string{"tokaido", "stability", "experimental"}, "project")
		reloadConfig()
		TokaidoMenu()
	case 3:
		TokaidoMenu()
	}
}

// TokaidoPhpversionMenu is exposes Tokaido-level config settings
func TokaidoPhpversionMenu() {
	menu := []ConfigGenericString{
		{
			Name:   "PHP 7.2",
			Type:   "value",
			Detail: "Enable the latest PHP 7.2 release at the time that this version of Tokaido was created",
		},
		{
			Name:   "PHP 7.3",
			Type:   "value",
			Detail: "Enable the latest PHP 7.3 release at the time that this version of Tokaido was created",
		},
		{
			Name:   "PHP 7.4",
			Type:   "value",
			Detail: "Enable the latest PHP 7.4 release at the time that this version of Tokaido was created",
		},
		{
			Name:   "« Tokaido Config",
			Type:   "menu",
			Detail: "Go back to the Main Menu",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }}`,
		Inactive: `   {{ .Name | cyan }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "menu" }}---------
{{ .Detail }}

{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » Tokaido Configuration » Stability Release Set",
		Items:     menu,
		Templates: templates,
		Size:      4,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		SetConfigValueByArgs([]string{"tokaido", "phpversion", "7.2"}, "project")
		reloadConfig()
		TokaidoMenu()
	case 1:
		SetConfigValueByArgs([]string{"tokaido", "phpversion", "7.3"}, "project")
		reloadConfig()
		TokaidoMenu()
	case 2:
		SetConfigValueByArgs([]string{"tokaido", "phpversion", "7.4"}, "project")
		reloadConfig()
		TokaidoMenu()
	case 3:
		TokaidoMenu()
	}
}

// DatabaseMenu ...
func DatabaseMenu() {
	gprj, err := GetGlobalProjectSettings()
	if err != nil {
		panic(err)
	}
	fmt.Println(aurora.BgRed(aurora.White("!! WARNING !!")))
	fmt.Println(aurora.BgRed(aurora.White("Changing your database configuration can completely break your database. We strongly recommend creating a `tok snapshot` first")))
	fmt.Println(aurora.BgRed(aurora.White("!! WARNING !!")))
	fmt.Println()
	menu := []ConfigGenericString{
		{
			Name:    "Database Engine",
			Type:    "menu",
			Current: GetConfig().Database.Engine,
			Default: "mysql",
			Detail:  "This is the current database engine. Wherever possible, you should set this to match what you're running in Production",
		},
		{
			Name:    "MySQL Version",
			Type:    "menu",
			Current: GetConfig().Database.Mysqlconfig.Version,
			Default: "5.7",
			Detail:  "Set the MySQL Version if MySQL is used",
		},
		{
			Name:    "MariaDB Version",
			Type:    "menu",
			Current: GetConfig().Database.Mariadbconfig.Version,
			Default: "10.3",
			Detail:  "Set the MariaDB Version if MariaDB is used",
		},
		{
			Name:    "Set static database port",
			Type:    "value",
			Current: strconv.Itoa(gprj.Database.Port),
			Detail:  "Instead of using a dynamic random port number, always expose this project's database locally on a pre-defined port",
		},
		{
			Name:    "« Main Menu",
			Type:    "mainmenu",
			Default: "",
			Current: "",
			Detail:  "Go back to the Main Menu",
		},
		{
			Name:    "Exit",
			Type:    "mainmenu",
			Default: "",
			Current: "",
			Detail:  "Stop editing your configuration",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }} {{ if ne .Type "mainmenu" }} {{if or (eq .Current .Default) (eq .Current "") }} Using default value [{{ .Default | cyan }}] {{ else }} Using custom value [{{ .Current | green }}] {{ end }} {{ end }}`,
		Inactive: `   {{ .Name | cyan }} {{ if ne .Type "mainmenu" }} {{if or (eq .Current .Default) (eq .Current "") }} Using default value [{{ .Default | cyan }}] {{ else }} Using custom value [{{ .Current | green }}] {{ end }} {{ end }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "mainmenu" }}---------
{{ .Detail }}

{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » Database Configuration",
		Items:     menu,
		Templates: templates,
		Size:      5,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		DatabaseEngineMenu()
		reloadConfig()
		DatabaseMenu()
	case 1:
		DatabaseMysqlconfigMenu()
		reloadConfig()
		DatabaseMenu()
	case 2:
		DatabaseMariadbconfigMenu()
		reloadConfig()
		DatabaseMenu()
	case 3:
		res := newStringValue("Specify the port to always use for this project:")
		SetGlobalConfigValueByArgs([]string{"global", "project", "database", "port", res})
		reloadConfig()
		DatabaseMenu()
	case 4:
		MainMenu()
	case 5:
		fmt.Println("Please note that if you have made config changes, you need to run `tok rebuild`")
		os.Exit(0)
	}
}

// DatabaseEngineMenu is exposes Tokaido-level config settings
func DatabaseEngineMenu() {

	menu := []ConfigGenericString{
		{
			Name:    "mysql",
			Type:    "value",
			Current: GetConfig().Database.Engine,
			Detail:  "Enable MySQL. This is the default provider for Tokaido.",
		},
		{
			Name:    "mariadb",
			Type:    "value",
			Current: GetConfig().Database.Engine,
			Detail:  "Enable MariaDB as an alternative to MySQL. Note that your database hostname will still be 'mysql'",
		},
		{
			Name:   "« Database Menu",
			Type:   "mainmenu",
			Detail: "Go back to the Database",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }} {{if eq .Current .Name }} [{{ "current setting" | green }}] {{ end }}`,
		Inactive: `   {{ .Name | cyan }} {{if eq .Current .Name }} [{{ "current setting" | green }}] {{ end }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "mainmenu" }}---------
{{ .Detail }}

{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » Database » Engine",
		Items:     menu,
		Templates: templates,
		Size:      4,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		SetConfigValueByArgs([]string{"database", "engine", "mysql"}, "project")
		reloadConfig()
		DatabaseMenu()
	case 1:
		SetConfigValueByArgs([]string{"database", "engine", "mariadb"}, "project")
		reloadConfig()
		DatabaseMenu()
	case 2:
		DatabaseMenu()
	}
}

// DatabaseMysqlconfigMenu enables the user to change the MySQL version
func DatabaseMysqlconfigMenu() {
	menu := []ConfigGenericString{
		{
			Name:    "MySQL 5.6",
			Current: GetConfig().Database.Mysqlconfig.Version,
			Type:    "value",
			Detail:  "Use the latest version of the MySQL 5.6 Docker image.\nSee https://hub.docker.com/_/mysql for more info.",
		},
		{
			Name:    "MySQL 5.7",
			Current: GetConfig().Database.Mysqlconfig.Version,
			Type:    "value",
			Detail:  "Use the latest version of the MySQL 5.7 Docker image.\nSee https://hub.docker.com/_/mysql for more info.\n\nThis is the Tokaido default",
		},
		{
			Name:    "MySQL 8.0",
			Current: GetConfig().Database.Mysqlconfig.Version,
			Type:    "value",
			Detail:  "Use the latest version of the MySQL 8.0 Docker image.\nSee https://hub.docker.com/_/mysql for more info.",
		},
		{
			Name:   "« Database Config",
			Type:   "menu",
			Detail: "Go back to the Database Config Menu",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }} {{if eq .Current .Name }} [{{ "current setting" | green }}] {{ end }}`,
		Inactive: `   {{ .Name | cyan }} {{if eq .Current .Name }} [{{ "current setting" | green }}] {{ end }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "menu" }}---------
{{ .Detail }}

{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » MySQL Configuration » Version",
		Items:     menu,
		Templates: templates,
		Size:      4,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		SetConfigValueByArgs([]string{"database", "mysqlconfig", "version", "5.6"}, "project")
		reloadConfig()
		DatabaseMenu()
	case 1:
		SetConfigValueByArgs([]string{"database", "mysqlconfig", "version", "5.7"}, "project")
		reloadConfig()
		DatabaseMenu()
	case 2:
		SetConfigValueByArgs([]string{"database", "mysqlconfig", "version", "8.0"}, "project")
		reloadConfig()
		DatabaseMenu()
	case 3:
		DatabaseMenu()
	}
}

// DatabaseMariadbconfigMenu enables the user to change the Mariadb version
func DatabaseMariadbconfigMenu() {
	menu := []ConfigGenericString{
		{
			Name:    "MariaDB 10.3",
			Type:    "value",
			Current: GetConfig().Database.Mysqlconfig.Version,
			Detail:  "Use the latest version of the MariaDB 10.3 Docker image.\nSee https://hub.docker.com/_/mysql for more info.\n\nThis is the Tokaido default",
		},
		{
			Name:    "MariaDB 10.4",
			Type:    "value",
			Current: GetConfig().Database.Mysqlconfig.Version,
			Detail:  "Use the latest version of the MariaDB 10.4 Docker image.\nSee https://hub.docker.com/_/mysql for more info.",
		},
		{
			Name:   "« Database Menu",
			Type:   "mainmenu",
			Detail: "Go back to the Main Menu",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }} {{if eq .Current .Name }} [{{ "current setting" | green }}] {{ end }}`,
		Inactive: `   {{ .Name | cyan }} {{if eq .Current .Name }} [{{ "current setting" | green }}] {{ end }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "menu" }}---------
{{ .Detail }}

{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » MariaDB Configuration » Version",
		Items:     menu,
		Templates: templates,
		Size:      3,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		SetConfigValueByArgs([]string{"database", "mariadbconfig", "version", "10.3"}, "project")
		reloadConfig()
		DatabaseMenu()
	case 1:
		SetConfigValueByArgs([]string{"database", "mariadbconfig", "version", "10.4"}, "project")
		reloadConfig()
		DatabaseMenu()
	case 2:
		DatabaseMenu()
	}
}

// DrupalMenu is exposes Tokaido-level config settings
func DrupalMenu() {
	menu := []ConfigGenericString{
		{
			Name:    "Drupal Root Path",
			Type:    "value",
			Current: GetConfig().Drupal.Path,
			Detail:  "This is the name of your Drupal document root, such as 'docroot' or 'web'.\nIf you change this value, you must run `tok destroy` and `tok rebuild`",
		},
		{
			Name:    "« Main Menu",
			Type:    "menu",
			Default: "",
			Current: "",
			Detail:  "Go back to the Main Menu",
		},
		{
			Name:    "Exit",
			Type:    "menu",
			Default: "",
			Current: "",
			Detail:  "Stop editing your configuration",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }} {{ if ne .Type "menu" }} Current setting: [{{ .Current | green }}] {{ end }}`,
		Inactive: `   {{ .Name | cyan }} {{ if ne .Type "menu" }} Current setting: [{{ .Current | green }}] {{ end }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "menu" }}---------
{{ .Detail }}

Current Setting: [{{ .Current | green }}]
{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » Tokaido Configuration",
		Items:     menu,
		Templates: templates,
		Size:      6,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		templates := &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		}

		prompt := promptui.Prompt{
			Label:     "Please enter the name of your Drupal root directory such as '/docroot' or '/web'.",
			Templates: templates,
		}

		dp, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		SetConfigValueByArgs([]string{"drupal", "path", dp}, "project")
		reloadConfig()
		DrupalMenu()
	case 1:
		MainMenu()
	case 2:
		fmt.Println("Please note that if you have made config changes, you need to run `tok rebuild`")
		os.Exit(0)
	}
}

// NginxMenu is exposes Nginx config settings that can be controlled from the codebase, rather than via the container env vars
func NginxMenu() {
	menu := []ConfigGenericString{
		{
			Name:    "Worker Connections",
			Type:    "value",
			Default: "1024",
			Current: GetConfig().Nginx.Workerconnections,
			Detail:  "Sets the maximum number of simultaneous connections that can be opened by a \nworker process.\n\nIt should be kept in mind that this number includes all connections \n(e.g. connections with proxied servers, among others), not only connections with\nclients. Another consideration is that the actual number of simultaneous\nconnections cannot exceed the current limit on the maximum number of\nopen files",
		},
		{
			Name:    "Client Body Max Size",
			Type:    "value",
			Default: "1024m",
			Current: GetConfig().Nginx.Clientmaxbodysize,
			Detail:  "Sets the maximum allowed size of the client request body, specified in the \n“Content-Length” request header field. If the size in a request exceeds the \nconfigured value, the 413 (Request Entity Too Large) error is returned to \nthe client.\n\nSet this value if you want to increase or decrease the max upload size in Drupal",
		},
		{
			Name:    "Keepalive Timeout",
			Type:    "value",
			Default: "65",
			Current: GetConfig().Nginx.Keepalivetimeout,
			Detail:  "Sets a timeout during which a keep-alive client connection will stay open on \nthe server side. The zero value disables keep-alive client connections.\n\nA lower or zero value may help mitigate DDoS attack traffic",
		},
		{
			Name:    "FastCGI Read Timeout",
			Type:    "value",
			Default: "300",
			Current: GetConfig().Nginx.Fastcgireadtimeout,
			Detail:  "This is how long Nginx will wait for the PHP FPM process to respond to a request\nbefore timing out. Note that increasing this value only impacts Nginx, and other\nservers in the request chain (such as a CDN) may impose their own timeout value\n\nAlso note that increasing this value and exposing long-running pages like public\nforms can easily result in a new DDoS vector that crashes your site. Wherever\npossible, you should set this value lower, not higher.",
		},
		{
			Name:    "FastCGI Buffers",
			Type:    "value",
			Default: "16 16k",
			Current: GetConfig().Nginx.Fastcgibuffers,
			Detail:  "Sets the number and size of the buffers used for reading a response from the\nPHP FPM server, for a single connection.\n\nSet this value as 'number size' such as '16 32k' for 16 x 32k buffers",
		},
		{
			Name:    "FastCGI Buffer Size",
			Type:    "value",
			Default: "32k",
			Current: GetConfig().Nginx.Fastcgibuffersize,
			Detail:  "When buffering of responses from the FastCGI server is enabled, limits the total\nsize of buffers that can be busy sending a response to the client while the\nresponse is not yet fully read. In the meantime, the rest of the buffers can be\nused for reading the response and, if needed, buffering part of the response to\na temporary file. By default, size is limited by the size of two buffers set by\nthe fastcgi_buffer_size and fastcgi_buffers directives.",
		},
		{
			Name:    "« Main Menu",
			Type:    "menu",
			Default: "",
			Current: "",
			Detail:  "Go back to the Main Menu",
		},
		{
			Name:    "Exit",
			Type:    "menu",
			Default: "",
			Current: "",
			Detail:  "Stop editing your configuration",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }} {{ if ne .Type "menu" }} {{if or (eq .Current .Default) (eq .Current "") }} Using default value [{{ .Default | cyan }}] {{ else }} Using custom value [{{ .Current | green }}] {{ end }} {{ end }}`,
		Inactive: `   {{ .Name | cyan }} {{ if ne .Type "menu" }} {{if or (eq .Current .Default) (eq .Current "") }} Using default value [{{ .Default | cyan }}] {{ else }} Using custom value [{{ .Current | green }}] {{ end }} {{ end }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "menu" }}---------
{{ .Detail }}

{{ if eq .Current "" }}Current Setting: Use Default Value {{ else }}Current Setting: [{{ .Current | green }}] {{ end }}
Default Setting: [{{ .Default | cyan }}]
{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » Tokaido Configuration",
		Items:     menu,
		Templates: templates,
		Size:      8,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		res := newStringValue("Specify the number of worker connections as an integer")
		SetConfigValueByArgs([]string{"nginx", "workerconnections", res}, "project")
		reloadConfig()
		NginxMenu()
	case 1:
		res := newStringValue("Specify the max body size in MB, with 'm' at the end (eg '64m' or '128m')")
		SetConfigValueByArgs([]string{"nginx", "clientmaxbodysize", res}, "project")
		reloadConfig()
		NginxMenu()
	case 2:
		res := newStringValue("Specify the keepalive timeout in seconds")
		SetConfigValueByArgs([]string{"nginx", "keepalivetimeout", res}, "project")
		reloadConfig()
		NginxMenu()
	case 3:
		res := newStringValue("Specify the fastcgi read timeout in seconds")
		SetConfigValueByArgs([]string{"nginx", "fastcgireadtimeout", res}, "project")
		reloadConfig()
		NginxMenu()
	case 4:
		res := newStringValue("Specify the fastcgi buffers config as 'number size' such as '16 16k'")
		SetConfigValueByArgs([]string{"nginx", "fastcgibuffers", res}, "project")
		reloadConfig()
		NginxMenu()
	case 5:
		res := newStringValue("Specify the fastcgi buffers size in kilobytes such as '16k' or '8k'")
		SetConfigValueByArgs([]string{"nginx", "fastcgibuffersize", res}, "project")
		reloadConfig()
		NginxMenu()
	case 6:
		MainMenu()
	case 7:
		fmt.Println("Please note that if you have made config changes, you need to run `tok rebuild`")
		os.Exit(0)
	}
}

// FpmMenu is exposes Nginx config settings that can be controlled from the codebase, rather than via the container env vars
func FpmMenu() {
	menu := []ConfigGenericString{
		{
			Name:    "Max Execution Time",
			Type:    "value",
			Default: "300",
			Current: GetConfig().Fpm.Maxexecutiontime,
			Detail:  "Sets the maximum time, in seconds, that a PHP request can run. Lower values are\nmuch better, as they improve performance and reduce DDoS attack risk\n\n\n\n",
		},
		{
			Name:    "Memory Limit",
			Type:    "value",
			Default: "256M",
			Current: GetConfig().Fpm.Phpmemorylimit,
			Detail:  "The maximum memory that a single PHP worker process can consume. Be sure to\nmatch this value to your hosting providers value to minimise the risk of \nsomething which works in Tokaido being busted in non-Tokaido production hosts",
		},
		{
			Name:    "Display Errors",
			Type:    "value",
			Default: "Off",
			Current: GetConfig().Fpm.Phpdisplayerrors,
			Detail:  "Display PHP errors in the browser",
		},
		{
			Name:    "Log Errors",
			Type:    "value",
			Default: "On",
			Current: GetConfig().Fpm.Phplogerrors,
			Detail:  "Log PHP errors to /tokaido/logs/nginx/errors.log",
		},
		{
			Name:    "Report Memory Leaks",
			Type:    "value",
			Default: "On",
			Current: GetConfig().Fpm.Phpreportmemleaks,
			Detail:  "When 'On', PHP will report memory leaks in the console and/or error log",
		},
		{
			Name:    "POST Max Size",
			Type:    "value",
			Default: "64M",
			Current: GetConfig().Fpm.Phppostmaxsize,
			Detail:  "The maximum allows upload size. Smaller is better, as large values\nincrease your DDoS attack risk",
		},
		{
			Name:    "Default Character Set",
			Type:    "value",
			Default: "UTF-8",
			Current: GetConfig().Fpm.Phpdefaultcharset,
			Detail:  "You almost always want this to be UTF-8",
		},
		{
			Name:    "Allow File Uploads",
			Type:    "value",
			Default: "On",
			Current: GetConfig().Fpm.Phpfileuploads,
			Detail:  "Turn this value 'Off' if you don't need file upload support, \nbut you almost certainly do",
		},
		{
			Name:    "Max Upload File Size",
			Type:    "value",
			Default: "64M",
			Current: GetConfig().Fpm.Phpuploadmaxfilesize,
			Detail:  "This is the maximum upload size and should match POST Max Size.\nSmaller is better, as large values increase your DDoS attack risk",
		},
		{
			Name:    "Max File Uploads",
			Type:    "value",
			Default: "20",
			Current: GetConfig().Fpm.Phpmaxfileuploads,
			Detail:  "Sets the maximum number of simultaneous file uploads",
		},
		{
			Name:    "Allow URL Fopen",
			Type:    "value",
			Default: "On",
			Current: GetConfig().Fpm.Phpallowurlfopen,
			Detail:  "This option enables the URL-aware fopen wrappers that enable accessing URL \nobject like files. Default wrappers are provided for the access of remote files\nusing the ftp or http protocol, some extensions like zlib may register\nadditional wrappers.",
		},
		{
			Name:    "« Main Menu",
			Type:    "menu",
			Default: "",
			Current: "",
			Detail:  "Go back to the Main Menu",
		},
		{
			Name:    "Exit",
			Type:    "menu",
			Default: "",
			Current: "",
			Detail:  "Stop editing your configuration",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }} {{ if ne .Type "menu" }} {{if or (eq .Current .Default) (eq .Current "") }} Using default value [{{ .Default | cyan }}] {{ else }} Using custom value [{{ .Current | green }}] {{ end }} {{ end }}`,
		Inactive: `   {{ .Name | cyan }} {{ if ne .Type "menu" }} {{if or (eq .Current .Default) (eq .Current "") }} Using default value [{{ .Default | cyan }}] {{ else }} Using custom value [{{ .Current | green }}] {{ end }} {{ end }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "menu" }}---------
{{ .Detail }}

{{ if eq .Current "" }}Current Setting: Use Default Value {{ else }}Current Setting: [{{ .Current | green }}] {{ end }}
Default Setting: [{{ .Default | cyan }}]
{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » Tokaido Configuration",
		Items:     menu,
		Templates: templates,
		Size:      13,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		res := newStringValue("Specify the max execution time in seconds")
		SetConfigValueByArgs([]string{"fpm", "maxexecutiontime", res}, "project")
		reloadConfig()
		FpmMenu()
	case 1:
		res := newStringValue("Specify memory limit in MB, with 'M' at the end (eg '1024M' or '512M')")
		SetConfigValueByArgs([]string{"fpm", "phpmemorylimit", res}, "project")
		reloadConfig()
		FpmMenu()
	case 2:
		if GetConfig().Fpm.Phpdisplayerrors == "On" {
			SetConfigValueByArgs([]string{"fpm", "phpdisplayerrors", "Off"}, "project")
		} else {
			SetConfigValueByArgs([]string{"fpm", "phpdisplayerrors", "On"}, "project")
		}
		reloadConfig()
		FpmMenu()
	case 3:
		if GetConfig().Fpm.Phplogerrors == "On" {
			SetConfigValueByArgs([]string{"fpm", "phplogerrors", "Off"}, "project")
		} else {
			SetConfigValueByArgs([]string{"fpm", "phplogerrors", "On"}, "project")
		}
		reloadConfig()
		FpmMenu()
	case 4:
		if GetConfig().Fpm.Phpreportmemleaks == "On" {
			SetConfigValueByArgs([]string{"fpm", "phpreportmemleaks", "Off"}, "project")
		} else {
			SetConfigValueByArgs([]string{"fpm", "phpreportmemleaks", "On"}, "project")
		}
		reloadConfig()
		FpmMenu()
	case 5:
		res := newStringValue("Specify the max POST size in MB, with 'M' at the end (eg '64M' or '10M')")
		SetConfigValueByArgs([]string{"fpm", "phppostmaxsize", res}, "project")
		reloadConfig()
		FpmMenu()
	case 6:
		res := newStringValue("Specify the default character set")
		SetConfigValueByArgs([]string{"fpm", "phpdefaultcharset", res}, "project")
		reloadConfig()
		FpmMenu()
	case 7:
		if GetConfig().Fpm.Phpfileuploads == "On" {
			SetConfigValueByArgs([]string{"fpm", "phpfileuploads", "Off"}, "project")
		} else {
			SetConfigValueByArgs([]string{"fpm", "phpfileuploads", "On"}, "project")
		}
		reloadConfig()
		FpmMenu()
	case 8:
		res := newStringValue("Specify the max POST size in MB, with 'M' at the end (eg '64M' or '10M')")
		SetConfigValueByArgs([]string{"fpm", "phpuploadmaxfilesize", res}, "project")
		reloadConfig()
		FpmMenu()
	case 9:
		res := newStringValue("Specify the maximum number of simultaneous file uploads")
		SetConfigValueByArgs([]string{"fpm", "phpmaxfileuploads", res}, "project")
		reloadConfig()
		FpmMenu()
	case 10:
		if GetConfig().Fpm.Phpallowurlfopen == "On" {
			SetConfigValueByArgs([]string{"fpm", "phpallowurlfopen", "Off"}, "project")
		} else {
			SetConfigValueByArgs([]string{"fpm", "phpallowurlfopen", "On"}, "project")
		}
		reloadConfig()
		FpmMenu()
	case 11:
		MainMenu()
	case 12:
		fmt.Println("Please note that if you have made config changes, you need to run `tok rebuild`")
		os.Exit(0)
	}
}

// ServicesMenu ...
func ServicesMenu() {
	menu := []ConfigGenericString{
		{
			Name:    "Solr Search 6.6",
			Type:    "value",
			Default: "false",
			Current: strconv.FormatBool(GetConfig().Services.Solr.Enabled),
			Detail:  "If true, Solr 6.6 will be provisioned at solr:8983.\nRun `tok open solr` to access it",
		},
		{
			Name:    "Memcache",
			Default: "true",
			Type:    "value",
			Current: strconv.FormatBool(GetConfig().Services.Memcache.Enabled),
			Detail:  "If true, Memcache will be provisioned at memcached:11211",
		},
		{
			Name:    "Redis",
			Type:    "value",
			Default: "false",
			Current: strconv.FormatBool(GetConfig().Services.Redis.Enabled),
			Detail:  "If true, Redis will be available at redis:6379",
		},
		{
			Name:    "Mailhog",
			Type:    "value",
			Default: "false",
			Current: strconv.FormatBool(GetConfig().Services.Mailhog.Enabled),
			Detail:  "Turn on Mailhog to easily test SMTP settings.\nRun `tok open mailhog` to access it",
		},
		{
			Name:    "PHP Adminer",
			Type:    "value",
			Default: "false",
			Current: strconv.FormatBool(GetConfig().Services.Adminer.Enabled),
			Detail:  "Turn on Adminer for a MySQL GUI tool, similar to phpmyadmin ",
		},
		{
			Name:    "« Main Menu",
			Type:    "menu",
			Default: "",
			Current: "",
			Detail:  "Go back to the Main Menu",
		},
		{
			Name:    "Exit",
			Type:    "menu",
			Default: "",
			Current: "",
			Detail:  "Stop editing your configuration",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   `🤔 {{ .Name | cyan }} {{ if ne .Type "menu" }} {{ if eq .Current .Default }} Using default value [{{ .Default | cyan }}] {{ else }} Using custom value [{{ .Current | green }}] {{ end }} {{ end }}`,
		Inactive: `   {{ .Name | cyan }} {{ if ne .Type "menu" }} {{ if eq .Current .Default }} Using default value [{{ .Default | cyan }}] {{ else }} Using custom value [{{ .Current | green }}] {{ end }} {{ end }}`,
		Selected: "{{ .Name | blue }}",
		Details: `
{{ if ne .Type "menu" }}---------
{{ .Detail }}

Default Setting: [{{ .Default | cyan }}]
Current Setting: [{{ .Current | green }}]
{{ end }}
`,
	}

	prompt := promptui.Select{
		Label:     "Main Menu » Tokaido Configuration",
		Items:     menu,
		Templates: templates,
		Size:      7,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch i {
	case 0:
		if GetConfig().Services.Solr.Enabled == true {
			SetConfigValueByArgs([]string{"services", "solr", "enabled", "false"}, "project")
		} else {
			SetConfigValueByArgs([]string{"services", "solr", "enabled", "true"}, "project")
		}
		reloadConfig()
		ServicesMenu()
	case 1:
		if GetConfig().Services.Memcache.Enabled == true {
			SetConfigValueByArgs([]string{"services", "memcache", "enabled", "false"}, "project")
		} else {
			SetConfigValueByArgs([]string{"services", "memcache", "enabled", "true"}, "project")
		}
		reloadConfig()
		ServicesMenu()
	case 2:
		if GetConfig().Services.Redis.Enabled == true {
			SetConfigValueByArgs([]string{"services", "redis", "enabled", "false"}, "project")
		} else {
			SetConfigValueByArgs([]string{"services", "redis", "enabled", "true"}, "project")
		}
		reloadConfig()
		ServicesMenu()
	case 3:
		if GetConfig().Services.Mailhog.Enabled == true {
			SetConfigValueByArgs([]string{"services", "mailhog", "enabled", "false"}, "project")
		} else {
			SetConfigValueByArgs([]string{"services", "mailhog", "enabled", "true"}, "project")
		}
		reloadConfig()
		ServicesMenu()
	case 4:
		if GetConfig().Services.Adminer.Enabled == true {
			SetConfigValueByArgs([]string{"services", "adminer", "enabled", "false"}, "project")
		} else {
			SetConfigValueByArgs([]string{"services", "adminer", "enabled", "true"}, "project")
		}
		reloadConfig()
		ServicesMenu()
	case 5:
		MainMenu()
	case 6:
		fmt.Println("Please note that if you have made config changes, you need to run `tok rebuild`")
		os.Exit(0)
	}
}
