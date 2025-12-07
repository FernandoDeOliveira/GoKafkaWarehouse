package config

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Config struct {
	OLTP DatabaseConfig
	OLAP DatabaseConfig
}

func Load() *Config {
	return &Config{
		OLTP: DatabaseConfig{
			Host:     getEnv("OLTP_DB_HOST", "localhost"),
			Port:     getEnv("OLTP_DB_PORT", "3306"),
			User:     getEnv("OLTP_DB_USER", "root"),
			Password: getEnv("OLTP_DB_PASSWORD", "password"),
			DBName:   getEnv("OLTP_DB_NAME", "oltp_db"),
		},
		OLAP: DatabaseConfig{
			Host:     getEnv("OLAP_DB_HOST", "localhost"),
			Port:     getEnv("OLAP_DB_PORT", "3307"),
			User:     getEnv("OLAP_DB_USER", "root"),
			Password: getEnv("OLAP_DB_PASSWORD", "password"),
			DBName:   getEnv("OLAP_DB_NAME", "olap_db"),
		},
	}
}


func (dc *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		dc.User, dc.Password, dc.Host, dc.Port, dc.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
