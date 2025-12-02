package services

import (
	"context"

	"github.com/Sebas3004tian/api-news/internal/clients/qdrant"
	"github.com/Sebas3004tian/api-news/internal/models"
)

type QdrantService struct {
	Client *qdrant.Client
}

func NewQdrantService(client *qdrant.Client) *QdrantService {
	return &QdrantService{Client: client}
}

func (s *QdrantService) EnsureCollection(cfg models.CollectionConfig) error {
	body := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     cfg.VectorSize,
			"distance": "Cosine",
		},
		"hnsw_config": map[string]interface{}{
			"m":            cfg.HnswM,
			"ef_construct": cfg.HnswEfConst,
			"on_disk":      false,
		},
		"on_disk_payload": true,
	}

	return s.Client.EnsureCollection(body)
}

func (s *QdrantService) Insert(vector []float32, payload map[string]string) error {
	return s.Client.InsertPoint(vector, payload)
}

func (s *QdrantService) Search(ctx context.Context, vector []float32, limit int) ([]SearchResult, error) {
	req := qdrant.QdrantSearchRequest{
		Vector:      vector,
		Limit:       limit,
		WithPayload: true,
	}

	points, err := s.Client.Search(ctx, req)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, len(points))
	for i, p := range points {
		results[i] = SearchResult{
			ID:      p.ID,
			Score:   p.Score,
			Payload: p.Payload,
		}
	}

	return results, nil
}

var _ QdrantVectorService = (*QdrantService)(nil)
