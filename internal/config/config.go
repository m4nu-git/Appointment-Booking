// Package config loads and validates application configuration from environment
// variables. It fails fast at startup rather than failing later at request time.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all runtime configuration for the service.
type Config struct {
	ServerPort  string
	Environment string // development | production | test

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret      string
	JWTExpiryHours int

	LogLevel string
}

// Load reads configuration from environment variables. Required variables that
// are missing cause an immediate error (fail fast at boot, not at request time).
func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", ""),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}

	expiryStr := getEnv("JWT_EXPIRY_HOURS", "24")
	expiry, err := strconv.Atoi(expiryStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS %q: %w", expiryStr, err)
	}
	cfg.JWTExpiryHours = expiry

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate ensures required secrets/credentials are present. This is what makes
// the service fail at boot instead of failing confusingly on the first request.
func (c *Config) validate() error {
	required := map[string]string{
		"DB_USER":     c.DBUser,
		"DB_PASSWORD": c.DBPassword,
		"DB_NAME":     c.DBName,
		"JWT_SECRET":  c.JWTSecret,
	}

	for name, value := range required {
		if value == "" {
			return fmt.Errorf("missing required environment variable: %s", name)
		}
	}

	return nil
}

// DSN builds the Postgres connection string used by GORM.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
