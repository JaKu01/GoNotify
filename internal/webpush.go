package internal

import (
	"encoding/json"
	"github.com/SherClockHolmes/webpush-go"
	"gorm.io/gorm"
	"log"
)

type WebPushSubscription struct {
	gorm.Model
	Subscription string
}

func SaveSubscription(subscription string) {
	if Connection == nil {
		log.Fatalf("Database connection not initialized")
	}

	Connection.Create(&WebPushSubscription{Subscription: subscription})
}

func SendToAllSubscribers(body []byte) {
	var subscriptions []WebPushSubscription
	Connection.Find(&subscriptions)

	for _, subscription := range subscriptions {
		sendNotificationToSubscription(subscription.Subscription, body)
	}
}

func sendNotificationToSubscription(subscription string, body []byte) {
	// Send Notification
	s := &webpush.Subscription{}
	err := json.Unmarshal([]byte(subscription), s)
	if err != nil {
		log.Printf("Error unmarshalling subscription: %v", err)
	}

	log.Printf("Sending notification to %v", s)
	log.Printf("Raw value notification to %v", subscription)

	resp, err := webpush.SendNotification(body, s, &webpush.Options{
		Subscriber:      "example@example.com", // Do not include "mailto:"
		VAPIDPublicKey:  VapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             30,
	})
	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}
	defer resp.Body.Close()
}
