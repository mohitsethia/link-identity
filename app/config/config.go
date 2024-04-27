package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
}

// Values ...
var Values config

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal(err, "Error loading .env file")
	}
	if Values.Server.Port = os.Getenv("server.port"); Values.Server.Port == "" {
		panic("server port cannot be empty")
	}
	if Values.Database.Host = os.Getenv("database.host"); Values.Database.Host == "" {
		panic("database host cannot be empty")
	}
	if Values.Database.Name = os.Getenv("database.name"); Values.Database.Name == "" {
		panic("database name cannot be empty")
	}
	if Values.Database.Password = os.Getenv("database.pass"); Values.Database.Password == "" {
		panic("database password cannot be empty")
	}
	if Values.Database.Username = os.Getenv("database.username"); Values.Database.Username == "" {
		panic("database username cannot be empty")
	}
	if Values.Database.Port = os.Getenv("database.port"); Values.Database.Port == "" {
		panic("database port cannot be empty")
	}
}
