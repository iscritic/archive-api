package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer HTTPServer `yaml:"http_server"`
	SMTP       SMTPConfig `yaml:"smtp"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"localhost"`
	Port        int           `yaml:"port" env:"HTTP_SERVER_PORT" env-default:"8082"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
}

type SMTPConfig struct {
	Host     string `yaml:"host" env:"SMTP_HOST" env-default:"smtp.example.com"`
	Port     int    `yaml:"port" env:"SMTP_PORT" env-default:"587"`
	User     string `yaml:"user" env:"SMTP_USER" env-default:"your_email@example.com"`
	Password string `yaml:"password" env:"SMTP_PASSWORD"`
	Sender   string `yaml:"sender" env:"SMTP_SENDER" env-default:"your_email@example.com"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Error reading environment variables: %s", err)
	}

	return &cfg
}
