package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL string
	RedisAddr   string
	RedisPass   string
	RedisDB     int
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		PostgresURL: os.Getenv("POSTGRES_URL"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		RedisPass:   os.Getenv("REDIS_PASS"),
		RedisDB:     getEnvAsInt("REDIS_DB", 0),
	}
}

func getEnvAsInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	var i int
	_, err := fmt.Sscanf(val, "%d", &i)
	if err != nil {
		return defaultVal
	}
	return i
}
