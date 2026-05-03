package config

import (
	"os"
)

type Config struct {
	GRPC   GRPCConfig
	GitHub GitHubConfig
}

type GRPCConfig struct {
	Port string
}

type GitHubConfig struct {
	Token string
	Owner string
	Repo  string
}

func NewConfig() *Config {
	return &Config{
		GRPC: GRPCConfig{
			Port: getEnv("GRPC_PORT", "50051"),
		},
		GitHub: GitHubConfig{
			Token: getEnv("GITHUB_TOKEN", ""),
			Owner: getEnv("GITHUB_OWNER", "malfoit"),
			Repo:  getEnv("GITHUB_REPO", "SimpleProject"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
