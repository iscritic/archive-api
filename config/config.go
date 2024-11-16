package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerAddress string `yaml:"server_address"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
