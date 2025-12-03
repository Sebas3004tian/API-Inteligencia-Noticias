package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Client struct {
	Host       string
	Port       int
	Collection string
}

func NewClient(host string, port int, collection string) *Client {
	return &Client{Host: host, Port: port, Collection: collection}
}

func (c *Client) EnsureCollection(body map[string]interface{}) error {
	url := fmt.Sprintf("http://%s:%d/collections/%s", c.Host, c.Port, c.Collection)

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

func (c *Client) InsertPoint(vector []float32, payload map[string]string) error {
	url := fmt.Sprintf("http://%s:%d/collections/%s/points",
		c.Host, c.Port, c.Collection)

	body := map[string]interface{}{
		"points": []map[string]interface{}{
			{
				"id":      uuid.New().String(),
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

func (c *Client) Search(ctx context.Context, reqBody QdrantSearchRequest) ([]QdrantPoint, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"http://%s:%d/collections/%s/points/search",
		c.Host, c.Port, c.Collection,
	)

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

	var searchResp QdrantSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, err
	}

	return searchResp.Result, nil
}

func (c *Client) GetBySourceName(ctx context.Context, source string) ([]QdrantPoint, error) {
	url := fmt.Sprintf(
		"http://%s:%d/collections/%s/points/scroll",
		c.Host, c.Port, c.Collection,
	)

	filter := map[string]interface{}{
		"must": []map[string]interface{}{
			{
				"key": "source_name",
				"match": map[string]interface{}{
					"value": source,
				},
			},
		},
	}

	body := map[string]interface{}{
		"with_payload": true,
		"limit":        500,
		"filter":       filter,
	}

	bodyBytes, _ := json.Marshal(body)

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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("scroll failed: HTTP %d", resp.StatusCode)
	}

	var scrollResp QdrantScrollResponse
	if err := json.NewDecoder(resp.Body).Decode(&scrollResp); err != nil {
		return nil, err
	}

	return scrollResp.Result.Points, nil
}

func (c *Client) Scroll(ctx context.Context, body map[string]interface{}) ([]QdrantPoint, error) {
	url := fmt.Sprintf(
		"http://%s:%d/collections/%s/points/scroll",
		c.Host, c.Port, c.Collection,
	)

	bodyBytes, _ := json.Marshal(body)

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

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("scroll failed: HTTP %d", resp.StatusCode)
	}

	var scrollResp QdrantScrollResponse
	if err := json.NewDecoder(resp.Body).Decode(&scrollResp); err != nil {
		return nil, err
	}

	return scrollResp.Result.Points, nil
}
