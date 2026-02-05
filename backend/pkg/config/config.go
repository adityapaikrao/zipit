package config

import (
	"fmt"
	"os"
	"strconv"
)

/*Defines a Struct to hold DB connection params*/
type DBConfig struct {
	Driver        string
	Host          string
	Port          int
	User          string
	Password      string
	DbName        string
	ConnectionURL string
	SSLMode       string
}

func NewDBConfig() (*DBConfig, error) {
	// 1. If DATABASE_URL is provided (e.g., in Railway/Neon), use it directly
	if connURL := os.Getenv("DATABASE_URL"); connURL != "" {
		return &DBConfig{
			Driver:        "postgres",
			ConnectionURL: connURL,
		}, nil
	}

	// 2. Fallback to individual components
	portStr := getEnvOrDefault("DB_PORT", "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %v", err)
	}

	host := getEnvOrDefault("DB_HOST", "localhost")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := getEnvOrDefault("DB_SSL_MODE", "disable")

	if user == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if dbName == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}

	return &DBConfig{
		Driver:   "postgres",
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DbName:   dbName,
		SSLMode:  sslMode,
		ConnectionURL: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			user, password, host, port, dbName, sslMode),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
