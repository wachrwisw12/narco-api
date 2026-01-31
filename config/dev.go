package config

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func loadDev() *Config {
	privKeyBytes, err := os.ReadFile("keys/private.pem")
	if err != nil {
		log.Fatal("cannot read private key:", err)
	}

	pubKeyBytes, err := os.ReadFile("keys/public.pem")
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

	return &Config{
		AppEnv:     "dev",
		JWTPrivKey: privateKey,
		JWTPubKey:  publicKey,
	}
}
