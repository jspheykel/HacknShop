package config

import "fmt"

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName string
}

func (c Config) DSN() string {
	// parseTime=true ensures DATETIME -> time.Time
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

func Default() Config {
	return Config{
		DBUser: "heykel",
		DBPass: "310103",
		DBHost: "127.0.0.1",
		DBPort: 3306,
		DBName: "hacknshop_db",
	}
}
