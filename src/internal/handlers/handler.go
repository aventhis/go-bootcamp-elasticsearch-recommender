package handlers

import (
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

	// Получаем данные из хранилища
	places, total, err := h.store.GetPlaces(limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка получения данных: %s", err), http.StatusInternalServerError)
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
		LastPage:    (total + limit - 1) / limit,
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
