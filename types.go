package main

import (
	"strconv"
	"time"
)

type WeatherRecord struct {
	ID string `json:"id"`
	Date          time.Time `json:"date"`
	Precipitation float64 `json:"precipitation"`
	TempMax       float64 `json:"temp_max"`
	TempMin       float64 `json:"temp_min"`
	Wind          float64 `json:"wind"`
	Weather       string `json:"weather"`
}

func NewWeatherRecord(data []string) *WeatherRecord {
	date, _ := time.Parse("2006-01-02", data[0])
	precipitation, _ := strconv.ParseFloat(data[1], 64)
	tempMax, _ := strconv.ParseFloat(data[2], 64)
	tempMin, _ := strconv.ParseFloat(data[3], 64)
	wind, _ := strconv.ParseFloat(data[4], 64)
	return &WeatherRecord{
		ID: data[0],
		Date:          date,
		Precipitation: precipitation,
		TempMax:       tempMax,
		TempMin:       tempMin,
		Wind:          wind,
		Weather:       data[5],
	}
}

type UploadResponse struct {
	InsertTime string				`json:"file_db_insert_time"`
	Result []map[string]interface{} `json:"result"`
	Status string                   `json:"status"`
	Time   string                   `json:"count_query_time"`
}
