package models

type Article struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Url         string `json:"url"`
	Image       string `json:"image"`
	PublishedAt string `json:"publishedAt"`
	SourceName  string `json:"sourceName"`
	SourceURL   string `json:"sourceUrl"`
}
