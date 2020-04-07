package conf

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/ironstar-io/tokaido/system/fs"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// GetConfig ...
func GetConfig() *Config {
	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal("Failed to retrieve configuration values\n", err)
	}

	// Create an empty global config file if on doesn't exist
	gcPath := getConfigPath("global")
	if !fs.CheckExists(gcPath) {
		WriteGlobalConfig(Global{})
	}

	// Load the global config in to the config struct without using Viper
	gcFile, err := ioutil.ReadFile(gcPath)
	if err != nil {
		log.Fatalf("There was an issue reading in your global config file\n%v", err)
	}

	err = yaml.Unmarshal(gcFile, &config.Global)
	if err != nil {
		log.Fatalf("There was an issue unpacking your global config file\n%v", err)
	}

	return config
}

// GetGlobalConfig ...
func GetGlobalConfig() *Global {
	global := Global{}
	// Create an empty global config file if on doesn't exist
	gcPath := getConfigPath("global")
	if !fs.CheckExists(gcPath) {
		WriteGlobalConfig(global)
	}

	// Load the global config in to the config struct without using Viper
	gcFile, err := ioutil.ReadFile(gcPath)
	if err != nil {
		log.Fatalf("There was an issue reading in your global config file\n%v", err)
	}
	err = yaml.Unmarshal(gcFile, &global)
	if err != nil {
		log.Fatalf("There was an issue unpacking your global config file\n%v", err)
	}

	return &global
}

// GetProjectPath returns the full system path to this project as it exists in the global.yml file
func GetProjectPath() (path string) {
	// If the path exists in the global config, return it
	for _, v := range GetConfig().Global.Projects {
		if v.Name == GetConfig().Tokaido.Project.Name {
			return v.Path
		}
	}

	// The path isn't in global.yml yet, so we'll just return it based on the current project context
	return fs.ProjectRoot()
}

// GetGlobalProjectSettings returns the current global conf object for the current project
func GetGlobalProjectSettings() (*Project, error) {
	c := GetConfig()
	for _, v := range c.Global.Projects {
		if v.Name == c.Tokaido.Project.Name {
			return &v, nil
		}
	}

	return &Project{}, fmt.Errorf("unable to find global project configuration")
}

// GetConfigValueByArgs - Get a config value based on the arguments sent from the command line
func GetConfigValueByArgs(args []string) (reflect.Value, error) {
	c := GetConfig()

	if len(args) == 0 {
		return reflect.ValueOf(nil), errors.New("No arguments were provided. See `tok config-get --help` for usage")
	}

	r, err := getField(c, normalizeFieldSlice(args))
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	if !r.IsValid() {
		return reflect.ValueOf(nil), errors.New("`" + strings.ToLower(strings.Join(args, " ")) + "` is not a valid Tokaido configuration path")
	}

	return r, nil
}

func normalizeFieldSlice(args []string) []string {
	var s []string
	for _, a := range args {
		f := strings.ToUpper(string([]rune(a)[0]))
		s = append(s, f+strings.ToLower(a[1:]))
	}

	return s
}

func getField(v *Config, fields []string) (reflect.Value, error) {
	r := reflect.ValueOf(v)
	iv := reflect.Indirect(r)

	f := iv.FieldByName(fields[0])
	if len(fields) == 1 {
		return f, nil
	}

	for i, a := range fields {
		if i == 0 {
			continue
		}
		if !f.IsValid() {
			return reflect.ValueOf(nil), errors.New("`" + strings.ToLower(strings.Join(fields, " ")) + "` is not a valid Tokaido configuration path")
		}

		switch f.Kind() {
		case reflect.String:
			fallthrough
		case reflect.Bool:
			fallthrough
		case reflect.Int:
			return reflect.ValueOf(nil), errors.New("`" + strings.ToLower(strings.Join(fields, " ")) + "` is a value and cannot have a value set against it as a key")
		case reflect.Map:
			return f, nil
		default:
			f = f.FieldByName(a)
		}
	}

	return f, nil
}
