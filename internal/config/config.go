package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName string
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

func Default() Config {
	// Fallback defaults
	dbUser := getenvDefault("DB_USER", "root")
	dbPass := getenvDefault("DB_PASS", "")
	dbHost := getenvDefault("DB_HOST", "127.0.0.1")
	dbName := getenvDefault("DB_NAME", "games_cli")

	// Parse port with default
	port := getenvIntDefault("DB_PORT", 3306)

	return Config{
		DBUser: dbUser,
		DBPass: dbPass,
		DBHost: dbHost,
		DBPort: port,
		DBName: dbName,
	}
}

// Helper to get env or default
func getenvDefault(key, def string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return def
}

// Helper to get env int or default
func getenvIntDefault(key string, def int) int {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			return n
		}
	}
	return def
}
