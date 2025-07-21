package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
