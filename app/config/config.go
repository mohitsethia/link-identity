package config

import (
	"github.com/spf13/viper"
)

type config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
}

var Values config

var DefaultServerPort = "8000"

func init() {
	// TODO : Add once
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Values); err != nil {
		panic(err)
	}
}
