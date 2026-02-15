package routers

import (
	handlers "api-naco/handlers"
	middlewares "api-naco/midleware"

	"github.com/gofiber/fiber/v2"
)

func SetupNarcoticsReport(nacoticRoute fiber.Router) {
	nacoticRoute.Post(
		"/sendreport",
		// middlewares.ReportLimiter(), // 🔒 กัน spam report
		handlers.SendReport,
	)

	nacoticRoute.Get(
		"/reportInit",
		middlewares.PublicLimiter(),
		handlers.ReportInit,
	)

	nacoticRoute.Get(
		"/case-reports",
		// middlewares.PublicLimiter(),
		middlewares.JWTMiddleware,
		handlers.ListReports,
	)
	nacoticRoute.Post("/update-status", middlewares.JWTMiddleware, handlers.UpdateStatusReport)
	nacoticRoute.Post(
		"/track",
		middlewares.TrackLimiter(), // 🔒 brute force tracking
		handlers.TrackReport,
	)

	// nacoticRoute.Get(
	// 	"/app-init",
	// 	middlewares.OptionalJWT(),
	// 	middlewares.PublicLimiter(),
	// 	handlers.AppInit,
	// )

	nacoticRoute.Get("/test", handlers.Test)
}

func SetupAuth(auth fiber.Router) {
	auth.Post(
		"/singin",
		middlewares.LoginLimiter(), // 🔒 brute force login
		handlers.Authhandler,
	)

	auth.Post(
		"/register",
		middlewares.RegisterLimiter(), // 🔒 spam account
		handlers.Registerhandler,
	)

	auth.Get(
		"/me",
		middlewares.JWTMiddleware,
		handlers.Me,
	)
}
