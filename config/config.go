package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// DebugMode is accessible from outside the package to make it easier to check.
var DebugMode bool

// There are some things we run at startup that relate to the runtime config, such as a debug mode
func init() {
	DebugMode = os.Getenv("DEBUG") == "true"
}

// Config represents what we pull in from the 'xero:' key in config.yml.
type Config struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	AppPort      int    `yaml:"app_port"`
}

// New - returns an instance of Config from the given filePath - If empty, will default to 'config.yml'.
func New(filePath string) *Config {
	if filePath == "" {
		filePath = "config.yml"
	}
	// Check that the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Println("Config file does not exist:", filePath)
		log.Fatalln(err)
	}
	configFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("An error occurred while trying to read the config file:", filePath)
		log.Fatalln(err)
	}

	var config Config

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Println("Unable to parse contents of YAML config file:", filePath)
		log.Fatalln(err)
	}
	if DebugMode {
		log.Println("Loaded config:")
		config.Print()
	}

	// Set some defaults if they weren't provided
	if config.AppPort == 0 {
		config.AppPort = 8000
		if DebugMode {
			log.Println("Using default app_port:", config.AppPort)
		}
	}

	// Validate that the config has the required values
	validateConfig(&config)

	// Attempt to parse the YAML file for the details we need.
	return &config
}

func validateConfig(c *Config) {
	var fieldsMissing []string
	if c.ClientID == "" {
		fieldsMissing = append(fieldsMissing, "client_id")
	}
	if c.ClientSecret == "" {
		fieldsMissing = append(fieldsMissing, "client_secret")
	}
	if len(fieldsMissing) > 0 {
		log.Println("The follow fields appear to be missing from your config file:")
		for _, configKey := range fieldsMissing {
			log.Println("-", configKey)
		}
		log.Fatalln("Please ensure all required config values are present.")
	}
}

// Print outputs the loaded client struct.
func (c *Config) Print() {
	log.Println("Client ID:    ", c.ClientID)
	log.Println("Client Secret:", c.ClientSecret)
}
