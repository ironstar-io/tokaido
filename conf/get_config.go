package conf

import (
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// GetConfig ...
func GetConfig() *Config {
	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal("Failed to retrieve configuration values\n", err)
	}

	return config
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
