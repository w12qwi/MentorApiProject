package config

import (
	"log"
	"os"
	"strconv"
)

type GRPCconfig struct {
	Host    string
	Port    string
	Timeout int
}

func LoadGRPCConfig() GRPCconfig {
	timeout, err := strconv.Atoi(os.Getenv("GRPC_TIMEOUT"))
	if err != nil {
		log.Fatal(err)
	}

	return GRPCconfig{
		Host:    os.Getenv("GRPC_SERVER_HOST"),
		Port:    os.Getenv("GRPC_SERVER_PORT"),
		Timeout: timeout,
	}
}
