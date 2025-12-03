package services

import (
	"context"

	"github.com/Sebas3004tian/api-news/internal/clients/qdrant"
	"github.com/Sebas3004tian/api-news/internal/http/dto"
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
		score := 0.0
		if p.Score != nil {
			score = *p.Score
		}

		results[i] = SearchResult{
			ID:      p.ID,
			Score:   score,
			Payload: p.Payload,
		}
	}

	return results, nil
}

func (s *QdrantService) SearchByVectorAndSource(vector []float32, source string, limit int) ([]dto.ArticleWithScore, error) {
	filter := map[string]interface{}{
		"must": []map[string]interface{}{
			{
				"key": "source_name",
				"match": map[string]string{
					"value": source,
				},
			},
		},
	}

	req := qdrant.QdrantSearchRequest{
		Vector:      vector,
		Limit:       limit,
		WithPayload: true,
		Filter:      &filter,
	}

	points, err := s.Client.Search(context.Background(), req)
	if err != nil {
		return nil, err
	}

	articles := make([]dto.ArticleWithScore, len(points))
	for i, p := range points {
		score := 0.0
		if p.Score != nil {
			score = *p.Score
		}

		pl := p.Payload

		articles[i] = dto.ArticleWithScore{
			Article: models.Article{
				ID:          p.ID,
				Title:       pl["title"],
				Description: pl["description"],
				Content:     pl["content"],
				Url:         pl["url"],
				Image:       pl["image"],
				PublishedAt: pl["published_at"],
				SourceName:  pl["source_name"],
				SourceURL:   pl["source_url"],
			},
			Score: score,
		}
	}

	return articles, nil
}

func (s *QdrantService) GetAllSources(ctx context.Context) ([]string, error) {
	body := map[string]interface{}{
		"limit":        5000,
		"with_payload": true,
	}

	points, err := s.Client.Scroll(ctx, body)
	if err != nil {
		return nil, err
	}

	unique := map[string]struct{}{}
	for _, p := range points {
		if src, ok := p.Payload["source_name"]; ok {
			unique[src] = struct{}{}
		}
	}

	var sources []string
	for src := range unique {
		sources = append(sources, src)
	}

	return sources, nil
}

func (s *QdrantService) GetArticlesBySourceName(ctx context.Context, source string) ([]models.Article, error) {
	points, err := s.Client.GetBySourceName(ctx, source)
	if err != nil {
		return nil, err
	}

	articles := make([]models.Article, len(points))
	for i, p := range points {
		pl := p.Payload
		articles[i] = models.Article{
			ID:          p.ID,
			Title:       pl["title"],
			Description: pl["description"],
			Content:     pl["content"],
			Url:         pl["url"],
			Image:       pl["image"],
			PublishedAt: pl["published_at"],
			SourceName:  pl["source_name"],
			SourceURL:   pl["source_url"],
		}
	}

	return articles, nil
}

var _ QdrantVectorService = (*QdrantService)(nil)
