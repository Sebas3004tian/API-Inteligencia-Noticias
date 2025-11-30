package models

type CollectionConfig struct {
	VectorSize     int
	HnswM          int
	HnswEfConst    int
	PayloadIndexes map[string]string
}
