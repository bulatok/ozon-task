package logic

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Port        string `yaml:"port"`
	DatabaseURL string `yaml:"database_url"`
}

func NewConfig(ConfigPath string) (*Config, error) {
	config := &Config{}

	// Opening config file ...
	f, err := os.Open(ConfigPath)
	if err != nil {
		return nil, err
	}

	// Initing file ...
	d := yaml.NewDecoder(f)

	// Parsing configs
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
