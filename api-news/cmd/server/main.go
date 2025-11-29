package main

import (
	"log"
	"os"

	appHttp "github.com/Sebas3004tian/api-news/internal/http"
	"github.com/Sebas3004tian/api-news/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.api"); err != nil {
		log.Println("No .env file found, continuing...")
	}

	embeddingURL := os.Getenv("EMBEDDING_SERVICE_URL")
	if embeddingURL == "" {
		log.Fatal("ERROR: EMBEDDING_SERVICE_URL not set in .env")
	}

	app := fiber.New()

	embeddingClient := services.NewHttpEmbeddingClient(embeddingURL)

	embedService := services.NewEmbedService(embeddingClient)

	handler := appHttp.NewArticleHandler(embedService)

	app.Post("/index", handler.Index)

	log.Println("Servidor iniciado en :8080")
	log.Fatal(app.Listen(":8080"))
}
