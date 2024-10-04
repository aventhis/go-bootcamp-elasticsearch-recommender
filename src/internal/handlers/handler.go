package handlers

import (
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/db"
	"net/http"
)

type Handler struct {
	store db.Store
}

func NewHandler(store db.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
