package routers

import "github.com/gofiber/fiber/v2"

func SetupRoute(app *fiber.App) {
	api := app.Group("/api")

	SetupNarcoticsReport(api.Group("/v1"))
	SetupAuth(api.Group("/auth"))
}
