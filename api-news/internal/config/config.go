package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	EmbeddingURL          string
	EmbeddingVectorLength int
	QdrantHost            string
	QdrantPort            int
	Collection            string
}

func Load() *Config {
	_ = godotenv.Load(".env.api")

	embeddingURL := getEnv("MY_EMBEDDING_SERVICE_URL", "")
	embeddingVectorLengthStr := getEnv("EMBEDDING_VECTOR_LENGTH", "")
	qdrantHost := getEnv("QDRANT_HOST", "")
	qdrantPortStr := getEnv("QDRANT_PORT", "")
	collection := getEnv("QDRANT_COLLECTION", "articles")

	if embeddingURL == "" {
		log.Fatal("ERROR: Missing embedding service url")
	}

	if embeddingVectorLengthStr == "" || qdrantHost == "" || qdrantPortStr == "" {
		log.Fatal("ERROR: Missing required environment variables")
	}

	qdrantPort, err := strconv.Atoi(qdrantPortStr)
	if err != nil {
		log.Fatalf("QDRANT_PORT must be a number: %v", err)
	}

	embeddingVectorLength, err := strconv.Atoi(embeddingVectorLengthStr)
	if err != nil {
		log.Fatalf("EMBEDDING_VECTOR_LENGTH must be a number: %v", err)
	}

	return &Config{
		EmbeddingURL:          embeddingURL,
		EmbeddingVectorLength: embeddingVectorLength,
		QdrantHost:            qdrantHost,
		QdrantPort:            qdrantPort,
		Collection:            collection,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
