package conf

import (
	"github.com/spf13/viper"
	"log"
)

// GetConfig ...
func GetConfig() *Config {
	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal("Failed to retrieve configuration values")
	}

	return config
}
