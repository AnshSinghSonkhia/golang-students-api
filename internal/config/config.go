package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// HTTPServer holds the configuration for the HTTP server.
type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

// Config holds the application configuration.
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

// MustLoad reads the configuration from a file specified by the CONFIG_PATH environment variable or command line flag.
// It returns a pointer to the Config struct or exits the program with an error message if the configuration cannot be loaded.
func MustLoad() *Config {
	var configPath string // Initialize configPath variable

	configPath = os.Getenv("CONFIG_PATH") // Check if the environment variable CONFIG_PATH is set

	if configPath == "" { // If CONFIG_PATH is not set, check for command line flags
		flags := flag.String("config", "", "Path to the configuration file")

		flag.Parse() // Parse command line flags

		configPath = *flags // Assign the value of the --config flag to configPath

		if configPath == "" {
			log.Fatal("Config path is not set. Please set the CONFIG_PATH environment variable or use the --config flag.")
		}
	}

	// Check if the configuration file exists
	// If it does not exist, log an error and exit the program
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file does not exist: %s", configPath)
	}

	var cfg Config // Initialize a Config struct to hold the configuration

	err := cleanenv.ReadConfig(configPath, &cfg) // Read the configuration from the specified file into the cfg struct
	// If there is an error reading the configuration file, log the error and exit the program

	if err != nil {
		log.Fatalf("Failed to read configuration file: %s", err.Error())
	}

	return &cfg // Return a pointer to the loaded configuration struct
}
