package main

import (
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/data"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/db"
	"log"
)

func main() {
	esClient, err := db.NewElasticsearchClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	err = db.CreateIndex(esClient)
	if err != nil {
		log.Fatalf("Ошибка при создании индекса: %s\n", err)
	}

	restaurants, err := data.LoadRestaurant("../../internal/data/data.csv")
	if err != nil {
		log.Fatalf("Ошибка загрузки данных из файла: %s\n", err)
	}
	err = db.BulkInsert(esClient, restaurants)
	if err != nil {
		log.Fatalf("Ошибка при выполнении Bulk-запроса: %s", err)
	}

	log.Println("Данные успешно загружены в Elasticsearch")
}
