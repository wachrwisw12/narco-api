package routers

import (
	handlers "api-naco/handlers"
	middlewares "api-naco/midleware"

	"github.com/gofiber/fiber/v2"
)

func SetupNarcoticsReport(nacoticRoute fiber.Router) {
	nacoticRoute.Post("/sendreport", handlers.SendReport)
	nacoticRoute.Get("/reportInit", handlers.ReportInit)
	nacoticRoute.Get("/reports", handlers.ListReports)
	nacoticRoute.Post("/track", handlers.TrackReport)

	nacoticRoute.Get("/app-init", middlewares.OptionalJWT(), handlers.AppInit)

	nacoticRoute.Get("/test", handlers.Test)
}

func SetupAuth(auth fiber.Router) {
	auth.Post("/singin", handlers.Authhandler)
	auth.Post("/register", handlers.Registerhandler)
}
