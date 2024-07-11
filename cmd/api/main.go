package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/shamil/weather/internal/infrastructure/http"
	"github.com/shamil/weather/internal/infrastructure/usecase/api"
	"github.com/shamil/weather/internal/repository"
	"net"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/shamil/weather/config"
	"github.com/shamil/weather/internal/service"
	"github.com/shamil/weather/pkg/log"
	"github.com/shamil/weather/pkg/signal"
)

func main() {
	application := cli.App{
		Name: "Api-Service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config-file",
				Required: true,
				Usage:    "YAML config filepath",
				EnvVars:  []string{"CONFIG_FILE"},
				FilePath: "/srv/secret/config_file",
			},
			&cli.StringFlag{
				Name:     "bind-address",
				Usage:    "IP и порт сервера, например: 0.0.0.0:3001",
				Required: false,
				Value:    "0.0.0.0:3003",
				EnvVars:  []string{"BIND_ADDRESS"},
			},
			&cli.StringFlag{
				Name:     "bind-socket",
				Usage:    "Путь к Unix сокет файлу",
				Required: false,
				Value:    "/tmp/api_service.sock",
				EnvVars:  []string{"BIND_SOCKET"},
			},
			&cli.IntFlag{
				Name:     "listener",
				Usage:    "Unix socket or TCP",
				Required: false,
				Value:    1,
				EnvVars:  []string{"LISTENER"},
			},
		},
		Action: Main,
		After: func(c *cli.Context) error {
			log.Info("stopped")
			return nil
		},
	}

	if err := application.Run(os.Args); err != nil {
		log.Error(err)
	}

}

func Main(ctx *cli.Context) error {
	appContext, cancel := context.WithCancel(ctx.Context)
	defer func() {
		cancel()
		log.Info("app context is canceled, Api-Service is down!")
	}()

	cfg, err := config.New(ctx.String("config-file"))
	if err != nil {
		return err
	}

	apis, err := service.New(appContext, &service.Options{
		Database: &cfg.Database,
	})
	if err != nil {
		return err
	}

	defer func() {
		apis.Shutdown(func(err error) {
			log.Warning(err)
		})
		apis.Stacktrace()
	}()

	await, stop := signal.Notifier(func() {
		log.Info("Api-Service, start shutdown process..")
	})

	repo := repository.New(apis.Pool.Builder())
	useCase := api.NewApiUseCase(repo)
	handler := http.New(useCase)

	// repository := repo.New(apis.Pool)
	app := fiber.New(fiber.Config{
		ServerHeader: "Worker-Service Server",
	})
	app.Get("/alive", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Alive")
	})
	handler.MountRoutes(app)

	go func() {
		var ln net.Listener
		if ln, err = signal.Listener(
			apis.Context(),
			ctx.Int("listener"),
			ctx.String("bind-socket"),
			ctx.String("bind-address"),
		); err != nil {
			stop(err)
			return
		}
		if err = app.Listener(ln); err != nil {
			stop(err)
		}
	}()

	log.Info("Api-Service is launched")
	return await()

}
