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

func (m *MockQdrantService) EnsureCollection(cfg models.CollectionConfig) error {
	args := m.Called(cfg)
	return args.Error(0)
}

func (m *MockQdrantService) GetAllSources(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockQdrantService) GetArticlesBySourceName(ctx context.Context, source string) ([]models.Article, error) {
	args := m.Called(ctx, source)
	return args.Get(0).([]models.Article), args.Error(1)
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

func TestArticleHandler_GetSources_Unit(t *testing.T) {
	mockEmb := new(MockEmbedService)
	mockQ := new(MockQdrantService)

	mockQ.On("GetAllSources", mock.Anything).Return([]string{
		"El Tiempo", "Semana", "El Espectador",
	}, nil)

	articleService := services.NewArticleService(mockEmb, mockQ)
	h := handler.NewArticleHandler(articleService, mockEmb, mockQ)

	app := fiber.New()
	app.Get("/sources", h.GetSources)

	req := httptest.NewRequest("GET", "/sources", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result []string
	json.NewDecoder(resp.Body).Decode(&result)

	assert.Contains(t, result, "El Tiempo")
	assert.Contains(t, result, "Semana")
	assert.Contains(t, result, "El Espectador")

	mockQ.AssertCalled(t, "GetAllSources", mock.Anything)
}

func TestArticleHandler_GetArticlesBySource_Unit(t *testing.T) {
	mockEmb := new(MockEmbedService)
	mockQ := new(MockQdrantService)

	mockQ.On(
		"GetArticlesBySourceName",
		mock.Anything,
		"infobae",
	).Return([]models.Article{
		{
			ID:          "123",
			Title:       "Noticia 1",
			Description: "Desc",
			Content:     "Contenido",
			Url:         "https://example.com",
			Image:       "",
			PublishedAt: "2025-12-02T01:11:09Z",
			SourceName:  "infobae",
			SourceURL:   "https://www.infobae.com",
		},
		{
			ID:          "456",
			Title:       "Noticia 2",
			Description: "Desc 2",
			Content:     "Contenido 2",
			Url:         "https://example2.com",
			Image:       "",
			PublishedAt: "2025-12-02T02:11:09Z",
			SourceName:  "infobae",
			SourceURL:   "https://www.infobae.com",
		},
	}, nil)

	articleService := services.NewArticleService(mockEmb, mockQ)
	h := handler.NewArticleHandler(articleService, mockEmb, mockQ)

	app := fiber.New()
	app.Get("/articles/source/:source", h.GetArticlesBySource)

	req := httptest.NewRequest("GET", "/articles/source/infobae", nil)
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var articles []models.Article
	err := json.NewDecoder(resp.Body).Decode(&articles)
	assert.Nil(t, err)

	assert.Len(t, articles, 2)
	assert.Equal(t, "infobae", articles[0].SourceName)

	mockQ.AssertCalled(t, "GetArticlesBySourceName", mock.Anything, "infobae")
}

func TestArticleHandler_PerspectiveAnalysis_Unit(t *testing.T) {
	mockEmb := new(MockEmbedService)
	mockQ := new(MockQdrantService)

	body := map[string]interface{}{
		"query":   "Casa emerito",
		"sources": []string{"El Economista", "Europa Press"},
	}

	jsonBody, _ := json.Marshal(body)

	vector := []float32{0.1, 0.2, 0.3}

	mockEmb.
		On("EmbedText", "Casa emerito").
		Return(vector, nil)

	mockQ.On(
		"SearchByVectorAndSource",
		vector,
		"El Economista",
		10,
	).Return([]dto.ArticleWithScore{
		{Score: 0.30},
		{Score: 0.28},
		{Score: 0.27},
	}, nil)

	mockQ.On(
		"SearchByVectorAndSource",
		vector,
		"Europa Press",
		10,
	).Return([]dto.ArticleWithScore{
		{Score: 0.25},
		{Score: 0.20},
		{Score: 0.18},
		{Score: 0.28},
	}, nil)

	articleService := services.NewArticleService(mockEmb, mockQ)
	h := handler.NewArticleHandler(articleService, mockEmb, mockQ)

	app := fiber.New()
	app.Post("/analysis/perspective", h.PerspectiveAnalysis)

	req := httptest.NewRequest("POST", "/analysis/perspective", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string][]map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	summary := result["summary"]
	assert.Len(t, summary, 2)

	assert.Equal(t, "El Economista", summary[0]["source"])
	assert.Equal(t, float64(3), summary[0]["articles_count"])

	assert.Equal(t, "Europa Press", summary[1]["source"])
	assert.Equal(t, float64(4), summary[1]["articles_count"])

	mockEmb.AssertCalled(t, "EmbedText", "Casa emerito")
	mockQ.AssertNumberOfCalls(t, "SearchByVectorAndSource", 2)
}
