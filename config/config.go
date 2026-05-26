package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppEnv         string
	AppPort        string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTAccessExp   time.Duration
	JWTRefreshExp  time.Duration
	TestDBHost     string
	TestDBPort     string
	TestDBUser     string
	TestDBPassword string
	TestDBName     string
}

var Config AppConfig

func LoadConfig() {
	_ = godotenv.Load()

	env := os.Getenv("APP_ENV")

	if env == "test" {
		_ = godotenv.Load(".env.test")
		err := godotenv.Load(".env.test")
		if err != nil {
			log.Println(".env.test not found, trying parent path")
			_ = godotenv.Load("../.env.test")
		}
	} else {
		_ = godotenv.Load()
	}

	Config = AppConfig{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", ""),
		DBPort:        getEnv("DB_PORT", ""),
		DBUser:        getEnv("DB_USER", ""),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", ""),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-secret-change-this"),
		JWTAccessExp:  parseDuration(getEnv("JWT_ACCESS_EXP", "15m")),
		JWTRefreshExp: parseDuration(getEnv("JWT_REFRESH_EXP", "168h")),

		TestDBHost:     getEnv("TEST_DB_HOST", ""),
		TestDBPort:     getEnv("TEST_DB_PORT", ""),
		TestDBUser:     getEnv("TEST_DB_USER", ""),
		TestDBPassword: getEnv("TEST_DB_PASSWORD", ""),
		TestDBName:     getEnv("TEST_DB_NAME", ""),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env variable: %s", key)
	}
	return val
}

func parseDuration(value string) time.Duration {
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("invalid duration format: %s", value)
	}
	return d
}
