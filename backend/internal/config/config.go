package config

import (
	"os"
)

type Config struct {
	Port        string
	DBPath      string
	JWTSecret   string
	UploadDir   string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "./stockapp.db"),
		JWTSecret: getEnv("JWT_SECRET", "stockapp-secret-key"),
		UploadDir: getEnv("UPLOAD_DIR", "./uploads"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}