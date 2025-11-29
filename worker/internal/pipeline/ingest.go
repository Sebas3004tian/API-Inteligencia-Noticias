package pipeline

import (
	"fmt"
	"log"

	"github.com/Sebas3004tian/api-inteligencia-noticias/worker/internal/gnews"
)

type Ingestor struct {
	Client *gnews.Client
	Query  string
	Lang   string
	Max    string
}

func NewIngestor(client *gnews.Client, q, lang, max string) *Ingestor {
	return &Ingestor{
		Client: client,
		Query:  q,
		Lang:   lang,
		Max:    max,
	}
}

func (i *Ingestor) Run() error {
	log.Println(" Consultando GNews...")

	resp, err := i.Client.Search(i.Query, i.Lang, i.Max)
	if err != nil {
		return err
	}

	log.Printf("Recibidos %d artículos\n", len(resp.Articles))

	for idx, a := range resp.Articles {
		fmt.Printf("\n [%d] %s\n", idx+1, a.Title)
		fmt.Printf("   %s\n", a.Description)
		fmt.Printf("   Fuente: %s (%s)\n", a.Source.Name, a.Source.URL)
		fmt.Printf("   URL: %s\n", a.URL)
		fmt.Printf("   Publicado en: %s\n", a.PublishedAt)
	}

	return nil
}
