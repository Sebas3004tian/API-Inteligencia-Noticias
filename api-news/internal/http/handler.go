package http

import (
	"log"

	"github.com/Sebas3004tian/api-news/internal/models"
	"github.com/Sebas3004tian/api-news/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	Embeds *services.EmbedService
}

func NewArticleHandler(e *services.EmbedService) *ArticleHandler {
	return &ArticleHandler{Embeds: e}
}
func (h *ArticleHandler) Index(c *fiber.Ctx) error {
	var articles []models.Article

	if err := c.BodyParser(&articles); err != nil {
		var single models.Article
		if err2 := c.BodyParser(&single); err2 != nil {
			return fiber.ErrBadRequest
		}
		articles = append(articles, single)
	}

	results := make([]fiber.Map, 0)

	for _, article := range articles {
		text := article.Title + " " + article.Description + " " + article.Content

		vector, err := h.Embeds.EmbedText(text)
		if err != nil {
			return err
		}

		log.Println("Artículo recibido:")
		log.Println("ID:", article.ID)
		log.Println("Title:", article.Title)
		log.Println("Vector:", vector)

		results = append(results, fiber.Map{
			"id":     article.ID,
			"status": "indexed (simulado)",
			"vector": vector,
		})
	}

	return c.JSON(fiber.Map{
		"count":   len(results),
		"results": results,
	})
}
