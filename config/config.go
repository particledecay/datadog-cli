package config

import (
	"github.com/spf13/viper"
)

var APIKey string

func Load() {
	viper.SetEnvPrefix("dd")
	viper.AutomaticEnv()

	APIKey = viper.GetString("api_key")
}
