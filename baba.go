package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dthomas/weatherbaba/model"
	"github.com/sajari/regression"
)

func main() {
	locations := []string{"adl", "mel", "syd"}
	lats := []float64{-34.71, -37.86, -33.86}
	lons := []float64{138.62, 144.76, 151.21}
	elvn := []float64{17, 20, 39}
	dats := []float64{1467243000, 1467313200, 1467313200}
	log.SetOutput(ioutil.Discard)
	for idx, loc := range locations {
		go fmt.Println(predictWeatherFor(loc, lats[idx], lons[idx], elvn[idx], dats[idx]))
	}
	var dummyInput string
	fmt.Scanln(&dummyInput)
}

func predictWeatherFor(loc string, lat, lon, eln, dttm float64) string {
	name := "./data/" + loc + ".csv"
	log.Println("Opening", name)
	// read file and scan all the data for computation
	csvfile, err := os.Open(name)
	defer csvfile.Close()

	if err != nil {
		log.Fatalln(err)
		return "ERR"
	}

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1

	data, err := reader.ReadAll()

	if err != nil {
		log.Fatalln(err)
		return "ERR"
	}

	temp := make(chan float64)
	humd := make(chan float64)
	pres := make(chan float64)
	rain := make(chan float64)

	go predictItem(data, lat, lon, eln, dttm, temp, "Temprature", 6)
	go predictItem(data, lat, lon, eln, dttm, humd, "Humidity", 7)
	go predictItem(data, lat, lon, eln, dttm, pres, "Pressure", 8)
	go predictItem(data, lat, lon, eln, dttm, rain, "Rain", 9)

	t, r, h, p := <-temp, <-rain, <-humd, <-pres

	cond := weatherCondition(t, r)

	return fmt.Sprintf("%s|%.2f,%.2f,%.0f|%s|%+f|%.1f|%.0f", strings.ToUpper(loc), lat, lon, eln, cond, t, p, h)
}

func predictItem(data [][]string, lat, lon, eln, dttm float64, item chan float64, desc string, idx int64) {
	log.Printf("Predicting %s at {%f,%f}\n", desc, lat, lon)
	r := new(regression.Regression)
	r.SetObserved(desc)
	r.SetVar(0, "Tim")
	r.SetVar(1, "Lat")
	r.SetVar(2, "Lon")
	r.SetVar(3, "Eln")

	for _, row := range data {
		var w model.Weather
		w.Parse(row)
		val, _ := strconv.ParseFloat(row[idx], 64)
		r.Train(regression.DataPoint(val, []float64{w.Date, w.Latitude, w.Longitude, w.Elevation}))
	}

	r.Run()

	// fmt.Printf("Regression:\n%s\n", r)
	// fmt.Printf("Regression formula:\n%v\n", r.Formula)

	pred, err := r.Predict([]float64{dttm, lat, lon, eln})

	if err != nil {
		log.Println(err)
		item <- -9999999.999
		return
	}

	item <- pred
}

func weatherCondition(temp, rain float64) string {
	if temp <= 0.00 {
		return "Snow"
	} else if rain > 0.00 {
		return "Rain"
	} else if temp > 10.00 {
		return "Sunny"
	}
	return "Cloudy"
}
