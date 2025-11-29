package services

import (
	"crypto/sha1"
	"encoding/binary"
)

type EmbedService struct{}

func NewEmbedService() *EmbedService {
	return &EmbedService{}
}

// Simulación de un embedding, mientras pienso como integrar un servicio real
func (s *EmbedService) EmbedText(text string) []float32 {
	h := sha1.Sum([]byte(text))
	vector := make([]float32, 8)

	for i := 0; i < 8; i++ {
		chunk := h[i*2 : i*2+2]
		vector[i] = float32(binary.BigEndian.Uint16(chunk)) / 65535.0
	}

	return vector
}
