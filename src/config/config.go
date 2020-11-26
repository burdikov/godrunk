package config

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strings"
)

// Config represents application configuration
type Config struct {
	Token          *string `yaml:"token" envconfig:"TOKEN"`
	Port           int     `yaml:"port" envconfig:"PORT"`
	WebhookAddress *string `yaml:"webhookAddress" envconfig:"WEBHOOK_ADDRESS"`
	Debug          bool    `yaml:"debug" envconfig:"DEBUG"`
}

// GetConfig retrieves application configuration
func GetConfig(filename string) *Config {
	config := &Config{}

	if _, err := os.Stat(filename); err == nil {
		f, err := os.Open(filename)
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()

		err = yaml.NewDecoder(f).Decode(config)
		if err != nil {
			log.Panic(err)
		}
	} else {
		log.Printf("Error reading %v: %v", filename, err)
		log.Print("Proceeding without configuration file.")
	}

	err := envconfig.Process("GODRUNK", config)
	if err != nil {
		log.Panic(err)
	}

	CheckConfig(config)

	tmp := strings.TrimSpace(*config.WebhookAddress)
	config.WebhookAddress = &tmp

	PrintConfig(config)

	return config
}

func PrintConfig(config *Config) {
	log.Printf("Port: %v", config.Port)
	log.Printf("Webhook address: %v", *config.WebhookAddress)
}

func CheckConfig(config *Config) {
	if config.Token == nil {
		log.Fatal("Bot token is not provided. Can't proceed.")
	}

	if config.WebhookAddress == nil {
		log.Fatal("Webhook address is not provided. Can't proceed.")
	}
}
