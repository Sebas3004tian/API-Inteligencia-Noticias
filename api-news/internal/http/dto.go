package http

type IndexResult struct {
	ID     string    `json:"id"`
	Status string    `json:"status,omitempty"`
	Error  string    `json:"error,omitempty"`
	Vector []float32 `json:"vector,omitempty"`
}

type SearchResponse struct {
	ID      string                 `json:"id"`
	Score   float32                `json:"score"`
	Payload map[string]interface{} `json:"item"`
}
