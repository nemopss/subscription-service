package config

import (
	"github.com/joho/godotenv"

	"github.com/nemopss/subscription-service/pkg/logger"
)

func LoadConfig(path string) {
	err := godotenv.Load(path)
	if err != nil {
		logger.Log.Warn("Error loading .env file")
	}
}
