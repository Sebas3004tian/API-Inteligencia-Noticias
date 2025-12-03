package dto

import (
	"github.com/Sebas3004tian/api-news/internal/models"
)

type ArticleWithScore struct {
	models.Article
	Score float64
}
