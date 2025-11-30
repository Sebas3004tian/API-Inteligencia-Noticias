package services

import (
	"github.com/Sebas3004tian/api-news/internal/clients/embed"
)

type EmbedService struct {
	client embed.EmbeddingClient
}

func NewEmbedService(client embed.EmbeddingClient) *EmbedService {
	return &EmbedService{client: client}
}

func (s *EmbedService) EmbedText(text string) ([]float32, error) {
	return s.client.GetEmbedding(text)
}
