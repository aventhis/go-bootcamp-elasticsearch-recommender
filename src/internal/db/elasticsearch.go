package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/types"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
)

type Store interface {
	// returns a list of items, a total number of hits and (or) an error in case of one
	GetPlaces(limit int, offset int) ([]types.Place, int, error)
}

type ElasticsearchStore struct {
	client *elasticsearch.Client
}

func NewElasticsearchStore(client *elasticsearch.Client) *ElasticsearchStore {
	return &ElasticsearchStore{client: client}
}

// limit — это количество записей
// offset — это смещение от начала списка записей
func (es *ElasticsearchStore) GetPlaces(limit int, offset int) ([]types.Place, int, error) {
	// Формируем запрос к Elasticsearch с использованием параметров `limit` и `offset`
	query := fmt.Sprintf(`{
    "from": %d,
    "size": %d,
    "sort": [{ "id": "asc" }]
	}`, offset, limit)

	// Выполняем запрос к Elasticsearch
	res, err := es.client.Search(
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex("places"),                // Указываем индекс, из которого будем получать данные
		es.client.Search.WithBody(strings.NewReader(query)), // Передаем запрос в формате JSON
		es.client.Search.WithTrackTotalHits(true),           // "hits" - общее количество найденных документов
	)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка выполнения запроса к Elasticsearch: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, 0, fmt.Errorf("ошибка в ответе от Elasticsearch: %s", res.String())
	}

	// Декодируем ответ от Elasticsearch
	var response map[string]interface{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка декодирования ответа от Elasticsearch: %w", err)
	}

	// Извлекаем "hits" и проверяем, является ли оно картой
	hits, ok := response["hits"].(map[string]interface{})
	if !ok {
		log.Fatal("Ошибка: 'hits' не является map[string]interface{}")
	}

	// Извлекаем "total" из hits и проверяем, является ли оно картой
	total, ok := hits["total"].(map[string]interface{})
	if !ok {
		log.Fatal("Ошибка: 'total' не является map[string]interface{}")
	}

	value, ok := total["value"].(float64)
	if !ok {
		// Обработка ошибки, если приведение типа не удалось
		log.Fatal("Ошибка: значение 'value' не является float64")
	}

	totalHits := int(value)

	// Получаем массив документов из "hits"
	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		log.Fatal("Ошибка: 'hits' не является []interface{}")
	}

	var places []types.Place
	for _, hit := range hitsArray {
		var place types.Place
		// Извлекаем "_source" из текущего документа
		source := hit.(map[string]interface{})["_source"]
		// Сериализуем карту "_source" в JSON-байты
		sourceBytes, err := json.Marshal(source)
		if err != nil {
			log.Println("Ошибка при сериализации '_source':", err)
			continue
		}
		if err := json.Unmarshal(sourceBytes, &place); err != nil {
			continue
		}
		places = append(places, place)
	}
	return places, totalHits, nil
}
