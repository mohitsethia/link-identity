package config

import (
	"log"
	"os"
	"path/filepath"

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
	// Define a slice of possible locations for the .env file
	possiblePaths := []string{
		"./.env",     // Current directory
		"../.env",    // Parent directory
		"$HOME/.env", // Home directory
		".env",
		"../../.env",
		filepath.Join(
			os.Getenv("GOPATH"),
			"src", "github.com", "mohitsethia", "link-identity", ".env"), // Go project directory
	}

	// Try each path and load the first .env file found
	var err error
	for _, path := range possiblePaths {
		if err = godotenv.Load(filepath.Clean(os.ExpandEnv(path))); err == nil {
			break
		}
	}
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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
