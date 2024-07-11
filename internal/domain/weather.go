package domain

import (
	"encoding/json"
	"github.com/shamil/weather/pkg/openweathermap"
)

type CityWeather struct {
	Name      string            `json:"name"`
	Country   string            `json:"country"`
	Lat       float64           `json:"lat"`
	Lon       float64           `json:"lon"`
	Forecasts []WeatherForecast `json:"forecasts"`
}

type WeatherForecast struct {
	Date string                     `json:"date"`
	Temp float64                    `json:"temp"`
	Data openweathermap.WeatherList `json:"data"`

	City    string `json:"city"`
	Country string `json:"country"`
}

func (a *WeatherForecast) DataStr() string {
	str, _ := json.Marshal(a.Data)
	return string(str)
}
