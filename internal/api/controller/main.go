package controller

import "github.com/gofiber/fiber/v2"

type MainController struct {
}

func NewMainController() *MainController {
	return &MainController{}
}

func (a *MainController) Health(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
