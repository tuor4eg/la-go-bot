package main

import (
	"log"

	"la-go-bot/internal/bot"
	"la-go-bot/internal/config"
)

func main() {
	cfg := config.Load()

	telegramBot, err := bot.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := telegramBot.Start(); err != nil {
		log.Fatal(err)
	}
}
