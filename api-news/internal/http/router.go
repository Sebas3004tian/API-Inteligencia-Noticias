package http

import (
	"github.com/Sebas3004tian/api-news/internal/http/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, articleHandler *handler.ArticleHandler, healthHandler *handler.HealthHandler) {
	app.Post("/index", articleHandler.Index)
	app.Get("/search", articleHandler.Search)
	app.Get("/", healthHandler.Health)
	app.Post("/analysis/perspective", articleHandler.PerspectiveAnalysis)
	app.Get("/sources", articleHandler.GetSources)
	app.Get("/articles/source/:source", articleHandler.GetArticlesBySource)

}
