package services

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

type QdrantService struct {
	Client     *qdrant.Client
	Collection string
}

func NewQdrantService(host string, port int, collection string) *QdrantService {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: host,
		Port: port,
	})
	if err != nil {
		log.Fatal("Error creando cliente Qdrant:", err)
	}
	return &QdrantService{
		Client:     client,
		Collection: collection,
	}
}

func (q *QdrantService) EnsureCollection(vectorSize int) error {
	err := q.Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: q.Collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(vectorSize),
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if err != nil && !strings.Contains(err.Error(), "AlreadyExists") {
		return err
	}

	log.Printf("Collection '%s' asegurada.", q.Collection)
	return nil
}

func (q *QdrantService) InsertPoint(vector []float32, payload map[string]string) error {
	id := uuid.New().String()
	_, err := q.Client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: q.Collection,
		Points: []*qdrant.PointStruct{
			{
				Id: &qdrant.PointId{
					PointIdOptions: &qdrant.PointId_Uuid{Uuid: id},
				},
				Vectors: &qdrant.Vectors{
					VectorsOptions: &qdrant.Vectors_Vector{
						Vector: &qdrant.Vector{Data: vector},
					},
				},
				Payload: func() map[string]*qdrant.Value {
					out := make(map[string]*qdrant.Value)
					for k, v := range payload {
						out[k] = &qdrant.Value{
							Kind: &qdrant.Value_StringValue{
								StringValue: v,
							},
						}
					}
					return out
				}(),
			},
		},
	})
	return err
}
