package main

import (
	"fmt"
	"spl-notification/internal/api/controller"
	"spl-notification/internal/api/middleware"
	"spl-notification/internal/config"
	"spl-notification/internal/database"
	"spl-notification/internal/repository"
	"spl-notification/internal/server"
	"spl-notification/internal/service"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			// Setup environment config
			config.NewEnviromentConfig,
			NewValidator,
			// Database connection
			database.CreateTursoConnection,
			// Middleware
			middleware.NewAuthMiddleware,
			// Controllers
			controller.NewMainController,
			controller.NewTrackController,
			// Services
			fx.Annotate(
				service.NewAccessServiceImpl,
				fx.As(new(service.AccessService)),
			),
			fx.Annotate(
				service.NewNotificationServiceImpl,
				fx.As(new(service.NotificationService)),
			),
			fx.Annotate(
				service.NewTrackServiceImpl,
				fx.As(new(service.TrackService)),
			),
			fx.Annotate(
				service.NewSourceServiceImpl,
				fx.As(new(service.SourceService)),
			),
			// Setup Repositories
			fx.Annotate(
				repository.NewTrackRepositoryImpl,
				fx.As(new(repository.TrackRepository)),
			),
		),
		// Setup Server
		fx.Invoke(server.CreateFiberServer),
		// Start Pub/Sub notification consumer
		fx.Invoke(func(notificationService service.NotificationService) {
			go notificationService.HandleNotification()
		}),
		fx.Invoke(func(accessService service.AccessService) {
			s, err := gocron.NewScheduler()
			if err != nil {
				fmt.Println("Error creating scheduler:", err)
				return
			}

			s.NewJob(
				gocron.DurationJob(5*time.Second),
				gocron.NewTask(func() {
					accesses, err := accessService.GetRecentlyAccess()
					if err != nil {
						fmt.Println("[CRON] Error fetching accesses:", err)
						return
					}

					if len(accesses) == 0 {
						return
					}

					err = accessService.CheckAccess(accesses)
					if err != nil {
						fmt.Println("[CRON] Error checking accesses:", err)
						return
					}
				}),
				gocron.WithSingletonMode(gocron.LimitModeWait),
			)

			s.Start()
		}),
	).Run()
}

func NewValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}
