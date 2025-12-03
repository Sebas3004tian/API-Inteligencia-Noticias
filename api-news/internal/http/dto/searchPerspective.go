package dto

type PerspectiveRequest struct {
	Query   string   `json:"query"`
	Sources []string `json:"sources"`
}

type PerspectiveResponseItem struct {
	Source        string  `json:"source"`
	ArticlesCount int     `json:"articles_count"`
	SimilarityAvg float64 `json:"similarity_avg"`
}

type PerspectiveResponse struct {
	Summary []PerspectiveResponseItem `json:"summary"`
}
