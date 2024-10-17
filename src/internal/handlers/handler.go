package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/db"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/types"
	"github.com/golang-jwt/jwt/v4"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	store db.Store
}

// секретный ключ для подписи токенов
var jwtKey = []byte("my_secret_key")

func NewHandler(store db.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр "page" из запроса
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, fmt.Sprintf("Некорректное значение 'page': '%s'", pageStr), http.StatusBadRequest)
		return
	}

	// Количество записей на одной странице
	limit := 10
	offset := (page - 1) * limit

	// Получаем общее количество записей для расчета страниц
	_, total, err := h.store.GetPlaces(1, 0) // Делаем запрос для получения только общего числа записей
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения данных: %s", err), http.StatusInternalServerError)
		return
	}

	// Вычисляем последнюю допустимую страницу
	lastPage := (total + limit - 1) / limit
	if page > lastPage {
		http.Error(w, fmt.Sprintf("Некорректное значение 'page': '%s'", pageStr), http.StatusBadRequest)
		return
	}

	// Теперь получаем данные для текущей страницы
	places, _, err := h.store.GetPlaces(limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения данных: %s", err), http.StatusInternalServerError)
		return
	}

	// Формируем HTML-ответ
	data := struct {
		Title       string
		Places      []types.Place
		Total       int
		CurrentPage int
		NextPage    int
		PrevPage    int
		LastPage    int
	}{
		Title:       "Список ресторанов",
		Places:      places,
		Total:       total,
		CurrentPage: page,
		NextPage:    page + 1,
		PrevPage:    page - 1,
		LastPage:    lastPage,
	}

	if page == lastPage {
		data.NextPage = 0
	}

	tmpl, err := template.ParseFiles("../internal/data/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка загрузки шаблона: %s", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка рендеринга: %s", err), http.StatusInternalServerError)
	}
}

func (h *Handler) JSONHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр "page" из запроса
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, fmt.Sprintf("Некорректное значение 'page': '%s'", pageStr), http.StatusBadRequest)
		return
	}

	// Количество записей на одной странице
	limit := 10
	offset := (page - 1) * limit

	// Получаем общее количество записей для расчета страниц
	_, total, err := h.store.GetPlaces(1, 0) // Делаем запрос для получения только общего числа записей
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения данных: %s", err), http.StatusInternalServerError)
		return
	}

	// Вычисляем последнюю допустимую страницу
	lastPage := (total + limit - 1) / limit
	if page > lastPage {
		http.Error(w, fmt.Sprintf("Некорректное значение 'page': '%s'", pageStr), http.StatusBadRequest)
		return
	}

	// Теперь получаем данные для текущей страницы
	places, _, err := h.store.GetPlaces(limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения данных: %s", err), http.StatusInternalServerError)
		return
	}

	// Формируем JSON-ответ
	data := struct {
		Name     string        `json:"name"`
		Total    int           `json:"total"`
		Places   []types.Place `json:"places"`
		PrevPage int           `json:"prev_page"`
		NextPage int           `json:"next_page"`
		LastPage int           `json:"last_page"`
	}{
		Name:     "Places",
		Total:    total,
		Places:   places,
		LastPage: (total + limit - 1) / limit,
	}

	if page > 1 {
		data.PrevPage = page - 1
	}

	if page < data.LastPage {
		data.NextPage = page + 1
	}

	// Устанавливаем заголовок Content-Type для ответа в формате JSON
	w.Header().Set("Content-Type", "application/json")

	// Отправляем JSON-ответ
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Ошибка при рендеринге данных: %s"}`, err), http.StatusInternalServerError)
	}
}

func (h *Handler) RecommendHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметр "page" из запроса
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	if latStr == "" || lonStr == "" {
		http.Error(w, fmt.Sprintf("Параметры lat и lon обязательны"), http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректное значение для 'lat': %s", latStr), http.StatusBadRequest)
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректное значение для 'lat': %s", lonStr), http.StatusBadRequest)
	}

	//fmt.Fprintf(w, "Получены координаты: lat = %f, lon = %f", lat, lon)

	// Формируем запрос для поиска ближайших мест
	query := fmt.Sprintf(`{
	  "size": 3,
	  "sort": [
		{
		  "_geo_distance": {
			"location": {
			  "lat": %f,
			  "lon": %f
			},
			"order": "asc",
			"unit": "km",
			"mode": "min",
			"distance_type": "arc",
			"ignore_unmapped": true
		  }
		}
	  ]
	}`, lat, lon)

	// Выполняем запрос к Elasticsearch
	res, err := h.store.Search(query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка поиска: %s", err), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		http.Error(w, fmt.Sprintf("Ошибка в ответе от Elasticsearch: %s", res.String()), http.StatusInternalServerError)
		return
	}

	// Обрабатываем ответ от Elasticsearch
	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка декодирования ответа: %s", err), http.StatusInternalServerError)
		return
	}

	// Извлекаем результаты из ответа
	hitsBlock, ok := response["hits"].(map[string]interface{})
	if !ok {
		http.Error(w, "Ошибка извлечения блока 'hits' из ответа", http.StatusInternalServerError)
		return
	}

	hitsArray, ok := hitsBlock["hits"].([]interface{})
	if !ok {
		http.Error(w, "Ошибка извлечения блока 'hits' из блока 'hits'", http.StatusInternalServerError)
		return
	}

	// Формируем список мест
	var places []types.Place
	for _, hit := range hitsArray {
		var place types.Place
		source := hit.(map[string]interface{})["_source"]
		sourceBytes, _ := json.Marshal(source)
		if err := json.Unmarshal(sourceBytes, &place); err != nil {
			continue
		}
		places = append(places, place)
	}

	// Формируем JSON-ответ
	data := struct {
		Name   string        `json:"name"`
		Places []types.Place `json:"places"`
	}{
		Name:   "Recommendation",
		Places: places,
	}

	// Устанавливаем заголовок Content-Type и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при рендеринге данных: %s", err), http.StatusInternalServerError)
	}
}

func (h *Handler) GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin": true,
		"name":  "User",
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	//Подписание токена
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"token": "` + tokenString + `"}`))
}
