package data

import (
	"bufio"
	"encoding/csv"
	"errors"
	"github.com/aventhis/go-bootcamp-elasticsearch-recommender/internal/types"
	"io"
	"os"
	"strconv"
)

func LoadPlaces(filePath string) ([]types.Place, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("ошибка при открытии файла: " + err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '\t'

	// Пропуск заголовка
	if _, err = reader.Read(); err != nil {
		return nil, err
	}

	var places []types.Place

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		id, _ := strconv.ParseInt(record[0], 10, 64)
		lat, _ := strconv.ParseFloat(record[5], 64)
		lon, _ := strconv.ParseFloat(record[4], 64)

		place := types.Place{
			ID:      id + 1,
			Name:    record[1],
			Address: record[2],
			Phone:   record[3],
		}
		place.Location.Lat = lat
		place.Location.Lon = lon

		places = append(places, place)
	}

	return places, nil
}
