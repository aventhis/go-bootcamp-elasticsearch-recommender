package main

import (
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/src/internal/data"
	elasticsearch_ "github.com/aventhis/go-bootcamp-elasticsearch-recommender/src/internal/elasticsearch"
	"log"
)

func main() {
	es, err := elasticsearch_.NewClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	err = elasticsearch_.CreateIndex(es)
	if err != nil {
		log.Fatalf("Ошибка при создании индекса: %s\n", err)
	}

	filepath := "../internal/data/data.csv"
	restaurants, err := data.LoadRestaurant(filepath)
	if err != nil {
		log.Fatalf("Ошибка загрузки данных из файла: %s\n", err)
	}
	err = elasticsearch_.BulkInsert(es, restaurants)
	if err != nil {
		log.Fatalf("Ошибка при выполнении Bulk-запроса: %s", err)
	}

	log.Println("Данные успешно загружены в Elasticsearch")
}
