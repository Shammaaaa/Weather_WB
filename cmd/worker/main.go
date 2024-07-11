package main

import (
	"context"

	"net"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"

	"github.com/shamil/weather/config"
	"github.com/shamil/weather/internal/infrastructure/usecase/updater"
	"github.com/shamil/weather/internal/repository"
	"github.com/shamil/weather/internal/service"
	"github.com/shamil/weather/pkg/log"
	"github.com/shamil/weather/pkg/signal"
)

func main() {
	// это просто штука для получения аргументов из командной строки
	// $ go run main.go --config-file ./config.yml --listener 1
	application := cli.App{
		Name: "Worker-Service",
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
				Value:    "0.0.0.0:3004",
				EnvVars:  []string{"BIND_ADDRESS"},
			},
			&cli.StringFlag{
				Name:     "bind-socket",
				Usage:    "Путь к Unix сокет файлу",
				Required: false,
				Value:    "/tmp/worker_service.sock",
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

	// запускаем нашу cli команду
	if err := application.Run(os.Args); err != nil {
		log.Error(err)
	}

}

// Main функция как main, чтобы разбить main на более мелкие части
func Main(ctx *cli.Context) error {
	appContext, cancel := context.WithCancel(ctx.Context)
	defer func() {
		cancel()
		log.Info("app context is canceled, Worker-Service is down!")
	}()

	// получаем наш конфиг файл
	cfg, err := config.New(ctx.String("config-file"))
	if err != nil {
		return err
	}

	// инициализируем сервис с нужными зависимостями
	srv, err := service.New(appContext, &service.Options{
		Database: &cfg.Database,
	})
	if err != nil {
		return err
	}

	// эта штука отрабатывает при завершении приложения по ctrl+c
	// тут закрываются соединения у всех зависимостей, которые реализуют метод Drop() из pkg/drop
	defer func() {
		srv.Shutdown(func(err error) {
			log.Warning(err)
		})
		srv.Stacktrace()
	}()

	// тут инициализируем ожидание сигнала системы, например, ctrl+c
	await, stop := signal.Notifier(func() {
		log.Info("Worker-Service, start shutdown process..")
	})

	// репозиторий с юз кей
	repo := repository.New(srv.Pool.Builder())
	useCase := updater.NewUpdaterUseCase(repo, cfg.Token)

	// запускаем background worker
	go func() {
		useCase.Work(srv.Context())
	}()

	// тут простой health check, localhost:3004/alive
	go func() {
		app := fiber.New(fiber.Config{
			ServerHeader: "Worker-Service Server",
		})
		app.Get("/alive", func(ctx *fiber.Ctx) error {
			return ctx.SendString("Alive")
		})

		var ln net.Listener
		if ln, err = signal.Listener(
			srv.Context(),
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

	log.Info("Worker-Service is launched")
	return await()

}
