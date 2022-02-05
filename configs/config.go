package configs

import (
	"gopkg.in/yaml.v2"
	"os"
)

const configPath string = "./configs/config.yml"

type Config struct {
	Port        string `yaml:"port"`
	DatabaseURL string `yaml:"database_url"`
}

func NewConfig() (*Config, error) {
	config := &Config{}

	// Opening config file ...
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	d := yaml.NewDecoder(f)

	// Parsing configs
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
