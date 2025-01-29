package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HBUsername string
	HBPassword string
	Port       string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		HBUsername: os.Getenv("HB_USERNAME"),
		HBPassword: os.Getenv("HB_PASSWORD"),
		Port:       "9090",
	}, nil
}
