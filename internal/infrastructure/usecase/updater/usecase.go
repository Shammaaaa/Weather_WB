package updater

import (
	"context"
	"sync"
	"time"

	"github.com/shamil/weather/config"
	"github.com/shamil/weather/internal/domain"
	"github.com/shamil/weather/pkg/log"
	"github.com/shamil/weather/pkg/openweathermap"
)

type WeatherRepository interface {
	CitySave(ctx context.Context, cities ...domain.City) error
	WeatherSave(ctx context.Context, weathers ...domain.WeatherForecast) error
}

type UseCase struct {
	weatherRepository WeatherRepository
	token             string
}

func NewUpdaterUseCase(weatherRepository WeatherRepository, token string) *UseCase {
	return &UseCase{
		weatherRepository: weatherRepository,
		token:             token,
	}
}

func (u *UseCase) Work(ctx context.Context) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	if err := u.doWork(ctx); err != nil {
		log.Warningf("Worker-Serice: doWork: %s", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := u.doWork(ctx); err != nil {
				log.Warningf("Worker-Serice: doWork: %s", err)
			}
		}
	}
}

func (u *UseCase) doWork(ctx context.Context) error {
	// const dtLayout = "2006-01-02 15:04:05"

	// insert or update cities on database
	if err := u.weatherRepository.CitySave(ctx, config.Cities...); err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	for _, city := range config.Cities {
		wg.Add(1)

		// async requests
		go func(city domain.City) {
			defer wg.Done()

			forecasts := make([]domain.WeatherForecast, 0, len(config.Cities))

			log.Infof("handle %s weather", city.Name)

			weather, err := openweathermap.Weather(
				ctx, city.LatitudeStr(), city.LongitudeStr(), u.token,
			)
			if err != nil {
				log.Warningf("failed to fetch city weather information: %s", err)
			}

			for _, w := range weather.List {
				forecasts = append(forecasts, domain.WeatherForecast{
					Date:    w.DtTxt,
					Temp:    w.Main.Temp,
					Data:    w,
					City:    city.Name,
					Country: city.Country,
				})
			}

			// update weather information for each city
			if err = u.weatherRepository.WeatherSave(ctx, forecasts...); err != nil {
				log.Warningf("failed to save city weather information: %s", err)
			}
		}(city)
	}

	// wait handle all cities
	wg.Wait()

	log.Info("successfully saves weather information")

	return nil
}
