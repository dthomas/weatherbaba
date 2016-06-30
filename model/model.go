package model

import "strconv"

// Weather Information
type Weather struct {
	IATA       string
	Date       float64
	Latitude   float64
	Longitude  float64
	Elevation  float64
	Gust       float64
	Temprature float64
	Humidity   float64
	Pressure   float64
	Rainfall   float64
}

// Parse string to struct assuming all data is formatted corectly. Not checking for conversion errors
func (w *Weather) Parse(data []string) {
	w.IATA = data[0]
	w.Date, _ = strconv.ParseFloat(data[1], 64)
	w.Latitude, _ = strconv.ParseFloat(data[2], 64)
	w.Longitude, _ = strconv.ParseFloat(data[3], 64)
	w.Elevation, _ = strconv.ParseFloat(data[4], 64)
	w.Gust, _ = strconv.ParseFloat(data[5], 64)
	w.Temprature, _ = strconv.ParseFloat(data[6], 64)
	w.Humidity, _ = strconv.ParseFloat(data[7], 64)
	w.Pressure, _ = strconv.ParseFloat(data[8], 64)
	w.Rainfall, _ = strconv.ParseFloat(data[9], 64)
}
