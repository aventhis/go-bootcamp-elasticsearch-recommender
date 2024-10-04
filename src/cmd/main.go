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

	// Запускаем HTTP-сервер на порту 8888
	log.Println("Сервер запущен на порту 8888...")

	handler := handlers.NewHandler(store)

	http.HandleFunc("/", handler.ServeHTTP)

	log.Println("Сервер запущен на http://localhost:8888")
	if err = http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %s", err)
	}

	//// Проверка работы функции GetPlaces
	//limit := 10 // Лимит количества возвращаемых записей
	//offset := 2 // Смещение
	//places, totalHits, err := store.GetPlaces(limit, offset)
	//if err != nil {
	//	log.Fatalf("Ошибка при получении мест из Elasticsearch: %s", err)
	//}
	//
	//// Выводим общее количество найденных документов
	//fmt.Printf("Найдено документов: %d\n", totalHits)
	//
	//// Выводим результаты
	//for i, place := range places {
	//	fmt.Printf("Место %d: %+v\n", i+1, place)
	//}

}
