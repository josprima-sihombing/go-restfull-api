package config

import (
	"log"
	"os"
)

type EnvConfig struct {
	DatabaseURL string
}

var Env EnvConfig

func LoadEnv() {
	Env = EnvConfig{
		DatabaseURL: getConfig("DATABASE_URL"),
	}
}

func getConfig(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("Environtment variable %s is not set", key)
	}

	return value
}
