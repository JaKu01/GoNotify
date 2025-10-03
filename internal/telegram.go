package internal

import (
	"context"
	"errors"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"os"
)

func SendTelegramMessage(apiToken string, chatId int64, subject string, message string) error {
	token, err := getTelegramAPIToken(apiToken)
	if err != nil {
		return err
	}

	telegramService, err := telegram.New(token)
	if err != nil {
		return err
	}

	telegramService.AddReceivers(chatId)

	notify.UseServices(telegramService)

	err = notify.Send(context.Background(), subject, message)
	return err
}

func getTelegramAPIToken(token string) (string, error) {
	if token != "" {
		return token, nil
	}
	envToken := os.Getenv("TELEGRAM_API_TOKEN")
	if envToken == "" {
		return "", errors.New("TELEGRAM_API_TOKEN is empty")
	}
	return envToken, nil
}
