package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"localhost:8080"`
}

func MustLoad() *Config {
	config := os.Getenv("CONFIG_PATH")
	if config == "" {
		log.Fatal("CONFIG PATH is not set")
	}

	if _, err := os.Stat(config); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", config)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(config, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
