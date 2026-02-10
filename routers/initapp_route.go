package routers

import (
	handlers "api-naco/handlers"
	middlewares "api-naco/midleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInnitApp(InnitApp fiber.Router) {
	InnitApp.Post(
		"/search-village",
		middlewares.ReportLimiter(), // 🔒 กัน spam report
		handlers.SearchVillage,
	)
}
