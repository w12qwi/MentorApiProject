package config

import "os"

type AppConfig struct {
	Host string
	Port string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Host: os.Getenv("APP_HOST"),
		Port: os.Getenv("APP_PORT"),
	}
}
