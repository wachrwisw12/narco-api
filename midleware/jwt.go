package middlewares

import (
	"errors"
	"strings"
	"time"

	"api-naco/config"
	"api-naco/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")

		// 🔓 guest route
		if auth == "" {
			if len(requiredRoles) == 0 {
				c.Locals("role", "guest")
				return c.Next()
			}
			return fiber.ErrUnauthorized
		}

		// Bearer token
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.ErrUnauthorized
		}

		tokenStr := parts[1]

		token, err := jwt.ParseWithClaims(
			tokenStr,
			&models.JWTClaims{},
			func(t *jwt.Token) (interface{}, error) {
				if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
					return nil, fiber.ErrUnauthorized
				}
				return config.Cfg.JWTPubKey, nil
			},
		)

		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		claims := token.Claims.(*models.JWTClaims)

		// 🧱 ตรวจ role
		if len(requiredRoles) > 0 {
			allowed := false
			for _, r := range requiredRoles {
				if claims.Role == r {
					allowed = true
					break
				}
			}
			if !allowed {
				return fiber.ErrForbidden
			}
		}

		// แนบ user info
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
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

func GenerateJWT(cfg *config.Config, userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"iss":  "api-naco",
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(cfg.JWTPrivKey)
}
