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
	var article models.Article

	if err := c.BodyParser(&article); err != nil {
		return fiber.ErrBadRequest
	}

	text := article.Title + " " + article.Description + " " + article.Content

	vector, err := h.Embeds.EmbedText(text)
	if err != nil {
		return err
	}

	log.Println("Artículo recibido:")
	log.Println("ID: ", article.ID)
	log.Println("Title:", article.Title)
	log.Println("Vector:", vector)

	return c.JSON(fiber.Map{
		"status": "indexed (simulado)",
		"vector": vector,
	})
}
