package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func LoginLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: time.Minute,
	})
}

func RegisterLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        3,
		Expiration: time.Minute,
	})
}

func ReportLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,
		Expiration: time.Minute,
	})
}

func TrackLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        30,
		Expiration: time.Minute,
	})
}

func PublicLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,
		Expiration: time.Minute,
	})
}
