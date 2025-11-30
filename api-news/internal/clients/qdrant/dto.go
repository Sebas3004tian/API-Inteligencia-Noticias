package qdrant

type QdrantSearchRequest struct {
	Vector      []float32               `json:"vector"`
	Limit       int                     `json:"limit"`
	WithPayload bool                    `json:"with_payload"`
	Params      *QdrantSearchParams     `json:"params,omitempty"`
	Filter      *map[string]interface{} `json:"filter,omitempty"`
}

type QdrantPoint struct {
	ID      string            `json:"id"`
	Payload map[string]string `json:"payload"`
	Score   float64           `json:"score"`
}

type QdrantSearchResponse struct {
	Result []QdrantPoint `json:"result"`
}

type QdrantSearchParams struct {
	HnswEf int `json:"hnsw_ef,omitempty"`
}
