package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type EmbeddingClient interface {
	GetEmbedding(text string) ([]float32, error)
}

type HttpEmbeddingClient struct {
	BaseURL string
}

func NewHttpEmbeddingClient(baseURL string) *HttpEmbeddingClient {
	return &HttpEmbeddingClient{BaseURL: baseURL}
}

func (c *HttpEmbeddingClient) GetEmbedding(text string) ([]float32, error) {
	payload := map[string]string{"text": text}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(c.BaseURL+"/embed", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Embedding []float32 `json:"embedding"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Embedding, nil
}
