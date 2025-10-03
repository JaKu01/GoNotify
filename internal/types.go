package internal

import (
	"gorm.io/gorm"
)

type WebPushSubscription struct {
	gorm.Model
	Endpoint     string `gorm:"unique"`
	Subscription string
}

type WebPushUnsubscriptionRequest struct {
	Endpoint string `json:"endpoint"`
}

type NotificationRequest struct {
	Subject     string `json:"subject"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

type TelegramRequest struct {
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	ApiToken string `json:"api_token"`
	ChatId   int64  `json:"chat_id"`
}

type NotificationResponse struct {
	Message string `json:"message"`
}
