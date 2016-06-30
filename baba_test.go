package main

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestPredictWeatherForSuccess(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	if msg := predictWeatherFor("adl", 0.0, 0.0, 0.0, 00); msg != "ADL|0.00,0.00,0|Rain|+2242.263801|-47684.6|-11678" {
		t.Errorf("Expected ADL|0.00,0.00,0|Rain|+2242.263801|-47684.6|-11678, but received %s", msg)
	}
}

func TestPredictItemSuccess(t *testing.T) {
	data := [][]string{
		[]string{"EQT", "0000000", "0.00", "0.00", "00", "00", "0.0", "0.0", "0.0", "0.0"},
		[]string{"EQS", "0000001", "0.00", "0.00", "00", "00", "0.0", "0.0", "0.0", "0.0"},
		[]string{"EQS", "0000002", "0.00", "0.00", "00", "00", "0.0", "0.0", "0.0", "0.0"},
	}

	sample := make(chan float64)
	go predictItem(data, 0.0, 0.00, 0, 3.0, sample, "Sample", 9)

	if item := <-sample; item != 0 {
		t.Errorf("Expected 0, but received %f", item)
	}
}

func TestWeatherConditionSnow(t *testing.T) {
	if weatherCondition(0, 0) != "Snow" {
		t.FailNow()
	}
}

func TestWeatherConditionRain(t *testing.T) {
	if weatherCondition(0.1, 1) != "Rain" {
		t.FailNow()
	}
}

func TestWeatherConditionCloudy(t *testing.T) {
	if weatherCondition(0.1, 0) != "Cloudy" {
		t.FailNow()
	}
}

func TestWeatherConditionSunny(t *testing.T) {
	if weatherCondition(10.1, 0) != "Sunny" {
		t.FailNow()
	}
}
