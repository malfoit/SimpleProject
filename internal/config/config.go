package config

import "os"

type Config struct {
	GRPC GRPCConfig
}

type GRPCConfig struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		GRPC: GRPCConfig{
			Port: getEnv("GRPC_PORT", "50051"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
