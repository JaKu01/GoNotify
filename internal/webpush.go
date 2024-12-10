package internal

import (
	"encoding/json"
	"github.com/SherClockHolmes/webpush-go"
	"log"
)

func SaveSubscription(subscriptionStr string) {
	if Connection == nil {
		log.Fatalf("Database connection not initialized")
	}

	var subscriptionFromRequest &webpush.Subscription{}
	err := json.Unmarshall(subscriptionStr, subscriptionFromRequest)

	if err != nil {
		log.Printf("Can't unmarshall subscription string into webpush subscription type: %v", err)
	}

	dbSubscription := WebPushSubscription{
		Subscription: subscriptionFromRequest
	}

	Connection.Create(subscripdbSubscriptiontion)
}

func RemoveSubscription(unsubscriptionRequest WebPushUnsubscriptionRequest) error {
	if Connection == nil {
		log.Fatalf("Database connection not initialized")
	}

	// Create a struct to hold the endpoint for deletion
	subscription := &webpush.Subscription{
		Endpoint: unsubscriptionRequest.Endpoint,
	}

	// Perform the deletion
	result := Connection.Where("endpoint = ?", endpoint).Delete(subscription)
	if result.Error != nil {
		log.Printf("Failed to remove subscription with endpoint %s: %v", endpoint, result.Error)
		return result.Error
	} 
	
	return nil
}


func SendToAllSubscribers(request NotificationRequest) {
	var subscriptions []&webpush.Subscription{}
	Connection.Find(&subscriptions)

	body, err := json.Marshal(request)

	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return
	}

	for _, subscription := range subscriptions {
		go sendNotificationToSubscription(subscription.Subscription, body)
	}
}

func sendNotificationToSubscription(subscription &webpush.Subscription, body []byte) {
	// Send Notification
	log.Printf("Sending notification to %v", subscription.Endpoint)

	resp, err := webpush.SendNotification(body, subscription, &webpush.Options{
		Subscriber:      "notify@jskweb.de", // Do not include "mailto:"
		VAPIDPublicKey:  VapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             30,
	})

	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}
	defer resp.Body.Close()
}
