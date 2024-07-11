package http

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/shamil/weather/internal/domain"
	"strconv"
)

type WeatherUsecase interface {
	GetCitiesWithWeatherForecasts(ctx context.Context) ([]domain.City, error)
	GetCityWeather(ctx context.Context, cityName string) (domain.CityWeather, error)
	GetCityWeatherForecasts(ctx context.Context, cityName, date string) ([]domain.WeatherForecast, error)
	GetWeatherWithDateTime(ctx context.Context, cityName, date, timestamp string) ([]domain.WeatherForecast, error)
	GetFavoriteCities(ctx context.Context, userID int) ([]domain.FavoriteCity, error)
	AddFavoriteCity(ctx context.Context, userID, cityID int, cityName string) error
}

type HandlerImpl struct {
	weatherUsecase WeatherUsecase
}

func (h *HandlerImpl) GetCitiesWithWeatherForecasts(ctx *fiber.Ctx) error {
	cities, err := h.weatherUsecase.GetCitiesWithWeatherForecasts(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(cities)
}

func (h *HandlerImpl) GetCityWeather(ctx *fiber.Ctx) error {
	cityName := ctx.Params("city")
	cityWeather, err := h.weatherUsecase.GetCityWeather(ctx.Context(), cityName)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(cityWeather)
}

func (h *HandlerImpl) GetCityWeatherForecasts(ctx *fiber.Ctx) error {
	cityName := ctx.Params("city")
	date := ctx.Params("date")
	forecasts, err := h.weatherUsecase.GetCityWeatherForecasts(ctx.Context(), cityName, date)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(forecasts)
}

func (h *HandlerImpl) GetWeatherWithDateTime(ctx *fiber.Ctx) error {
	cityName := ctx.Params("city")
	date := ctx.Params("date")
	timestamp := ctx.Params("time")

	cityWeather, err := h.weatherUsecase.GetWeatherWithDateTime(ctx.Context(), cityName, date, timestamp)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(cityWeather)
}

func (h *HandlerImpl) GetFavoriteCities(ctx *fiber.Ctx) error {
	userID, _ := strconv.Atoi(ctx.Query("user_id"))
	favoriteCities, err := h.weatherUsecase.GetFavoriteCities(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(favoriteCities)
}

func (h *HandlerImpl) AddFavoriteCity(ctx *fiber.Ctx) error {
	userID, _ := strconv.Atoi(ctx.Query("user_id"))
	cityID, _ := strconv.Atoi(ctx.Query("city_id"))
	cityName := ctx.Query("city_name")

	if err := h.weatherUsecase.AddFavoriteCity(ctx.Context(), userID, cityID, cityName); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "City added to favorites",
	})
}

func New(useCase WeatherUsecase) *HandlerImpl {
	return &HandlerImpl{weatherUsecase: useCase}
}
func (h *HandlerImpl) MountRoutes(app *fiber.App) {
	app.Get("/cities", h.GetCitiesWithWeatherForecasts)
	app.Get("/weather/:city", h.GetCityWeather)
	app.Get("/weather/forecast/:city/:date", h.GetCityWeatherForecasts)
	app.Get("/weather/:city/:date/:time", h.GetWeatherWithDateTime)
	app.Get("/favorites", h.GetFavoriteCities)
	app.Post("/favorites", h.AddFavoriteCity)
}
