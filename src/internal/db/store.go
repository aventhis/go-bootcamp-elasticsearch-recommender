package db

import "github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/types"

type Store interface {
	// returns a list of items, a total number of hits and (or) an error in case of one
	GetPlaces(limit int, offset int) ([]types.Place, int, error)
}
