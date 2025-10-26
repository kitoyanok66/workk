package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret   string
	JWTTTLHours time.Duration

	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	DBSSLMode  string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	jwtTTLStr := os.Getenv("JWT_TTL_HOURS")
	if jwtTTLStr == "" {
		jwtTTLStr = "24"
	}
	jwtTTL, err := strconv.Atoi(jwtTTLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_TTL_HOURS: %w", err)
	}

	cfg := &Config{
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTTTLHours: time.Duration(jwtTTL) * time.Hour,

		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required database environment variables")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("missing JWT_SECRET in environment")
	}

	return cfg, nil
}
