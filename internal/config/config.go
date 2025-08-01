package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Db               string
	DbPassword       string
	DbUsername       string
	DbPort           string
	DbHost           string
	Port             int
	JWTSecret        string
	JWTExpiryMinutes time.Duration
}

// New returns a config object from the env and a non-nil error if the env value is not present
func New() (*Config, error) {
	// database configs
	db := getEnv("DB_DATABASE", "chow")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbUsername := getEnv("DB_USERNAME", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbHost := getEnv("DB_HOST", "localhost")

	// server configs
	port := getEnvInt("PORT", 8000)
	jwtSecret := getEnv("JWT_SECRET", "random-token")
	jwtExpiry := getEnvInt("JWT_EXPIRY_MINUTES", 60)

	return &Config{
		Db:               db,
		DbPassword:       dbPassword,
		DbUsername:       dbUsername,
		DbPort:           dbPort,
		DbHost:           dbHost,
		Port:             port,
		JWTSecret:        jwtSecret,
		JWTExpiryMinutes: time.Duration(jwtExpiry),
	}, nil
}

func getEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	intVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return fallback
	}
	return int(intVal)
}
