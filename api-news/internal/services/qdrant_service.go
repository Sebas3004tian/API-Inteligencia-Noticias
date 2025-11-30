package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sebas3004tian/api-news/internal/models"
	"github.com/google/uuid"
)

type QdrantService struct {
	Host       string
	Port       int
	Collection string
}

func NewQdrantService(host string, port int, collection string) *QdrantService {
	return &QdrantService{
		Host:       host,
		Port:       port,
		Collection: collection,
	}
}

func (q *QdrantService) EnsureCollection(config models.CollectionConfig) error {
	url := fmt.Sprintf("http://%s:%d/collections/%s", q.Host, q.Port, q.Collection)

	body := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     config.VectorSize,
			"distance": "Cosine",
		},
		"hnsw_config": map[string]interface{}{
			"m":            config.HnswM,
			"ef_construct": config.HnswEfConst,
			"on_disk":      false,
		},
		"on_disk_payload": true,
	}

	if len(config.PayloadIndexes) > 0 {
		payloadSchema := make(map[string]interface{})
		for field, indexType := range config.PayloadIndexes {
			payloadSchema[field] = map[string]interface{}{
				"field_index_type": indexType,
			}
		}
		body["payload_schema"] = payloadSchema
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

func (q *QdrantService) SearchHTTP(ctx context.Context, vector []float32, limit int) ([]models.QdrantPoint, error) {
	reqBody := models.QdrantSearchRequest{
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

	var searchResp models.QdrantSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return searchResp.Result, nil
}
