package config

import (
	"os"
)

type Config struct {
	DBHost                string
	DBPort                string
	DBUser                string
	DBPassword            string
	DBName                string
	DBSSLMode             string
	JWTSecret             string
	Port                  string
	UserAttemptLimit      int
	UserWindowMinutes     int
	UserSuspensionMinutes int
	IPAttemptLimit        int
	IPWindowMinutes       int
}

func LoadConfig() *Config {
	return &Config{
		DBHost:                getEnv("DB_HOST", "localhost"),
		DBPort:                getEnv("DB_PORT", "5432"),
		DBUser:                getEnv("DB_USER", "postgres"),
		DBPassword:            getEnv("DB_PASSWORD", "postgres"),
		DBName:                getEnv("DB_NAME", "brute_force_login"),
		DBSSLMode:             getEnv("DB_SSLMODE", "disable"),
		JWTSecret:             getEnv("JWT_SECRET", ""),
		Port:                  getEnv("PORT", "8080"),
		UserAttemptLimit:      5,
		UserWindowMinutes:     5,
		UserSuspensionMinutes: 15,
		IPAttemptLimit:        100,
		IPWindowMinutes:       5,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
