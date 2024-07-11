package repository

import (
	"context"

	"github.com/shamil/weather/internal/domain"
	"github.com/shamil/weather/internal/infrastructure/database"
)

func (r *Repository) CitySave(ctx context.Context, cities ...domain.City) error {
	const query = `INSERT INTO cities (name, country, lat, lon) VALUES ($1, $2, $3, $4)
			       ON CONFLICT(name) DO UPDATE SET country = excluded.country, lat = excluded.lat, lon = excluded.lon`

	err := database.WithTransaction(ctx, r.db, func(transaction database.Transaction) error {
		for _, city := range cities {
			_, err := transaction.Exec(query, city.Name, city.Country, city.Lat, city.Lon)
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
