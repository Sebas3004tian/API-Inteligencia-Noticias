package dto

type IndexResult struct {
	ID     string    `json:"id"`
	Status string    `json:"status,omitempty"`
	Error  string    `json:"error,omitempty"`
	Vector []float32 `json:"vector,omitempty"`
}
