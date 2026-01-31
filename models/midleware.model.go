package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"` // guest, officer, admin
	jwt.RegisteredClaims
}
