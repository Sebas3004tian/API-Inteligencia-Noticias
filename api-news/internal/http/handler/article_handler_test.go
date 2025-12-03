package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"net/url"

	"github.com/Sebas3004tian/api-news/internal/http/dto"
	"github.com/Sebas3004tian/api-news/internal/http/handler"
	"github.com/Sebas3004tian/api-news/internal/models"
	"github.com/Sebas3004tian/api-news/internal/services"
)

type MockEmbedService struct{ mock.Mock }

func (m *MockEmbedService) EmbedText(text string) ([]float32, error) {
	args := m.Called(text)
	return args.Get(0).([]float32), args.Error(1)
}

type MockQdrantService struct{ mock.Mock }

func (m *MockQdrantService) Insert(vector []float32, payload map[string]string) error {
	args := m.Called(vector, payload)
	return args.Error(0)
}

func (m *MockQdrantService) Search(ctx context.Context, vector []float32, limit int) ([]services.SearchResult, error) {
	args := m.Called(ctx, vector, limit)
	return args.Get(0).([]services.SearchResult), args.Error(1)
}

func (m *MockQdrantService) SearchByVectorAndSource(vector []float32, source string, limit int) ([]dto.ArticleWithScore, error) {
	args := m.Called(vector, source, limit)
	return args.Get(0).([]dto.ArticleWithScore), args.Error(1)
}

func TestArticleHandler_Index_Unit(t *testing.T) {
	mockEmb := new(MockEmbedService)
	mockQ := new(MockQdrantService)

	articles := []models.Article{
		{
			Title:       "Noticia de prueba",
			Description: "Descripción corta",
			Content:     "Contenido completo",
			Url:         "https://example.com",
		},
	}

	mockEmb.On("EmbedText", "Noticia de prueba Descripción corta Contenido completo").Return([]float32{0.1, 0.2, 0.3}, nil)
	mockQ.On("Insert", []float32{0.1, 0.2, 0.3}, mock.Anything).Return(nil)

	articleService := services.NewArticleService(mockEmb, mockQ)
	h := handler.NewArticleHandler(articleService, mockEmb, mockQ)

	app := fiber.New()
	app.Post("/index", h.Index)

	body, _ := json.Marshal(articles)
	req := httptest.NewRequest("POST", "/index", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockEmb.AssertCalled(t, "EmbedText", "Noticia de prueba Descripción corta Contenido completo")
	mockQ.AssertCalled(t, "Insert", []float32{0.1, 0.2, 0.3}, mock.Anything)
}

func TestArticleHandler_Search_Unit(t *testing.T) {
	mockEmb := new(MockEmbedService)
	mockQ := new(MockQdrantService)

	mockEmb.On("EmbedText", "consulta de prueba").Return([]float32{0.5, 0.1, -0.2}, nil)
	mockQ.On("Search", mock.Anything, []float32{0.5, 0.1, -0.2}, 10).Return([]services.SearchResult{
		{ID: "1", Score: 0.9, Payload: map[string]string{"title": "Resultado de prueba"}},
	}, nil)

	articleService := services.NewArticleService(mockEmb, mockQ)
	h := handler.NewArticleHandler(articleService, mockEmb, mockQ)

	app := fiber.New()
	app.Get("/search", h.Search)

	req := httptest.NewRequest("GET", "/search?query="+url.QueryEscape("consulta de prueba"), nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockEmb.AssertCalled(t, "EmbedText", "consulta de prueba")
	mockQ.AssertCalled(t, "Search", mock.Anything, []float32{0.5, 0.1, -0.2}, 10)
}
