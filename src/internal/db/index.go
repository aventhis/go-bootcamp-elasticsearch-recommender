package db

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
)

func CreateIndex(es *elasticsearch.Client) error {
	mapping := `{
		"mappings": {
			"properties": {
				"id": {
					"type": "long"
				},
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

	// Удаляем индекс, если он уже существует
	res, err := es.Indices.Delete([]string{"places"}, es.Indices.Delete.WithIgnoreUnavailable(true))
	if err != nil {
		return fmt.Errorf("ошибка при удалении индекса: %s", err)
	}
	defer res.Body.Close()

	// Создаем индекс с маппингом
	res, err = es.Indices.Create(
		"places",
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("ошибка при создании индекса: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ошибка ответа от Elasticsearch: %s", err)
	}

	log.Println("Индекс 'places' и маппинг созданы успешно")
	return nil
}
