package handler

import (
	"github.com/Sebas3004tian/api-news/internal/http/dto"
	"github.com/Sebas3004tian/api-news/internal/models"
	"github.com/Sebas3004tian/api-news/internal/services"
	"github.com/gofiber/fiber/v2"
)

type ArticleHandler struct {
	Articles      *services.ArticleService
	EmbedService  services.EmbeddingService
	QdrantService services.QdrantVectorService
}

func NewArticleHandler(a *services.ArticleService, e services.EmbeddingService, q services.QdrantVectorService) *ArticleHandler {
	return &ArticleHandler{
		Articles:      a,
		EmbedService:  e,
		QdrantService: q,
	}
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

	hasErrors := false
	for _, r := range results {
		if r["error"] != nil {
			hasErrors = true
			break
		}
	}

	if hasErrors {
		return c.Status(207).JSON(results)
	}

	return c.Status(200).JSON(results)
}

func (h *ArticleHandler) Search(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		return fiber.NewError(fiber.StatusBadRequest, "query param is required")
	}

	results, err := h.Articles.SearchArticles(c.Context(), query)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(results)
}
func (h *ArticleHandler) PerspectiveAnalysis(c *fiber.Ctx) error {
	var req dto.PerspectiveRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	queryVector, err := h.EmbedService.EmbedText(req.Query)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "embedding failed")
	}

	var results []dto.PerspectiveResponseItem

	for _, source := range req.Sources {
		articles, err := h.QdrantService.SearchByVectorAndSource(queryVector, source, 10)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "search failed")
		}

		var sum float64
		for _, a := range articles {
			sum += a.Score
		}
		avg := 0.0
		if len(articles) > 0 {
			avg = sum / float64(len(articles))
		}

		results = append(results, dto.PerspectiveResponseItem{
			Source:        source,
			ArticlesCount: len(articles),
			SimilarityAvg: avg,
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.PerspectiveResponse{Summary: results})
}

func (h *ArticleHandler) GetSources(c *fiber.Ctx) error {
	sources, err := h.QdrantService.GetAllSources(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(sources)
}

func (h *ArticleHandler) GetArticlesBySource(c *fiber.Ctx) error {
	source := c.Params("source")

	articles, err := h.QdrantService.GetArticlesBySourceName(c.Context(), source)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(articles)
}
