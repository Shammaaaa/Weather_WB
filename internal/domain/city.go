package domain

import "fmt"

type City struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

func (c *City) LatitudeStr() string {
	return fmt.Sprintf("%f", c.Lat)
}

func (c *City) LongitudeStr() string {
	return fmt.Sprintf("%f", c.Lon)
}
