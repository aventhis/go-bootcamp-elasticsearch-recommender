package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/db"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/types"
	"html/template"
	"net/http"
	"strconv"
)

type Handler struct {
	store db.Store
}

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
		http.Error(w, fmt.Sprintf("Параметры lat и lon обязательны", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr,64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректное значение для 'lat': %s", latStr), http.StatusBadRequest)
	}

	lon, err := strconv.ParseFloat(lonStr,64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Некорректное значение для 'lat': %s", lonStr), http.StatusBadRequest)
	}

	// Здесь вы можете использовать полученные значения lat и lon для выполнения логики рекомендации
	fmt.Fprintf(w, "Получены координаты: lat = %f, lon = %f", lat, lon)

	//// Количество записей на одной странице
	//limit := 10
	//offset := (page - 1) * limit
	//
	//// Получаем общее количество записей для расчета страниц
	//_, total, err := h.store.GetPlaces(1, 0) // Делаем запрос для получения только общего числа записей
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Ошибка получения данных: %s", err), http.StatusInternalServerError)
	//	return
	//}
	//
	//// Вычисляем последнюю допустимую страницу
	//lastPage := (total + limit - 1) / limit
	//if page > lastPage {
	//	http.Error(w, fmt.Sprintf("Некорректное значение 'page': '%s'", pageStr), http.StatusBadRequest)
	//	return
	//}
	//
	//// Теперь получаем данные для текущей страницы
	//places, _, err := h.store.GetPlaces(limit, offset)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf("Ошибка получения данных: %s", err), http.StatusInternalServerError)
	//	return
	//}
	//
	//// Формируем JSON-ответ
	//data := struct {
	//	Name     string        `json:"name"`
	//	Total    int           `json:"total"`
	//	Places   []types.Place `json:"places"`
	//	PrevPage int           `json:"prev_page"`
	//	NextPage int           `json:"next_page"`
	//	LastPage int           `json:"last_page"`
	//}{
	//	Name:     "Places",
	//	Total:    total,
	//	Places:   places,
	//	LastPage: (total + limit - 1) / limit,
	//}
	//
	//if page > 1 {
	//	data.PrevPage = page - 1
	//}
	//
	//if page < data.LastPage {
	//	data.NextPage = page + 1
	//}
	//
	//// Устанавливаем заголовок Content-Type для ответа в формате JSON
	//w.Header().Set("Content-Type", "application/json")
	//
	//// Отправляем JSON-ответ
	//err = json.NewEncoder(w).Encode(data)
	//if err != nil {
	//	http.Error(w, fmt.Sprintf(`{"error": "Ошибка при рендеринге данных: %s"}`, err), http.StatusInternalServerError)
	//}
}
