package http

import (
	"github.com/Sebas3004tian/api-news/internal/models"
	"github.com/Sebas3004tian/api-news/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	Articles *services.ArticleService
}

func NewArticleHandler(a *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{Articles: a}
}

func (h *ArticleHandler) Index(c *fiber.Ctx) error {
	var articles []models.Article
	if err := c.BodyParser(&articles); err != nil {
		return fiber.ErrBadRequest
	}

	results, err := h.Articles.IndexArticles(articles)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(results)
}

func (h *ArticleHandler) Search(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query param is required")
	}

	results, err := h.Articles.SearchArticles(query)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(results)
}
