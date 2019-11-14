package conf

import (
	"strconv"

	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/ironstar-io/tokaido/system/hash"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/logrusorgru/aurora"
)

// SetConfigValueByArgs updates the config file by merging a **single** new value with the
// current in memory configuration. Once merged, it writes the updated config to disk
// - args are a slice of new values such as `[]string{"nginx", "workerconnections", "30"}`
// - configType is either 'project' or 'global' and will determine which file is updated
func SetConfigValueByArgs(args []string, configType string) {
	if configType != "project" && configType != "global" {
		fmt.Println(aurora.Sprintf("The config file %s is unknown", aurora.Bold(configType)))
		os.Exit(1)
	}

	validateArgs(args)

	newYaml := argsToYaml(args) // carries yaml-formatted string of singular args slice
	configPath := getConfigPath(configType)

	// 'runningConfig' initially carries in-memory config from Viper, which does not differentiate
	// between our "project" and "global" config files
	// later on we merge our yaml config into this in-memory config
	runningConfig := GetConfig()

	// merge our newYaml into our runningConfig
	newConfig := mergeConfigInMemory(newYaml, configPath, runningConfig)

	// Viper doesn't split config in memory so 'runningConfig' now contains merged
	// project and global config settings. We need to split them out.
	if configType == "project" {
		// These values must not be written to the project config file so we reset them to nil or empty
		runningConfig.Global.Syncservice = ""
		runningConfig.Global.Proxy.Enabled = false
		runningConfig.Global.Projects = nil
		runningConfig.Global.Telemetry = Telemetry{}

		// Stop our debug flag from leaking into project config
		emptyDebug := new(bool)
		runningConfig.Tokaido.Debug = *emptyDebug
	}

	writeConfig(runningConfig, configPath)

	// validate that our config was saved successfully
	compareFiles(newConfig, configPath)
}

// SetGlobalConfigValueByArgs enables the update of global config values by a string slice
// When working project-level global config, it defaults to updating the active project
func SetGlobalConfigValueByArgs(args []string) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("Error: too few arguments for global config")
	}

	if args[1] == "syncservice" {
		if len(args) < 3 {
			return fmt.Errorf("Error: too few arguments for global project config. Did you want 'project'?")
		}

		g := GetGlobalConfig()
		if err != nil {
			return err
		}

		g.Syncservice = args[2]
		WriteGlobalConfig(*g)
		return nil
	}

	if args[1] == "project" {
		if len(args) < 3 {
			return fmt.Errorf("Error: too few arguments for global project config. Did you want 'xdebug'?")
		}

		p, err := GetGlobalProjectSettings()
		if err != nil {
			return err
		}

		if args[2] == "database" {
			if len(args) < 4 {
				return fmt.Errorf("Error: too few arguments for database port. Please specify 'port {number}'")
			}

			if args[3] != "port" {
				return fmt.Errorf("Error: unknown argument '%s'. Expected 'port'", args[3])
			}

			dbPort, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}
			if dbPort < 1024 || dbPort > 65535 {
				return fmt.Errorf("Error: you must specify an static database port between 1025 and 65535")
			}

			p.Database.Port = dbPort
			WriteGlobalProjectSettings(p)
			return nil
		}

		if args[2] == "xdebug" {
			if len(args) < 4 {
				return fmt.Errorf("Error: too few arguments for xdebug config. Please specify 'enabled {true/false}' or 'port {number}'")
			}

			if args[3] == "enabled" {
				if len(args) < 5 {
					return fmt.Errorf("Error: too few arguments for xdebug setting. Please specify 'enabled {true/false}'")
				}
				enabled, err := strconv.ParseBool(args[4])
				if err != nil {
					return err
				}

				p.Xdebug.Enabled = enabled
				WriteGlobalProjectSettings(p)
				return nil
			}

			if args[3] == "fpmport" {
				if len(args) < 5 {
					return fmt.Errorf("Error: too few arguments for xdebug port. Please specify 'fpmport {number}'")
				}
				fpmPort, err := strconv.Atoi(args[4])
				if err != nil {
					return err
				}
				if fpmPort < 1024 || fpmPort > 65535 {
					return fmt.Errorf("Error: you must specify an xdebug fpmport between 1025 and 65535")
				}

				p.Xdebug.FpmPort = fpmPort
				WriteGlobalProjectSettings(p)
				return nil
			}

			return nil
		}
	}

	return fmt.Errorf("requested global config path is not known")
}

// mergeConfigInMemory takes our saved config from disk, the new yaml string, and the running
// config and merges all three into a new byte slice that can be saved to disk
func mergeConfigInMemory(newYaml, configPath string, runningConfig *Config) (newConfig []byte) {
	// Read the saved config file from disk
	newConfig, err := ioutil.ReadFile(configPath) // newConfig will eventually be written to disk
	if err != nil {
		log.Fatalf("There was an issue reading in your config file\n%v", err)
	}

	// Unmarshal the new config from disk into our running config
	err = yaml.Unmarshal(newConfig, &runningConfig)
	if err != nil {
		log.Fatalf("There was an issue parsing your config file\n%v", err)
	}

	// Merge the new yaml with our in-memory config
	err = yaml.Unmarshal([]byte(newYaml), &runningConfig)
	if err != nil {
		log.Fatalf("There was an issue updating your config file\n%v", err)
	}

	return newConfig
}

// writeConfig writes the provided in-memory config to the configPath
func writeConfig(runningConfig *Config, configPath string) {
	// Now that we've merged the config, we'll write that merged config to disk
	newMarhsalled, err := yaml.Marshal(runningConfig)
	if err != nil {
		log.Fatalf("There was an issue building your config file\n%v", err)
	}

	fs.Replace(configPath, newMarhsalled)
}

// RegisterProject adds a project to the global config file
func RegisterProject(name, path string) {
	gcPath := getConfigPath("global")

	// Read the global config from file
	// Using the in-memory config from Viper isn't an option here because it would
	// contain _all_ of the project-level config, and rubbing those out is too
	// verbose and difficult to scale. Thankfully nothing modifies global config
	// at run time so this mechanism is safe.
	gcFile, err := ioutil.ReadFile(gcPath)
	if err != nil {
		log.Fatalf("There was an issue reading in your global config file\n%v", err)
	}

	gc := &Global{}
	err = yaml.Unmarshal(gcFile, gc)
	if err != nil {
		log.Fatalf("There was an issue unpacking your global config file\n%v", err)
	}

	// Add this project to the global list if it isn't already there
	found := false
	for _, v := range gc.Projects {
		if v.Name == name {
			found = true
		}
	}

	if !found {
		project := Project{
			Name: name,
			Path: path,
		}
		gc.Projects = append(gc.Projects, project)
	}

	// Write the updated global config back to file
	newMarhsalled, err := yaml.Marshal(gc)
	if err != nil {
		log.Fatalf("There was a fatal issue updating your global config file\n%v", err)
	}

	fs.Replace(gcPath, newMarhsalled)
}

// DeregisterProject removes a project from the global config file
func DeregisterProject(name string) {
	gcPath := getConfigPath("global")

	// Read the global config from file
	// Using the in-memory config from Viper isn't an option here because it would
	// contain _all_ of the project-level config, and rubbing those out is too
	// verbose and difficult to scale. Thankfully nothing modifies global config
	// at run time so this mechanism is safe.
	gcFile, err := ioutil.ReadFile(gcPath)
	if err != nil {
		log.Fatalf("There was an issue reading in your global config file\n%v", err)
	}

	gc := &Global{}
	err = yaml.Unmarshal(gcFile, gc)
	if err != nil {
		log.Fatalf("There was an issue unpacking your global config file\n%v", err)
	}

	// Remove this project to the global list if it isn't already there
	var index int
	for i, v := range gc.Projects {
		if v.Name == name {
			index = i
		}
	}
	if index > 0 {
		gc.Projects = append(gc.Projects[:index], gc.Projects[index+1:]...)
	}

	// Write the updated global config back to file
	WriteGlobalConfig(*gc)

}

// WriteGlobalConfig overwrites the existing global.yml with the supplied config
func WriteGlobalConfig(ng Global) {
	gcPath := getConfigPath("global")

	// Marhsal the supplied global config into yaml
	newMarhsalled, err := yaml.Marshal(ng)
	if err != nil {
		log.Fatalf("There was a fatal issue updating your global config file\n%v", err)
	}

	fs.Replace(gcPath, newMarhsalled)
}

// WriteGlobalProjectSettings takes an updated Project object and merges it with all of the
// Projects defined in the global config, finally saving the whole global config back to disk
func WriteGlobalProjectSettings(p *Project) {
	g := GetConfig().Global

	for k, v := range g.Projects {
		if v.Name == p.Name {
			g.Projects[k] = *p
		}
	}

	WriteGlobalConfig(g)
}

// compareFiles checks original and new config files to identify if any values were changed
func compareFiles(original []byte, newPath string) {
	o, err := hash.BytesMD5(original)
	if err != nil {
		log.Fatalf("There was an issue opening the new config file\n%v", err)
	}

	n, err := hash.FileMD5(newPath)
	if err != nil {
		log.Fatalf("There was an issue opening the new config file:\n%v", err)
	}

	if o == n {
		fmt.Println("Action resulted in no change to config")
		return
	}

}

func unmarshalConfig(cp string) *Config {
	c := &Config{}

	yf, err := ioutil.ReadFile(cp)
	if err != nil {
		log.Fatalf("There was an issue reading in your config file\n%v", err)
	}

	err = yaml.Unmarshal(yf, c)
	if err != nil {
		log.Fatalf("There was an issue parsing your config file\n%v", err)
	}

	return c
}

func getConfigPath(configFile string) string {
	var cp string
	if configFile == "project" {
		cp = filepath.Join(GetProjectPath(), ".tok", "config.yml")
	} else if configFile == "global" {
		cp = filepath.Join(fs.HomeDir(), ".tok", "global.yml")
	} else {
		fmt.Println(aurora.Sprintf("The config file %s is unknown", aurora.Bold(configFile)))
	}

	// Initialise the config file if it doesn't exist
	var _, errf = os.Stat(cp)
	if os.IsNotExist(errf) {
		// The global .tok path requires appropriate permissions
		gp := filepath.Dir(cp)
		if configFile == "global" && !fs.CheckExists(gp) {
			err := os.MkdirAll(gp, 0700)
			if err != nil {
				log.Fatalf("Unexpected error creating global config directory")
			}
		}
		fs.TouchEmpty(cp)
	}

	return cp
}

// argsToYaml converts a string slice to a single yaml formatted-string
func argsToYaml(args []string) string {
	var y string
	for i, a := range args {
		if i == len(args)-1 {
			y = y + " " + a
			continue
		}
		y = y + calcWhitespace(i) + mapEdgeKeys(a) + ":"
	}

	return y
}

func calcWhitespace(i int) string {
	if i == 0 {
		return ""
	}

	w := "\n"
	for x := 1; x <= i; x++ {
		w = w + "  "
	}

	return w
}

func mapEdgeKeys(a string) string {
	var keyMap = map[string]string{
		"volumesfrom": "volumes_from",
		"dependson":   "depends_on",
		"workingdir":  "working_dir",
	}

	if keyMap[a] != "" {
		return keyMap[a]
	}

	return a
}

func validateArgs(args []string) {
	if len(args) < 2 {
		log.Fatal("At least two arguments are required in order to set a config value")
	}

	ca := args[:len(args)-1]

	_, err := GetConfigValueByArgs(ca)
	if err != nil {
		log.Fatal(err)
	}
}
