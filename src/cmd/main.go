package main

import (
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/data"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/db"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/handlers"
	"log"
	"net/http"
)

func main() {
	// Создаем клиент Elasticsearch
	esClient, err := db.NewElasticsearchClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Создаем индекс и маппинг
	err = db.CreateIndex(esClient)
	if err != nil {
		log.Fatalf("Ошибка при создании индекса: %s\n", err)
	}

	// Загружаем данные из CSV файла
	restaurants, err := data.LoadPlaces("../internal/data/data.csv")
	if err != nil {
		log.Fatalf("Ошибка загрузки данных из файла: %s\n", err)
	}

	// Выполняем Bulk-загрузку данных в Elasticsearch
	err = db.BulkInsert(esClient, restaurants)
	if err != nil {
		log.Fatalf("Ошибка при выполнении Bulk-запроса: %s", err)
	}
	log.Println("Данные успешно загружены в Elasticsearch")

	// Создаем хранилище
	store := db.NewElasticsearchStore(esClient)

	// Создаем обработчики
	handler := handlers.NewHandler(store)

	// Регистрируем маршруты
	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/api/places", handler.JSONHandler)
	http.HandleFunc("/api/recommend", handler.RecommendHandler)
	http.HandleFunc("/api/get_token", handler.GetTokenHandler)

	log.Println("Сервер запущен на http://localhost:8888")
	if err = http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %s", err)
	}
}
