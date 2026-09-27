package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerPort       string
	JWTSecret        string
	HospitalABaseURL string
	HospitalBBaseURL string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
}

func Load() *Config {
	return &Config{
		ServerPort:       getEnv("SERVER_PORT", ":8080"),
		JWTSecret:        getEnv("JWT_SECRET", "super-secret-key-5678"),
		HospitalABaseURL: getEnv("HOSPITAL_A_BASE_URL", "https://hospital-a.api.co.th"),
		HospitalBBaseURL: getEnv("HOSPITAL_B_BASE_URL", "https://hospital-b.api.co.th"),
		DBHost:           getEnv("DB_HOST", "127.0.0.1"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "KhonNaRak5555"),
		DBName:           getEnv("DB_NAME", "hospital_middleware"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
