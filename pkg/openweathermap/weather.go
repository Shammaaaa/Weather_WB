package openweathermap

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (w WeatherList) Value() (driver.Value, error) {
	return json.Marshal(w)
}

func (w *WeatherList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &w)
}

type WeatherResponse struct {
	Cod     string        `json:"cod"`
	Message int           `json:"message"`
	Cnt     int           `json:"cnt"`
	List    []WeatherList `json:"list"`
	City    struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}

type WeatherList struct {
	Dt   int `json:"dt"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
		Humidity  int     `json:"humidity"`
		TempKf    float64 `json:"temp_kf"`
	} `json:"main"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Visibility int     `json:"visibility"`
	Pop        float64 `json:"pop"`
	Sys        struct {
		Pod string `json:"pod"`
	} `json:"sys"`
	DtTxt string `json:"dt_txt"`
	Rain  struct {
		ThreeH float64 `json:"3h"`
	} `json:"rain,omitempty"`
}

func Weather(ctx context.Context, lat, lon, token string) (WeatherResponse, error) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/forecast?lat=%s&lon=%s&appid=%s&lang=ru&units=metric",
		lat, lon, token,
	)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return WeatherResponse{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return WeatherResponse{}, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var response WeatherResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return WeatherResponse{}, err
	}
	return response, nil
}
