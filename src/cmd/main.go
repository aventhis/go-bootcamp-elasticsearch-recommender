package main

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s\n", err)
	}

	//res, err := es.Info()
	//if err != nil {
	//	log.Fatalf("Error getting the response: %s\n", err)
	//}
	//
	//defer res.Body.Close()
	//
	//if res.IsError() {
	//	log.Fatalf("Error response from Elasticsearch: %s\n", res.String())
	//}
	//
	//// Выводим статус ответа
	//fmt.Printf("Elasticsearch response status: %s\n", res.Status())

	mapping := `{
		"mappings" {
			"properties": {
				"name": {
					"type":  "text"
				},
				"address": {
					"type":  "text"
				},
				"phone": {
					"type":  "text"
				},
				"location": {
				  "type": "geo_point"
				}
  			}
		}
	}`

	res, err := es.Indices.Create("places", es.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		log.Fatalf("Ошибка при создании индекса: %s", err)
	}
	defer res.Body.Close()

	log.Println("Индекс и маппинг созданы успешно")
}
