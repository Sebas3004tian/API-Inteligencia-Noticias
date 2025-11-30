package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	EmbeddingURL string
	QdrantHost   string
	QdrantPort   int
	Collection   string
}

func Load() *Config {
	_ = godotenv.Load(".env.api")

	embeddingURL := getEnv("EMBEDDING_SERVICE_URL", "")
	qdrantHost := getEnv("QDRANT_HOST", "")
	qdrantPortStr := getEnv("QDRANT_PORT", "")
	collection := getEnv("QDRANT_COLLECTION", "articles")

	if embeddingURL == "" || qdrantHost == "" || qdrantPortStr == "" {
		log.Fatal("ERROR: Missing required environment variables")
	}

	qdrantPort, err := strconv.Atoi(qdrantPortStr)
	if err != nil {
		log.Fatalf("QDRANT_PORT must be a number: %v", err)
	}

	return &Config{
		EmbeddingURL: embeddingURL,
		QdrantHost:   qdrantHost,
		QdrantPort:   qdrantPort,
		Collection:   collection,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
