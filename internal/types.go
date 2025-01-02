package internal

import (
	"gorm.io/gorm"
)

type WebPushSubscription struct {
	gorm.Model
	Endpoint     string `gorm:"unique"`
	AuthSecret   string
	P256dhKey    string
}

type WebPushUnsubscriptionRequest struct {
	Endpoint string `json:"endpoint"`
}

type NotificationRequest struct {
	Subject     string `json:"subject"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

type NotificationResponse struct {
	Message string `json:"message"`
}
