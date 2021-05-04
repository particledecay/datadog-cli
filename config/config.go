package config

import (
	"github.com/spf13/viper"
)

var APIKey string
var AppKey string

func Load() {
	viper.SetEnvPrefix("dd")
	viper.AutomaticEnv()

	APIKey = viper.GetString("api_key")
	AppKey = viper.GetString("app_key")
}
