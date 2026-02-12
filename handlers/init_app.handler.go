package handlers

import (
	"log"
	"os"

	"api-naco/models"
	"api-naco/services"

	"github.com/gofiber/fiber/v2"
)

func SearchVillage(c *fiber.Ctx) error {
	provinceID := os.Getenv("PROVICE")
	println("provicer", provinceID)
	q := c.Query("q") // รับ query param

	if q == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "q is required",
		})
	}

	// ตัวอย่าง log
	log.Println("search:", q)

	// TODO: query DB
	result, err := services.DistrictsService(q, provinceID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error",
		})
	}
	if result == nil {
		result = []models.District{}
	}

	return c.JSON(result)
}
