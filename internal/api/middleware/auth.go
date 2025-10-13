package middleware

import (
	"spl-notification/internal/config"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	config *config.EnvironmentConfig
}

func NewAuthMiddleware(config *config.EnvironmentConfig) *AuthMiddleware {
	return &AuthMiddleware{config: config}
}

func (u *AuthMiddleware) ValidateAuthHeader(c *fiber.Ctx) error {
	headerValue := c.Get("X-Auth-Token")

	if headerValue == "" || u.config.AuthString != headerValue {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}
