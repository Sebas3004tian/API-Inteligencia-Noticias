package pipeline

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/Sebas3004tian/api-inteligencia-noticias/worker/internal/gnews"
)

type Ingestor struct {
	Client      *gnews.Client
	Country     string
	Lang        string
	Max         string
	ApiEndpoint string
}

type ArticlePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

func NewIngestor(client *gnews.Client, country, lang, max string, endpoint string) *Ingestor {
	return &Ingestor{
		Client:      client,
		Country:     country,
		Lang:        lang,
		Max:         max,
		ApiEndpoint: endpoint,
	}
}

func (i *Ingestor) Run() error {
	log.Println("Consultando GNews con múltiples queries...")

	letters := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}

	articleMap := make(map[string]ArticlePayload)

	for _, q := range letters {
		log.Printf(" - Consultando con query='%s'", q)

		resp, err := i.Client.Search(q, i.Country, i.Lang, i.Max)
		if err != nil {
			log.Printf("Error con query '%s': %v", q, err)
			continue
		}

		for _, a := range resp.Articles {
			id := uuid.New().String()

			key := a.URL
			if key == "" {
				key = a.Title
			}

			if _, exists := articleMap[key]; exists {
				continue
			}

			articleMap[key] = ArticlePayload{
				ID:          id,
				Title:       a.Title,
				Description: a.Description,
				Content:     a.Content,
			}
		}
	}

	var finalPayload []ArticlePayload
	for _, v := range articleMap {
		finalPayload = append(finalPayload, v)
	}

	log.Printf("Total final sin duplicados: %d artículos", len(finalPayload))

	data, err := json.Marshal(finalPayload)
	if err != nil {
		return fmt.Errorf("error serializando payload: %w", err)
	}

	log.Println("Enviando artículos al API:", i.ApiEndpoint)

	req, err := http.NewRequest("POST", i.ApiEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	respApi, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error enviando datos al API: %w", err)
	}
	defer respApi.Body.Close()

	if respApi.StatusCode < 200 || respApi.StatusCode > 299 {
		return fmt.Errorf("API respondió con código %d", respApi.StatusCode)
	}

	log.Println("Ingesta completada exitosamente")
	return nil
}
