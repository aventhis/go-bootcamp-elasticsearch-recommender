package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/types"
	"github.com/elastic/go-elasticsearch/v8"
	"strings"
)

type ElasticsearchStore struct {
	client *elasticsearch.Client
}

func NewElasticsearchStore() (*ElasticsearchStore, error) {
	ctg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(ctg)
	if err != nil {
		return nil, err
	}
	return &ElasticsearchStore{client: es}, nil
}

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

	var places []types.Place

	return nil, 0, nil
}
