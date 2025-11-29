package main

import (
	"log"
	"os"
	"strconv"

	appHttp "github.com/Sebas3004tian/api-news/internal/http"
	"github.com/Sebas3004tian/api-news/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.api")

	embeddingURL := os.Getenv("EMBEDDING_SERVICE_URL")
	qdrantHOST := os.Getenv("QDRANT_HOST")
	qdrantPORTStr := os.Getenv("QDRANT_PORT")

	if embeddingURL == "" || qdrantHOST == "" || qdrantPORTStr == "" {
		log.Fatal("ERROR: Missing env vars")
	}

	qdrantPORT, err := strconv.Atoi(qdrantPORTStr)
	if err != nil {
		log.Fatal("QDRANT_PORT must be a number")
	}

	app := fiber.New()

	// Embedding service
	embeddingClient := services.NewHttpEmbeddingClient(embeddingURL)
	embedService := services.NewEmbedService(embeddingClient)

	// Qdrant service
	qdrantService := services.NewQdrantService(qdrantHOST, qdrantPORT, "articles")
	if err := qdrantService.EnsureCollection(384); err != nil {
		log.Fatal("No se pudo crear/asegurar colección:", err)
	}

	handler := appHttp.NewArticleHandler(embedService, qdrantService)

	app.Post("/index", handler.Index)

	log.Println("Servidor iniciado en :8081")
	log.Fatal(app.Listen(":8081"))
}
