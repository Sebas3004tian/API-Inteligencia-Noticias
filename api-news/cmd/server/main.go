package main

import (
	"log"

	"github.com/Sebas3004tian/api-news/internal/http"
	"github.com/Sebas3004tian/api-news/internal/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	embedService := services.NewEmbedService()
	handler := http.NewArticleHandler(embedService)

	app.Post("/index", handler.Index)

	log.Println("Servidor iniciado en :8080")
	log.Fatal(app.Listen(":8080"))
}
