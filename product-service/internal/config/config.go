package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	GRPCPort    int
}

func Load() (*Config, error) {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "gorm.db"),
		GRPCPort:    getEnvAsInt("GRPC_PORT", 50052),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if val, err := strconv.Atoi(value); err != nil {
			return val
		}
	}
	return defaultVal
}
