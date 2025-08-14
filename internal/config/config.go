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
		DBUser: "avnadmin",
		DBPass: "AVNS_bOzy9VzKp8zsu1AyOen",
		DBHost: "mysql-jhp3101-geeksheykel-a001.h.aivencloud.com",
		DBPort: 21861,
		DBName: "hacknshopdb",
	}
}
