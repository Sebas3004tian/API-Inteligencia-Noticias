package services

import (
	"context"

	"github.com/Sebas3004tian/api-news/internal/http/dto"
)

type EmbeddingService interface {
	EmbedText(text string) ([]float32, error)
}

type QdrantVectorService interface {
	Insert(vector []float32, payload map[string]string) error
	Search(ctx context.Context, vector []float32, limit int) ([]SearchResult, error)
	SearchByVectorAndSource(vector []float32, source string, limit int) ([]dto.ArticleWithScore, error)
}

type SearchResult struct {
	ID      string
	Score   float64
	Payload map[string]string
}
