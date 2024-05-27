package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DatabasePath    string `yaml:"database_path" env-required:"true"`
	KinopoiskClient `yaml:"kinopoisk_client"`
}

type KinopoiskClient struct {
	Timeout  time.Duration `yaml:"timeout" env-default:"10s"`
	Uri      string        `yaml:"uri" env-required:"true"`
	ApiToken string        `yaml:"api_token" env-required:"true"`
}

func MustLoad() *Config {
	configPath := "./config/config.yaml"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
