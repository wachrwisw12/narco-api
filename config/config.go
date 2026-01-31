package config

import (
	"crypto/rsa"

	"api-naco/models"
)

type Config struct {
	AppEnv     string
	JWTPrivKey *rsa.PrivateKey
	JWTPubKey  *rsa.PublicKey
	DB         models.DBConfig
}

var Cfg *Config
