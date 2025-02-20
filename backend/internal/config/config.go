package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	QPAppKey    string
	QPAppSecret string
	QPAuthBase  string
	QPAPIBase   string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:        "9090",
		QPAppKey:    os.Getenv("QP_APP_KEY"),
		QPAppSecret: os.Getenv("QP_APP_SECRET"),
		QPAuthBase:  os.Getenv("QP_AUTH_API_BASE"),
		QPAPIBase:   os.Getenv("QP_API_BASE"),
	}, nil
}
