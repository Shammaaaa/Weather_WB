package api

import (
	"context"
	"github.com/shamil/weather/internal/domain"
)

type WeatherRepository interface {
	GetCitiesWithWeatherForecasts(ctx context.Context) ([]domain.City, error)
	GetCityWeather(ctx context.Context, cityName string) (domain.CityWeather, error)
	GetCityWeatherForecasts(ctx context.Context, cityName, date string) ([]domain.WeatherForecast, error)
	GetWeatherWithDateTime(ctx context.Context, cityName, date, timestamp string) ([]domain.WeatherForecast, error)
	GetFavoriteCities(ctx context.Context, userID int) ([]domain.FavoriteCity, error)
	AddFavoriteCity(ctx context.Context, userID, cityID int, cityName string) error
}

type UseCase struct {
	WeatherRepository WeatherRepository
}

func NewApiUseCase(weatherRepository WeatherRepository) *UseCase {
	return &UseCase{
		WeatherRepository: weatherRepository,
	}
}

func (u *UseCase) GetCitiesWithWeatherForecasts(ctx context.Context) ([]domain.City, error) {
	return u.WeatherRepository.GetCitiesWithWeatherForecasts(ctx)
}

func (u *UseCase) GetCityWeather(ctx context.Context, cityName string) (domain.CityWeather, error) {
	return u.WeatherRepository.GetCityWeather(ctx, cityName)
}
func (u *UseCase) GetCityWeatherForecasts(ctx context.Context, cityName, date string) ([]domain.WeatherForecast, error) {
	return u.WeatherRepository.GetCityWeatherForecasts(ctx, cityName, date)
}

func (u *UseCase) GetWeatherWithDateTime(ctx context.Context, cityName, date, timestamp string) ([]domain.WeatherForecast, error) {
	return u.WeatherRepository.GetWeatherWithDateTime(ctx, cityName, date, timestamp)
}

func (u *UseCase) GetFavoriteCities(ctx context.Context, userID int) ([]domain.FavoriteCity, error) {
	return u.WeatherRepository.GetFavoriteCities(ctx, userID)
}

func (u *UseCase) AddFavoriteCity(ctx context.Context, userID, cityID int, cityName string) error {
	return u.WeatherRepository.AddFavoriteCity(ctx, userID, cityID, cityName)
}
