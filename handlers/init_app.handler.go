package handlers

import (
	"api-naco/services"

	"github.com/gofiber/fiber/v2"
)

func AppInit(c *fiber.Ctx) error {
	role := "guest"
	// userID := ""

	if r := c.Locals("role"); r != nil {
		role = r.(string)
	}

	// if uid := c.Locals("user_id"); uid != nil {
	// 	userID = uid.(string)
	// }

	menus, err := services.GetMenusByRole(role)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"role":  role,
		"menus": menus,
	})
}
