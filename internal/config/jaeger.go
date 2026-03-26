package config

import "os"

type JaegerConfig struct {
	Port           string
	Host           string
	TracesEndpoint string
}

func LoadJaegerConfig() JaegerConfig {
	return JaegerConfig{
		Port:           os.Getenv("JAEGER_PORT"),
		Host:           os.Getenv("JAEGER_HOST"),
		TracesEndpoint: os.Getenv("JAEGER_TRACES_ENDPOINT"),
	}
}
