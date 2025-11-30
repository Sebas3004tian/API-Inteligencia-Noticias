package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GNewsAPIKey  string
	GNewsLang    string
	GNewsMax     string
	GNewsCountry string
	ApiEndpoint  string
}

func Load() *Config {
	err := godotenv.Load(".env.worker")
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		GNewsAPIKey:  getEnv("GNEWS_API_KEY", ""),
		GNewsLang:    getEnv("GNEWS_LANG", "es"),
		GNewsMax:     getEnv("GNEWS_MAX", "10"),
		GNewsCountry: getEnv("GNEWS_COUNTRY", "co"),
		ApiEndpoint:  getEnv("API_ENDPOINT", "http://localhost:8080/index"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
