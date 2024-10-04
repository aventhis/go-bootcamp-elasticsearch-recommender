package handlers

import (
	"fmt"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/db"
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
		http.Error(w, fmt.Sprintf("Invalid 'page' value: '%s'", pageStr), http.StatusBadRequest)
		return
	}
}
