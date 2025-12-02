package main

import (
	"log"

	"github.com/Sebas3004tian/api-inteligencia-noticias/worker/internal/config"
	"github.com/Sebas3004tian/api-inteligencia-noticias/worker/internal/gnews"
	"github.com/Sebas3004tian/api-inteligencia-noticias/worker/internal/pipeline"
)

func main() {
	cfg := config.Load()

	client := gnews.NewClient(cfg.GNewsAPIKey)

	ingestor := pipeline.NewIngestor(
		client,
		cfg.GNewsCountry,
		cfg.GNewsLang,
		cfg.GNewsMax,
		cfg.ApiEndpoint,
	)

	if err := ingestor.Run(); err != nil {
		log.Fatal("Error ejecutando ingesta:", err)
	}
}
