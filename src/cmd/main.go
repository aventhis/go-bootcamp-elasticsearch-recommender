package main

import (
	"fmt"
	elasticsearch_ "github.com/aventhis/go-bootcamp-elasticsearch-recommender/src/internal/elasticsearch"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Ошибка при создании клиента Elasticsearch: %s\n", err)
	}

	err = elasticsearch_.CreateIndex(es)
	if err != nil {
		log.Fatalf("Ошибка при создании индекса: %s\n", err)
	}

	fmt.Println("the end")
}
