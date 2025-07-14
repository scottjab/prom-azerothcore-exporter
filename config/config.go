package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the exporter
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DSN      string
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port string
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		Database: DatabaseConfig{
			User:     os.Getenv("WOW_DB_USER"),
			Password: os.Getenv("WOW_DB_PASS"),
			Host:     os.Getenv("WOW_DB_HOST"),
			Port:     getEnvOrDefault("WOW_DB_PORT", "3306"),
		},
		Server: ServerConfig{
			Port: getEnvOrDefault("PORT", "7000"),
		},
	}

	// Check if DSN is provided directly
	if dsn := os.Getenv("WOW_DB_DSN"); dsn != "" {
		cfg.Database.DSN = dsn
	} else {
		cfg.Database.DSN = cfg.Database.buildDSN()
	}

	return cfg
}

// buildDSN builds the database connection string from individual components
func (c *DatabaseConfig) buildDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/acore_characters?parseTime=true",
		c.User, c.Password, c.Host, c.Port)
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
