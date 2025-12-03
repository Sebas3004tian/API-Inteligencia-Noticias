package http

import (
	"github.com/Sebas3004tian/api-news/internal/http/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, articleHandler *handler.ArticleHandler, healthHandler *handler.HealthHandler) {
	app.Post("/index", articleHandler.Index)
	app.Get("/search", articleHandler.Search)
	app.Get("/", healthHandler.Health)
	app.Get("/analysis/perspective", articleHandler.PerspectiveAnalysis)

}
