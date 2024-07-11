package openweathermap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CityResponse struct {
	Name       string            `json:"name"`
	LocalNames map[string]string `json:"local_names"`
	Lat        float64           `json:"lat"`
	Lon        float64           `json:"lon"`
	Country    string            `json:"country"`
	State      string            `json:"state"`
}

func City(ctx context.Context, city, token string) ([]CityResponse, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", city, token)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var response []CityResponse
	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
