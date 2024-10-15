package db

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func NewElasticsearchClient() (*elasticsearch.Client, error) {
	ctg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}

	es, err := elasticsearch.NewClient(ctg)
	if err != nil {
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	log.Println("Подключено к Elasticsearch")
	return es, nil
}
