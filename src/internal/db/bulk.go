package db

import (
	"encoding/json"
	"fmt"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/types"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
)

func BulkInsert(es *elasticsearch.Client, places []types.Place) error {
	// Начинаем формировать Bulk-запрос для загрузки данных
	var bulkRequest strings.Builder
	for _, place := range places {
		//{ "index": { "_index": "places", "_id": "1" } }
		//{ "name": "Cafe 123", "address": "123 Main St", "location": { "lat": 40.73, "lon": -73.93 } }
		meta := fmt.Sprintf(`{ "index" : { "_index" : "places", "_id" : "%d" } }%s`, place.ID, "\n")
		dataJSON, err := json.Marshal(place)
		if err != nil {
			return fmt.Errorf("ошибка при сериализации данных: %s", err)
		}
		bulkRequest.WriteString(meta)
		bulkRequest.WriteString(string(dataJSON))
		bulkRequest.WriteString("\n")
	}
	// Выполняем Bulk-запрос для загрузки данных в Elasticsearch
	res, err := es.Bulk(strings.NewReader(bulkRequest.String()), es.Bulk.WithIndex("places"))
	if err != nil {
		return fmt.Errorf("ошибка при выполнении Bulk-запроса: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Ошибка в ответе на Bulk-запрос: %s", res.String())
	}
	return nil
}
