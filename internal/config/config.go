package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	ApiSecretKey  string
	ApiBaseURL    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	apiSecretKey := os.Getenv("API_SECRET_KEY")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	apiBaseURL := os.Getenv("API_BASE_URL")
	if apiBaseURL == "" {
		log.Fatal("API_BASE_URL is not set")
	}

	return &Config{
		TelegramToken: token,
		ApiSecretKey:  apiSecretKey,
		ApiBaseURL:    apiBaseURL,
	}
}
