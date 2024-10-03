package main

import (
	"fmt"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/db"
	"log"
)

func main() {
	store, err := db.NewElasticsearchStore()
	if err != nil {
		log.Fatalf("Ошибка при создании Elasticsearch хранилища: %s", err)
	}

	// Проверка работы функции GetPlaces
	limit := 10 // Лимит количества возвращаемых записей
	offset := 0 // Смещение
	places, totalHits, err := store.GetPlaces(limit, offset)
	if err != nil {
		log.Fatalf("Ошибка при получении мест из Elasticsearch: %s", err)
	}

	// Выводим общее количество найденных документов
	fmt.Printf("Найдено документов: %d\n", totalHits)

	// Выводим результаты
	for i, place := range places {
		fmt.Printf("Место %d: %+v\n", i+1, place)
	}
}
