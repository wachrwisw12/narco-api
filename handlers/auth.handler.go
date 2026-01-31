package handlers

import (
	"api-naco/config"
	"api-naco/models"
	"api-naco/services"

	"github.com/gofiber/fiber/v2"
)

func Registerhandler(c *fiber.Ctx) error {
	var body models.User
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request json body", "nil": body})
	}
	user, err := services.AuthRegisterService(body)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"body": user})
}

func Authhandler(c *fiber.Ctx) error {
	var body models.AuthRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request json body"})
	}

	result, err := services.AuthLoginService(config.Cfg, body)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
