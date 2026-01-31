package config

import (
	"log"
	"os"
)

func Load() {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "dev"
	}

	switch env {
	case "prod":
		Cfg = loadProd()
	default:
		Cfg = loadDev()
	}

	log.Printf("🚀 running in %s mode\n", env)
}
