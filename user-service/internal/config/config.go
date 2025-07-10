package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	GRPCPort    int
}

func Load() (*Config, error) {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "gorm.db"),
		JWTSecret:   getEnv("JWT_SECRET", "secret-key"),
		GRPCPort:    getEnvAsInt("GRPC_PORT", 50051),
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
		var i int
		if _, err := fmt.Scanf(value, "%d", &i); err == nil {
			return i
		}
	}
	return defaultVal
}
