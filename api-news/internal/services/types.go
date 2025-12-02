package services

import "context"

type EmbeddingService interface {
	EmbedText(text string) ([]float32, error)
}

type QdrantVectorService interface {
	Insert(vector []float32, payload map[string]string) error
	Search(ctx context.Context, vector []float32, limit int) ([]SearchResult, error)
}

type SearchResult struct {
	ID      string
	Score   float64
	Payload map[string]string
}
