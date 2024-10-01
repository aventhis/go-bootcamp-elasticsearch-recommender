package main

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s\n", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting the response: %s\n", err)
	}

	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error response from Elasticsearch: %s\n", res.String())
	}

	// Выводим статус ответа
	fmt.Printf("Elasticsearch response status: %s\n", res.Status())
}
