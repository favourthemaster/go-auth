package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func Init() error {
	// Load environment variables from .env file
	if err := LoadEnv(); err != nil {
		return err
	}

	// Additional initialization logic can be added here if needed

	return nil
}

// LoadEnv loads environment variables from a .env file
func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

// GetEnv retrieves an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetMailerConfig() MailerConfig {
	return MailerConfig{
		Sender: getEnv("MAIL_SENDER", "noreply@project.com"),
	}
}

func GetRedisConfig() RedisConfig {
	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		db = 0 // Default to 0 if conversion fails
	}
	return RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		Database: db,
		SSLMode:  getEnv("REDIS_SSLMODE", "disable"),
	}
}

// GetDBConfig returns the database configuration from environment variables
func GetDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		Username: getEnv("DB_USERNAME", "postgres"),
		Password: getEnv("DB_PASSWORD", "admin"),
		Database: getEnv("DB_DATABASE", "testdb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}
