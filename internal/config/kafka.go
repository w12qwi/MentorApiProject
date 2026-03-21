package config

import "os"

type KafkaConfig struct {
	BrokerHost string
	BrokerPort string
	Topic      string
}

func LoadKafkaConfig() KafkaConfig {
	return KafkaConfig{
		BrokerHost: os.Getenv("KAFKA_BROKER_HOST"),
		BrokerPort: os.Getenv("KAFKA_BROKER_PORT"),
		Topic:      os.Getenv("KAFKA_TOPIC"),
	}
}
