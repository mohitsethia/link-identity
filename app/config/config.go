package config

import (
	"github.com/spf13/viper"
)

type config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
}

// Values ...
var Values config

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
	if Values.Server.Port == "" {
		panic("server port cannot be empty")
	}
	if Values.Database.Host == "" {
		panic("database host cannot be empty")
	}
	if Values.Database.Name == "" {
		panic("database name cannot be empty")
	}
	if Values.Database.Password == "" {
		panic("database password cannot be empty")
	}
	if Values.Database.Username == "" {
		panic("database username cannot be empty")
	}
	if Values.Database.Port == "" {
		panic("database port cannot be empty")
	}
}
