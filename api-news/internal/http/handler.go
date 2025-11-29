package http

import (
	"log"

	"github.com/Sebas3004tian/api-news/internal/models"
	"github.com/Sebas3004tian/api-news/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	Embeds *services.EmbedService
	Qdrant *services.QdrantService
}

func NewArticleHandler(e *services.EmbedService, q *services.QdrantService) *ArticleHandler {
	return &ArticleHandler{Embeds: e, Qdrant: q}
}

func (h *ArticleHandler) Index(c *fiber.Ctx) error {
	var articles []models.Article

	if err := c.BodyParser(&articles); err != nil {
		return fiber.ErrBadRequest
	}

	var results []map[string]interface{}

	for _, article := range articles {
		text := article.Title + " " + article.Description + " " + article.Content

		vector, err := h.Embeds.EmbedText(text)
		if err != nil {
			log.Println("Error embedding text:", err)
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
		}

		if err := h.Qdrant.InsertPoint(vector, payload); err != nil {
			log.Println("Error inserting point:", err)
			results = append(results, map[string]interface{}{
				"id":    article.ID,
				"error": "Failed to index",
			})
			continue
		}

		log.Printf("Inserted article %s into Qdrant", article.ID)
		results = append(results, map[string]interface{}{
			"id":     article.ID,
			"status": "indexed",
			"vector": vector,
		})
	}

	return c.JSON(results)
}
