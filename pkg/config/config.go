package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PostgresDSN string
	RedisAddr   string
	RedisPass   string
	RedisDB     string
	JWTSecret   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system env vars")
	}
	cfg := &Config{
		PostgresDSN: mustGetEnv("POSTGRES_DSN"),
		RedisAddr:   mustGetEnv("REDIS_ADDR"),
		RedisPass:   os.Getenv("REDIS_PASSWORD"),
		RedisDB:     os.Getenv("REDIS_DB"),
		JWTSecret:   mustGetEnv("JWT_SECRET"),
	}
	return cfg
}

func mustGetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	log.Fatalf("Required environment variable %s not set", key)
	return ""
}
