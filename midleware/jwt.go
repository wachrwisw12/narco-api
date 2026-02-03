package middlewares

import (
	"errors"
	"strings"
	"time"

	"api-naco/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(c *fiber.Ctx) error {
	auth := c.Get("Authorization")

	// 1. ต้องมี Authorization header
	if auth == "" {
		return fiber.ErrUnauthorized
	}

	// 2. ต้องเป็น Bearer token
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return fiber.ErrUnauthorized
	}

	tokenStr := parts[1]

	// 3. Parse + Verify token
	token, err := jwt.Parse(
		tokenStr,
		func(t *jwt.Token) (interface{}, error) {
			// 3.1 เช็ค algorithm
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fiber.ErrUnauthorized
			}
			// 3.2 ใช้ Public key
			return config.Cfg.JWTPubKey, nil
		},
	)

	// 4. token ไม่ valid หรือ error (รวมหมดอายุ)
	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Next()
	}
	c.Locals("username", claims["username"])
	// ✅ ผ่าน = token ถูก + ยังไม่หมดอายุ
	return c.Next()
}

func OptionalJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")

		// ไม่มี token → guest
		if auth == "" {
			return c.Next()
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return config.Cfg.JWTPubKey, nil
		})

		if err != nil || !token.Valid {
			// token พัง → ถือเป็น guest (ไม่ throw)
			return c.Next()
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Next()
		}

		// 🔑 set context
		c.Locals("user_id", claims["sub"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

func GenerateJWT(cfg *config.Config, userID int, roleID int, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub":      userID,
		"role":     roleID,
		"iss":      "api-naco",
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(cfg.JWTPrivKey)
}
