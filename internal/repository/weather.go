package repository

import (
	"context"
	"github.com/shamil/weather/pkg/log"

	"github.com/shamil/weather/internal/domain"
	"github.com/shamil/weather/internal/infrastructure/database"
)

func (r *Repository) WeatherSave(ctx context.Context, weathers ...domain.WeatherForecast) error {
	const query = `INSERT INTO weather_forecasts (city_id, datetime, temp, data) 
				   VALUES ((SELECT id FROM cities WHERE name = $1 AND country = $2), $3, $4, $5)
				   ON CONFLICT(city_id, datetime) DO UPDATE SET temp = excluded.temp, data = excluded.data`

	err := database.WithTransaction(ctx, r.db, func(transaction database.Transaction) error {
		for _, w := range weathers {
			log.Infof("save city %s weather", w.City)

			_, err := transaction.Exec(query, w.City, w.Country, w.Date, w.Temp, w.DataStr())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
