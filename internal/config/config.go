package config

import (
	"defig-stats-artifact/internal/logger"
	gdot "github.com/joho/godotenv"
)

func LoadEnv() {
	err := gdot.Load(".env")
	if err != nil {
		logger.ErrorLogger.Fatal("Error loading .env file", err)
	}
}
