package main

import (
	"fmt"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/src/internal/data"
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

	filepath := "../internal/data/data.csv"
	restaurants, err := data.LoadRestaurant(filepath)
	//fmt.Println(restaurants[0])
	fmt.Println(restaurants[1])
	fmt.Println("the end")
}
