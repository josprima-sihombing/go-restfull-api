package config

import (
	"log"
	"os"
)

type EnvConfig struct {
	DatabaseURL string
	JWTSecret   string
}

var Env EnvConfig

func LoadEnv() {
	Env = EnvConfig{
		DatabaseURL: getConfig("DATABASE_URL"),
		JWTSecret:   getConfig("JWT_SECRET"),
	}
}

func getConfig(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("Environtment variable %s is not set", key)
	}

	return value
}
