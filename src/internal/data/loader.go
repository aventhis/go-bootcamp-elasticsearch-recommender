package data

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
)

func LoadRestaurant(filePath string) ([]Restaurant, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("ошибка при открытии файла: " + err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '\t' // разделитель поля - табуляция

	// Пропуск заголовка
	if _, err = reader.Read(); err != nil {
		return nil, err
	}

	var restaurants []Restaurant

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		id, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			// Пропускаем строки с неверным ID
			continue
		}
		lat, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			continue
		}
		lon, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			continue
		}

		restaurant := Restaurant{
			ID:      id,
			Name:    record[1],
			Address: record[2],
			Phone:   record[3],
		}
		restaurant.Location.Lat = lat
		restaurant.Location.Lon = lon

		restaurants = append(restaurants, restaurant)

	}

	return restaurants, nil

}
