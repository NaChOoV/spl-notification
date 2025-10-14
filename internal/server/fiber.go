package server

import (
	"context"
	"spl-notification/internal/api/controller"
	"spl-notification/internal/api/middleware"
	"spl-notification/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/fx"
)

func CreateFiberServer(
	lc fx.Lifecycle,
	mainController *controller.MainController,
	trackController *controller.TrackController,
	authMiddleware *middleware.AuthMiddleware,
	config *config.EnvironmentConfig,
) {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	// Setup routes
	app.Get("/health", mainController.Health)
	// Track
	app.Get("/track/:chatId", authMiddleware.ValidateAuthHeader, trackController.GetAllFollowTracks)
	app.Get("/track/send/:chatId", authMiddleware.ValidateAuthHeader, trackController.SendAllFollowTracks)
	app.Post("/track", authMiddleware.ValidateAuthHeader, trackController.CreateTrack)
	app.Delete("/track", authMiddleware.ValidateAuthHeader, trackController.DeleteTrack)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			port := config.Port
			go app.Listen(":" + port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}
