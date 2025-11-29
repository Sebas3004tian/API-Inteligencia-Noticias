package services

type EmbedService struct {
	client EmbeddingClient
}

func NewEmbedService(client EmbeddingClient) *EmbedService {
	return &EmbedService{client: client}
}

func (s *EmbedService) EmbedText(text string) ([]float32, error) {
	return s.client.GetEmbedding(text)
}
