package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/Sebas3004tian/api-news/internal/config"
	"github.com/Sebas3004tian/api-news/internal/http"
	"github.com/Sebas3004tian/api-news/internal/models"
	"github.com/Sebas3004tian/api-news/internal/services"
)

func main() {
	cfg := config.Load()

	app := fiber.New()

	// Embedding service
	embeddingClient := services.NewHttpEmbeddingClient(cfg.EmbeddingURL)
	embedService := services.NewEmbedService(embeddingClient)

	// Qdrant service
	qdrantService := services.NewQdrantService(cfg.QdrantHost, cfg.QdrantPort, cfg.Collection)

	qdrantConfig := models.CollectionConfig{
		VectorSize:     384,
		HnswM:          16,
		HnswEfConst:    100,
		PayloadIndexes: nil,
	}

	if err := qdrantService.EnsureCollection(qdrantConfig); err != nil {
		log.Fatal("No se pudo crear/asegurar colección:", err)
	}

	articleService := services.NewArticleService(embedService, qdrantService)
	articleHandler := http.NewArticleHandler(articleService)

	app.Post("/index", articleHandler.Index)
	app.Get("/search", articleHandler.Search)

	log.Println("Servidor iniciado en :8081")
	log.Fatal(app.Listen(":8081"))
}
