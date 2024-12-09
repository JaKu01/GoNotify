package internal

import "gorm.io/gorm"

type WebPushSubscription struct {
	gorm.Model
	Subscription string
}

type NotificationRequest struct {
	Subject     string `json:"subject"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

type NotificationResponse struct {
	Message string `json:"message"`
}
