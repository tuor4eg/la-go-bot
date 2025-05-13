package bot

import (
	"encoding/json"
	"fmt"
	"log"

	"la-go-bot/internal/config"
	"la-go-bot/internal/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	config    *config.Config
	api       *tgbotapi.BotAPI
	apiClient *ApiClient
}

type BotButtons struct {
	Info tgbotapi.KeyboardButton
	Help tgbotapi.KeyboardButton
}

// Метод для получения всех кнопок в виде слайса
func (b BotButtons) ToSlice() []tgbotapi.KeyboardButton {
	return []tgbotapi.KeyboardButton{
		b.Info,
		b.Help,
	}
}

func New(cfg *config.Config) (*Bot, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	log.Printf("Initializing bot with token: %s...", cfg.TelegramToken[:10])

	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	apiClient := NewApiClient(cfg.ApiSecretKey, cfg.ApiBaseURL)
	if apiClient == nil {
		return nil, fmt.Errorf("failed to create API client")
	}

	bot := &Bot{
		config:    cfg,
		api:       api,
		apiClient: apiClient,
	}

	log.Printf("Bot initialized successfully")
	return bot, nil
}

func (b *Bot) getUserLanguage(message *tgbotapi.Message) string {
	// Сначала пробуем получить язык из API
	response, err := b.apiClient.GetUserInfo(fmt.Sprintf("%d", message.From.ID))
	if err == nil {
		var userResponse handler.UserResponse
		if err := json.Unmarshal([]byte(response), &userResponse); err == nil {
			return userResponse.User.Settings.Language
		}
	}

	// Если не получилось, используем язык из настроек клиента
	if message.From.LanguageCode == "ru" || message.From.LanguageCode == "en" {
		return message.From.LanguageCode
	}

	return "ru" // По умолчанию русский
}

func (b *Bot) getClosestCameras(lat, lng float64, message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	var err error

	response, err := b.apiClient.getClosestCameras(lat, lng, 500)
	if err != nil {
		return msg, fmt.Errorf("error getting closest cameras: %w", err)
	}

	formattedCameras, err := handler.FormatCameras(response, b.getUserLanguage(message))
	if err != nil {
		return msg, fmt.Errorf("error formatting cameras: %w", err)
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, formattedCameras)

	return msg, nil
}

func (b *Bot) handleCoordinates(message *tgbotapi.Message, lat, lng float64) (tgbotapi.MessageConfig, error) {
	return b.getClosestCameras(lat, lng, message)
}

func (b *Bot) handleLocation(message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	return b.getClosestCameras(message.Location.Latitude, message.Location.Longitude, message)
}

func (b *Bot) handleMessage(message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	userLang := b.getUserLanguage(message)
	buttons := b.getButtons(userLang)

	var msg tgbotapi.MessageConfig
	var err error

	switch message.Text {
	case "/start":
		msg, err = b.handleStartCommand(message, userLang, buttons)
	case buttons.Info.Text, "/info":
		msg, err = b.handleInfoCommand(message)
	case buttons.Help.Text, "/help":
		msg, err = b.handleHelpCommand(message, userLang)
	default:
		msg = tgbotapi.NewMessage(message.Chat.ID, handler.GetTranslation(userLang, "unknown_cmd"))
	}

	if err != nil {
		return msg, fmt.Errorf("error handling message: %w", err)
	}

	return msg, nil
}

func (b *Bot) getButtons(userLang string) BotButtons {
	return BotButtons{
		Info: tgbotapi.NewKeyboardButton(handler.GetTranslation(userLang, "info_cmd")),
		Help: tgbotapi.NewKeyboardButton(handler.GetTranslation(userLang, "help_cmd")),
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message, userLang string, buttons BotButtons) (tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, handler.GetTranslation(userLang, "choose_action"))

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			buttons.ToSlice()...,
		),
	)
	keyboard.ResizeKeyboard = true

	msg.ReplyMarkup = keyboard

	return msg, nil
}

func (b *Bot) handleInfoCommand(message *tgbotapi.Message) (tgbotapi.MessageConfig, error) {
	userID := fmt.Sprintf("%d", message.From.ID)
	log.Printf("Getting info for user: %s", userID)

	response, err := b.apiClient.GetUserInfo(userID)
	if err != nil {
		return tgbotapi.MessageConfig{}, fmt.Errorf("error getting user info: %w", err)
	}

	formattedInfo, err := handler.FormatUserInfo(response)

	if err != nil {
		return tgbotapi.MessageConfig{}, fmt.Errorf("error formatting user info: %w", err)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, formattedInfo)

	return msg, nil
}

func (b *Bot) handleHelpCommand(message *tgbotapi.Message, userLang string) (tgbotapi.MessageConfig, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, handler.FormatCommandList(userLang))
	return msg, nil
}

func (b *Bot) Start() error {
	if b.api == nil {
		return fmt.Errorf("telegram API is nil")
	}

	if b.apiClient == nil {
		return fmt.Errorf("API client is nil")
	}

	log.Printf("Authorized on account %s", b.api.Self.UserName)

	// Configure updates
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Start receiving updates
	updates := b.api.GetUpdatesChan(updateConfig)

	var msg tgbotapi.MessageConfig
	var err error

	// Handle updates
	for update := range updates {
		if update.Message == nil {
			continue
		}

		lat, lng, parseError := handler.ParseCoordinates(update.Message.Text)

		switch {
		case update.Message.Location != nil:
			msg, err = b.handleLocation(update.Message)
		case parseError == nil:
			msg, err = b.handleCoordinates(update.Message, lat, lng)
		default:
			msg, err = b.handleMessage(update.Message)
		}

		if err != nil {
			log.Printf("Error handling message: %v", err)
		}

		if _, err := b.api.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}

	return nil
}
