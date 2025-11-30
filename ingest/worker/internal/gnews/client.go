package gnews

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	APIKey string
}

func NewClient(apiKey string) *Client {
	return &Client{APIKey: apiKey}
}

func (c *Client) Search(q string, country string, lang string, max string) (*Response, error) {
	url := fmt.Sprintf(
		"https://gnews.io/api/v4/search?q=%s&lang=%s&max=%s&apikey=%s",
		q, lang, max, c.APIKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
