package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type QdrantService struct {
	Host       string
	Port       int
	Collection string
}

type qdrantSearchRequest struct {
	Vector      []float32 `json:"vector"`
	Limit       int       `json:"limit"`
	WithPayload bool      `json:"with_payload"`
}

type qdrantPoint struct {
	ID      string            `json:"id"`
	Payload map[string]string `json:"payload"`
	Score   float64           `json:"score"`
}

type qdrantSearchResponse struct {
	Result []qdrantPoint `json:"result"`
}

func NewQdrantService(host string, port int, collection string) *QdrantService {
	return &QdrantService{
		Host:       host,
		Port:       port,
		Collection: collection,
	}
}

func (q *QdrantService) EnsureCollection(vectorSize int) error {
	url := fmt.Sprintf("http://%s:%d/collections/%s", q.Host, q.Port, q.Collection)

	body := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     vectorSize,
			"distance": "Cosine",
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 409 {
		return fmt.Errorf("failed to create collection: HTTP %d", resp.StatusCode)
	}

	return nil
}

func (q *QdrantService) InsertPoint(vector []float32, payload map[string]string) error {
	url := fmt.Sprintf("http://%s:%d/collections/%s/points", q.Host, q.Port, q.Collection)

	id := uuid.New().String()

	body := map[string]interface{}{
		"points": []map[string]interface{}{
			{
				"id":      id,
				"vector":  vector,
				"payload": payload,
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to insert point: HTTP %d", resp.StatusCode)
	}

	return nil
}

func (q *QdrantService) SearchHTTP(ctx context.Context, vector []float32, limit int) ([]qdrantPoint, error) {
	reqBody := qdrantSearchRequest{
		Vector:      vector,
		Limit:       limit,
		WithPayload: true,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%d/collections/%s/points/search", q.Host, q.Port, q.Collection)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("qdrant returned status %d", resp.StatusCode)
	}

	var searchResp qdrantSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return searchResp.Result, nil
}
