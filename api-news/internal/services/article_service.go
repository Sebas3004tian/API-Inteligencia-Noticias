package services

import (
	"log"

	"github.com/Sebas3004tian/api-news/internal/models"
)

type ArticleService struct {
	Embeds *EmbedService
	Qdrant *QdrantService
}

func NewArticleService(e *EmbedService, q *QdrantService) *ArticleService {
	return &ArticleService{Embeds: e, Qdrant: q}
}

func (s *ArticleService) IndexArticles(articles []models.Article) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	for _, article := range articles {

		text := article.Title + " " + article.Description + " " + article.Content

		vector, err := s.Embeds.EmbedText(text)
		if err != nil {
			log.Println("Error embedding:", err)
			results = append(results, map[string]interface{}{
				"id":    article.ID,
				"error": err.Error(),
			})
			continue
		}

		payload := map[string]string{
			"id":          article.ID,
			"title":       article.Title,
			"description": article.Description,
			"content":     article.Content,
			"link":        article.Link,
		}

		err = s.Qdrant.InsertPoint(vector, payload)
		if err != nil {
			log.Println("Error inserting into Qdrant:", err)
			results = append(results, map[string]interface{}{
				"id":    article.ID,
				"error": "Failed to index",
			})
			continue
		}

		log.Println("Inserted:", article.ID)

		results = append(results, map[string]interface{}{
			"id":     article.ID,
			"status": "indexed",
			"vector": vector,
		})
	}

	return results, nil
}

func (s *ArticleService) SearchArticles(query string) ([]map[string]interface{}, error) {

	vector, err := s.Embeds.EmbedText(query)
	if err != nil {
		return nil, err
	}

	results, err := s.Qdrant.SearchHTTP(nil, vector, 10)
	if err != nil {
		return nil, err
	}

	var resp []map[string]interface{}

	for _, r := range results {
		resp = append(resp, map[string]interface{}{
			"id":    r.ID,
			"score": r.Score,
			"item":  r.Payload,
		})
	}

	return resp, nil
}
