package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/Sebas3004tian/api-news/internal/clients/embed"
	"github.com/Sebas3004tian/api-news/internal/clients/qdrant"
	"github.com/Sebas3004tian/api-news/internal/config"
	apiHttp "github.com/Sebas3004tian/api-news/internal/http"
	"github.com/Sebas3004tian/api-news/internal/http/handler"
	"github.com/Sebas3004tian/api-news/internal/models"
	"github.com/Sebas3004tian/api-news/internal/services"
)

func main() {
	cfg := config.Load()

	app := fiber.New()

	// Clients
	embeddingClient := embed.NewHttpEmbeddingClient(cfg.EmbeddingURL)
	qdrantClient := qdrant.NewClient(cfg.QdrantHost, cfg.QdrantPort, cfg.Collection)

	// Services
	embedService := services.NewEmbedService(embeddingClient)
	qdrantService := services.NewQdrantService(qdrantClient)

	qdrantConfig := models.CollectionConfig{
		VectorSize:     cfg.EmbeddingVectorLength,
		HnswM:          16,
		HnswEfConst:    100,
		PayloadIndexes: nil,
	}
	if err := qdrantService.EnsureCollection(qdrantConfig); err != nil {
		log.Fatal("No se pudo crear/asegurar colección:", err)
	}
	articleService := services.NewArticleService(embedService, qdrantService)

	// Handlers
	articleHandler := handler.NewArticleHandler(articleService, embedService, qdrantService)
	healthHandler := handler.NewHealthHandler()

	// Routes
	apiHttp.SetupRoutes(app, articleHandler, healthHandler)

	log.Println("Servidor iniciado en :8081")
	log.Fatal(app.Listen(":8081"))
}
