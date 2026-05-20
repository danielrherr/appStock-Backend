package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port          string
	DBPath        string
	DatabaseURL   string
	JWTSecret    string
	UploadDir    string
	// Connection pool settings for Supabase
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLifetime int // in seconds
}

func Load() *Config {
	// DATABASE_URL takes precedence over DBPath (for PostgreSQL)
	databaseURL := getEnv("DATABASE_URL", "")
	
	// Connection pool defaults (Supabase friendly)
	maxOpenConns := 25 // Supabase free tier limit
	maxIdleConns := 5
	connMaxLifetime := 300 // 5 minutes in seconds
	
	return &Config{
		Port:            getEnv("PORT", "8080"),
		DBPath:          getEnv("DB_PATH", "./stockapp.db"),
		DatabaseURL:     databaseURL,
		JWTSecret:       getEnv("JWT_SECRET", "stockapp-secret-key"),
		UploadDir:       getEnv("UPLOAD_DIR", "./backend/uploads"),
		MaxOpenConns:   getEnvInt("DB_MAX_OPEN_CONNS", maxOpenConns),
		MaxIdleConns:   getEnvInt("DB_MAX_IDLE_CONNS", maxIdleConns),
		ConnMaxLifetime: getEnvInt("DB_CONN_MAX_LIFETIME", connMaxLifetime),
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}