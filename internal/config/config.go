package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	NatsURL  string
	Stream   string
	Subject  string
	HTTPPort string
}

func Load() *Config {
	godotenv.Load()

	cfg := &Config{
		NatsURL:  getEnv("NATS_URL", "nats://localhost:4222"),
		Stream:   getEnv("NATS_STREAM", "GEO_STREAM"),
		Subject:  getEnv("NATS_SUBJECT", "geo.coordinates"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
