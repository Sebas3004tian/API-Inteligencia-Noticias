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
	Query       string
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

func NewIngestor(client *gnews.Client, q, lang, max string, endpoint string) *Ingestor {
	return &Ingestor{
		Client:      client,
		Query:       q,
		Lang:        lang,
		Max:         max,
		ApiEndpoint: endpoint,
	}
}

func (i *Ingestor) Run() error {
	log.Println("Consultando GNews...")

	resp, err := i.Client.Search(i.Query, i.Lang, i.Max)
	if err != nil {
		return err
	}

	log.Printf("Recibidos %d artículos\n", len(resp.Articles))

	var payload []ArticlePayload

	for _, a := range resp.Articles {
		payload = append(payload, ArticlePayload{
			ID:          uuid.New().String(),
			Title:       a.Title,
			Description: a.Description,
			Content:     a.Content,
		})
	}

	data, err := json.Marshal(payload)
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
