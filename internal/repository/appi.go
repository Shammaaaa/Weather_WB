package repository

import (
	"context"
	"fmt"
	"github.com/shamil/weather/internal/domain"
)

func (r *Repository) GetCitiesWithWeatherForecasts(ctx context.Context) ([]domain.City, error) {
	query := `
        SELECT c.name, c.country, c.lat, c.lon
        FROM cities c
        INNER JOIN weather_forecasts wf ON c.id = wf.city_id
        GROUP BY c.id
        ORDER BY c.name
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cities := make([]domain.City, 0)
	for rows.Next() {
		var city domain.City
		if err := rows.Scan(&city.Name, &city.Country, &city.Lat, &city.Lon); err != nil {
			return nil, err
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cities, nil
}

func (r *Repository) GetCityWeather(ctx context.Context, cityName string) (domain.CityWeather, error) {
	query := `
        SELECT c.name, c.country, c.lat, c.lon, wf.datetime, wf.temp, wf.data
        FROM cities c
        INNER JOIN weather_forecasts wf ON c.id = wf.city_id
        WHERE c.name = $1
        ORDER BY wf.datetime
    `

	rows, err := r.db.QueryContext(ctx, query, cityName)
	if err != nil {
		return domain.CityWeather{}, err
	}
	defer rows.Close()

	var cityWeather domain.CityWeather
	for rows.Next() {
		var forecast domain.WeatherForecast
		if err := rows.Scan(&cityWeather.Name, &cityWeather.Country, &cityWeather.Lat, &cityWeather.Lon,
			&forecast.Date, &forecast.Temp, &forecast.Data); err != nil {
			return domain.CityWeather{}, err
		}
		cityWeather.Forecasts = append(cityWeather.Forecasts, forecast)
	}

	if err := rows.Err(); err != nil {
		return domain.CityWeather{}, err
	}

	return cityWeather, nil
}

func (r *Repository) GetCityWeatherForecasts(ctx context.Context, cityName, date string) ([]domain.WeatherForecast, error) {
	query := `
        SELECT c.country, c.name, AVG(wf.temp), wf.datetime
        FROM cities c
        INNER JOIN weather_forecasts wf ON c.id = wf.city_id
        WHERE c.name = $1 AND wf.datetime::date = $2
        GROUP BY 1, 2, 4
        ORDER BY wf.datetime
        
    `

	rows, err := r.db.QueryContext(ctx, query, cityName, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	forecasts := make([]domain.WeatherForecast, 0)
	for rows.Next() {
		var forecast domain.WeatherForecast
		if err := rows.Scan(&forecast.Country, &forecast.City, &forecast.Temp, &forecast.Date); err != nil {
			return nil, err
		}
		forecasts = append(forecasts, forecast)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return forecasts, nil
}

func (r *Repository) GetWeatherWithDateTime(ctx context.Context, cityName, date, timestamp string) ([]domain.WeatherForecast, error) {
	withTime := ""
	if timestamp != "" {
		withTime = fmt.Sprintf(" AND wf.datetime::time = '%s'", timestamp)
	}

	query := fmt.Sprintf(`
        SELECT c.name, wf.datetime, wf.data
        FROM cities c
        INNER JOIN weather_forecasts wf ON c.id = wf.city_id
        WHERE c.name = $1 AND wf.datetime::date = $2 %s`, withTime)

	rows, err := r.db.QueryContext(ctx, query, cityName, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cityWeather := make([]domain.WeatherForecast, 0)
	for rows.Next() {
		var forecast domain.WeatherForecast
		if err := rows.Scan(&forecast.City, &forecast.Date, &forecast.Data); err != nil {
			return nil, err
		}
		cityWeather = append(cityWeather, forecast)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cityWeather, nil
}

func (r *Repository) GetFavoriteCities(ctx context.Context, userID int) ([]domain.FavoriteCity, error) {
	query := `
        SELECT fc.id, fc.city_id, c.name, fc.user_id 
        FROM favorite_cities fc
        JOIN cities c ON fc.city_id = c.id
        WHERE fc.user_id = $1
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	favoriteCities := make([]domain.FavoriteCity, 0)
	for rows.Next() {
		var favoriteCity domain.FavoriteCity
		if err := rows.Scan(&favoriteCity.ID, &favoriteCity.CityID, &favoriteCity.CityName, &favoriteCity.UserID); err != nil {
			return nil, err
		}
		favoriteCities = append(favoriteCities, favoriteCity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favoriteCities, nil
}

func (r *Repository) AddFavoriteCity(ctx context.Context, userID, cityID int, cityName string) error {
	query := `
        INSERT INTO favorite_cities (user_id, city_id, city_name)
        VALUES ($1, $2, $3)
    `

	_, err := r.db.Exec(query, userID, cityID, cityName)
	return err
}
