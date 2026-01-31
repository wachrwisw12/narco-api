package config

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func loadProd() *Config {
	privKeyBytes, err := os.ReadFile("/run/keys/private.pem")
	if err != nil {
		log.Fatal("cannot read private key:", err)
	}

	pubKeyBytes, err := os.ReadFile("/run/keys/public.pem")
	if err != nil {
		log.Fatal("cannot read public key:", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		log.Fatal("invalid private key:", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		log.Fatal("invalid public key:", err)
	}
	// dbConfig, err := os.ReadFile(".env.api")
	// Load DB config from environment variables

	// ตรวจสอบว่ามีค่าว่างหรือไม่

	return &Config{
		AppEnv:     "dev",
		JWTPrivKey: privateKey,
		JWTPubKey:  publicKey,
	}
}
